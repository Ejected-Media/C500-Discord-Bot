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
