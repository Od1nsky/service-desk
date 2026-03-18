package repository

import (
	"database/sql"

	"github.com/servicedesk/backend/internal/models"
)

// GradeTypeRepo handles all grade_type-related database operations.
type GradeTypeRepo struct {
	db *sql.DB
}

func NewGradeTypeRepo(db *sql.DB) *GradeTypeRepo {
	return &GradeTypeRepo{db: db}
}

// List returns all grade types ordered by id.
func (r *GradeTypeRepo) List() ([]*models.GradeType, error) {
	const q = `SELECT id, name, color, max_score FROM grade_types ORDER BY id`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gradeTypes []*models.GradeType
	for rows.Next() {
		gt := &models.GradeType{}
		if err := rows.Scan(&gt.ID, &gt.Name, &gt.Color, &gt.MaxScore); err != nil {
			return nil, err
		}
		gradeTypes = append(gradeTypes, gt)
	}
	return gradeTypes, rows.Err()
}
