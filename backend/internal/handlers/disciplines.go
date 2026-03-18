package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/servicedesk/backend/internal/models"
	"github.com/servicedesk/backend/internal/repository"
)

// DisciplineHandler handles discipline-related endpoints.
type DisciplineHandler struct {
	repo *repository.DisciplineRepo
}

func NewDisciplineHandler(repo *repository.DisciplineRepo) *DisciplineHandler {
	return &DisciplineHandler{repo: repo}
}

// List returns all disciplines.
func (h *DisciplineHandler) List(c *gin.Context) {
	disciplines, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch disciplines"})
		return
	}
	if disciplines == nil {
		disciplines = []*models.Discipline{}
	}
	c.JSON(http.StatusOK, disciplines)
}

// disciplineRequest is shared body for Create and Update.
type disciplineRequest struct {
	Name        string `json:"name"        binding:"required"`
	Description string `json:"description"`
}

// Create adds a new discipline. Admin only.
func (h *DisciplineHandler) Create(c *gin.Context) {
	var req disciplineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discipline, err := h.repo.Create(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create discipline"})
		return
	}

	c.JSON(http.StatusCreated, discipline)
}

// Update modifies an existing discipline. Admin only.
func (h *DisciplineHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req disciplineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	discipline, err := h.repo.Update(id, req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update discipline"})
		return
	}
	if discipline == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "discipline not found"})
		return
	}

	c.JSON(http.StatusOK, discipline)
}

// Delete removes a discipline. Admin only. Returns 409 if grades reference it.
func (h *DisciplineHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
