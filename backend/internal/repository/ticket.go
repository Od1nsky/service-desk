package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/servicedesk/backend/internal/models"
)

// TicketFilter holds parameters for listing tickets.
type TicketFilter struct {
	Status     string
	PriorityID int
	CategoryID int
	AuthorID   string
	Page       int
	Limit      int
}

// TicketRepo handles all ticket-related database operations.
type TicketRepo struct {
	db *sql.DB
}

func NewTicketRepo(db *sql.DB) *TicketRepo {
	return &TicketRepo{db: db}
}

// Create inserts a new ticket and returns the created record.
func (r *TicketRepo) Create(title, description, authorID string, categoryID, priorityID int) (*models.Ticket, error) {
	const q = `
		INSERT INTO tickets (title, description, author_id, category_id, priority_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, number, title, description, status, author_id, assignee_id,
		          category_id, priority_id, rating, created_at, updated_at, closed_at`

	t := &models.Ticket{}
	var assigneeID sql.NullString
	var rating sql.NullInt64
	var closedAt sql.NullTime

	err := r.db.QueryRow(q, title, description, authorID, categoryID, priorityID).
		Scan(
			&t.ID, &t.Number, &t.Title, &t.Description, &t.Status,
			&t.AuthorID, &assigneeID,
			&t.CategoryID, &t.PriorityID,
			&rating, &t.CreatedAt, &t.UpdatedAt, &closedAt,
		)
	if err != nil {
		return nil, err
	}

	if assigneeID.Valid {
		t.AssigneeID = &assigneeID.String
	}
	if rating.Valid {
		v := int(rating.Int64)
		t.Rating = &v
	}
	if closedAt.Valid {
		t.ClosedAt = &closedAt.Time
	}
	return t, nil
}

// GetByID returns a fully-populated ticket including nested Author, Assignee, Category, Priority.
func (r *TicketRepo) GetByID(id string) (*models.Ticket, error) {
	const q = `
		SELECT
			t.id, t.number, t.title, t.description, t.status,
			t.author_id, t.assignee_id,
			t.category_id, t.priority_id,
			t.rating, t.created_at, t.updated_at, t.closed_at,
			-- author
			ua.id, ua.full_name, ua.email,
			-- assignee (nullable)
			ub.id, ub.full_name, ub.email,
			-- category
			c.id, c.name, c.description, c.created_at,
			-- priority
			p.id, p.name, p.color, p.sla_hours
		FROM tickets t
		JOIN users ua ON ua.id = t.author_id
		LEFT JOIN users ub ON ub.id = t.assignee_id
		JOIN categories c ON c.id = t.category_id
		JOIN priorities p ON p.id = t.priority_id
		WHERE t.id = $1`

	t := &models.Ticket{}
	author := &models.UserInfo{}
	category := &models.Category{}
	priority := &models.Priority{}

	var assigneeID sql.NullString
	var rating sql.NullInt64
	var closedAt sql.NullTime

	// Nullable assignee fields
	var assigneeUserID sql.NullString
	var assigneeFullName sql.NullString
	var assigneeEmail sql.NullString

	err := r.db.QueryRow(q, id).Scan(
		&t.ID, &t.Number, &t.Title, &t.Description, &t.Status,
		&t.AuthorID, &assigneeID,
		&t.CategoryID, &t.PriorityID,
		&rating, &t.CreatedAt, &t.UpdatedAt, &closedAt,
		&author.ID, &author.FullName, &author.Email,
		&assigneeUserID, &assigneeFullName, &assigneeEmail,
		&category.ID, &category.Name, &category.Description, &category.CreatedAt,
		&priority.ID, &priority.Name, &priority.Color, &priority.SLAHours,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if assigneeID.Valid {
		t.AssigneeID = &assigneeID.String
	}
	if rating.Valid {
		v := int(rating.Int64)
		t.Rating = &v
	}
	if closedAt.Valid {
		t.ClosedAt = &closedAt.Time
	}

	t.Author = author
	t.Category = category
	t.Priority = priority

	if assigneeUserID.Valid {
		t.Assignee = &models.UserInfo{
			ID:       assigneeUserID.String,
			FullName: assigneeFullName.String,
			Email:    assigneeEmail.String,
		}
	}

	return t, nil
}

// List returns a paginated slice of tickets matching the filter, plus the total count.
func (r *TicketRepo) List(f TicketFilter) ([]*models.Ticket, int, error) {
	args := []interface{}{}
	argN := 1
	where := " WHERE 1=1"

	if f.Status != "" {
		where += fmt.Sprintf(" AND t.status = $%d", argN)
		args = append(args, f.Status)
		argN++
	}
	if f.PriorityID > 0 {
		where += fmt.Sprintf(" AND t.priority_id = $%d", argN)
		args = append(args, f.PriorityID)
		argN++
	}
	if f.CategoryID > 0 {
		where += fmt.Sprintf(" AND t.category_id = $%d", argN)
		args = append(args, f.CategoryID)
		argN++
	}
	if f.AuthorID != "" {
		where += fmt.Sprintf(" AND t.author_id = $%d", argN)
		args = append(args, f.AuthorID)
		argN++
	}

	// Count query
	countQuery := "SELECT COUNT(*) FROM tickets t" + where
	var total int
	if err := r.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 20
	}
	offset := (f.Page - 1) * f.Limit

	listQuery := fmt.Sprintf(`
		SELECT
			t.id, t.number, t.title, t.status,
			t.author_id, t.assignee_id,
			t.category_id, t.priority_id,
			t.rating, t.created_at, t.updated_at, t.closed_at,
			ua.id, ua.full_name, ua.email,
			c.id, c.name, c.description, c.created_at,
			p.id, p.name, p.color, p.sla_hours
		FROM tickets t
		JOIN users ua ON ua.id = t.author_id
		LEFT JOIN users ub ON ub.id = t.assignee_id
		JOIN categories c ON c.id = t.category_id
		JOIN priorities p ON p.id = t.priority_id
		%s
		ORDER BY t.created_at DESC
		LIMIT $%d OFFSET $%d`, where, argN, argN+1)

	args = append(args, f.Limit, offset)

	rows, err := r.db.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tickets []*models.Ticket
	for rows.Next() {
		t := &models.Ticket{}
		author := &models.UserInfo{}
		category := &models.Category{}
		priority := &models.Priority{}

		var assigneeID sql.NullString
		var rating sql.NullInt64
		var closedAt sql.NullTime

		if err := rows.Scan(
			&t.ID, &t.Number, &t.Title, &t.Status,
			&t.AuthorID, &assigneeID,
			&t.CategoryID, &t.PriorityID,
			&rating, &t.CreatedAt, &t.UpdatedAt, &closedAt,
			&author.ID, &author.FullName, &author.Email,
			&category.ID, &category.Name, &category.Description, &category.CreatedAt,
			&priority.ID, &priority.Name, &priority.Color, &priority.SLAHours,
		); err != nil {
			return nil, 0, err
		}

		if assigneeID.Valid {
			t.AssigneeID = &assigneeID.String
		}
		if rating.Valid {
			v := int(rating.Int64)
			t.Rating = &v
		}
		if closedAt.Valid {
			t.ClosedAt = &closedAt.Time
		}

		t.Author = author
		t.Category = category
		t.Priority = priority

		tickets = append(tickets, t)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return tickets, total, nil
}

// UpdateStatus changes a ticket's status inside a transaction and records the change in status_history.
func (r *TicketRepo) UpdateStatus(id, oldStatus, newStatus, changedByID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var closedAtExpr string
	if newStatus == string(models.StatusClosed) {
		closedAtExpr = ", closed_at = NOW()"
	}

	updateQ := fmt.Sprintf(
		`UPDATE tickets SET status = $1, updated_at = NOW()%s WHERE id = $2`,
		closedAtExpr,
	)
	if _, err = tx.Exec(updateQ, newStatus, id); err != nil {
		return err
	}

	const histQ = `
		INSERT INTO status_history (ticket_id, changed_by, old_status, new_status)
		VALUES ($1, $2, $3, $4)`
	if _, err = tx.Exec(histQ, id, changedByID, oldStatus, newStatus); err != nil {
		return err
	}

	return tx.Commit()
}

// AssignTicket sets the assignee for a ticket.
func (r *TicketRepo) AssignTicket(ticketID, assigneeID string) error {
	const q = `UPDATE tickets SET assignee_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(q, assigneeID, ticketID)
	return err
}

// UpdateRating sets the rating on a ticket.
func (r *TicketRepo) UpdateRating(ticketID string, rating int) error {
	const q = `UPDATE tickets SET rating = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(q, rating, ticketID)
	return err
}

// GetStatusHistory returns all status changes for a ticket ordered chronologically.
func (r *TicketRepo) GetStatusHistory(ticketID string) ([]*models.StatusHistory, error) {
	const q = `
		SELECT sh.id, sh.ticket_id, sh.changed_by, sh.old_status, sh.new_status, sh.changed_at,
		       u.id, u.full_name, u.email
		FROM status_history sh
		JOIN users u ON u.id = sh.changed_by
		WHERE sh.ticket_id = $1
		ORDER BY sh.changed_at ASC`

	rows, err := r.db.Query(q, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.StatusHistory
	for rows.Next() {
		sh := &models.StatusHistory{}
		user := &models.UserInfo{}
		if err := rows.Scan(
			&sh.ID, &sh.TicketID, &sh.ChangedBy, &sh.OldStatus, &sh.NewStatus, &sh.ChangedAt,
			&user.ID, &user.FullName, &user.Email,
		); err != nil {
			return nil, err
		}
		sh.User = user
		history = append(history, sh)
	}
	return history, rows.Err()
}

// ensure time import is used (used via sql.NullTime indirectly through time.Time)
var _ = time.Now
