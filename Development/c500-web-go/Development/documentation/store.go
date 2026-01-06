package handlers

import (
	"html/template"
	"log"
	"net/http"

	// IMPORT your models package here
	"c500-web-go/models"
)

// StoreHandler holds dependencies for store-related pages.
type StoreHandler struct {
	Templates *template.Template
}

// NewStoreHandler creates a new instance of StoreHandler.
func NewStoreHandler(t *template.Template) *StoreHandler {
	return &StoreHandler{
		Templates: t,
	}
}

// StorePageData is the specific structure of data we send to the store.html template.
type StorePageData struct {
	Title    string
	Products []models.Product
}

// Index handles requests for the main store page (/store).
// It fetches products and renders templates/store.html.
func (h *StoreHandler) Index(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch the data (currently mock data)
	products := models.GetMockProducts()

	// 2. Prepare the data for the template
	data := StorePageData{
		Title:    "Browse Items - c500-web-go Store",
		Products: products,
	}

	// 3. Render the template
	tmplName := "store.html"
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
