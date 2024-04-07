package handler

import (
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/service"
	"log/slog"
)

type Handler struct {
	config  *config.Config
	log     *slog.Logger
	service *service.Service
}

func New(cfg *config.Config, log *slog.Logger, s *service.Service) *Handler {
	return &Handler{
		config:  cfg,
		log:     log,
		service: s,
	}
}
