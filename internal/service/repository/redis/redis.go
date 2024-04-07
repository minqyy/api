package redis

import (
	"context"
	"github.com/minqyy/api/internal/config"
	"github.com/redis/go-redis/v9"
)

func New(cfg config.Redis) (*redis.Client, error) {
	db := redis.NewClient(&redis.Options{
		Addr:     cfg.Host,
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       0,
	})

	status := db.Ping(context.Background())
	if status.Err() != nil {
		return nil, status.Err()
	}

	return db, nil
}
