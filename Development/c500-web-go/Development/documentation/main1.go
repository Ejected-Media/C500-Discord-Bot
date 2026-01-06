package main

import (
    // ... other imports remain the same ...
	"log"
	"net/http"
	// ...
    
	"c500-web-go/handlers"
)

func main() {
	// 1. Parse Templates (Same as before)
	tmpl, err := parseTemplates("./templates")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}
	// log.Println("Templates parsed successfully.") // Optional noise reduction

	// 2. Initialize Handlers
	// --- CHANGE: Initialize the new HomeHandler ---
	homeHandler := handlers.NewHomeHandler(tmpl)
	docsHandler := handlers.NewDocsHandler(tmpl)

	// 3. Create a new router (ServeMux) (Same as before)
	mux := http.NewServeMux()

	// --- Static File Serving --- (Same as before)
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// --- Application Routes ---

	// --- CHANGE: Replace temporary redirect with real Home handler ---
    // The "/" matches everything not matched by other routes, so the handler
    // must check path == "/"
	mux.HandleFunc("/", homeHandler.Landing)

	// Main Docs Hub (Same as before)
	mux.HandleFunc("/docs/", docsHandler.Index)

    // ... rest of the docs routes remain the same ...
	mux.HandleFunc("/docs/buyer", docsHandler.Buyer)
	mux.HandleFunc("/docs/seller", docsHandler.Seller)
	mux.HandleFunc("/docs/developer", docsHandler.Developer)
	mux.HandleFunc("/docs/admin", docsHandler.Admin)


	// 4. Start the Server (Same as before)
    // ...
}
// ... parseTemplates function remains the same ...
