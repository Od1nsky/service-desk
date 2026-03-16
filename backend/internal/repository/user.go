package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/servicedesk/backend/internal/models"
)

// UserRepo handles all user-related database operations.
type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user and returns the created record.
func (r *UserRepo) Create(email, passwordHash, fullName string, role models.Role) (*models.User, error) {
	const q = `
		INSERT INTO users (email, password, full_name, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, email, full_name, role, is_active, created_at, updated_at`

	u := &models.User{}
	err := r.db.QueryRow(q, email, passwordHash, fullName, string(role)).
		Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetByEmail returns the user with the given email, or nil if not found.
func (r *UserRepo) GetByEmail(email string) (*models.User, error) {
	const q = `
		SELECT id, email, password, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE email = $1`

	u := &models.User{}
	err := r.db.QueryRow(q, email).
		Scan(&u.ID, &u.Email, &u.Password, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetByID returns the user with the given UUID, or nil if not found.
func (r *UserRepo) GetByID(id string) (*models.User, error) {
	const q = `
		SELECT id, email, password, full_name, role, is_active, created_at, updated_at
		FROM users
		WHERE id = $1`

	u := &models.User{}
	err := r.db.QueryRow(q, id).
		Scan(&u.ID, &u.Email, &u.Password, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// List returns all users without the password field.
func (r *UserRepo) List() ([]*models.User, error) {
	const q = `
		SELECT id, email, full_name, role, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC`

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.FullName, &u.Role, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// UpdateRole sets a new role for the user identified by id.
func (r *UserRepo) UpdateRole(id string, role models.Role) error {
	const q = `UPDATE users SET role = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(q, string(role), id)
	return err
}

// UpdateActive sets is_active for the user identified by id.
func (r *UserRepo) UpdateActive(id string, isActive bool) error {
	const q = `UPDATE users SET is_active = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(q, isActive, id)
	return err
}

// SaveRefreshToken persists a new refresh token record.
func (r *UserRepo) SaveRefreshToken(userID, token string, expiresAt time.Time) error {
	const q = `
		INSERT INTO refresh_tokens (user_id, token, expires_at)
		VALUES ($1, $2, $3)`
	_, err := r.db.Exec(q, userID, token, expiresAt)
	return err
}

// GetRefreshToken looks up a refresh token by its value; returns nil if not found.
func (r *UserRepo) GetRefreshToken(token string) (*models.RefreshToken, error) {
	const q = `
		SELECT id, user_id, token, expires_at, created_at
		FROM refresh_tokens
		WHERE token = $1`

	rt := &models.RefreshToken{}
	err := r.db.QueryRow(q, token).
		Scan(&rt.ID, &rt.UserID, &rt.Token, &rt.ExpiresAt, &rt.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return rt, nil
}

// DeleteRefreshToken removes a specific refresh token.
func (r *UserRepo) DeleteRefreshToken(token string) error {
	const q = `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(q, token)
	return err
}

// DeleteUserRefreshTokens removes all refresh tokens belonging to a user.
func (r *UserRepo) DeleteUserRefreshTokens(userID string) error {
	const q = `DELETE FROM refresh_tokens WHERE user_id = $1`
	_, err := r.db.Exec(q, userID)
	return err
}
