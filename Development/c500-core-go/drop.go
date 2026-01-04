package domain

import "time"

// DropType Enum to differentiate sales models
type DropType string

const (
	DropTypeCommission DropType = "commission"
	DropTypeRTS        DropType = "ready_to_ship"
)

// DropStatus Enum for lifecycle state
type DropStatus string

const (
	DropStatusAvailable DropStatus = "available"
	DropStatusPending   DropStatus = "pending_payment"
	DropStatusSold      DropStatus = "sold"
)

// Drop represents a single listing in the marketplace.
// The struct tags (e.g., `json:"title"`) define how it looks in API calls
// and how it's stored in Firestore.
type Drop struct {
	ID              string                 `json:"id" firestore:"id"` // Firestore Document ID
	SellerDiscordID string                 `json:"seller_discord_id" firestore:"seller_discord_id"`
	Title           string                 `json:"title" firestore:"title" binding:"required"`
	Price           float64                `json:"price" firestore:"price" binding:"required,gt=0"`
	ImageURL        string                 `json:"image_url" firestore:"image_url" binding:"required,url"`
	Description     string                 `json:"description" firestore:"description"`
	Type            DropType               `json:"type" firestore:"type"`
	Status          DropStatus             `json:"status" firestore:"status"`
	StructuredData  map[string]interface{} `json:"structured_data" firestore:"structured_data"`
	CreatedAt       time.Time              `json:"created_at" firestore:"created_at"`
}

// CreateDropRequest defines the exact JSON payload the Python Bot must send
type CreateDropRequest struct {
	SellerDiscordID string   `json:"seller_discord_id" binding:"required"`
	Title           string   `json:"title" binding:"required,max=100"`
	Price           float64  `json:"price" binding:"required,gt=0"`
	// ... other fields ...
}
