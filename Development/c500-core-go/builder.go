package domain

import (
	"time"
)

// ProfileData holds the custom "Geocities-style" content for a builder's public profile.
// This data is stored nested within the Builder document in Firestore.
type ProfileData struct {
	// Raw HTML submitted by the user.
	// SECURITY CRITICAL: This must be securely sanitized before rendering on the web client.
	HTML string `json:"html" firestore:"html"`

	// Raw CSS submitted by the user.
	// SECURITY CRITICAL: This must be securely sanitized/scoped before rendering.
	CSS string `json:"css" firestore:"css"`
}

// Builder represents a user in the C500 ecosystem.
// While we call them "Builders," this struct represents buyers as well.
// The struct tags (json:"...") define how data is mapped when talking to the API or Firestore.
type Builder struct {
	// ID is the unique identifier used as the Firestore document ID.
	// In our design, this is identical to the DiscordID.
	ID string `json:"id" firestore:"id"`

	// DiscordID is the user's unique snowflake ID from Discord.
	// This is our primary way of identifying users across services.
	DiscordID string `json:"discord_id" firestore:"discord_id" binding:"required"`

	// DisplayName is their current Discord username (cached for UI display).
	DisplayName string `json:"display_name" firestore:"display_name"`

	// StripeAccountID is the Connected Express Account ID (e.g., "acct_1GSE7...").
	// This will be empty/nil if they have not completed seller onboarding.
	// 'omitempty' means it won't be sent in JSON if it's empty.
	StripeAccountID string `json:"stripe_account_id,omitempty" firestore:"stripe_account_id,omitempty"`

	// IsVerifiedBuilder is a flag set by community admins allowing access to selling commands.
	IsVerifiedBuilder bool `json:"is_verified_builder" firestore:"is_verified_builder"`

	// Profile contains their custom public profile customizations.
	Profile ProfileData `json:"profile_data" firestore:"profile_data"`

	// Standard timestamps for record keeping.
	CreatedAt time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt time.Time `json:"updated_at" firestore:"updated_at"`
}

// NewBuilder is a constructor function to create a basic, empty user instance.
// It sets the default timestamps and ensures IDs are synced.
func NewBuilder(discordID, displayName string) *Builder {
	now := time.Now().UTC()
	return &Builder{
		ID:                discordID, // We use the DiscordID as the primary document key
		DiscordID:         discordID,
		DisplayName:       displayName,
		IsVerifiedBuilder: false, // Users start unverified by default
		CreatedAt:         now,
		UpdatedAt:         now,
		// Initialize empty profile data so it's not nil in the DB
		Profile: ProfileData{
			HTML: "",
			CSS:  "",
		},
	}
}

// CanSell is a domain helper method to check if all requirements to create a drop are met.
// Business logic shouldn't rely on just one flag; they need Stripe set up too.
func (b *Builder) CanSell() bool {
	// They must be marked verified by admins AND have finished Stripe setup.
	return b.IsVerifiedBuilder && b.StripeAccountID != ""
}
