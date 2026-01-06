This is the final step in closing the loop on the basic purchase flow. You've successfully sent the user to Stripe to pay, and Stripe has sent them back to your specified SuccessUrl. Now we just need a page to welcome them back.
Here is the plan to create the "Order Success" page:
 * Create an Order Handler: A new handler file (handlers/order.go) to manage post-purchase pages.
 * Create a Success Template: An HTML file (templates/success.html) to display a "Thank You" message.
 * Update main.go: Register the new /success route.
Step 1: Create the Order Handler (handlers/order.go)
Create a new file named c500-web-go/handlers/order.go.
For now, this handler will just render a static success page. In a real production app, you would use the session_id query parameter in the URL to ask Stripe for the order details and verify that the payment actually succeeded before showing this page.
package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// OrderHandler holds dependencies for order confirmation pages.
type OrderHandler struct {
	Templates *template.Template
}

// NewOrderHandler creates a new instance.
func NewOrderHandler(t *template.Template) *OrderHandler {
	return &OrderHandler{
		Templates: t,
	}
}

// Success handles the landing page after a successful Stripe payment.
// Route: GET /success?session_id={CHECKOUT_SESSION_ID}
func (h *OrderHandler) Success(w http.ResponseWriter, r *http.Request) {
	// NOTE: In a real application, you would verify the session_id here.
	// sessionID := r.URL.Query().Get("session_id")
	// Currently, we just assume if they reached this URL, it's fine for a demo.

	data := PageData{
		Title: "Order Confirmed! - c500-web-go",
	}

	// Render template logic (copying this pattern again for expediency)
	tmplName := "success.html"
	if h.Templates.Lookup(tmplName) == nil {
		log.Printf("Error: Template '%s' not found.", tmplName)
		http.Error(w, "Template missing", http.StatusInternalServerError)
		return
	}

	err := h.Templates.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Printf("Error executing template '%s': %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

Step 2: Create the Success Template (templates/success.html)
Create a new HTML file in the templates/ directory named success.html.
We'll add some specific green styling to make it look like a successful confirmation page.
{{extends "base.html"}}

{{define "title"}}{{ .Title }}{{end}}

{{define "head_extra"}}
<style>
    .success-container {
        text-align: center;
        padding: 4rem 2rem;
        background-color: #f0fff4; /* Light green background */
        border: 1px solid #c6f6d5; /* Slightly darker green border */
        border-radius: 8px;
        margin-top: 2rem;
    }

    .success-icon {
        font-size: 5rem;
        color: #38a169; /* Nice green color */
        margin-bottom: 1rem;
    }

    .success-container h1 {
        margin-bottom: 1rem;
        color: #2f855a;
    }

    .success-container p {
        font-size: 1.2rem;
        margin-bottom: 2rem;
        color: #4a5568;
    }

    /* Re-using button styles but making it behave like a link */
    .btn-primary {
        display: inline-block;
        padding: 12px 30px;
        background-color: var(--color-primary);
        color: white;
        border-radius: 4px;
        font-weight: bold;
        transition: background-color 0.2s;
    }
    .btn-primary:hover {
        background-color: #357abd;
        text-decoration: none;
    }
</style>
{{end}}

{{define "content"}}
<section class="success-container">
    <div class="success-icon">✓</div>
    <h1>Thank You for Your Order!</h1>
    <p>Your payment has been processed successfully.</p>
    <p>You will receive a confirmation email shortly.</p>
    
    <a href="/store" class="btn-primary">Continue Shopping</a>
</section>
{{end}}

Step 3: Update main.go
Finally, register the new handler and route in main.go.
// ... imports ...
// ...

func main() {
    // ... load env and parse templates ...

	// 2. Initialize Handlers
	homeHandler := handlers.NewHomeHandler(tmpl)
	docsHandler := handlers.NewDocsHandler(tmpl)
	storeHandler := handlers.NewStoreHandler(tmpl)
	checkoutHandler := handlers.NewCheckoutHandler()
    // --- NEW: Init OrderHandler ---
	orderHandler := handlers.NewOrderHandler(tmpl)

	// 3. Create router
	mux := http.NewServeMux()

    // ... other routes ...
	mux.HandleFunc("POST /create-checkout-session/{id}", checkoutHandler.CreateCheckoutSession)

	// --- NEW ROUTE: Order Success Page ---
    // Stripe makes a GET request to this URL after payment.
	mux.HandleFunc("GET /success", orderHandler.Success)

    // ... docs routes and server start ...
}

Final Test Run
 * Restart your server: go run main.go
 * Go to the store: http://localhost:8080/store
 * Select a product: e.g., "Vintage Discord Wumpus Plush".
 * Click "Buy Now with Stripe".
 * Complete payment on Stripe: Use test card 4242424242424242, any future date (e.g., 12/34), and any CVC (e.g., 123).
 * Success! Stripe will redirect you back to your local server, and you should now see the green "Thank You for Your Order!" page instead of a 404 error.
Congratulations! You have built a functional end-to-end e-commerce flow with Go, HTML/CSS, and Stripe.
You now have a solid foundation. Future steps could include setting up a real PostgreSQL database instead of mock data, building the Discord bot to listen for commands, or setting up Stripe Webhooks to automatically fulfill orders when a payment succeeds.
