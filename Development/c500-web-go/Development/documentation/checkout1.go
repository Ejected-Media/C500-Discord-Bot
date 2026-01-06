package handlers

import (
	// ... imports ...
	"c500-web-go/models"
    // ...
)

// CheckoutHandler dependencies.
type CheckoutHandler struct{
    // NEW: Add model dependency
    Products *models.ProductModel
}

// Update constructor
func NewCheckoutHandler(p *models.ProductModel) *CheckoutHandler {
	return &CheckoutHandler{
        Products: p,
    }
}

func (h *CheckoutHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
    // ... get ID from URL (same) ...
	idStr := r.PathValue("id")
    id, _ := strconv.Atoi(idStr) // (simplified error check for brevity here)

	// 2. Fetch real product details from DB
	product, err := h.Products.Get(id)
	if err != nil {
        log.Printf("DB Error: %v", err)
		http.Error(w, "Internal Error", 500)
		return
	}
    if product == nil {
        http.NotFound(w, r)
        return
    }

	// ... Stripe configuration (same) ...

	// 4. Create session parameters using REAL product data
	params := &stripe.CheckoutSessionParams{
		// ... types ...
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String(string(stripe.CurrencyUSD)),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
                        // Use the data from the DB product
						Name:        stripe.String(product.Name),
						Description: stripe.String(product.Description),
						Images:      stripe.StringSlice([]string{product.ImageURL}),
					},
                    // Use the real price from DB
					UnitAmount: stripe.Int64(int64(product.Price * 100)),
				},
				Quantity: stripe.Int64(1),
			},
		},
        // ... rest of Stripe params (same) ...
	}
    
    // ... create session and redirect (same) ...
}
