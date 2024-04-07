package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/minqyy/api/internal/app/request"
	"github.com/minqyy/api/internal/app/response"
	"github.com/minqyy/api/internal/lib/log/sl"
	"github.com/minqyy/api/internal/service/repository/postgres/user"
	"github.com/minqyy/api/pkg/requestid"
	"log/slog"
	"net/http"
	"net/mail"
)

func (h *Handler) SignUp(ctx *gin.Context) {
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

func (h *Handler) Login(ctx *gin.Context) {
	log := h.log.With(
		slog.String("op", "handler.Login"),
		slog.String("request_id", requestid.Get(ctx)),
	)

	var body request.UserLogin

	if err := ctx.BindJSON(&body); err != nil {
		log.Debug("error occurred while decode request body", sl.Err(err))
		response.SendInvalidRequestBodyError(ctx)
		return
	}

	passwordHash := h.service.Hasher.Create([]byte(body.Password))

	usr, err := h.service.Repository.User.GetByCredentials(ctx, body.Email, passwordHash)
	if errors.Is(err, user.ErrUserNotFound) {
		log.Debug("user not found with provided credentials", slog.String("email", body.Email), slog.String("password_hash", passwordHash))
		response.SendError(ctx, http.StatusNotFound, "user not found")
		return
	}
	if err != nil {
		log.Error("could not get user by credentials", sl.Err(err))
		response.SendError(ctx, http.StatusInternalServerError, "could not get user")
		return
	}

	tokenPair, err := h.service.TokenManager.GenerateTokenPair(usr.ID)
	if err != nil {
		log.Error("error occurred while generating token pair",
			slog.String("user_id", usr.ID),
			sl.Err(err),
		)
		response.SendError(ctx, http.StatusInternalServerError, "can't create token pair")
		return
	}

	err = h.service.Repository.Session.Create(ctx, tokenPair.RefreshToken, usr.ID, ctx.ClientIP(), ctx.Request.UserAgent())
	if err != nil {
		log.Error("error occurred while creating refresh session in database", sl.Err(err))
		response.SendError(ctx, http.StatusInternalServerError, "can't create refresh session")
		return
	}

	ctx.JSON(http.StatusOK, response.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})
}
