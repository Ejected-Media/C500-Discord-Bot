package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"c500-core-go/internal/domain"
)

// Define the custom errors the Handler layer is expecting.
var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrUnauthorizedSeller    = errors.New("user is not the seller of this order")
	ErrOrderAlreadyFulfilled = errors.New("order is not in held status")
	ErrStripePayoutFailed    = errors.New("failed to release funds via stripe")
)

// OrderRepository defines DB operations needed for fulfillment.
// (Implemented in internal/database/firestore.go)
type OrderRepository interface {
	GetOrderByID(ctx context.Context, orderID string) (*domain.Order, error)
	// UpdateOrderFulfillment performs a partial update on specific fields.
	UpdateOrderFulfillment(ctx context.Context, orderID string, updates map[string]interface{}) error
}

// Use existing BuilderRepository interface to fetch seller's Stripe ID.
// (Defined in builder_service.go)
// type BuilderRepository interface { ...GetByID... }

// StripeIntegration needs a new method for payouts.
// (Implemented in internal/integrations/stripe/client.go)
type StripeIntegration interface {
	// ... existing CreateCheckoutSession ...

	// ReleaseEscrowFunds transfers held funds to the destination account.
	ReleaseEscrowFunds(ctx context.Context, paymentIntentID, destinationStripeAcctID string, amountCents int64) error
}

// FulfillmentService is the concrete implementation.
type fulfillmentService struct {
	orderRepo   OrderRepository
	builderRepo BuilderRepository
	stripe      StripeIntegration
}

// NewFulfillmentService constructor.
func NewFulfillmentService(or OrderRepository, br BuilderRepository, si StripeIntegration) *fulfillmentService {
	return &fulfillmentService{
		orderRepo:   or,
		builderRepo: br,
		stripe:      si,
	}
}

// ==========================================
// Business Logic
// ==========================================

// FulfillOrderWithShipping handles RTS orders.
func (s *fulfillmentService) FulfillOrderWithShipping(ctx context.Context, orderID, sellerDiscordID, tracking, carrier string) error {
	// Prepare the specific updates for shipping
	updates := map[string]interface{}{
		"tracking_number": tracking,
		"carrier":         carrier,
	}
	// Call shared logic helper
	return s.processFulfillment(ctx, orderID, sellerDiscordID, updates)
}

// FulfillOrderWithVOD handles Commission orders.
func (s *fulfillmentService) FulfillOrderWithVOD(ctx context.Context, orderID, sellerDiscordID, vodURL string) error {
	// Prepare the specific updates for VOD
	updates := map[string]interface{}{
		"vod_link": vodURL,
	}
	// Call shared logic helper
	return s.processFulfillment(ctx, orderID, sellerDiscordID, updates)
}

// processFulfillment is the shared Core Logic responsible for checks and payouts.
func (s *fulfillmentService) processFulfillment(ctx context.Context, orderID, requestingSellerID string, specificUpdates map[string]interface{}) error {
	// 1. Fetch the Order
	order, err := s.orderRepo.GetOrderByID(ctx, orderID)
	if err != nil {
		// Assume repo maps DB not found to standard error, or check here.
		return fmt.Errorf("failed to fetch order: %w", ErrOrderNotFound)
	}

	// 2. SECURITY CHECK: Is the requester actually the seller?
	if order.SellerDiscordID != requestingSellerID {
		return ErrUnauthorizedSeller
	}

	// 3. STATE CHECK: Is the money actually held right now?
	if order.EscrowStatus != domain.EscrowHeld {
		return ErrOrderAlreadyFulfilled
	}

	// 4. Fetch Seller's Stripe destination account ID.
	// We need to look up the builder profile to get this.
	seller, err := s.builderRepo.GetByID(ctx, order.SellerDiscordID)
	if err != nil || seller.StripeAccountID == "" {
		return fmt.Errorf("critical: cannot find seller stripe account for payout: %v", err)
	}

	// 5. THE BIG MOMENT: Call Stripe to release the funds.
	// We move money from the PaymentIntent to the Seller's connected account.
	err = s.stripe.ReleaseEscrowFunds(ctx, order.StripePaymentIntentID, seller.StripeAccountID, order.PriceInCents)
	if err != nil {
		// This is bad. Database says held, but Stripe refused to pay out.
		// Log heavily. Do NOT update DB status. Let user try again or contact support.
		return fmt.Errorf("%w: %v", ErrStripePayoutFailed, err)
	}

	// 6. Update Database Record.
	// Money moved successfully. Now we update our records.
	// Add standard fields to the specific updates map.
	specificUpdates["escrow_status"] = domain.EscrowReleased
	specificUpdates["updated_at"] = time.Now().UTC()

	err = s.orderRepo.UpdateOrderFulfillment(ctx, orderID, specificUpdates)
	if err != nil {
		// DANGER: Money moved in Stripe, but DB update failed.
		// System is in inconsistent state. Needs high priority log/alert.
		return fmt.Errorf("CRITICAL: Funds released but DB update failed for order %s: %v", orderID, err)
	}

	return nil
}
