package handlers

import (
	"html/template"
	"log"
	"net/http"
    "strconv"

	"c500-web-go/models"
)

// StoreHandler holds dependencies.
type StoreHandler struct {
	Templates *template.Template
    // NEW: Add the product model dependency
	Products  *models.ProductModel
}

// Update constructor to accept the model
func NewStoreHandler(t *template.Template, p *models.ProductModel) *StoreHandler {
	return &StoreHandler{
		Templates: t,
		Products:  p,
	}
}

// ... StorePageData struct remains the same ...

// Index handles the main store page.
func (h *StoreHandler) Index(w http.ResponseWriter, r *http.Request) {
	// --- CHANGE: Fetch data from real DB instead of mock functions ---
	products, err := h.Products.All()
	if err != nil {
		log.Printf("Error fetching products from DB: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := StorePageData{
		Title:    "Browse Items - c500-web-go Store",
		Products: products,
	}
    // ... rendering logic remains the same ...
	h.render(w, "store.html", data) // Assuming you refactored render helper, otherwise use ExecuteTemplate
}

// ProductDetail handles single product page.
func (h *StoreHandler) ProductDetail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

    // --- CHANGE: Use DB model to get product ---
	product, err := h.Products.Get(id)
	if err != nil {
        // Log DB errors
		log.Printf("Error fetching product %d: %v", id, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
    // If product is nil, it means it wasn't found in DB
	if product == nil {
		http.NotFound(w, r)
		return
	}

    // ... rendering logic remains the same, passing *product ...
     h.render(w, "product_detail.html", product)
}

// (Optional: If you haven't refactored the 'render' helper from docs.go into a shared place, 
// you'll need to copy/paste it here or use ExecuteTemplate directly like before. 
// For brevity, I assumed a helper function like in docs.go)
func (h *StoreHandler) render(w http.ResponseWriter, tmplName string, data interface{}) {
     // ... (implementation same as in docs.go but taking interface{}) ...
	if h.Templates.Lookup(tmplName) == nil {
		log.Printf("Error: Template '%s' not found.", tmplName)
		http.Error(w, "Template missing", http.StatusInternalServerError)
		return
	}
	err := h.Templates.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Printf("Error executing template '%s': %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
