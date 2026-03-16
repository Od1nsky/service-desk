package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/servicedesk/backend/internal/models"
)

// CategoryRepo handles all category-related database operations.
type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepo(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db: db}
}

// List returns all categories ordered by id.
func (r *CategoryRepo) List() ([]*models.Category, error) {
	const q = `SELECT id, name, description, created_at FROM categories ORDER BY id`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cats []*models.Category
	for rows.Next() {
		c := &models.Category{}
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt); err != nil {
			return nil, err
		}
		cats = append(cats, c)
	}
	return cats, rows.Err()
}

// GetByID returns the category with the given id, or nil if not found.
func (r *CategoryRepo) GetByID(id int) (*models.Category, error) {
	const q = `SELECT id, name, description, created_at FROM categories WHERE id = $1`

	c := &models.Category{}
	err := r.db.QueryRow(q, id).Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Create inserts a new category.
func (r *CategoryRepo) Create(name, description string) (*models.Category, error) {
	const q = `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at`

	c := &models.Category{}
	err := r.db.QueryRow(q, name, description).
		Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Update modifies name and description for an existing category and returns the updated record.
func (r *CategoryRepo) Update(id int, name, description string) (*models.Category, error) {
	const q = `
		UPDATE categories SET name = $1, description = $2
		WHERE id = $3
		RETURNING id, name, description, created_at`

	c := &models.Category{}
	err := r.db.QueryRow(q, name, description, id).
		Scan(&c.ID, &c.Name, &c.Description, &c.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Delete removes a category. Returns an error if any tickets reference this category.
func (r *CategoryRepo) Delete(id int) error {
	var count int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM tickets WHERE category_id = $1`, id).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("cannot delete category: %d ticket(s) are referencing it", count)
	}
	_, err := r.db.Exec(`DELETE FROM categories WHERE id = $1`, id)
	return err
}
