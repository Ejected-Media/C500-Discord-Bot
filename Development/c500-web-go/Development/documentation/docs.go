package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// DocsHandler holds the dependencies needed for documentation requests,
// specifically the parsed HTML templates.
type DocsHandler struct {
	Templates *template.Template
}

// NewDocsHandler creates a new instance of DocsHandler with the given templates.
// This is typically called in main.go during server startup.
func NewDocsHandler(t *template.Template) *DocsHandler {
	return &DocsHandler{
		Templates: t,
	}
}

// PageData is a generic struct to pass data to the templates.
// You can expand this based on what dynamic content your pages need.
type PageData struct {
	Title string
	// Add other fields here if needed (e.g., CurrentUser, ActiveNavHighlight)
}

// --- Handler Methods ---

// Index handles requests for the main documentation hub page (/docs/).
// It renders templates/docs/index.html.
func (h *DocsHandler) Index(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Documentation Hub - c500-web-go",
	}
	h.render(w, "docs/index.html", data)
}

// Buyer handles requests for the buyer's manual (/docs/buyer).
// It renders templates/docs/buyer.html.
func (h *DocsHandler) Buyer(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Buyer's Manual",
	}
	h.render(w, "docs/buyer.html", data)
}

// Seller handles requests for the seller's manual (/docs/seller).
// It renders templates/docs/seller.html.
func (h *DocsHandler) Seller(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Seller's Manual",
	}
	h.render(w, "docs/seller.html", data)
}

// Developer handles requests for the developer's manual (/docs/developer).
// It renders templates/docs/developer.html.
func (h *DocsHandler) Developer(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Developer's Documentation",
	}
	h.render(w, "docs/developer.html", data)
}

// Admin handles requests for the admin/mod manual (/docs/admin).
// It renders templates/docs/admin.html.
func (h *DocsHandler) Admin(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "Admin & Moderator Manual",
	}
	h.render(w, "docs/admin.html", data)
}

// --- Helper Functions ---

// render is a private helper method to execute templates uniformly.
// It handles looking up the template name and logging errors if execution fails.
func (h *DocsHandler) render(w http.ResponseWriter, tmplName string, data PageData) {
	// Check if the template exists in the parsed set.
	if h.Templates.Lookup(tmplName) == nil {
		log.Printf("Error: Template '%s' not found. Ensure it is in the templates directory and parsed correctly in main.go.", tmplName)
		http.Error(w, "Page not found (Template missing)", http.StatusNotFound)
		return
	}

	// Execute the template with the provided data.
	err := h.Templates.ExecuteTemplate(w, tmplName, data)
	if err != nil {
		log.Printf("Error executing template '%s': %v", tmplName, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
