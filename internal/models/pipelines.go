package models

import "time"

type Pipeline struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Tasks       []Task         `json:"tasks,omitempty"`
	Schedule    string         `json:"schedule,omitempty"` // Cron expression for time-based triggers
	CreatedAt   time.Time      `json:"created_at,omitempty"`
	UpdatedAt   time.Time      `json:"updated_at,omitempty"`
	Metadata    map[string]any `json:"metadat,omitempty"`
	Status      string         `json:"status,omitempty"`
}

type Task struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Execute      func() error `json:"execute,omitempty"`
	Dependencies []string     `json:"dependencies,omitempty"`
	Status       string       `json:"status"`
}
