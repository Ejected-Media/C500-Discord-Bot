package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"c500-core-go/internal/domain"
	"c500-core-go/internal/service"
)

type DropHandler struct {
	dropService service.DropService // Dependency injection of business logic
}

// NewDropHandler is a constructor
func NewDropHandler(ds service.DropService) *DropHandler {
	return &DropHandler{dropService: ds}
}

// CreateDrop handles POST /api/internal/drops/create
func (h *DropHandler) CreateDrop(c *gin.Context) {
	var req domain.CreateDropRequest

	// 1. Bind JSON from request body to struct and validate inputs
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 2. Call the Service layer to perform business logic
	// (The service layer handles talking to Firestore, adding timestamps, setting initial status, etc.)
	newDrop, err := h.dropService.CreateNewDrop(c.Request.Context(), req)
	if err != nil {
		// In a real app, check error type to decide between 500 vs 400 errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create drop"})
		return
	}

	// 3. Return success response with the newly created drop data
	c.JSON(http.StatusCreated, newDrop)
}
