// ... inside handlers/checkout.go ...
	// ... imports, including strconv ...
    "strconv"
    // ...

func (h *CheckoutHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
    // ... fetch product from DB ...
	product, err := h.Products.Get(id)
    // ... checks ...

	// ... stripe setup ...

	// 4. Create the checkout session parameters
	params := &stripe.CheckoutSessionParams{
		// --- NEW SECTION: Attach Metadata ---
		// We attach our internal Product ID so we know what was bought later.
		Metadata: map[string]string{
			"product_id": strconv.Itoa(product.ID),
		},
		// ------------------------------------

		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
        // ... rest of the params remain the same ...
    
