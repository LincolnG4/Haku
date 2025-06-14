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
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Pipelines: &PipelinesStore{db},
		Users:     &UsersStore{db},
	}
}
