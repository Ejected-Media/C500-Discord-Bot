package main

import (
	"context"
	"log"
	"os"

	"c500-core-go/config"
	"c500-core-go/internal/database"
	"c500-core-go/internal/integrations/stripe"
	"c500-core-go/internal/service"
	"c500-core-go/internal/transport"
	"c500-core-go/internal/transport/handlers"
)

func main() {
	// 1. Load configuration (Env variables)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	ctx := context.Background()

	// 2. Initialize Infrastructure Clients
	// Connect to Firestore
	firestoreClient, err := database.NewFirestoreClient(ctx, cfg.GoogleProjectID)
	if err != nil {
		log.Fatalf("Failed to init Firestore: %v", err)
	}
	defer firestoreClient.Close()

	// Init Stripe
	stripeClient := stripe.NewClient(cfg.StripeSecretKey)

	// 3. Initialize Service Layer (Business Logic)
	// We inject the database client into the service layer
	dropService := service.NewDropService(firestoreClient)
	// builderService := service.NewBuilderService(firestoreClient, stripeClient)

	// 4. Initialize HTTP Handlers
	// We inject the services into the handlers
	dropHandler := handlers.NewDropHandler(dropService)
	webhookHandler := handlers.NewWebhookHandler(dropService, stripeClient, cfg.StripeWebhookSecret)

	// 5. Setup Router & Start Server
	router := transport.SetupRouter(dropHandler, webhookHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("C500 Core API starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
