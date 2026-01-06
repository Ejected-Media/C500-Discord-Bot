This is an exciting pivot point. You have built a complete, self-contained web application.
Now, it's time to fulfill the "integration" part of the project's name. Remember the original developer architecture diagram? It showed a Python Discord Bot talking to your Go API.
Right now, your Go app only speaks HTML intended for humans in browsers. We need to teach it to speak JSON intended for other programs (like bots).
The next step is to build a JSON API endpoint to list products.
This will allow an external script (which we will write later) to ask your server, "Hey, what are you selling right now?" and get a structured answer back.
Here is the plan to create your first API endpoint:
 * Create an API Handler: A new handler file dedicated to serving JSON data instead of HTML templates.
 * Register Route: Add an /api/products route to main.go.
 * Test: Verify the endpoint returns correct JSON data.
Step 1: Create the API Handler (handlers/api.go)
Create a new file named c500-web-go/handlers/api.go.
This handler looks different from the others. Notice:
 * It doesn't import html/template.
 * It uses encoding/json.
 * It sets the Content-Type header to application/json.
<!-- end list -->
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"c500-web-go/models"
)

// APIHandler holds dependencies for REST API endpoints.
type APIHandler struct {
	Products *models.ProductModel
	// We might add OrderModel here later for admin API endpoints
}

// NewAPIHandler creates a new instance.
func NewAPIHandler(p *models.ProductModel) *APIHandler {
	return &APIHandler{
		Products: p,
	}
}

// --- JSON Helper ---

// jsonResponse is a helper to write JSON data and status codes to the response.
func (h *APIHandler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	// NewEncoder writes directly to the response stream
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("JSON encode error: %v", err)
		// If we can't even encode JSON, something is seriously wrong.
		http.Error(w, "Internal Server Error trying to encode JSON", http.StatusInternalServerError)
	}
}

// errorResponse is a helper to send standardized JSON error messages.
func (h *APIHandler) errorResponse(w http.ResponseWriter, status int, message string) {
	// Create a quick anonymous struct for the error shape
	data := struct {
		Error string `json:"error"`
	}{
		Error: message,
	}
	h.jsonResponse(w, status, data)
}

// --- API Endpoints ---

// ListProducts handles GET /api/products.
// It returns a JSON array of all available products.
func (h *APIHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// 1. Ensure this is a GET request
	if r.Method != http.MethodGet {
		h.errorResponse(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// 2. Fetch data from the database using existing model
	products, err := h.Products.All()
	if err != nil {
		log.Printf("DB Error in API: %v", err)
		h.errorResponse(w, http.StatusInternalServerError, "Database error fetching products")
		return
	}

	// 3. If no products exist, ensure we return an empty array [] instead of null
	if products == nil {
		products = []models.Product{}
	}

	// 4. Send the JSON response
	h.jsonResponse(w, http.StatusOK, products)
}

Step 2: Register Route in main.go
Open main.go to initialize the new API handler and register its route.
We will use a convention of starting API routes with /api/.
// ... imports ...

func main() {
    // ... load env, connect DB, init models, parse templates ...

	// --- Initialize Handlers ---
	homeHandler := handlers.NewHomeHandler(tmpl)
	docsHandler := handlers.NewDocsHandler(tmpl)
	storeHandler := handlers.NewStoreHandler(tmpl, productModel)
	checkoutHandler := handlers.NewCheckoutHandler(productModel)
	orderHandler := handlers.NewOrderHandler(tmpl)
	webhookHandler := handlers.NewWebhookHandler(productModel, orderModel)
    // --- NEW: Init API Handler ---
	apiHandler := handlers.NewAPIHandler(productModel)


	// 3. Create router
	mux := http.NewServeMux()

	// ... static files and web routes ...

	// --- NEW ROUTE: API Endpoints ---
    // We explicitly require GET method for listing products.
	mux.HandleFunc("GET /api/products", apiHandler.ListProducts)


    // ... webhook route and server start ...
}

Step 3: Test the API
 * Restart your server: go run main.go
 * Open your browser (or a tool like Postman, or use curl in the terminal).
 * Navigate to: http://localhost:8080/api/products
Instead of an HTML page with navigation bars and CSS, you should see raw, structured data like this:
[
  {
    "id": 1,
    "name": "DB Vintage Wumpus Plush",
    "description": "From the database! A rare 2018 plushie.",
    "price": 45,
    "image_url": "https://placehold.co/600x400/png?text=DB+Wumpus",
    "seller_id": "discord_user_123"
  },
  {
    "id": 2,
    "name": "DB Emoji Pack",
    "description": "10 custom emojis stored in SQL.",
    "price": 15.5,
    "image_url": "https://placehold.co/600x400/png?text=DB+Emoji",
    "seller_id": "artist_jane"
  },
  // ... etc ...
]

You have successfully opened a door for other applications to interact with your data! The next logical step will be to create the Python script (the "bot") that uses this API.
