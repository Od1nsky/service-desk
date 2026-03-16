package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/servicedesk/backend/internal/middleware"
	"github.com/servicedesk/backend/internal/models"
	"github.com/servicedesk/backend/internal/repository"
)

// TicketHandler handles ticket-related endpoints.
type TicketHandler struct {
	repo *repository.TicketRepo
}

func NewTicketHandler(repo *repository.TicketRepo) *TicketHandler {
	return &TicketHandler{repo: repo}
}

// List returns a paginated list of tickets.
// Employees only see their own tickets; support/admin see all.
func (h *TicketHandler) List(c *gin.Context) {
	role, _ := c.Get(middleware.ContextRole)
	userID, _ := c.Get(middleware.ContextUserID)

	f := repository.TicketFilter{}

	if s := c.Query("status"); s != "" {
		f.Status = s
	}
	if p := c.Query("priority_id"); p != "" {
		if id, err := strconv.Atoi(p); err == nil {
			f.PriorityID = id
		}
	}
	if cat := c.Query("category_id"); cat != "" {
		if id, err := strconv.Atoi(cat); err == nil {
			f.CategoryID = id
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	f.Page = page
	f.Limit = limit

	// Employees may only view their own tickets
	if role.(string) == string(models.RoleEmployee) {
		f.AuthorID = userID.(string)
	}

	tickets, total, err := h.repo.List(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch tickets"})
		return
	}
	if tickets == nil {
		tickets = []*models.Ticket{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tickets,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// createTicketRequest is the expected body for POST /tickets.
type createTicketRequest struct {
	Title       string `json:"title"       binding:"required"`
	Description string `json:"description"`
	CategoryID  int    `json:"category_id" binding:"required,min=1"`
	PriorityID  int    `json:"priority_id" binding:"required,min=1"`
}

// Create creates a new ticket. Only employees may do this.
func (h *TicketHandler) Create(c *gin.Context) {
	var req createTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID, _ := c.Get(middleware.ContextUserID)

	ticket, err := h.repo.Create(req.Title, req.Description, authorID.(string), req.CategoryID, req.PriorityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ticket"})
		return
	}

	c.JSON(http.StatusCreated, ticket)
}

// GetByID returns a single ticket with its full details and status history.
// Employees may only view their own tickets.
func (h *TicketHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get(middleware.ContextRole)
	userID, _ := c.Get(middleware.ContextUserID)

	ticket, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	// Employees can only see their own tickets
	if role.(string) == string(models.RoleEmployee) && ticket.AuthorID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	history, err := h.repo.GetStatusHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch history"})
		return
	}
	if history == nil {
		history = []*models.StatusHistory{}
	}

	c.JSON(http.StatusOK, gin.H{
		"ticket":  ticket,
		"history": history,
	})
}

// updateStatusRequest is the expected body for PATCH /tickets/:id/status.
type updateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus transitions a ticket to a new status. Only support/admin.
func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	changedByID, _ := c.Get(middleware.ContextUserID)

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	newStatus := models.Status(req.Status)
	if !ticket.Status.CanTransitionTo(newStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status transition from " + string(ticket.Status) + " to " + req.Status,
		})
		return
	}

	if err := h.repo.UpdateStatus(id, string(ticket.Status), req.Status, changedByID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

// Assign assigns the ticket to the currently authenticated support user.
func (h *TicketHandler) Assign(c *gin.Context) {
	id := c.Param("id")
	assigneeID, _ := c.Get(middleware.ContextUserID)

	ticket, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	if err := h.repo.AssignTicket(id, assigneeID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket assigned"})
}

// rateRequest is the expected body for PATCH /tickets/:id/rate.
type rateRequest struct {
	Rating int `json:"rating" binding:"required,min=1,max=5"`
}

// Rate lets the ticket author rate a resolved ticket.
func (h *TicketHandler) Rate(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get(middleware.ContextUserID)

	var req rateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}

	if ticket.AuthorID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only the ticket author can rate"})
		return
	}
	if ticket.Status != models.StatusResolved {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can only rate tickets in 'resolved' status"})
		return
	}

	if err := h.repo.UpdateRating(id, req.Rating); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update rating"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "rating saved"})
}
