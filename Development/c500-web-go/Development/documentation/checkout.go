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
