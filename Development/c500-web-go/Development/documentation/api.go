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
