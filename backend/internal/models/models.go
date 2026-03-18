package models

import "time"

type Role string

const (
	RoleStudent Role = "student"
	RoleTeacher Role = "teacher"
	RoleAdmin   Role = "admin"
)

type Status string

const (
	StatusPending        Status = "pending"
	StatusGraded         Status = "graded"
	StatusNotCertified   Status = "not_certified"
	StatusCertified      Status = "certified"
	StatusClosed         Status = "closed"
)

var allowedTransitions = map[Status][]Status{
	StatusPending:      {StatusGraded},
	StatusGraded:       {StatusNotCertified, StatusCertified, StatusPending},
	StatusNotCertified: {StatusGraded},
	StatusCertified:    {StatusClosed},
	StatusClosed:       {},
}

func (s Status) CanTransitionTo(next Status) bool {
	for _, a := range allowedTransitions[s] {
		if a == next {
			return true
		}
	}
	return false
}

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FullName  string    `json:"full_name"`
	Role      Role      `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserInfo struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
}

type Discipline struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type GradeType struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	MaxScore int    `json:"max_score"`
}

type Grade struct {
	ID           string     `json:"id"`
	Number       int        `json:"number"`
	TaskName     string     `json:"task_name"`
	Notes        string     `json:"notes"`
	Status       Status     `json:"status"`
	StudentID    string     `json:"student_id"`
	TeacherID    *string    `json:"teacher_id"`
	DisciplineID int        `json:"discipline_id"`
	GradeTypeID  int        `json:"grade_type_id"`
	Score        *int       `json:"score"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	GradedAt     *time.Time `json:"graded_at"`

	Student    *UserInfo   `json:"student,omitempty"`
	Teacher    *UserInfo   `json:"teacher,omitempty"`
	Discipline *Discipline `json:"discipline,omitempty"`
	GradeType  *GradeType  `json:"grade_type,omitempty"`
}

type Comment struct {
	ID         string    `json:"id"`
	GradeID    string    `json:"grade_id"`
	AuthorID   string    `json:"author_id"`
	Body       string    `json:"body"`
	IsInternal bool      `json:"is_internal"`
	CreatedAt  time.Time `json:"created_at"`
	Author     *UserInfo `json:"author,omitempty"`
}

type GradeHistory struct {
	ID        string    `json:"id"`
	GradeID   string    `json:"grade_id"`
	ChangedBy string    `json:"changed_by"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
	User      *UserInfo `json:"user,omitempty"`
}

type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type TeacherStat struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Count    int    `json:"count"`
}

type DashboardStats struct {
	TotalGrades    int            `json:"total_grades"`
	PendingGrades  int            `json:"pending_grades"`
	AvgScore       *float64       `json:"avg_score"`
	GradesByStatus map[string]int `json:"grades_by_status"`
	TopTeachers    []TeacherStat  `json:"top_teachers"`
}
