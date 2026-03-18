package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/servicedesk/backend/internal/middleware"
	"github.com/servicedesk/backend/internal/models"
	"github.com/servicedesk/backend/internal/repository"
)

// UserHandler handles user management endpoints (admin only).
type UserHandler struct {
	repo *repository.UserRepo
}

func NewUserHandler(repo *repository.UserRepo) *UserHandler {
	return &UserHandler{repo: repo}
}

// List returns all users.
func (h *UserHandler) List(c *gin.Context) {
	users, err := h.repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	if users == nil {
		users = []*models.User{}
	}
	c.JSON(http.StatusOK, users)
}

// updateRoleRequest is the expected body for PATCH /users/:id/role.
type updateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

// UpdateRole changes the role of a user. Admin cannot change their own role.
func (h *UserHandler) UpdateRole(c *gin.Context) {
	targetID := c.Param("id")
	callerID, _ := c.Get(middleware.ContextUserID)

	if targetID == callerID.(string) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot change your own role"})
		return
	}

	var req updateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validRoles := map[string]bool{
		string(models.RoleStudent): true,
		string(models.RoleTeacher): true,
		string(models.RoleAdmin):   true,
	}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role; must be one of: student, teacher, admin"})
		return
	}

	user, err := h.repo.GetByID(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := h.repo.UpdateRole(targetID, models.Role(req.Role)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

// blockRequest is the expected body for PATCH /users/:id/block.
type blockRequest struct {
	IsActive bool `json:"is_active"`
}

// Block enables or disables a user account. Admin cannot block themselves.
func (h *UserHandler) Block(c *gin.Context) {
	targetID := c.Param("id")
	callerID, _ := c.Get(middleware.ContextUserID)

	if targetID == callerID.(string) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot block yourself"})
		return
	}

	var req blockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.repo.GetByID(targetID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if err := h.repo.UpdateActive(targetID, req.IsActive); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}
