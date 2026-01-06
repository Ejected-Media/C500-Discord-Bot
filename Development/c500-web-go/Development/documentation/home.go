package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// HomeHandler holds dependencies for general site pages.
type HomeHandler struct {
	Templates *template.Template
}

// NewHomeHandler creates a new instance of HomeHandler.
func NewHomeHandler(t *template.Template) *HomeHandler {
	return &HomeHandler{
		Templates: t,
	}
}

// Landing handles requests for the root path (/).
// It renders templates/home.html.
func (h *HomeHandler) Landing(w http.ResponseWriter, r *http.Request) {
	// Important: ServeMux treats "/" as a catch-all.
	// We must ensure this handler only responds to the exact root path.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	data := PageData{
		Title: "Home - c500-web-go",
	}
	
	// Re-use the render helper logic (you could refactor the helper to a shared utility file later)
	tmplName := "home.html"
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
