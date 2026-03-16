package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/servicedesk/backend/internal/repository"
)

// StatsHandler handles analytics/dashboard endpoints.
type StatsHandler struct {
	repo *repository.StatsRepo
}

func NewStatsHandler(repo *repository.StatsRepo) *StatsHandler {
	return &StatsHandler{repo: repo}
}

// Dashboard returns aggregated metrics for the admin/support dashboard.
func (h *StatsHandler) Dashboard(c *gin.Context) {
	data, err := h.repo.Dashboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch dashboard data"})
		return
	}
	c.JSON(http.StatusOK, data)
}
