package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/servicedesk/backend/internal/models"
	"github.com/servicedesk/backend/internal/repository"
)

// GradeTypeHandler handles grade_type-related endpoints.
type GradeTypeHandler struct {
	repo *repository.GradeTypeRepo
}

func NewGradeTypeHandler(repo *repository.GradeTypeRepo) *GradeTypeHandler {
	return &GradeTypeHandler{repo: repo}
}

// List returns all grade types.
func (h *GradeTypeHandler) List(c *gin.Context) {
	gradeTypes, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch grade types"})
		return
	}
	if gradeTypes == nil {
		gradeTypes = []*models.GradeType{}
	}
	c.JSON(http.StatusOK, gradeTypes)
}
