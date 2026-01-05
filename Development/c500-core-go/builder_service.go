package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"c500-core-go/internal/domain"
	// We would likely have a custom error package, e.g., "c500-core-go/pkg/errs"
)

var (
	ErrBuilderNotFound = errors.New("builder not found")
	ErrStripeError     = errors.New("stripe integration error")
)

// BuilderRepository defines the interface used to persist Builder data.
// This service doesn't care *how* it's saved (Firestore, SQL, memory), only that these methods exist.
// This implementation lives in internal/database/firestore.go
type BuilderRepository interface {
	GetByID(ctx context.Context, discordID string) (*domain.Builder, error)
	Create(ctx context.Context, builder *domain.Builder) error
	UpdateStripeID(ctx context.Context, discordID, stripeAccountID string) error
	UpdateProfileData(ctx context.Context, discordID string, profile domain.ProfileData) error
}

// StripeIntegration defines the interface for talking to the Stripe API.
// This implementation lives in internal/integrations/stripe/client.go
type StripeIntegration interface {
	// CreateAccountLink generates a one-time URL for a user to onboard with Stripe Express.
	CreateAccountLink(discordID string) (string, error)
}

// BuilderService is the concrete implementation containing business logic.
type builderService struct {
	repo   BuilderRepository
	stripe StripeIntegration
}

// NewBuilderService is the constructor used in main.go to inject dependencies.
func NewBuilderService(repo BuilderRepository, stripe StripeIntegration) *builderService {
	return &builderService{
		repo:   repo,
		stripe: stripe,
	}
}

// --- Business Logic Methods ---

// GetOrCreateBuilder gets an existing user or creates a skeleton record if it's their first time.
// Used when a user runs any command in Discord to ensure they exist in our DB.
func (s *builderService) GetOrCreateBuilder(ctx context.Context, discordID, displayName string) (*domain.Builder, error) {
	// 1. Try to find them first
	builder, err := s.repo.GetByID(ctx, discordID)
	if err == nil {
		// Found them, return existing data
		return builder, nil
	}

	// If err is generic, return it. If it's specifically "Not Found", proceed to create.
	// (Simplified error checking for this example)
	if !errors.Is(err, ErrBuilderNotFound) {
		 return nil, fmt.Errorf("error searching for builder: %w", err)
	}

	// 2. Not found, so create a new domain object
	newBuilder := domain.NewBuilder(discordID, displayName)

	// 3. Persist it
	if err := s.repo.Create(ctx, newBuilder); err != nil {
		return nil, fmt.Errorf("failed to create new builder record: %w", err)
	}

	return newBuilder, nil
}

// GetStripeOnboardingLink orchestrates the process of letting a user become a seller.
func (s *builderService) GetStripeOnboardingLink(ctx context.Context, discordID string) (string, error) {
	// 1. Ensure user exists in our DB first.
	builder, err := s.repo.GetByID(ctx, discordID)
	if err != nil {
		// In reality, you might auto-create them here, but for now, fail if they don't exist.
		return "", ErrBuilderNotFound
	}

	// 2. If they already have a Stripe ID connected, they don't need an onboarding link.
	if builder.StripeAccountID != "" {
		// Depending on business logic, maybe return a "dashboard login" link instead.
		// For now, just note they are already set up.
		return "", errors.New("user is already connected to Stripe")
	}

	// 3. Call Stripe integration to generate the secure, ephemeral link.
	url, err := s.stripe.CreateAccountLink(discordID)
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrStripeError, err)
	}

	// 4. Return the URL to be sent privately via Discord DM.
	return url, nil
}

// CompleteStripeOnboarding is called via Webhook when Stripe tells us someone finished setup.
func (s *builderService) CompleteStripeOnboarding(ctx context.Context, discordID, newStripeAccountID string) error {
	// The repository handles the database specific logic of finding the document and patching just this field.
	// We update the UpdatedAt timestamp here as part of the business logic.
	err := s.repo.UpdateStripeID(ctx, discordID, newStripeAccountID)
	if err != nil {
		return fmt.Errorf("failed to link stripe account ID in DB: %w", err)
	}
	return nil
}

// UpdateBuilderProfile handles saving the custom HTML/CSS.
// CRITICAL: This is where security sanitization must happen.
func (s *builderService) UpdateBuilderProfile(ctx context.Context, discordID string, rawHTML, rawCSS string) (*domain.Builder, error) {
	// 1. Validate inputs (e.g., enforce character limits)
	if len(rawHTML) > 10000 || len(rawCSS) > 5000 {
		return nil, errors.New("profile content exceeds size limits")
	}

	// 2. TODO: SECURITY SANITIZATION
	// Before saving, we MUST pass rawHTML through a strict sanitizer library like 'bluemonday'.
	// We strip dangerous tags (<script>, <iframe>, onmouseover attributes, etc.).
	// sanitizedHTML := bluemonday.UGCPolicy().Sanitize(rawHTML)
	// For this example code, we assume raw is okay, but DO NOT do this in prod.
	sanitizedHTML := rawHTML
	sanitizedCSS := rawCSS // CSS also needs validation against malicious external url() calls.

	// 3. Create the data update struct
	newData := domain.ProfileData{
		HTML: sanitizedHTML,
		CSS:  sanitizedCSS,
	}

	// 4. Save to repository
	if err := s.repo.UpdateProfileData(ctx, discordID, newData); err != nil {
		return nil, fmt.Errorf("failed to update profile data: %w", err)
	}

	// 5. Fetch and return the updated builder info
	return s.repo.GetByID(ctx, discordID)
}

// GetPublicBuilder is used by the Web frontend to render profiles.
// It might apply different logic than GetByID (e.g., only showing verified builders).
func (s *builderService) GetPublicBuilder(ctx context.Context, discordID string) (*domain.Builder, error) {
	builder, err := s.repo.GetByID(ctx, discordID)
	if err != nil {
		return nil, err
	}

	// Optional business logic: Maybe we don't want to show profiles for banned users
	// or users who aren't verified builders yet.
	// if !builder.IsVerifiedBuilder { ... }

	return builder, nil
}
