package transport

import (
	"github.com/gin-gonic/gin"
	"c500-core-go/internal/transport/handlers"
	"c500-core-go/internal/transport/middleware"
)

// SetupRouter initializes the Gin engine and defines routes
func SetupRouter(dropHandler *handlers.DropHandler, webhookHandler *handlers.WebhookHandler) *gin.Engine {
	r := gin.Default() // Creates router with default logger and recovery middleware

	// --- Health Check (for Cloud Run) ---
	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

	// --- Public Webhooks (Stripe talks to these) ---
	webhooks := r.Group("/webhooks")
	{
		// Stripe verifies its own signatures, so no internal auth needed here
		webhooks.POST("/stripe", webhookHandler.HandleStripeEvent)
	}

	// --- Internal APIs (Only the Python Bot talks to these) ---
	// We protect these with a shared secret middleware so the public can't access them.
	internal := r.Group("/api/internal")
	internal.Use(middleware.InternalAuthCheck()) 
	{
		drops := internal.Group("/drops")
		{
			drops.POST("/create", dropHandler.CreateDrop)
			// drops.GET("/:id", dropHandler.GetDrop)
		}

		// sellers := internal.Group("/seller") { ... }
	}

	return r
}
