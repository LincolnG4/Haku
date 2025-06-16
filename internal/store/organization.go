package store

import (
	"context"
	"database/sql"
	"errors"
)

type OrganizationStore struct {
	db *sql.DB
}

type Organization struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

func (s *OrganizationStore) Create(ctx context.Context, org *Organization) error {
	query := `
	INSERT INTO organizations (name, description)
	VALUES ($1, $2) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		org.Name,
		org.Description,
	).Scan(
		&org.ID,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *OrganizationStore) GetByID(ctx context.Context, orgID int64) (*Organization, error) {
	query := `
		SELECT id, name, description, created_at, updated_at 
		FROM organizations WHERE id=$1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	org := Organization{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		orgID,
	).Scan(
		&org.ID,
		&org.Name,
		&org.Description,
		&org.CreatedAt,
		&org.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return &org, nil
}
