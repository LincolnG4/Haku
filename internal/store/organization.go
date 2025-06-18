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

type OrganizationMember struct {
	ID             int64  `json:"id"`
	OrganizationID int64  `json:"organization_id"`
	UserID         int64  `json:"user_id"`
	RoleID         int64  `json:"role_id"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type Role struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

const (
	AdminRole int64 = iota + 1
	DeveloperRole
	ViewerRole
)

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

func (s *OrganizationStore) GetByID(ctx context.Context, orgID int64) (Organization, error) {
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
			return Organization{}, ErrNotFound
		default:
			return Organization{}, err
		}
	}

	return org, nil
}

func (s *OrganizationStore) GetMembers(ctx context.Context, orgID int64) ([]OrganizationMember, error) {
	query := `
		SELECT id, user_id, organization_id, role_id, created_at, updated_at 
		FROM organization_members WHERE organization_id=$1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	orgs := []OrganizationMember{}
	rows, err := s.db.QueryContext(
		ctx,
		query,
		orgID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	defer rows.Close()

	for rows.Next() {
		var o OrganizationMember
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.OrganizationID,
			&o.RoleID,
			&o.CreatedAt,
			&o.UpdatedAt); err != nil {
			return nil, err
		}
		orgs = append(orgs, o)
	}

	return orgs, nil
}

func (s *OrganizationStore) AddMember(ctx context.Context, member *OrganizationMember) error {
	query := `
		INSERT INTO organization_members (user_id, organization_id, role_id)
		VALUES ($1, $2, $3) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		member.UserID,
		member.OrganizationID,
		member.RoleID,
	).Scan(
		&member.ID,
		&member.CreatedAt,
		&member.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *OrganizationStore) GetMember(ctx context.Context, orgID, userID int64) (OrganizationMember, error) {
	query := `
		SELECT id, user_id, organization_id, role_id, created_at, updated_at 
		FROM organization_members WHERE organization_id=$1 AND user_id=$2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	org := OrganizationMember{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		orgID,
		userID,
	).Scan(
		&org.ID,
		&org.UserID,
		&org.OrganizationID,
		&org.RoleID,
		&org.CreatedAt,
		&org.UpdatedAt)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return OrganizationMember{}, ErrNotFound
		default:
			return OrganizationMember{}, err
		}
	}

	return org, nil
}
