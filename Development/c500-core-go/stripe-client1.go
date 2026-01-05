package stripe

import (
	// ... previous imports go here
	"fmt"

	"github.com/stripe/stripe-go/v74"
	// "github.com/stripe/stripe-go/v74/checkout/session" // Already imported
	"github.com/stripe/stripe-go/v74/transfer" // NEW: Needed for payouts
	// "c500-core-go/internal/domain" // Already imported
)

// ... (Client struct and CreateCheckoutSession method remain above) ...

// =================================================================
// NEW METHOD: ReleaseEscrowFunds
// This fulfills the interface needed by fulfillment_service.go
// =================================================================

// ReleaseEscrowFunds moves money from the platform's balance to the seller's connected account.
func (c *Client) ReleaseEscrowFunds(ctx context.Context, paymentIntentID, destinationStripeAcctID string, amountCents int64) error {

	// 1. Configure the Transfer parameters.
	params := &stripe.TransferParams{
		// Amount to transfer in cents.
		Amount: stripe.Int64(amountCents),
		// Currency must match the original payment.
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		// The destination is the seller's connected Express account ID.
		Destination: stripe.String(destinationStripeAcctID),
		// CRITICAL: Link this transfer to the original charge.
		// This ensures that if the buyer charges back, the funds are pulled from
		// the seller's account, not the platform's main balance.
		SourceTransaction: stripe.String(paymentIntentID),
	}

	// 2. Set up Idempotency.
	// We generate a unique key for this specific payout operation.
	// A simple way is to combine the source payment ID and destination ID.
	idempotencyKey := fmt.Sprintf("payout_%s_%s", paymentIntentID, destinationStripeAcctID)
	params.IdempotencyKey = stripe.String(idempotencyKey)

	// 3. Perform the network call to Stripe.
	// The .New() function creates the transfer. If successful, funds move instantly.
	_, err := transfer.New(params)
	if err != nil {
		// The error from Stripe is usually detailed. We wrap it to give context.
		// e.g. "stripe transfer failed: card_error: Your balance is insufficient."
		return fmt.Errorf("stripe transfer api failed: %w", err)
	}

	// If err is nil, the money has officially moved.
	return nil
}
