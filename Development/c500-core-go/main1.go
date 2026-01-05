package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"

	"c500-core-go/internal/database"
	stripeintegration "c500-core-go/internal/integrations/stripe"
	"c500-core-go/internal/service"
	transport "c500-core-go/internal/transport/http"
)

func main() {
	// 1. CRITICAL: Load Configuration from Environment Variables.
	// In production (Cloud Run), these are set in the deployment configuration.
	// Locally, you'd set them in your terminal before running.
	gcpProjectID := os.Getenv("GCP_PROJECT_ID")
	stripeSecretKey := os.Getenv("STRIPE_SECRET_KEY")
	stripeWebhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	port := os.Getenv("PORT")

	if gcpProjectID == "" || stripeSecretKey == "" || stripeWebhookSecret == "" {
		log.Fatal("Missing required environment variables (GCP_PROJECT_ID, STRIPE_SECRET_KEY, STRIPE_WEBHOOK_SECRET)")
	}
	if port == "" {
		port = "8080" // Default for Cloud Run
	}

	ctx := context.Background()

	// 2. Initialize Infrastructure Clients
	log.Println("Initializing Firestore client...")
	firestoreClient, err := database.NewFirestoreClient(ctx, gcpProjectID)
	if err != nil {
		log.Fatalf("Failed to init firestore: %v", err)
	}
	defer firestoreClient.Close() // Ensure connection closes on shutdown

	log.Println("Initializing Stripe global key...")
	// The official Stripe SDK uses a global variable for the secret key.
	stripe.Key = stripeSecretKey
	// Create our wrapper client.
	stripeClient := stripeintegration.NewClient()

	// =====================================================================
	// 3. THE WIRING PHASE (Dependency Injection)
	// We build the layers from the bottom up: Repo -> Service -> Handler
	// =====================================================================

	log.Println("Wiring application layers...")

	// --- Layer 3: Repositories & Integrations (Bottom) ---
	// firestoreClient already implements BuilderRepo, DropRepo, and OrderRepo interfaces.
	// stripeClient already implements StripeIntegration interface.

	// --- Layer 2: Services (Middle) ---
	// Inject repos/clients into services.
	// Note: We reuse firestoreClient wherever a Repo interface is needed.
	builderService := service.NewBuilderService(firestoreClient, stripeClient)
	// (We haven't written DropService yet, but it would go here)
	// dropService := service.NewDropService(firestoreClient, builderService)
	checkoutService := service.NewCheckoutService(firestoreClient, firestoreClient, stripeClient)
	fulfillmentService := service.NewFulfillmentService(firestoreClient, firestoreClient, stripeClient)

	// --- Layer 1: Handlers (Top) ---
	// Inject services into HTTP handlers.
	// builderHandler := transport.NewBuilderHandler(builderService) (Not written yet)
	// dropHandler := transport.NewDropHandler(dropService) (Not written yet)
	checkoutHandler := transport.NewCheckoutHandler(checkoutService)
	webhookHandler := transport.NewWebhookHandler(checkoutService, stripeWebhookSecret)
	fulfillmentHandler := transport.NewFulfillmentHandler(fulfillmentService)


	// 4. Setup HTTP Server (Gin Router)
	log.Println("Setting up HTTP router...")
	// Gin in release mode for production, debug mode locally.
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default() // Default includes basic logger and recovery middleware

	// A simple health check endpoint for Cloud Run load balancer.
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Register API Route Groups
	apiV1 := router.Group("/api/v1")
	{
		// Tell each handler to register its own paths under this group.
		// builderHandler.RegisterRoutes(apiV1)
		// dropHandler.RegisterRoutes(apiV1)
		checkoutHandler.RegisterRoutes(apiV1)
		fulfillmentHandler.RegisterRoutes(apiV1)
	}

	// Register Webhook Route (usually at root level or distinct path)
	// Note: It's NOT under /api/v1 because it's an external callback, not our internal API.
	webhookHandler.RegisterRoutes(router.Group("/"))


	// 5. Start the Engine
	log.Printf("🚀 C500 Core API starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
