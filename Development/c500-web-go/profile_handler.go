package handlers

import (
	"html/template"
	"net/http"
	"github.com/gin-gonic/gin"
	"c500-web-go/internal/clients"
)

type ProfileHandler struct {
	coreClient *clients.CoreAPIClient
}

// GetBuilderProfile handles GET /builder/:username
func (h *ProfileHandler) GetBuilderProfile(c *gin.Context) {
	username := c.Param("username")

	// 1. Call Core API to fetch builder data (this involves an internal HTTP request)
	builderData, err := h.coreClient.GetPublicBuilderData(c.Request.Context(), username)
	if err != nil {
		// Handle 404 Not Found or 500 Internal Error gracefully
		c.HTML(http.StatusNotFound, "error.html", gin.H{"Message": "Builder not found."})
		return
	}

	// 2. Prepare data for the template.
	// CRITICAL STEP: Convert string data to template.HTML/CSS types.
	// This tells the template engine: "Trust me, don't escape this."
	data := gin.H{
		"Title":          builderData.DisplayName + "'s Profile",
		"BuilderName":    builderData.DisplayName,
		// The Core API has already sanitized these strings.
		"SafeCustomHTML": template.HTML(builderData.ProfileHTMLRaw),
		"SafeCustomCSS":  template.CSS(builderData.ProfileCSSRaw),
	}

	// 3. Render the profile.html template with the prepared data
	c.HTML(http.StatusOK, "profile.html", data)
}
