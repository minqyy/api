package service

import (
	"github.com/minqyy/api/internal/service/hasher"
	"github.com/minqyy/api/internal/service/repository"
	"github.com/minqyy/api/internal/service/token"
)

type Service struct {
	Repository   *repository.Repository
	Hasher       *hasher.Hasher
	TokenManager *token.Manager
}

func New(r *repository.Repository, h *hasher.Hasher, t *token.Manager) *Service {
	return &Service{
		Repository:   r,
		Hasher:       h,
		TokenManager: t,
	}
}
