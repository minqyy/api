package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/service/repository/postgres/user"
)

type Repository struct {
	User *user.Postgres
}

func New(postgres *sqlx.DB, cfg *config.Config) *Repository {
	return &Repository{
		User: user.New(postgres),
	}
}
