package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// OrderHandler holds dependencies for order confirmation pages.
type OrderHandler struct {
	Templates *template.Template
}

// NewOrderHandler creates a new instance.
func NewOrderHandler(t *template.Template) *OrderHandler {
	return &OrderHandler{
		Templates: t,
	}
}

// Success handles the landing page after a successful Stripe payment.
// Route: GET /success?session_id={CHECKOUT_SESSION_ID}
func (h *OrderHandler) Success(w http.ResponseWriter, r *http.Request) {
	// NOTE: In a real application, you would verify the session_id here.
	// sessionID := r.URL.Query().Get("session_id")
	// Currently, we just assume if they reached this URL, it's fine for a demo.

	data := PageData{
		Title: "Order Confirmed! - c500-web-go",
	}

	// Render template logic (copying this pattern again for expediency)
	tmplName := "success.html"
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
