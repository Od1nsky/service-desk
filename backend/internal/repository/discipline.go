package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/servicedesk/backend/internal/models"
)

// DisciplineRepo handles all discipline-related database operations.
type DisciplineRepo struct {
	db *sql.DB
}

func NewDisciplineRepo(db *sql.DB) *DisciplineRepo {
	return &DisciplineRepo{db: db}
}

// List returns all disciplines ordered by id.
func (r *DisciplineRepo) List() ([]*models.Discipline, error) {
	const q = `SELECT id, name, description, created_at FROM disciplines ORDER BY id`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var disciplines []*models.Discipline
	for rows.Next() {
		d := &models.Discipline{}
		if err := rows.Scan(&d.ID, &d.Name, &d.Description, &d.CreatedAt); err != nil {
			return nil, err
		}
		disciplines = append(disciplines, d)
	}
	return disciplines, rows.Err()
}

// GetByID returns the discipline with the given id, or nil if not found.
func (r *DisciplineRepo) GetByID(id int) (*models.Discipline, error) {
	const q = `SELECT id, name, description, created_at FROM disciplines WHERE id = $1`

	d := &models.Discipline{}
	err := r.db.QueryRow(q, id).Scan(&d.ID, &d.Name, &d.Description, &d.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Create inserts a new discipline.
func (r *DisciplineRepo) Create(name, description string) (*models.Discipline, error) {
	const q = `
		INSERT INTO disciplines (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at`

	d := &models.Discipline{}
	err := r.db.QueryRow(q, name, description).
		Scan(&d.ID, &d.Name, &d.Description, &d.CreatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Update modifies name and description for an existing discipline and returns the updated record.
func (r *DisciplineRepo) Update(id int, name, description string) (*models.Discipline, error) {
	const q = `
		UPDATE disciplines SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description, created_at`

	d := &models.Discipline{}
	err := r.db.QueryRow(q, name, description, id).
		Scan(&d.ID, &d.Name, &d.Description, &d.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return d, nil
}

// Delete removes a discipline. Returns an error if any grades reference this discipline.
func (r *DisciplineRepo) Delete(id int) error {
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM grades WHERE discipline_id = $1`, id).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete discipline: %d grade(s) are referencing it", count)
	}
	_, err := r.db.Exec(`DELETE FROM disciplines WHERE id = $1`, id)
	return err
}
