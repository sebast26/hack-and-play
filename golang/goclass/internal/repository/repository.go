package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"goclass/internal/domain"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) Repository {
	return Repository{
		db: db,
	}
}

func (r Repository) Store(posts []domain.Post) error {
	return nil
}
