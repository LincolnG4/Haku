package models

import "time"

type Pipeline struct {
	ID          string         `json:"id"`
	UserID      string         `json:"user_id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Tasks       []Task         `json:"tasks"`
	Schedule    string         `json:"schedule,omitempty"` // Cron expression for time-based triggers
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
	Status      string         `json:"status,omitempty"`
}

type Task struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Activity     string         `json:"activity"` // Name of the activity (e.g., "copy-file", "send-email", "run-script")
	Config       map[string]any `json:"config"`   // Configuration for the activity
	Dependencies []string       `json:"dependencies,omitempty"`
	Status       string         `json:"status"`
}
