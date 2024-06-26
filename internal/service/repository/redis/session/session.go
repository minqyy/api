package session

import (
	"context"
	"encoding/json"
	"github.com/minqyy/api/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

const RefreshTokenCookie = "_refreshToken"

type Redis struct {
	client *redis.Client
	config *config.Config
}

type Session struct {
	UserID    string
	IP        string
	UserAgent string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func New(client *redis.Client, cfg *config.Config) *Redis {
	return &Redis{
		client: client,
		config: cfg,
	}
}

func (r *Redis) Create(ctx context.Context, refreshToken string, userID string, ip string, userAgent string) error {
	exists, err := r.client.Exists(ctx, refreshToken).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		return ErrRefreshTokenAlreadyExists
	}

	session := Session{
		UserID:    userID,
		IP:        ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(r.config.Token.Refresh.TTL),
	}

	marshalledSession, err := json.Marshal(session)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, refreshToken, marshalledSession, r.config.Token.Refresh.TTL).Err()
}

func (r *Redis) Close(ctx context.Context, refreshToken string) error {
	exists, err := r.client.Exists(ctx, refreshToken).Result()
	if err != nil {
		return err
	}
	if exists == 0 {
		return ErrSessionNotExists
	}
	return r.client.Del(ctx, refreshToken).Err()
}

func (r *Redis) Get(ctx context.Context, refreshToken string) (*Session, error) {
	exists, err := r.client.Exists(ctx, refreshToken).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, ErrSessionNotExists
	}

	marshalledData, err := r.client.Get(ctx, refreshToken).Result()
	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal([]byte(marshalledData), &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}
