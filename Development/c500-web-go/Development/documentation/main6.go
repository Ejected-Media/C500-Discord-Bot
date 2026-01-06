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
