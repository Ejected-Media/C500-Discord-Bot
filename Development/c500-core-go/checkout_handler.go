package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin" // We use Gin as our HTTP web framework

	"c500-core-go/internal/service"
)

// CheckoutHandler holds dependencies needed to process checkout requests.
type CheckoutHandler struct {
	checkoutService service.CheckoutService
}

// NewCheckoutHandler is the constructor.
func NewCheckoutHandler(cs service.CheckoutService) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutService: cs,
	}
}

// RegisterRoutes connects the HTTP URLs to the handler functions.
// This is called in main.go.
func (h *CheckoutHandler) RegisterRoutes(router *gin.RouterGroup) {
	// This maps the POST request from the Python bot to the handler below.
	router.POST("/checkout/session", h.CreateSession)
	// We would also register the Stripe Webhook endpoint here later:
	// router.POST("/webhooks/stripe", h.HandleStripeWebhook)
}


// ==========================================
// Request/Response Structs (Data Contracts)
// ==========================================

// createSessionRequest defines exactly what we expect the Python bot to send us.
type createSessionRequest struct {
	DropID         string `json:"drop_id" binding:"required"`
	BuyerDiscordID string `json:"buyer_discord_id" binding:"required"`
}

// createSessionResponse defines what we send back to Python.
type createSessionResponse struct {
	URL string `json:"url"`
}

// ==========================================
// Handler Functions
// ==========================================

// CreateSession is the entrypoint for the "Buy Now" button click.
// It receives JSON, validates it, and calls the business logic.
func (h *CheckoutHandler) CreateSession(c *gin.Context) {
	var req createSessionRequest

	// 1. Parse and Validate JSON input
	// Gin's ShouldBindJSON automatically checks that fields marked 'required' exist.
	if err := c.ShouldBindJSON(&req); err != nil {
		// If Python sent bad data, return HTTP 400 Bad Request.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call the Service Layer (The Business Brain)
	// This is where the actual work (DB checks, talking to Stripe) happens.
	checkoutURL, err := h.checkoutService.CreateCheckoutSession(c.Request.Context(), req.DropID, req.BuyerDiscordID)

	// 3. Handle Errors from the business logic
	if err != nil {
		// We check the type of error to return the correct HTTP status code.
		switch {
		case errors.Is(err, service.ErrDropNotFound):
			// The drop ID was invalid. Return HTTP 404.
			c.JSON(http.StatusNotFound, gin.H{"error": "Drop not found"})
		case errors.Is(err, service.ErrDropNotAvailable):
			// The drop is already sold or pending. Return HTTP 409 Conflict.
			// This triggers the "Too late!" message in the Python bot.
			c.JSON(http.StatusConflict, gin.H{"error": "Drop is no longer available"})
		default:
			// Something unexpected happened (DB down, Stripe down). Return HTTP 500.
			// In production, log the actual error securely and send a generic message to client.
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error processing checkout"})
		}
		return
	}

	// 4. Success! Return the Stripe URL to the Python bot.
	// Return HTTP 200 OK.
	c.JSON(http.StatusOK, createSessionResponse{
		URL: checkoutURL,
	})
}
