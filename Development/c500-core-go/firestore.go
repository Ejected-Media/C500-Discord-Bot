package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"cloud.google.com/go/firestore" // The official Google SDK
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"c500-core-go/internal/domain"
	"c500-core-go/internal/service" // Import service to reference custom errors like ErrBuilderNotFound
)

const (
	usersCollection = "users"
	// dropsCollection = "drops"
)

// FirestoreClient is our concrete implementation of the repository interfaces.
// It holds the actual pointer to the Google SDK client.
type FirestoreClient struct {
	client *firestore.Client
}

// NewFirestoreClient creates the connection to Google Cloud.
// This is called once in main.go at startup.
func NewFirestoreClient(ctx context.Context, projectID string) (*FirestoreClient, error) {
	// This call performs the actual network handshake with Google Cloud.
	// It relies on GOOGLE_APPLICATION_CREDENTIALS env var available in the cloud environment.
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to create firestore client: %w", err)
	}

	return &FirestoreClient{client: client}, nil
}

// Close ensures the connection is shut down properly when the app stops.
func (f *FirestoreClient) Close() error {
	return f.client.Close()
}

// =================================================================
// BuilderRepository Implementation
// These methods satisfy the interface defined in builder_service.go
// =================================================================

// GetByID fetches a single document from the "users" collection.
func (f *FirestoreClient) GetByID(ctx context.Context, discordID string) (*domain.Builder, error) {
	// 1. Define the path to the document: users/{discordID}
	docRef := f.client.Collection(usersCollection).Doc(discordID)

	// 2. Perform the network call to get the snapshot
	docSnap, err := docRef.Get(ctx)
	if err != nil {
		// 3. Translate Firestore-specific "not found" error into our domain "not found" error.
		// This ensures the service layer doesn't need to import gRPC packages.
		if status.Code(err) == codes.NotFound {
			return nil, service.ErrBuilderNotFound
		}
		return nil, fmt.Errorf("firestore get error: %w", err)
	}

	// 4. Deserialize: Convert raw Firestore data into our nice Go struct.
	var builder domain.Builder
	if err := docSnap.DataTo(&builder); err != nil {
		return nil, fmt.Errorf("failed to map data to builder struct: %w", err)
	}

	return &builder, nil
}

// Create saves a new builder structure as a document.
func (f *FirestoreClient) Create(ctx context.Context, builder *domain.Builder) error {
	// We use .Create() instead of .Set() to ensure we don't accidentally overwrite
	// an existing user if our service logic failed. It fails if doc exists.
	_, err := f.client.Collection(usersCollection).Doc(builder.ID).Create(ctx, builder)
	if err != nil {
		// In a real app, check here if error is "AlreadyExists"
		return fmt.Errorf("firestore create error: %w", err)
	}
	return nil
}

// UpdateStripeID performs a partial update on a document.
// It only changes the fields we specify, leaving the rest alone.
func (f *FirestoreClient) UpdateStripeID(ctx context.Context, discordID, stripeAccountID string) error {
	docRef := f.client.Collection(usersCollection).Doc(discordID)

	// We use []firestore.Update for atomic patch operations.
	updates := []firestore.Update{
		{Path: "stripe_account_id", Value: stripeAccountID},
		{Path: "updated_at", Value: time.Now().UTC()},
	}

	_, err := docRef.Update(ctx, updates)
	if err != nil {
		// Handle "not found" if trying to update a non-existent user
		if status.Code(err) == codes.NotFound {
			return service.ErrBuilderNotFound
		}
		return fmt.Errorf("firestore update stripe id error: %w", err)
	}
	return nil
}

// UpdateProfileData handles updating nested data structures.
func (f *FirestoreClient) UpdateProfileData(ctx context.Context, discordID string, profile domain.ProfileData) error {
	docRef := f.client.Collection(usersCollection).Doc(discordID)

	// Firestore allows updating nested fields using dot notation in the path.
	// This replaces the entire "profile_data" map with the new one.
	updates := []firestore.Update{
		{Path: "profile_data", Value: profile},
		{Path: "updated_at", Value: time.Now().UTC()},
	}

	_, err := docRef.Update(ctx, updates)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return service.ErrBuilderNotFound
		}
		return fmt.Errorf("firestore update profile error: %w", err)
	}
	return nil
}
