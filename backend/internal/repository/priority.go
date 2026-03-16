package repository

import (
	"database/sql"

	"github.com/servicedesk/backend/internal/models"
)

// PriorityRepo handles all priority-related database operations.
type PriorityRepo struct {
	db *sql.DB
}

func NewPriorityRepo(db *sql.DB) *PriorityRepo {
	return &PriorityRepo{db: db}
}

// List returns all priorities ordered by sla_hours descending (lowest SLA = highest urgency last).
func (r *PriorityRepo) List() ([]*models.Priority, error) {
	const q = `SELECT id, name, color, sla_hours FROM priorities ORDER BY id`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var priorities []*models.Priority
	for rows.Next() {
		p := &models.Priority{}
		if err := rows.Scan(&p.ID, &p.Name, &p.Color, &p.SLAHours); err != nil {
			return nil, err
		}
		priorities = append(priorities, p)
	}
	return priorities, rows.Err()
}
