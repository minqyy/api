package service

import (
	"github.com/minqyy/api/internal/service/hasher"
	"github.com/minqyy/api/internal/service/repository"
)

type Service struct {
	Repository *repository.Repository
	Hasher     *hasher.Hasher
}

func New(r *repository.Repository, h *hasher.Hasher) *Service {
	return &Service{
		Repository: r,
		Hasher:     h,
	}
}
