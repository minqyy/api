package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/service/repository/postgres/user"
	"github.com/minqyy/api/internal/service/repository/redis/session"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	User    *user.Postgres
	Session *session.Redis
}

func New(cfg *config.Config, postgresDB *sqlx.DB, redisDB *redis.Client) *Repository {
	return &Repository{
		User:    user.New(postgresDB),
		Session: session.New(redisDB, cfg),
	}
}
