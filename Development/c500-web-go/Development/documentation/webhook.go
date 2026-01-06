package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"c500-web-go/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// WebhookHandler holds dependencies (like DB access for later use).
type WebhookHandler struct {
	Products *models.ProductModel
}

func NewWebhookHandler(p *models.ProductModel) *WebhookHandler {
	return &WebhookHandler{Products: p}
}

// Handle is the endpoint that Stripe servers will send POST requests to.
// Route: POST /webhook
func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// 1. Read the request body payload
	const MaxBodyBytes = int64(65536) // Limit body size to prevent abuse
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	// 2. Get the signature from the header and the secret from .env
	sigHeader := r.Header.Get("Stripe-Signature")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		log.Fatal("STRIPE_WEBHOOK_SECRET is not set in .env")
	}

	// 3. Verify the event signature using Stripe's library.
	// This ensures the request actually came from Stripe and hasn't been tampered with.
	event, err := webhook.ConstructEvent(payload, sigHeader, endpointSecret)
	if err != nil {
		log.Printf("⚠️  Webhook signature verification failed. %v", err)
		w.WriteHeader(http.StatusBadRequest) // Return 400 Bad Request
		return
	}

	// 4. Handle the specific event type
	// We are only listening for 'checkout.session.completed' right now.
	if event.Type == "checkout.session.completed" {
		var session stripe.CheckoutSession
		// Unmarshal the raw JSON event data into the Stripe struct
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Printf("Error parsing webhook JSON: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// --- BUSINESS LOGIC GOES HERE ---
		// For now, we just print the success to the terminal.
		// In the next step, we will use this info to update our PostgreSQL database.
		log.Printf("\n💰 >>> PAYMENT RECEIVED! <<< 💰\n")
		log.Printf("Session ID: %s\n", session.ID)
		log.Printf("Amount Total: %d cents\n", session.AmountTotal)
		log.Printf("Customer Email: %s\n", session.CustomerDetails.Email)
		// You can extract metadata here to know WHICH product ID was bought.
		log.Println("---------------------------------\n")
	}

	// 5. Acknowledge receipt by returning a 200 OK
	w.WriteHeader(http.StatusOK)
}
