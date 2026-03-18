package repository

import (
	"database/sql"
	"errors"

	"github.com/servicedesk/backend/internal/models"
)

// CommentRepo handles all comment-related database operations.
type CommentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{db: db}
}

// Create inserts a new comment and returns the populated record including author info.
func (r *CommentRepo) Create(gradeID, authorID, body string, isInternal bool) (*models.Comment, error) {
	const insertQ = `
		INSERT INTO comments (grade_id, author_id, body, is_internal)
		VALUES ($1, $2, $3, $4)
		RETURNING id, grade_id, author_id, body, is_internal, created_at`

	c := &models.Comment{}
	err := r.db.QueryRow(insertQ, gradeID, authorID, body, isInternal).
		Scan(&c.ID, &c.GradeID, &c.AuthorID, &c.Body, &c.IsInternal, &c.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Populate author
	const userQ = `SELECT id, full_name, email FROM users WHERE id = $1`
	author := &models.UserInfo{}
	err = r.db.QueryRow(userQ, authorID).Scan(&author.ID, &author.FullName, &author.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if err == nil {
		c.Author = author
	}

	return c, nil
}

// ListByGrade returns all comments for a grade. When includeInternal is false,
// internal comments are excluded.
func (r *CommentRepo) ListByGrade(gradeID string, includeInternal bool) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.grade_id, c.author_id, c.body, c.is_internal, c.created_at,
		       u.id, u.full_name, u.email
		FROM comments c
		JOIN users u ON u.id = c.author_id
		WHERE c.grade_id = $1`

	if !includeInternal {
		query += " AND c.is_internal = FALSE"
	}
	query += " ORDER BY c.created_at ASC"

	rows, err := r.db.Query(query, gradeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{}
		author := &models.UserInfo{}
		if err := rows.Scan(
			&c.ID, &c.GradeID, &c.AuthorID, &c.Body, &c.IsInternal, &c.CreatedAt,
			&author.ID, &author.FullName, &author.Email,
		); err != nil {
			return nil, err
		}
		c.Author = author
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
