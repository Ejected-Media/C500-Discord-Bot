package domain

import (
	"time"
)

// DropStatus defines the lifecycle stages of a listing.
type DropStatus string

const (
	StatusDraft     DropStatus = "draft"      // Seller is still editing
	StatusAvailable DropStatus = "available"  // Live in the shop
	StatusPending   DropStatus = "pending"    // Buyer is in checkout flow (locked)
	StatusSold      DropStatus = "sold"       // Transaction complete
)

// Drop represents a marketplace listing.
type Drop struct {
	ID string `json:"id" firestore:"id"`

	// SellerDiscordID links this drop back to the user who created it.
	// This is a foreign key reference to the 'users' collection.
	SellerDiscordID string `json:"seller_discord_id" firestore:"seller_discord_id"`

	Title       string  `json:"title" firestore:"title" binding:"required,min=5"`
	Description string  `json:"description" firestore:"description"`
	// PriceInCents is stored as an integer to avoid floating point math errors.
	// e.g., $450.00 is stored as 45000.
	PriceInCents int64 `json:"price_in_cents" firestore:"price_in_cents" binding:"required,gt=0"`

	// Type is either "rts" (ready-to-ship) or "commission".
	Type string `json:"type" firestore:"type"`

	// Status tracks if it can be bought.
	Status DropStatus `json:"status" firestore:"status"`

	// Images would likely be a slice of URLs pointing to Cloud Storage buckets.
	ImageURLs []string `json:"image_urls" firestore:"image_urls"`

	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

// Helper to create a new, ready-to-save drop
func NewDrop(sellerID, title string, priceCents int64, dropType string) *Drop {
    now := time.Now().UTC()
    // In reality, ID would be generated here using a UUID library (e.g. google/uuid)
    // id := uuid.New().String()
    id := "temp_placeholder_id" 
    return &Drop{
        ID:              id,
        SellerDiscordID: sellerID,
        Title:           title,
        PriceInCents:    priceCents,
        Type:            dropType,
        Status:          StatusDraft, // Start as draft by default
        CreatedAt:       now,
        UpdatedAt:       now,
    }
}
