package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv" // Need strconv to convert string ID from URL to int

	"c500-web-go/models"
)

// ... (StoreHandler struct, NewStoreHandler, StorePageData, and Index method remain the same) ...

// --- NEW HANDLER METHOD ---

// ProductDetail handles requests for a single product page (e.g., /store/product/1).
func (h *StoreHandler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	// 1. Extract the ID from the URL path.
	// This relies on the route being defined as "/store/product/{id}" in main.go.
	idStr := r.PathValue("id")

	// 2. Convert string ID to integer.
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r) // Return 404 if ID is not a valid positive integer
		return
	}

	// 3. Fetch the specific product using our new model function.
	product, err := models.GetProductByID(id)
	if err != nil {
		// If product isn't found in our "database", return 404.
		http.NotFound(w, r)
		return
	}

	// 4. Render the template specifically for detail view.
	// We pass the product directly as the data dot (.).
	tmplName := "product_detail.html"
	if h.Templates.Lookup(tmplName) == nil {
		log.Printf("Error: Template '%s' not found.", tmplName)
		http.Error(w, "Template missing", http.StatusInternalServerError)
		return
	}

	err = h.Templates.ExecuteTemplate(w, tmplName, product)
	if err != nil {
		log.Printf("Error executing template '%s': %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
