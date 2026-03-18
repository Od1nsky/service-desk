package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/servicedesk/backend/internal/middleware"
	"github.com/servicedesk/backend/internal/models"
	"github.com/servicedesk/backend/internal/repository"
)

// GradeHandler handles grade-related endpoints.
type GradeHandler struct {
	repo *repository.GradeRepo
}

func NewGradeHandler(repo *repository.GradeRepo) *GradeHandler {
	return &GradeHandler{repo: repo}
}

// List returns a paginated list of grades.
// Students only see their own grades; teacher/admin see all.
func (h *GradeHandler) List(c *gin.Context) {
	role, _ := c.Get(middleware.ContextRole)
	userID, _ := c.Get(middleware.ContextUserID)

	f := repository.GradeFilter{}

	if s := c.Query("status"); s != "" {
		f.Status = s
	}
	if gt := c.Query("grade_type_id"); gt != "" {
		if id, err := strconv.Atoi(gt); err == nil {
			f.GradeTypeID = id
		}
	}
	if dis := c.Query("discipline_id"); dis != "" {
		if id, err := strconv.Atoi(dis); err == nil {
			f.DisciplineID = id
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

	// Students may only view their own grades
	if role.(string) == string(models.RoleStudent) {
		f.StudentID = userID.(string)
	}

	grades, total, err := h.repo.List(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch grades"})
		return
	}
	if grades == nil {
		grades = []*models.Grade{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  grades,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// createGradeRequest is the expected body for POST /grades.
type createGradeRequest struct {
	TaskName     string `json:"task_name"     binding:"required"`
	Notes        string `json:"notes"`
	DisciplineID int    `json:"discipline_id" binding:"required,min=1"`
	GradeTypeID  int    `json:"grade_type_id" binding:"required,min=1"`
}

// Create creates a new grade. Only students may do this.
func (h *GradeHandler) Create(c *gin.Context) {
	var req createGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentID, _ := c.Get(middleware.ContextUserID)

	grade, err := h.repo.Create(req.TaskName, req.Notes, studentID.(string), req.DisciplineID, req.GradeTypeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create grade"})
		return
	}

	c.JSON(http.StatusCreated, grade)
}

// GetByID returns a single grade with its full details and grade history.
// Students may only view their own grades.
func (h *GradeHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	role, _ := c.Get(middleware.ContextRole)
	userID, _ := c.Get(middleware.ContextUserID)

	grade, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if grade == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}

	// Students can only see their own grades
	if role.(string) == string(models.RoleStudent) && grade.StudentID != userID.(string) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	history, err := h.repo.GetGradeHistory(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch history"})
		return
	}
	if history == nil {
		history = []*models.GradeHistory{}
	}

	c.JSON(http.StatusOK, gin.H{
		"grade":   grade,
		"history": history,
	})
}

// updateStatusRequest is the expected body for PATCH /grades/:id/status.
type updateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus transitions a grade to a new status. Only teacher/admin.
func (h *GradeHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	changedByID, _ := c.Get(middleware.ContextUserID)

	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if grade == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}

	newStatus := models.Status(req.Status)
	if !grade.Status.CanTransitionTo(newStatus) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid status transition from " + string(grade.Status) + " to " + req.Status,
		})
		return
	}

	if err := h.repo.UpdateStatus(id, string(grade.Status), req.Status, changedByID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

// Assign assigns the grade to the currently authenticated teacher.
func (h *GradeHandler) Assign(c *gin.Context) {
	id := c.Param("id")
	teacherID, _ := c.Get(middleware.ContextUserID)

	grade, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if grade == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}

	if err := h.repo.AssignGrade(id, teacherID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign grade"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "grade assigned"})
}

// acknowledgeRequest is the expected body for PATCH /grades/:id/acknowledge.
type acknowledgeRequest struct {
	Score int `json:"score" binding:"required,min=0,max=100"`
}

// Acknowledge lets the teacher set the score on a graded record.
func (h *GradeHandler) Acknowledge(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get(middleware.ContextUserID)
	role, _ := c.Get(middleware.ContextRole)

	var req acknowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	grade, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}
	if grade == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "grade not found"})
		return
	}

	// Teachers can only set scores on grades assigned to them; admin can set any
	if role.(string) == string(models.RoleTeacher) {
		if grade.TeacherID == nil || *grade.TeacherID != userID.(string) {
			c.JSON(http.StatusForbidden, gin.H{"error": "only the assigned teacher can set the score"})
			return
		}
	}

	if err := h.repo.UpdateScore(id, req.Score); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update score"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "score saved"})
}
