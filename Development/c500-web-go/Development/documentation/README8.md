This is the exciting part: connecting your application to Stripe to actually process a payment.
We are going to create a new handler that will be triggered when a user clicks the "Buy Now" button. This handler's job is to tell Stripe what product is being bought and how much it costs. Stripe will then give us a special, secure URL to redirect the user to, where they can enter their payment information.
Here is the plan for implementing the checkout flow:
 * Create a new handler file: handlers/checkout.go to manage payment logic.
 * Implement the CreateCheckoutSession function: This function will fetch product details and use the Stripe API to generate a checkout session link.
 * Update main.go: Register the new route for this handler.
 * Update the Product Detail Template: Change the "Buy Now" button to submit a POST request to our new route.
Step 1 & 2: Create the Checkout Handler (handlers/checkout.go)
Create a new file named c500-web-go/handlers/checkout.go.
This code does the heavy lifting. It reads the Stripe secret key from the environment variable we set up, fetches the correct product from our mock data, and then makes an API call to Stripe to create a "Checkout Session".
package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	// Import our models and the Stripe libraries
	"c500-web-go/models"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
)

// CheckoutHandler holds dependencies for checkout-related actions.
// It doesn't need templates, as it redirects, doesn't render.
type CheckoutHandler struct{}

// NewCheckoutHandler creates a new instance.
func NewCheckoutHandler() *CheckoutHandler {
	return &CheckoutHandler{}
}

// CreateCheckoutSession handles POST requests to initiate a purchase.
// Route: POST /create-checkout-session/{id}
func (h *CheckoutHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	// 1. Get the Product ID from the URL path
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// 2. Fetch the product details from our mock "database"
	product, err := models.GetProductByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 3. Configure Stripe with your Secret Key
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Get your local domain from the .env file for redirect URLs
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		log.Fatal("DOMAIN environment variable not set in .env")
	}

	// 4. Create the checkout session parameters
	// This tells Stripe what the user is buying.
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name:        stripe.String(product.Name),
						Description: stripe.String(product.Description),
						// You can also pass the product image URL here so it shows up on the Stripe page!
						Images: stripe.StringSlice([]string{product.ImageURL}),
					},
					// Stripe expects amounts in cents (e.g., $10.00 = 1000 cents)
					UnitAmount: stripe.Int64(int64(product.Price * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		// URLs to redirect to after payment success or cancellation
		SuccessUrl: stripe.String(domain + "/success?session_id={CHECKOUT_SESSION_ID}"),
		CancelUrl:  stripe.String(domain + "/store/product/" + idStr),
	}

	// 5. Make the API call to create the session
	s, err := session.New(params)
	if err != nil {
		log.Printf("Error creating Stripe session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 6. Redirect the user's browser to the Stripe-hosted checkout page
	http.Redirect(w, r, s.URL, http.StatusSeeOther)
}

Step 3: Update main.go
Now, register the new handler and route in main.go.
// ... imports ...
	"c500-web-go/handlers"
)

func main() {
    // ... load env and parse templates ...

	// 2. Initialize Handlers
	homeHandler := handlers.NewHomeHandler(tmpl)
	docsHandler := handlers.NewDocsHandler(tmpl)
	storeHandler := handlers.NewStoreHandler(tmpl)
    // --- NEW: Init CheckoutHandler ---
	checkoutHandler := handlers.NewCheckoutHandler()

	// 3. Create router
	mux := http.NewServeMux()

    // ... other routes ...
	mux.HandleFunc("GET /store/product/{id}", storeHandler.ProductDetail)

	// --- NEW ROUTE: Create Checkout Session ---
    // Note: This is a POST request.
	mux.HandleFunc("POST /create-checkout-session/{id}", checkoutHandler.CreateCheckoutSession)

    // ... docs routes and server start ...
}

Step 4: Update the Product Detail Template (templates/product_detail.html)
Finally, change the "Buy Now" link into a form that sends a POST request to our new route.
Open templates/product_detail.html and find the <a> tag for the "Buy Now" button. Replace it with the following <form> block:
        <p class="seller-info">Sold by: {{.SellerID}}</p>

        <form action="/create-checkout-session/{{.ID}}" method="POST">
            <button type="submit" class="buy-now-btn">Buy Now with Stripe</button>
        </form>
    </div>
</article>
{{end}}

Test It Out
 * Restart your server: go run main.go
 * Go to a product page: e.g., http://localhost:8080/store/product/1
 * Click "Buy Now with Stripe".
You should be magically redirected to a Stripe-hosted checkout page that shows the correct product name, image, and price. You can use Stripe's test card numbers (like 4242424242424242) to complete a test payment.
Note: After a successful payment, Stripe will redirect you to http://localhost:8080/success..., which will currently give a "404 Not Found" error because we haven't built the success page yet. That will be our next step!
