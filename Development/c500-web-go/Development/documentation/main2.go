// ... imports remain the same ...

func main() {
	// 1. Parse Templates (same)
	tmpl, err := parseTemplates("./templates")
    // ...

	// 2. Initialize Handlers
	homeHandler := handlers.NewHomeHandler(tmpl)
	docsHandler := handlers.NewDocsHandler(tmpl)
    // --- CHANGE: Init store handler ---
	storeHandler := handlers.NewStoreHandler(tmpl)

	// 3. Create router (same)
	mux := http.NewServeMux()

	// ... static files and home routes same ...

	// --- CHANGE: Add Store Route ---
	mux.HandleFunc("/store", storeHandler.Index)

	// ... docs routes same ...
    // ... start server same ...
}
