package service

import (
	"context"
	"errors"
	"fmt"

	"c500-core-go/internal/domain"
)

// Define custom errors that the Handler layer will look for to determine HTTP status codes.
var (
	ErrDropNotFound     = errors.New("drop not found")
	ErrDropNotAvailable = errors.New("drop is not available for purchase")
	ErrStripeFailure    = errors.New("upstream stripe api failure")
)

// DropRepository defines the DB operations we need for checkout.
// Implemented in internal/database/firestore.go
type DropRepository interface {
	GetDropByID(ctx context.Context, dropID string) (*domain.Drop, error)
	UpdateDropStatus(ctx context.Context, dropID string, newStatus domain.DropStatus) error
	// In a real app, we'd also need a method to create an Order record here.
}

// StripeIntegration defines how we talk to Stripe.
// Implemented in internal/integrations/stripe/client.go
type StripeIntegration interface {
	// CreateCheckoutSession generates the url and the session ID.
	// We pass the whole drop object so Stripe knows the title and price.
	CreateCheckoutSession(ctx context.Context, drop *domain.Drop, buyerDiscordID string) (string, string, error)
}

// CheckoutService is the concrete implementation holding business logic.
type checkoutService struct {
	dropRepo DropRepository
	stripe   StripeIntegration
	// In the future, this would also need an OrderRepository.
}

// NewCheckoutService constructor used in main.go.
func NewCheckoutService(dr DropRepository, si StripeIntegration) *checkoutService {
	return &checkoutService{
		dropRepo: dr,
		stripe:   si,
	}
}

// ==========================================
// Business Logic
// ==========================================

// CreateCheckoutSession is the conductor for starting a purchase.
func (s *checkoutService) CreateCheckoutSession(ctx context.Context, dropID, buyerDiscordID string) (string, error) {
	// 1. Fetch the Drop info from the Database
	drop, err := s.dropRepo.GetDropByID(ctx, dropID)
	if err != nil {
		// We assume the repo returns a specific error if not found
		// In reality you'd check errors.Is(err, database.ErrNotFound)
		return "", fmt.Errorf("failed to fetch drop: %w", ErrDropNotFound)
	}

	// 2. CRITICAL BUSINESS RULE: Is it actually available?
	// This prevents race conditions where two people click buy at the exact same second.
	if drop.Status != domain.StatusAvailable {
		return "", ErrDropNotAvailable
	}

	// 3. Call Stripe to generate the payment link.
	// We pass buyer ID so Stripe attaches it as metadata to the transaction.
	checkoutURL, stripeSessionID, err := s.stripe.CreateCheckoutSession(ctx, drop, buyerDiscordID)
	if err != nil {
		// Log the actual error internally, return generic error upstream
		return "", fmt.Errorf("stripe session creation failed: %w", ErrStripeFailure)
	}

	// 4. Lock the Drop in the Database.
	// Mark it as 'pending' so no one else can start a checkout for it.
	// NOTE: In a production environment, steps 1, 2, and 4 should ideally happen
	// inside a Firestore Transaction to ensure perfect atomicity.
	err = s.dropRepo.UpdateDropStatus(ctx, dropID, domain.StatusPending)
	if err != nil {
		// DANGER ZONE: We created a Stripe link but failed to update our DB.
		// A user might pay for an item the system thinks is still available.
		//
		// Mitigation strategy for production code:
		// 1. Log a high-priority alert.
		// 2. Attempt to call Stripe API immediately to expire/cancel the sessionID we just created.
		return "", fmt.Errorf("CRITICAL: failed to mark drop pending after stripe link generation: %w", err)
	}

	// TODO: Step 4b - Create an initial 'Order' record in DB with status 'awaiting_payment'.

	// 5. Return the URL to be sent to the user.
	return checkoutURL, nil
}
