package store

import (
	"context"
	"database/sql"
	"errors"
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
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

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

func (s *PipelinesStore) GetByID(ctx context.Context, pipelineID int64) (*Pipelines, error) {
	query := `
		SELECT id, user_id, name, status, version, created_at, updated_at 
		FROM pipelines WHERE id=$1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	pipeline := Pipelines{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		pipelineID,
	).Scan(
		&pipeline.ID,
		&pipeline.UserID,
		&pipeline.Name,
		&pipeline.Status,
		&pipeline.Version,
		&pipeline.CreatedAt,
		&pipeline.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &pipeline, nil
}
func (s *PipelinesStore) Delete(ctx context.Context, pipelineID int64) error {
	query := `
		DELETE FROM pipelines WHERE id=$1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := s.db.ExecContext(
		ctx,
		query,
		pipelineID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *PipelinesStore) Update(ctx context.Context, pipeline *Pipelines) error {
	query := `
		UPDATE pipelines
		SET name = $1, version = version + 1 
		WHERE id = $2 AND version = $3
		RETURNING version
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		pipeline.Name,
		pipeline.ID,
		pipeline.Version,
	).Scan(&pipeline.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrNotFound
		default:
			return err
		}
	}

	return nil
}
