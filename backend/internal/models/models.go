package models

import "time"

type Role string

const (
	RoleEmployee Role = "employee"
	RoleSupport  Role = "support"
	RoleAdmin    Role = "admin"
)

type Status string

const (
	StatusNew        Status = "new"
	StatusInProgress Status = "in_progress"
	StatusWaiting    Status = "waiting"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

var allowedTransitions = map[Status][]Status{
	StatusNew:        {StatusInProgress},
	StatusInProgress: {StatusWaiting, StatusResolved},
	StatusWaiting:    {StatusInProgress, StatusResolved},
	StatusResolved:   {StatusClosed, StatusInProgress},
	StatusClosed:     {},
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

type Category struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Priority struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
	SLAHours int    `json:"sla_hours"`
}

type Ticket struct {
	ID          string     `json:"id"`
	Number      int        `json:"number"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	AuthorID    string     `json:"author_id"`
	AssigneeID  *string    `json:"assignee_id"`
	CategoryID  int        `json:"category_id"`
	PriorityID  int        `json:"priority_id"`
	Rating      *int       `json:"rating"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ClosedAt    *time.Time `json:"closed_at"`

	Author   *UserInfo `json:"author,omitempty"`
	Assignee *UserInfo `json:"assignee,omitempty"`
	Category *Category `json:"category,omitempty"`
	Priority *Priority `json:"priority,omitempty"`
}

type Comment struct {
	ID         string    `json:"id"`
	TicketID   string    `json:"ticket_id"`
	AuthorID   string    `json:"author_id"`
	Body       string    `json:"body"`
	IsInternal bool      `json:"is_internal"`
	CreatedAt  time.Time `json:"created_at"`
	Author     *UserInfo `json:"author,omitempty"`
}

type StatusHistory struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	ChangedBy string    `json:"changed_by"`
	OldStatus string    `json:"old_status"`
	NewStatus string    `json:"new_status"`
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
