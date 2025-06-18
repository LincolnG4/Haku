package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	QueryTimeoutDuration = 5 * time.Second
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Storage struct {
	Pipelines interface {
		Create(context.Context, *Pipelines) error
		GetByID(context.Context, int64) (*Pipelines, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Pipelines) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetByID(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
	}
	Organization interface {
		Create(context.Context, *Organization) error
		GetByID(context.Context, int64) (Organization, error)
		AddMember(context.Context, *OrganizationMember) error
		GetMembers(context.Context, int64) ([]OrganizationMember, error)
		GetMember(context.Context, int64, int64) (OrganizationMember, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Pipelines:    &PipelinesStore{db},
		Users:        &UsersStore{db},
		Organization: &OrganizationStore{db},
	}
}
