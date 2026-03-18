package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/servicedesk/backend/internal/middleware"
	"github.com/servicedesk/backend/internal/models"
	"github.com/servicedesk/backend/internal/repository"
)

// CommentHandler handles comment-related endpoints.
type CommentHandler struct {
	repo *repository.CommentRepo
}

func NewCommentHandler(repo *repository.CommentRepo) *CommentHandler {
	return &CommentHandler{repo: repo}
}

// List returns comments for a grade.
// Students see only non-internal comments; teacher/admin see all.
func (h *CommentHandler) List(c *gin.Context) {
	gradeID := c.Param("id")
	role, _ := c.Get(middleware.ContextRole)

	includeInternal := role.(string) == string(models.RoleTeacher) ||
		role.(string) == string(models.RoleAdmin)

	comments, err := h.repo.ListByGrade(gradeID, includeInternal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
		return
	}
	if comments == nil {
		comments = []*models.Comment{}
	}

	c.JSON(http.StatusOK, comments)
}

// createCommentRequest is the expected body for POST /grades/:id/comments.
type createCommentRequest struct {
	Body       string `json:"body"        binding:"required"`
	IsInternal bool   `json:"is_internal"`
}

// Create adds a new comment to the grade.
// Students cannot create internal comments.
func (h *CommentHandler) Create(c *gin.Context) {
	gradeID := c.Param("id")
	authorID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Students cannot post internal comments
	if role.(string) == string(models.RoleStudent) {
		req.IsInternal = false
	}

	comment, err := h.repo.Create(gradeID, authorID.(string), req.Body, req.IsInternal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}
