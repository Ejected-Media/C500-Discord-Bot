package http

import (
	"encoding/json"
	"io"
	"net/http"
	"os" // Needed to get the webhook secret from environment variables

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"

	"c500-core-go/internal/service"
)

// WebhookHandler holds dependencies needed to process incoming webhooks.
// It reuses the CheckoutService because handling a payment is part of the checkout lifecycle.
type WebhookHandler struct {
	checkoutService service.CheckoutService
	// In production, the webhook secret should be injected via config, not read directly from OS here.
	webhookSecret string
}

// NewWebhookHandler constructor.
// We pass the secret in here so it's only read from env vars once at startup.
func NewWebhookHandler(cs service.CheckoutService, secret string) *WebhookHandler {
	return &WebhookHandler{
		checkoutService: cs,
		webhookSecret:   secret,
	}
}

// RegisterRoutes connects the URL to the handler function.
// Called in main.go.
func (h *WebhookHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Stripe requests must be POST.
	router.POST("/webhooks/stripe", h.HandleStripeWebhook)
}

// ==========================================
// Handler Function
// ==========================================

// HandleStripeWebhook is the entrypoint for requests coming FROM Stripe's servers.
func (h *WebhookHandler) HandleStripeWebhook(c *gin.Context) {
	// 1. Read the raw body of the request.
	// We need the exact, unmodified bytes for signature verification.
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusServiceUnavailable)
		return
	}

	// 2. Get the signature header sent by Stripe.
	sigHeader := c.GetHeader("Stripe-Signature")

	// 3. CRITICAL SECURITY STEP: Verify the signature.
	// This function cryptographically checks that the payload was signed with our secret key.
	// If this fails, it means the request is fake or tampered with.
	event, err := webhook.ConstructEvent(payload, sigHeader, h.webhookSecret)
	if err != nil {
		// Invalid signature. Do not process.
		// Log this as a potential security incident in production.
		c.Status(http.StatusBadRequest)
		return
	}

	// 4. Event Routing: Switch based on what kind of event happened.
	switch event.Type {
	case "checkout.session.completed":
		// This is the event we want! Money has been successfully captured.
		var session stripe.CheckoutSession
		// Unmarshal the event data into a Stripe Session struct.
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return
		}

		// 5. Extract the crucial Metadata we attached back in stripe/client.go
		dropID := session.Metadata["drop_id"]
		buyerDiscordID := session.Metadata["buyer_discord_id"]
		// We also get the actual Stripe Transaction ID for our records.
		stripePaymentIntentID := session.PaymentIntent.ID

		// 6. Call the Service Layer to do the business logic (Update DB, notify user).
		// Note: We run this in a goroutine so we can immediately respond 200 OK to Stripe.
		// Stripe expects a quick response or it will keep retrying sending the webhook.
		go func() {
			// Create a background context, as the request context will be cancelled when this handler returns.
			ctx := context.Background()
			err := h.checkoutService.ProcessSuccessfulPayment(ctx, dropID, buyerDiscordID, stripePaymentIntentID)
			if err != nil {
				// Log this error heavily. It means we got paid but failed to update our system.
				// This requires manual intervention.
				// logger.Error("FAILED TO PROCESS PAYMENT", "error", err, "dropID", dropID)
			}
		}()

	default:
		// Handle other event types we don't care about (e.g., 'payment_intent.created').
		// Just return 200 OK so Stripe knows we received it.
	}

	// Acknowledge receipt to Stripe.
	c.Status(http.StatusOK)
}
