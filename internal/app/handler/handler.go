package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/minqyy/api/internal/app/request"
	"github.com/minqyy/api/internal/app/response"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/lib/log/sl"
	"github.com/minqyy/api/internal/service"
	"github.com/minqyy/api/internal/service/repository/postgres/user"
	"github.com/minqyy/api/pkg/requestid"
	"log/slog"
	"net/http"
	"net/mail"
)

type Handler struct {
	config  *config.Config
	log     *slog.Logger
	service *service.Service
}

// New returns a new instance of Handler.
func New(cfg *config.Config, log *slog.Logger, s *service.Service) *Handler {
	return &Handler{
		config:  cfg,
		log:     log,
		service: s,
	}
}

// Register      Creates a user in database
// @Summary      User registration
// @Description  Creates a user in database
// @Tags         auth
// @Router       /auth/signup       [post]
func (h *Handler) Register(ctx *gin.Context) {
	log := h.log.With(
		slog.String("op", "handler.Register"),
		slog.String("request_id", requestid.Get(ctx)),
	)

	var body request.UserCreate

	if err := ctx.BindJSON(&body); err != nil {
		log.Debug("Error occurred while decoding request body", sl.Err(err))
		response.SendInvalidRequestBodyError(ctx)
		return
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		log.Debug("Email is invalid", slog.String("email", body.Email))
		response.SendError(ctx, http.StatusBadRequest, "invalid email")
		return
	}

	passwordHash := h.service.Hasher.Create([]byte(body.Password))

	exists, err := h.service.Repository.User.IsUserExists(ctx, body.Email)
	if err != nil {
		log.Error("Error occurred while checking user existence", sl.Err(err))
		response.SendError(ctx, http.StatusInternalServerError, "can't create user")
		return
	}
	if exists {
		log.Debug("User already exists", slog.String("email", body.Email))
		response.SendError(ctx, http.StatusBadRequest, "user already exists")
		return
	}

	createdUser, err := h.service.Repository.User.Create(ctx, body.Email, passwordHash)
	if errors.Is(err, user.ErrUserAlreadyExists) {
		log.Debug("User already exists")
		response.SendError(ctx, http.StatusBadRequest, "user already exists")
		return
	}
	if err != nil {
		log.Error("Error occurred while creating user", sl.Err(err))
		response.SendError(ctx, http.StatusInternalServerError, "can't create user")
		return
	}

	log.Info("User created",
		slog.String("id", createdUser.ID),
		slog.String("email", createdUser.Email),
	)

	ctx.JSON(http.StatusCreated, response.User{
		ID:    createdUser.ID,
		Email: createdUser.Email,
	})
}
