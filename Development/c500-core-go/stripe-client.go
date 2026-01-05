package stripe

import (
	"context"
	"fmt"

	// The official Stripe Go SDK
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"

	"c500-core-go/internal/domain"
)

// Client is our wrapper around the Stripe API.
type Client struct {
	// In a real app, you might hold per-request context here if needed.
	// Stripe's Go SDK uses a global API key set in main.go, so this struct is simple.
}

// NewClient creates a new instance.
// The actual secret key is set globally in main.go using stripe.Key = "sk_test_..."
func NewClient() *Client {
	return &Client{}
}

// CreateCheckoutSession fulfills the interface defined in the Service layer.
func (c *Client) CreateCheckoutSession(ctx context.Context, drop *domain.Drop, buyerDiscordID string) (string, string, error) {

	// 1. Define where the user goes after they finish on the Stripe page.
	// These point to our Go Web Frontend (`c500-web-go`).
	// In production, these domains would come from config/env vars.
	successURL := "http://localhost:3000/checkout/success?session_id={CHECKOUT_SESSION_ID}"
	cancelURL := "http://localhost:3000/checkout/cancel"

	// 2. Construct the parameters for Stripe API.
	params := &stripe.CheckoutSessionParams{
		// We use 'payment' mode for one-time purchases.
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		// Tell Stripe where to redirect upon completion.
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
		// Define what is being bought.
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(drop.Title),
						// We could add the main image URL here too so it shows on checkout page
						// Images: []*string{stripe.String(drop.ImageURLs[0])},
					},
					// IMPORTANT: Stripe expects amounts in cents (e.g., 45000 for $450.00).
					// Our domain model already stores it this way, so it's a direct mapping.
					UnitAmount: stripe.Int64(drop.PriceInCents),
				},
				Quantity: stripe.Int64(1),
			},
		},
		// =====================================================================
		// CRITICAL SECTION: Metadata
		// This hidden data is attached to the transaction. Stripe doesn't use it,
		// but it sends it back to us in the 'payment_intent.succeeded' webhook.
		// This is how we link the money back to our database records.
		// =====================================================================
		Metadata: map[string]string{
			"drop_id":           drop.ID,
			"buyer_discord_id":  buyerDiscordID,
			"seller_discord_id": drop.SellerDiscordID,
			// We track drop type so we know if fulfillment needs shipping or a VOD link later.
			"drop_type":         drop.Type,
		},
	}

	// 3. Perform the network call to Stripe's servers.
	s, err := session.New(params)
	if err != nil {
		// This error means something went wrong talking to Stripe (auth error, API down, etc.)
		// We wrap it so our logs know exactly where it failed.
		return "", "", fmt.Errorf("stripe API call failed: %w", err)
	}

	// 4. Return the public URL for the frontend, and the private Session ID for our records.
	return s.URL, s.ID, nil
}
