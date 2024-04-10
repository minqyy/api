package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/service/repository/postgres/user"
	"github.com/minqyy/api/internal/service/repository/redis/session"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	User    IUser
	Session ISession
}

type IUser interface {
	Create(ctx context.Context, email, passwordHash string) (*user.User, error)
	IsUserExists(ctx context.Context, email string) (bool, error)
	GetByCredentials(ctx context.Context, email string, passwordHash string) (*user.User, error)
}

type ISession interface {
	Create(ctx context.Context, refreshToken string, userID string, ip string, userAgent string) error
	Close(ctx context.Context, refreshToken string) error
	Get(ctx context.Context, refreshToken string) (*session.Session, error)
}

func New(cfg *config.Config, postgresDB *sqlx.DB, redisDB *redis.Client) *Repository {
	return &Repository{
		User:    user.New(postgresDB),
		Session: session.New(redisDB, cfg),
	}
}
