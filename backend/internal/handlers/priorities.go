package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/servicedesk/backend/internal/models"
	"github.com/servicedesk/backend/internal/repository"
)

// PriorityHandler handles priority-related endpoints.
type PriorityHandler struct {
	repo *repository.PriorityRepo
}

func NewPriorityHandler(repo *repository.PriorityRepo) *PriorityHandler {
	return &PriorityHandler{repo: repo}
}

// List returns all priorities.
func (h *PriorityHandler) List(c *gin.Context) {
	priorities, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch priorities"})
		return
	}
	if priorities == nil {
		priorities = []*models.Priority{}
	}
	c.JSON(http.StatusOK, priorities)
}
