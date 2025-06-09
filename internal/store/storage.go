package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Pipelines interface {
		Create(context.Context, Pipelines) error
	}
	Users interface {
		Create(context.Context, User) error
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Pipelines: &PipelinesStore{db},
		Users:     &UsersStore{db},
	}
}
