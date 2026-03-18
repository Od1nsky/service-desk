package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/servicedesk/backend/internal/models"
)

// GradeFilter holds parameters for listing grades.
type GradeFilter struct {
	Status       string
	GradeTypeID  int
	DisciplineID int
	StudentID    string
	Page         int
	Limit        int
}

// GradeRepo handles all grade-related database operations.
type GradeRepo struct {
	db *sql.DB
}

func NewGradeRepo(db *sql.DB) *GradeRepo {
	return &GradeRepo{db: db}
}

// Create inserts a new grade and returns the created record.
func (r *GradeRepo) Create(taskName, notes, studentID string, disciplineID, gradeTypeID int) (*models.Grade, error) {
	const q = `
		INSERT INTO grades (task_name, notes, student_id, discipline_id, grade_type_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, number, task_name, notes, status, student_id, teacher_id,
		          discipline_id, grade_type_id, score, created_at, updated_at, graded_at`

	g := &models.Grade{}
	var teacherID sql.NullString
	var score sql.NullInt64
	var gradedAt sql.NullTime

	err := r.db.QueryRow(q, taskName, notes, studentID, disciplineID, gradeTypeID).
		Scan(
			&g.ID, &g.Number, &g.TaskName, &g.Notes, &g.Status,
			&g.StudentID, &teacherID,
			&g.DisciplineID, &g.GradeTypeID,
			&score, &g.CreatedAt, &g.UpdatedAt, &gradedAt,
		)
	if err != nil {
		return nil, err
	}

	if teacherID.Valid {
		g.TeacherID = &teacherID.String
	}
	if score.Valid {
		v := int(score.Int64)
		g.Score = &v
	}
	if gradedAt.Valid {
		g.GradedAt = &gradedAt.Time
	}
	return g, nil
}

// GetByID returns a fully-populated grade including nested Student, Teacher, Discipline, GradeType.
func (r *GradeRepo) GetByID(id string) (*models.Grade, error) {
	const q = `
		SELECT
			g.id, g.number, g.task_name, g.notes, g.status,
			g.student_id, g.teacher_id,
			g.discipline_id, g.grade_type_id,
			g.score, g.created_at, g.updated_at, g.graded_at,
			-- student
			us.id, us.full_name, us.email,
			-- teacher (nullable)
			ut.id, ut.full_name, ut.email,
			-- discipline
			d.id, d.name, d.description, d.created_at,
			-- grade_type
			gt.id, gt.name, gt.color, gt.max_score
		FROM grades g
		JOIN users us ON us.id = g.student_id
		LEFT JOIN users ut ON ut.id = g.teacher_id
		JOIN disciplines d ON d.id = g.discipline_id
		JOIN grade_types gt ON gt.id = g.grade_type_id
		WHERE g.id = $1`

	g := &models.Grade{}
	student := &models.UserInfo{}
	discipline := &models.Discipline{}
	gradeType := &models.GradeType{}

	var teacherID sql.NullString
	var score sql.NullInt64
	var gradedAt sql.NullTime

	// Nullable teacher fields
	var teacherUserID sql.NullString
	var teacherFullName sql.NullString
	var teacherEmail sql.NullString

	err := r.db.QueryRow(q, id).Scan(
		&g.ID, &g.Number, &g.TaskName, &g.Notes, &g.Status,
		&g.StudentID, &teacherID,
		&g.DisciplineID, &g.GradeTypeID,
		&score, &g.CreatedAt, &g.UpdatedAt, &gradedAt,
		&student.ID, &student.FullName, &student.Email,
		&teacherUserID, &teacherFullName, &teacherEmail,
		&discipline.ID, &discipline.Name, &discipline.Description, &discipline.CreatedAt,
		&gradeType.ID, &gradeType.Name, &gradeType.Color, &gradeType.MaxScore,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if teacherID.Valid {
		g.TeacherID = &teacherID.String
	}
	if score.Valid {
		v := int(score.Int64)
		g.Score = &v
	}
	if gradedAt.Valid {
		g.GradedAt = &gradedAt.Time
	}

	g.Student = student
	g.Discipline = discipline
	g.GradeType = gradeType

	if teacherUserID.Valid {
		g.Teacher = &models.UserInfo{
			ID:       teacherUserID.String,
			FullName: teacherFullName.String,
			Email:    teacherEmail.String,
		}
	}

	return g, nil
}

// List returns a paginated slice of grades matching the filter, plus the total count.
func (r *GradeRepo) List(f GradeFilter) ([]*models.Grade, int, error) {
	args := []interface{}{}
	argN := 1
	where := " WHERE 1=1"

	if f.Status != "" {
		where += fmt.Sprintf(" AND g.status = $%d", argN)
		args = append(args, f.Status)
		argN++
	}
	if f.GradeTypeID > 0 {
		where += fmt.Sprintf(" AND g.grade_type_id = $%d", argN)
		args = append(args, f.GradeTypeID)
		argN++
	}
	if f.DisciplineID > 0 {
		where += fmt.Sprintf(" AND g.discipline_id = $%d", argN)
		args = append(args, f.DisciplineID)
		argN++
	}
	if f.StudentID != "" {
		where += fmt.Sprintf(" AND g.student_id = $%d", argN)
		args = append(args, f.StudentID)
		argN++
	}

	// Count query
	countQuery := "SELECT COUNT(*) FROM grades g" + where
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
			g.id, g.number, g.task_name, g.status,
			g.student_id, g.teacher_id,
			g.discipline_id, g.grade_type_id,
			g.score, g.created_at, g.updated_at, g.graded_at,
			us.id, us.full_name, us.email,
			d.id, d.name, d.description, d.created_at,
			gt.id, gt.name, gt.color, gt.max_score
		FROM grades g
		JOIN users us ON us.id = g.student_id
		LEFT JOIN users ut ON ut.id = g.teacher_id
		JOIN disciplines d ON d.id = g.discipline_id
		JOIN grade_types gt ON gt.id = g.grade_type_id
		%s
		ORDER BY g.created_at DESC
		LIMIT $%d OFFSET $%d`, where, argN, argN+1)

	args = append(args, f.Limit, offset)

	rows, err := r.db.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var grades []*models.Grade
	for rows.Next() {
		g := &models.Grade{}
		student := &models.UserInfo{}
		discipline := &models.Discipline{}
		gradeType := &models.GradeType{}

		var teacherID sql.NullString
		var score sql.NullInt64
		var gradedAt sql.NullTime

		if err := rows.Scan(
			&g.ID, &g.Number, &g.TaskName, &g.Status,
			&g.StudentID, &teacherID,
			&g.DisciplineID, &g.GradeTypeID,
			&score, &g.CreatedAt, &g.UpdatedAt, &gradedAt,
			&student.ID, &student.FullName, &student.Email,
			&discipline.ID, &discipline.Name, &discipline.Description, &discipline.CreatedAt,
			&gradeType.ID, &gradeType.Name, &gradeType.Color, &gradeType.MaxScore,
		); err != nil {
			return nil, 0, err
		}

		if teacherID.Valid {
			g.TeacherID = &teacherID.String
		}
		if score.Valid {
			v := int(score.Int64)
			g.Score = &v
		}
		if gradedAt.Valid {
			g.GradedAt = &gradedAt.Time
		}

		g.Student = student
		g.Discipline = discipline
		g.GradeType = gradeType

		grades = append(grades, g)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return grades, total, nil
}

// UpdateStatus changes a grade's status inside a transaction and records the change in grade_history.
func (r *GradeRepo) UpdateStatus(id, oldStatus, newStatus, changedByID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var gradedAtExpr string
	if newStatus == string(models.StatusGraded) {
		gradedAtExpr = ", graded_at = NOW()"
	}

	updateQ := fmt.Sprintf(
		`UPDATE grades SET status = $1, updated_at = NOW()%s WHERE id = $2`,
		gradedAtExpr,
	)
	if _, err = tx.Exec(updateQ, newStatus, id); err != nil {
		return err
	}

	const histQ = `
		INSERT INTO grade_history (grade_id, changed_by, old_value, new_value)
		VALUES ($1, $2, $3, $4)`
	if _, err = tx.Exec(histQ, id, changedByID, oldStatus, newStatus); err != nil {
		return err
	}

	return tx.Commit()
}

// AssignGrade sets the teacher for a grade.
func (r *GradeRepo) AssignGrade(gradeID, teacherID string) error {
	const q = `UPDATE grades SET teacher_id = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(q, teacherID, gradeID)
	return err
}

// UpdateScore sets the score on a grade.
func (r *GradeRepo) UpdateScore(gradeID string, score int) error {
	const q = `UPDATE grades SET score = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(q, score, gradeID)
	return err
}

// GetGradeHistory returns all status changes for a grade ordered chronologically.
func (r *GradeRepo) GetGradeHistory(gradeID string) ([]*models.GradeHistory, error) {
	const q = `
		SELECT gh.id, gh.grade_id, gh.changed_by, gh.old_value, gh.new_value, gh.changed_at,
		       u.id, u.full_name, u.email
		FROM grade_history gh
		JOIN users u ON u.id = gh.changed_by
		WHERE gh.grade_id = $1
		ORDER BY gh.changed_at ASC`

	rows, err := r.db.Query(q, gradeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []*models.GradeHistory
	for rows.Next() {
		gh := &models.GradeHistory{}
		user := &models.UserInfo{}
		if err := rows.Scan(
			&gh.ID, &gh.GradeID, &gh.ChangedBy, &gh.OldValue, &gh.NewValue, &gh.ChangedAt,
			&user.ID, &user.FullName, &user.Email,
		); err != nil {
			return nil, err
		}
		gh.User = user
		history = append(history, gh)
	}
	return history, rows.Err()
}

// ensure time import is used
var _ = time.Now
