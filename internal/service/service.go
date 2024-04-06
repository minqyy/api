package service

import "github.com/minqyy/api/internal/service/repository"

type Service struct {
	Repository *repository.Repository
}

func New(r *repository.Repository) *Service {
	return &Service{
		Repository: r,
	}
}
