package repository

import (
	"database/sql"
)

// StatsRepo handles analytics queries.
type StatsRepo struct {
	db *sql.DB
}

func NewStatsRepo(db *sql.DB) *StatsRepo {
	return &StatsRepo{db: db}
}

// Dashboard returns a map of aggregated metrics for the dashboard.
func (r *StatsRepo) Dashboard() (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// -- tickets_by_status --
	statusRows, err := r.db.Query(`SELECT status, COUNT(*) FROM tickets GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer statusRows.Close()

	byStatus := make(map[string]int)
	totalTickets := 0
	for statusRows.Next() {
		var status string
		var count int
		if err := statusRows.Scan(&status, &count); err != nil {
			return nil, err
		}
		byStatus[status] = count
		totalTickets += count
	}
	if err := statusRows.Err(); err != nil {
		return nil, err
	}
	result["tickets_by_status"] = byStatus
	result["total_tickets"] = totalTickets

	// -- avg_resolution_hours --
	var avgResolution sql.NullFloat64
	err = r.db.QueryRow(`
		SELECT AVG(EXTRACT(EPOCH FROM (closed_at - created_at)) / 3600)
		FROM tickets
		WHERE status = 'closed' AND closed_at IS NOT NULL`).Scan(&avgResolution)
	if err != nil {
		return nil, err
	}
	if avgResolution.Valid {
		result["avg_resolution_hours"] = avgResolution.Float64
	} else {
		result["avg_resolution_hours"] = nil
	}

	// -- open_tickets --
	var openCount int
	err = r.db.QueryRow(`
		SELECT COUNT(*) FROM tickets
		WHERE status NOT IN ('closed')`).Scan(&openCount)
	if err != nil {
		return nil, err
	}
	result["open_tickets"] = openCount

	// -- top_assignees --
	assigneeRows, err := r.db.Query(`
		SELECT u.id, u.full_name, COUNT(*) AS count
		FROM tickets t
		JOIN users u ON u.id = t.assignee_id
		WHERE t.status IN ('resolved', 'closed')
		GROUP BY u.id, u.full_name
		ORDER BY count DESC
		LIMIT 5`)
	if err != nil {
		return nil, err
	}
	defer assigneeRows.Close()

	type assigneeStat struct {
		ID       string `json:"id"`
		FullName string `json:"full_name"`
		Count    int    `json:"count"`
	}
	var topAssignees []assigneeStat
	for assigneeRows.Next() {
		var a assigneeStat
		if err := assigneeRows.Scan(&a.ID, &a.FullName, &a.Count); err != nil {
			return nil, err
		}
		topAssignees = append(topAssignees, a)
	}
	if err := assigneeRows.Err(); err != nil {
		return nil, err
	}
	result["top_assignees"] = topAssignees

	return result, nil
}
