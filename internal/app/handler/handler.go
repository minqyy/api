package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/minqyy/api/internal/config"
	"log/slog"
	"time"
)

type Handler struct {
	config *config.Config
	log    *slog.Logger
}

// New returns a new instance of Handler.
func New(cfg *config.Config, log *slog.Logger) *Handler {
	return &Handler{
		config: cfg,
		log:    log,
	}
}

// Register      Creates a user in database.
// @Summary      User registration
// @Description  Creates a user in database
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body         request.UserCreate true "User data"
// @Success      201  {object}      response.User
// @Failure      400  {object}      response.Error
// @Failure      409  {object}      response.Error
// @Failure      500  {object}      response.Error
// @Router       /auth/signup       [post]
func (h *Handler) Register(ctx *gin.Context) {
	time.Sleep(10 * time.Second)
	return
}
