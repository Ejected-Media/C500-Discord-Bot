// ... inside main.go ...

	// --- Application Routes ---
	mux.HandleFunc("/", homeHandler.Landing)

	// Main Store Listings
	mux.HandleFunc("/store", storeHandler.Index)

	// --- NEW ROUTE: Product Detail ---
    // The {id} is a path value wildcard introduced in Go 1.22
	mux.HandleFunc("GET /store/product/{id}", storeHandler.ProductDetail)

	// ... docs routes remain the same ...
