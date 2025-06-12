package store

import (
	"context"
	"database/sql"
)

type Pipelines struct {
	// Basic Info
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	// Creation Info
	Version   int    `json:"version"`
	CreatedAt string `json:"create_at"`
	UpdatedAt string `json:"update_at"`
}

type PipelinesStore struct {
	db *sql.DB
}

func (s *PipelinesStore) Create(ctx context.Context, pipeline Pipelines) error {
	query := `
	INSERT INTO pipelines (user_id, name )
	VALUES  ($1,$2,$3) RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRowContext(
		ctx,
		query,
		pipeline.UserID,
		pipeline.Name,
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
