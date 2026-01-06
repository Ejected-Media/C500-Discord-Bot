package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv" // Needed for converting metadata string to int

	"c500-web-go/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// WebhookHandler dependencies updated to include Orders.
type WebhookHandler struct {
	Products *models.ProductModel
	Orders   *models.OrderModel // <-- Add this
}

// Update constructor to accept OrderModel.
func NewWebhookHandler(p *models.ProductModel, o *models.OrderModel) *WebhookHandler {
	return &WebhookHandler{
		Products: p,
		Orders:   o, // <-- Store it
	}
}

// Handle endpoint... (beginning of function remains the same) ...
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
    // ... read body, verify signature, parse event JSON ...
    // ... down to the event type check:

	// 4. Handle the specific event type
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// --- BUSINESS LOGIC UPDATED ---

		// 1. Extract the Product ID from the Metadata we sent earlier
		productIDStr := session.Metadata["product_id"]
		productID, err := strconv.Atoi(productIDStr)
		if err != nil {
			log.Printf("Error: Malformed product_id in metadata: %v", err)
			// We still return 200 OK to Stripe so they don't keep retrying,
			// but we log the error for our own investigation.
			w.WriteHeader(http.StatusOK)
			return
		}

		// 2. Create the Order struct payload
		order := &models.Order{
			StripeSessionID: session.ID,
			ProductID:       productID,
			AmountTotal:     int(session.AmountTotal), // Stripe sends int64, cast to int for DB
			CustomerEmail:   session.CustomerDetails.Email,
			Status:          string(session.PaymentStatus), // e.g., "paid"
		}

		// 3. Insert into PostgreSQL
		err = h.Orders.Insert(order)
		if err != nil {
			log.Printf("DATABASE ERROR: Could not save order: %v", err)
			// In a real app, you might return 500 here to make Stripe retry later,
			// depending on your idempotency strategy.
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Printf("✅ Order successfully saved to DB! ID: %d, Email: %s", order.ID, order.CustomerEmail)
	}

	// 5. Acknowledge receipt
	w.WriteHeader(http.StatusOK)
}
