package models

import "time"

type Pipelines struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Tasks       []Task         `json:"tasks"`
	Schedule    string         `json:"schedule"` // Cron expression for time-based triggers
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Metadata    map[string]any `json:"metadata"`
}

type Task struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Execute      func() error `json:"execute"`
	Dependencies []string     `json:"dependencies"`
	Status       string       `json:"status"`
}
