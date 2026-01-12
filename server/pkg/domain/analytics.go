package domain

import "time"

// UserVisit представляет запись о заходе пользователя
type UserVisit struct {
	ID        string    `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	SessionID string    `json:"session_id" db:"session_id"`
	Source    *string   `json:"source" db:"source"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// CreateUserVisit для создания новой записи о заходе
type CreateUserVisit struct {
	UserID    string  `json:"user_id"`
	SessionID string  `json:"session_id"`
	Source    *string `json:"source,omitempty"`
}

// UserVisitStats статистика заходов пользователя
type UserVisitStats struct {
	UserID           string    `json:"user_id" db:"user_id"`
	TotalVisits      int       `json:"total_visits" db:"total_visits"`
	UniqueDays       int       `json:"unique_days" db:"unique_days"`
	FirstVisit       time.Time `json:"first_visit" db:"first_visit"`
	LastVisit        time.Time `json:"last_visit" db:"last_visit"`
	DaysSpan         int       `json:"days_span" db:"days_span"`
	AvgVisitsPerDay  float64   `json:"avg_visits_per_day" db:"avg_visits_per_day"`
}
