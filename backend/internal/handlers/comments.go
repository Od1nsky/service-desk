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

// List returns comments for a ticket.
// Employees see only non-internal comments; support/admin see all.
func (h *CommentHandler) List(c *gin.Context) {
	ticketID := c.Param("id")
	role, _ := c.Get(middleware.ContextRole)

	includeInternal := role.(string) == string(models.RoleSupport) ||
		role.(string) == string(models.RoleAdmin)

	comments, err := h.repo.ListByTicket(ticketID, includeInternal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch comments"})
		return
	}
	if comments == nil {
		comments = []*models.Comment{}
	}

	c.JSON(http.StatusOK, comments)
}

// createCommentRequest is the expected body for POST /tickets/:id/comments.
type createCommentRequest struct {
	Body       string `json:"body"        binding:"required"`
	IsInternal bool   `json:"is_internal"`
}

// Create adds a new comment to the ticket.
// Employees cannot create internal comments.
func (h *CommentHandler) Create(c *gin.Context) {
	ticketID := c.Param("id")
	authorID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)

	var req createCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Employees cannot post internal comments
	if role.(string) == string(models.RoleEmployee) {
		req.IsInternal = false
	}

	comment, err := h.repo.Create(ticketID, authorID.(string), req.Body, req.IsInternal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}
