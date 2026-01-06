This is a great juncture. You have a database-driven frontend that can initiate payments.
However, right now, your application has a kind of "amnesia" regarding completed payments. When a user pays on Stripe and gets redirected to your /success page, your Go server doesn't actually know that the payment was successful. It just displays a generic thank-you page.
If you were selling digital goods, you wouldn't know when to send the download link. If you were managing inventory, you wouldn't know when to decrement the stock count.
The next critical step is to implement Stripe Webhooks.
A webhook is a way for Stripe to send a background message (a POST request) to your server saying, "Hey! Checkout Session X just completed successfully, and here are the details."
Here is the plan to close the payment loop:
 * Prerequisite: Setup Stripe CLI. We need a way for Stripe's cloud servers to talk to your local machine (localhost:8080). The Stripe CLI is the easiest tool for this.
 * Create Webhook Handler: Write Go code to receive, verify, and process the incoming webhook event.
 * Register Route: Add the webhook endpoint to main.go.
 * Test: Simulate a payment and watch your server react in real-time.
Step 1: Prerequisite - Setup Stripe CLI
Since Stripe can't send requests directly to your laptop's localhost, we need to open a secure tunnel.
 * Install the Stripe CLI: Follow the official instructions for your OS here: https://stripe.com/docs/stripe-cli
   * macOS (Homebrew): brew install stripe/stripe-cli/stripe
   * Windows/Linux: Follow links on the page above.
 * Login: Open a new terminal window and run:
   stripe login

   Follow the browser prompt to authenticate with your Stripe account.
 * Start Listening: In that same terminal, tell Stripe to forward webhook events to your local Go server. We are listening for the specific event that means a payment succeeded.
   stripe listen --forward-to localhost:8080/webhook --events checkout.session.completed

   Keep this terminal window running. You will see output like this:
   > Ready! Your webhook signing secret is whsec_...
   > 
 * Update .env: Copy that "signing secret" (starts with whsec_) and add it to your .env file. Your app needs this to verify that incoming requests are genuinely from Stripe.
<!-- end list -->
# c500-web-go/.env
# ... existing keys ...
STRIPE_SECRET_KEY=sk_test_...
DATABASE_URL=...

# Add this new key found in the 'stripe listen' terminal output:
STRIPE_WEBHOOK_SECRET=whsec_YOUR_ACTUAL_SECRET_HERE

Step 2: Create Webhook Handler (handlers/webhook.go)
Create a new file named c500-web-go/handlers/webhook.go.
This handler is different from the others. It doesn't render HTML. It reads a raw JSON payload, verifies its cryptographic signature using Stripe's library (crucial for security!), and then processes the event.
For this step, our "processing" will just be logging the success to the console. In the next stage, we'll update the database here.
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

Step 3: Register Route in main.go
Open main.go and register the new handler.
// ... imports ...

func main() {
    // ... load env, connect DB, init models, parse templates ...

	// 3. Initialize Handlers
	// ... existing handlers ...
	checkoutHandler := handlers.NewCheckoutHandler(productModel)
	orderHandler := handlers.NewOrderHandler(tmpl)
    // --- NEW: Init WebhookHandler ---
	// We pass productModel in case we need DB access later
	webhookHandler := handlers.NewWebhookHandler(productModel)

	// 4. Create router
	mux := http.NewServeMux()

    // ... existing routes ...
	mux.HandleFunc("GET /success", orderHandler.Success)

	// --- NEW ROUTE: Stripe Webhook Endpoint ---
    // Crucial: This must be a POST request.
	mux.HandleFunc("POST /webhook", webhookHandler.Handle)

    // ... server start ...
}

Test the Full Loop
This is the moment of truth.
 * Terminal 1 (Stripe CLI): Make sure stripe listen ... is still running.
 * Terminal 2 (Go Server): Restart your server: go run main.go
 * Browser: Go to your store, pick a product, click "Buy Now," and complete the test payment on Stripe's page.
 * Watch Terminal 2 (Go Server):
After you see the "Success" page in your browser, look at the terminal where your Go server is running. You should instantly see the webhook trigger:
2026/01/06 16:30:00 Environment variables loaded successfully.
...
Starting server on http://localhost:8080 ...

💰 >>> PAYMENT RECEIVED! <<< 💰
Session ID: cs_test_a1b2c3d4...
Amount Total: 4500 cents
Customer Email: test@example.com
---------------------------------

Your Go server is now automatically being notified by Stripe whenever a purchase completes. The next logical step is to stop just printing this to the log and start recording it in a new orders table in PostgreSQL.
