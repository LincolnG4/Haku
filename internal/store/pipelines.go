package store

import (
	"context"
	"database/sql"
)

type PipelineState int

const (
	StateCreated  string = "created"
	StateQueued   string = "queued"
	StateRunning  string = "running"
	StateError    string = "error"
	StateRetrying string = "retrying"
)

type Pipelines struct {
	// Basic Info
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Status string `json:"status"`
	// Creation Info
	Version   int    `json:"version"`
	CreatedAt string `json:"create_at"`
	UpdatedAt string `json:"update_at"`
}

type PipelinesStore struct {
	db *sql.DB
}

func (s *PipelinesStore) Create(ctx context.Context, pipeline *Pipelines) error {
	query := `
	INSERT INTO pipelines (user_id, name, status, version)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	pipeline.Status = StateCreated
	pipeline.Version = 1

	err := s.db.QueryRowContext(
		ctx,
		query,
		pipeline.UserID,
		pipeline.Name,
		pipeline.Status,
		pipeline.Version,
	).Scan(
		&pipeline.ID,
		&pipeline.CreatedAt,
		&pipeline.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PipelinesStore) GetByID(ctx context.Context, pipelineID int64) (Pipelines, error) {
	query := `
		SELECT * FROM pipelines WHERE id=$1
	`
	pipeline := Pipelines{}

	err := s.db.QueryRowContext(
		ctx,
		query,
		pipelineID,
	).Scan(
		&pipeline.ID,
		&pipeline.Name,
		&pipeline.Status,

		&pipeline.CreatedAt,
		&pipeline.UpdatedAt,
	)

	if err != nil {
		return pipeline, err
	}

	return pipeline, nil
}
