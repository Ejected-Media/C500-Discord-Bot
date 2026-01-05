package service

import (
	"context"
	"errors"
	"fmt"
	// "time" // removed unused import

	"c500-core-go/internal/domain"
)

// ... (previous error definitions)
var (
    // Add a new error for this specific failure case.
	ErrOrderCreationFailed = errors.New("failed to create final order record")
)

// DropRepository interface... (stays the same)
type DropRepository interface {
	GetDropByID(ctx context.Context, dropID string) (*domain.Drop, error)
	UpdateDropStatus(ctx context.Context, dropID string, newStatus domain.DropStatus) error
}

// NEW: OrderRepository defines how we save orders to the DB.
// This will be implemented in internal/database/firestore.go later.
type OrderRepository interface {
	CreateOrder(ctx context.Context, order *domain.Order) error
}

// StripeIntegration interface... (stays the same)

// CheckoutService holds the business logic dependencies.
type checkoutService struct {
	dropRepo  DropRepository
	orderRepo OrderRepository // NEW dependency added here.
	stripe    StripeIntegration
}

// NewCheckoutService constructor updated to accept the new repo.
func NewCheckoutService(dr DropRepository, or OrderRepository, si StripeIntegration) *checkoutService {
	return &checkoutService{
		dropRepo:  dr,
		orderRepo: or,
		stripe:    si,
	}
}

// ==========================================
// Business Logic
// ==========================================

// CreateCheckoutSession... (stays the same)

// NEW METHOD: ProcessSuccessfulPayment is called by the Webhook Handler.
// This is the "finisher" that makes the sale official in our system.
func (s *checkoutService) ProcessSuccessfulPayment(ctx context.Context, dropID, buyerDiscordID, stripePaymentIntentID string) error {
	// 1. Fetch the Drop details. We need the price and seller ID to create the order.
	drop, err := s.dropRepo.GetDropByID(ctx, dropID)
	if err != nil {
		// This is a critical data inconsistency error. Stripe has money for a drop we can't find.
		return fmt.Errorf("CRITICAL: failed to fetch drop %s for order finalization: %w", dropID, err)
	}

	// 2. Prepare the new Order domain object.
	newOrder := domain.NewOrder(
		drop.ID,
		buyerDiscordID,
		drop.SellerDiscordID,
		stripePaymentIntentID,
		drop.PriceInCents,
	)

	// 3. CRITICAL DB UPDATES.
	// In a production Firestore implementation, these two calls MUST be wrapped
	// in a single Transaction to ensure they either both succeed or both fail.
	// For this example, we show them as sequential steps.

	// 3a. Finalize the listing status from 'pending' to 'sold'.
	err = s.dropRepo.UpdateDropStatus(ctx, dropID, domain.StatusSold)
	if err != nil {
		// Danger: Stripe has money, but our DB still thinks the item is pending.
		return fmt.Errorf("failed to mark drop as sold: %w", err)
	}

	// 3b. Save the permanent Order record.
	err = s.orderRepo.CreateOrder(ctx, newOrder)
	if err != nil {
		// Major Danger: Drop is marked sold, but we have no record of who bought it.
		// This requires manual admin intervention to fix.
		return fmt.Errorf("%w: %v", ErrOrderCreationFailed, err)
	}

	// 4. (Optional Future Step) Trigger notifications.
	// e.g. Publish an event that the Python bot listens to:
	// events.Publish("order.created", newOrder) -> Bot DMs buyer & seller.

	return nil
}
