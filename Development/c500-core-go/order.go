package domain

import (
	"time"
)

// EscrowStatus defines the state of the funds held by Stripe.
type EscrowStatus string

const (
	EscrowHeld     EscrowStatus = "held"     // Funds secured, waiting for fulfillment.
	EscrowReleased EscrowStatus = "released" // Seller fulfilled, funds paid out.
	EscrowRefunded EscrowStatus = "refunded" // Something went wrong, buyer got money back.
)

// Order represents a finalized, paid-for transaction.
type Order struct {
	ID string `json:"id" firestore:"id"`

	// Foreign Keys linking the core entities.
	DropID          string `json:"drop_id" firestore:"drop_id"`
	BuyerDiscordID  string `json:"buyer_discord_id" firestore:"buyer_discord_id"`
	SellerDiscordID string `json:"seller_discord_id" firestore:"seller_discord_id"`

	// Financial record.
	PriceInCents int64        `json:"price_in_cents" firestore:"price_in_cents"`
	EscrowStatus EscrowStatus `json:"escrow_status" firestore:"escrow_status"`

	// The crucial link back to Stripe for potential refunds or disputes.
	StripePaymentIntentID string `json:"stripe_payment_intent_id" firestore:"stripe_payment_intent_id"`

	// Fulfillment details (added later by the seller).
	TrackingNumber string `json:"tracking_number,omitempty" firestore:"tracking_number,omitempty"`
	Carrier        string `json:"carrier,omitempty" firestore:"carrier,omitempty"`
	VODLink        string `json:"vod_link,omitempty" firestore:"vod_link,omitempty"`

	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

// NewOrder is a helper to create a new order object with default "held" status.
func NewOrder(dropID, buyerID, sellerID, paymentIntentID string, price int64) *Order {
	now := time.Now().UTC()
	// In a real app, use a UUID library here: id := uuid.New().String()
	id := "order_" + dropID + "_" + buyerID // Placeholder ID generation
	return &Order{
		ID:                    id,
		DropID:                dropID,
		BuyerDiscordID:        buyerID,
		SellerDiscordID:       sellerID,
		PriceInCents:          price,
		EscrowStatus:          EscrowHeld, // Funds start as held.
		StripePaymentIntentID: paymentIntentID,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
}
