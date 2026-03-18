package repository

import (
	"database/sql"

	"github.com/servicedesk/backend/internal/models"
)

// StatsRepo handles analytics queries.
type StatsRepo struct {
	db *sql.DB
}

func NewStatsRepo(db *sql.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

// Dashboard returns aggregated metrics for the dashboard.
func (r *StatsRepo) Dashboard() (*models.DashboardStats, error) {
	stats := &models.DashboardStats{
		GradesByStatus: make(map[string]int),
	}

	// -- grades_by_status --
	statusRows, err := r.db.Query(`SELECT status, COUNT(*) FROM grades GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer statusRows.Close()

	totalGrades := 0
	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats.GradesByStatus[status] = count
		totalGrades += count
	}
	if err := statusRows.Err(); err != nil {
		return nil, err
	}
	stats.TotalGrades = totalGrades

	// -- pending_grades --
	var pendingCount int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM grades WHERE status = 'pending'`).Scan(&pendingCount); err != nil {
		return nil, err
	}
	stats.PendingGrades = pendingCount

	// -- avg_score --
	var avgScore sql.NullFloat64
	if err := r.db.QueryRow(`
		SELECT AVG(score)
		FROM grades
		WHERE score IS NOT NULL`).Scan(&avgScore); err != nil {
		return nil, err
	}
	if avgScore.Valid {
		stats.AvgScore = &avgScore.Float64
	}

	// -- top_teachers --
	teacherRows, err := r.db.Query(`
		SELECT u.id, u.full_name, COUNT(*) AS count
		FROM grades g
		JOIN users u ON u.id = g.teacher_id
		WHERE g.status IN ('certified', 'closed')
		GROUP BY u.id, u.full_name
		ORDER BY count DESC
		LIMIT 5`)
	if err != nil {
		return nil, err
	}
	defer teacherRows.Close()

	var topTeachers []models.TeacherStat
	for teacherRows.Next() {
		var t models.TeacherStat
		if err := teacherRows.Scan(&t.ID, &t.FullName, &t.Count); err != nil {
			return nil, err
		}
		topTeachers = append(topTeachers, t)
	}
	if err := teacherRows.Err(); err != nil {
		return nil, err
	}
	if topTeachers == nil {
		topTeachers = []models.TeacherStat{}
	}
	stats.TopTeachers = topTeachers

	return stats, nil
}
