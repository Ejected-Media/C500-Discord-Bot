package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"c500-core-go/internal/service"
)

// FulfillmentHandler holds dependencies needed to process fulfillment requests.
type FulfillmentHandler struct {
	fulfillmentService service.FulfillmentService
}

// NewFulfillmentHandler is the constructor.
func NewFulfillmentHandler(fs service.FulfillmentService) *FulfillmentHandler {
	return &FulfillmentHandler{
		fulfillmentService: fs,
	}
}

// RegisterRoutes connects the HTTP URLs to the handler functions.
// This is called in main.go.
func (h *FulfillmentHandler) RegisterRoutes(router *gin.RouterGroup) {
	// Map endpoints to handle shipping and live fulfillment methods.
	// The ':orderID' part is a path parameter that Gin extracts for us.
	router.POST("/orders/:orderID/fulfill/ship", h.FulfillShip)
	router.POST("/orders/:orderID/fulfill/live", h.FulfillLive)
}


// ==========================================
// Request/Response Structs (Data Contracts)
// ==========================================

// fulfillShipRequest defines the expected JSON body from the Python bot's /fulfill ship command.
type fulfillShipRequest struct {
	TrackingNumber  string `json:"tracking_number" binding:"required"`
	Carrier         string `json:"carrier" binding:"required"`
	// We need the seller's ID to ensure the person trying to ship it actually owns the order.
	SellerDiscordID string `json:"seller_discord_id" binding:"required"`
}

// fulfillLiveRequest defines the expected JSON body from /fulfill live.
type fulfillLiveRequest struct {
	VODLink         string `json:"vod_url" binding:"required,url"` # Use 'url' validator
	SellerDiscordID string `json:"seller_discord_id" binding:"required"`
}


// ==========================================
// Handler Functions
// ==========================================

// FulfillShip processes requests for Ready-to-Ship items.
func (h *FulfillmentHandler) FulfillShip(c *gin.Context) {
	orderID := c.Param("orderID")
	var req fulfillShipRequest

	// 1. Parse and Validate JSON input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call the Service Layer
	err := h.fulfillmentService.FulfillOrderWithShipping(c.Request.Context(), orderID, req.SellerDiscordID, req.TrackingNumber, req.Carrier)

	// 3. Handle Errors
	if err != nil {
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		case errors.Is(err, service.ErrUnauthorizedSeller):
			// The person trying to ship it is not the seller on record.
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not the seller of this order"})
		case errors.Is(err, service.ErrOrderAlreadyFulfilled):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order is already fulfilled"})
		default:
			// Log actual error in production
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process fulfillment"})
		}
		return
	}

	// 4. Success
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Order fulfilled and funds released."})
}

// FulfillLive processes requests for Commission items.
func (h *FulfillmentHandler) FulfillLive(c *gin.Context) {
	orderID := c.Param("orderID")
	var req fulfillLiveRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.fulfillmentService.FulfillOrderWithVOD(c.Request.Context(), orderID, req.SellerDiscordID, req.VODLink)

	if err != nil {
		// (Same error handling logic as above)
		switch {
		case errors.Is(err, service.ErrOrderNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		case errors.Is(err, service.ErrUnauthorizedSeller):
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not the seller of this order"})
		case errors.Is(err, service.ErrOrderAlreadyFulfilled):
			c.JSON(http.StatusBadRequest, gin.H{"error": "Order is already fulfilled"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process fulfillment"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Order fulfilled and funds released."})
}
