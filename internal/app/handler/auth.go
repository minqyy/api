package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/minqyy/api/internal/app/request"
	"github.com/minqyy/api/internal/app/response"
	"github.com/minqyy/api/internal/lib/log/sl"
	"github.com/minqyy/api/internal/service/repository/postgres/user"
	"github.com/minqyy/api/internal/service/repository/redis/session"
	"github.com/minqyy/api/pkg/requestid"
	"log/slog"
	"net/http"
	"net/mail"
)

func (h *Handler) SignUp(ctx *gin.Context) {
	log := h.log.With(
		slog.String("op", "handler.SignUp"),
		slog.String("request_id", requestid.Get(ctx)),
	)

	var body request.UserCreate

	if err := ctx.BindJSON(&body); err != nil {
		log.Debug("error occurred while decoding request body", sl.Err(err))
		response.SendInvalidRequestBodyError(ctx)
		return
	}

	if _, err := mail.ParseAddress(body.Email); err != nil {
		log.Debug("email is invalid", slog.String("email", body.Email))
		response.SendError(ctx, http.StatusBadRequest, "invalid email")
		return
	}

	passwordHash := h.service.Hasher.Create([]byte(body.Password))

	exists, err := h.service.Repository.User.IsUserExists(ctx, body.Email)
	if err != nil {
		log.Error("error occurred while checking user existence", sl.Err(err))
		response.SendError(ctx, http.StatusInternalServerError, "can't create user")
		return
	}
	if exists {
		log.Debug("user already exists", slog.String("email", body.Email))
		response.SendError(ctx, http.StatusBadRequest, "user already exists")
		return
	}

	createdUser, err := h.service.Repository.User.Create(ctx, body.Email, passwordHash)
	if errors.Is(err, user.ErrUserAlreadyExists) {
		log.Debug("user already exists")
		response.SendError(ctx, http.StatusBadRequest, "user already exists")
		return
	}
	if err != nil {
		log.Error("error occurred while creating user", sl.Err(err))
		response.SendError(ctx, http.StatusInternalServerError, "can't create user")
		return
	}

	log.Info("user created",
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

func (h *Handler) Logout(ctx *gin.Context) {
	log := h.log.With(
		slog.String("op", "handler.Logout"),
		slog.String("request_id", requestid.Get(ctx)),
	)

	refreshToken, err := ctx.Cookie(session.RefreshTokenCookie)
	if err != nil {
		log.Debug("no refresh token cookie found", sl.Err(err))
		ctx.JSON(http.StatusBadRequest, response.Error{Message: "no refresh token cookie found"})
		return
	}

	err = h.service.Repository.Session.Close(ctx, refreshToken)
	if errors.Is(err, session.ErrSessionNotExists) {
		log.Debug("refresh session not found")
		response.SendError(ctx, http.StatusNotFound, "refresh session not found")
		return
	}
	if err != nil {
		response.SendError(ctx, http.StatusInternalServerError, "can't delete refresh session")
		log.Error("error occurred while deleting refresh session", sl.Err(err))
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *Handler) RefreshTokens(ctx *gin.Context) {
	log := h.log.With(
		slog.String("op", "handler.RefreshTokens"),
		slog.String("request_id", requestid.Get(ctx)),
	)

	refreshToken, err := ctx.Cookie(session.RefreshTokenCookie)
	if err != nil {
		log.Debug("no refresh token provided", sl.Err(err))
		response.SendError(ctx, http.StatusForbidden, "no refresh token cookie found")
		return
	}

	userSession, err := h.service.Repository.Session.Get(ctx, refreshToken)
	if err != nil {
		log.Debug("invalid refresh token")
		response.SendError(ctx, http.StatusForbidden, "invalid refresh token")
		return
	}

	err = h.service.Repository.Session.Close(ctx, refreshToken)
	if err != nil {
		log.Error("error occurred while deleting refresh session", sl.Err(err))
		response.SendError(ctx, http.StatusInternalServerError, "can't delete refresh session")
	}

	tokenPair, err := h.service.TokenManager.GenerateTokenPair(userSession.UserID)
	if err != nil {
		log.Error("error occurred while generating token pair",
			slog.String("user_id", userSession.UserID),
			sl.Err(err),
		)
		response.SendError(ctx, http.StatusInternalServerError, "can't create token pair")
		return
	}

	err = h.service.Repository.Session.Create(ctx, tokenPair.RefreshToken, userSession.UserID, ctx.ClientIP(), ctx.Request.UserAgent())
	if err != nil {
		log.Error("error occurred while creating refresh session",
			slog.String("user_id", userSession.UserID),
			sl.Err(err),
		)
		response.SendError(ctx, http.StatusInternalServerError, "can't create refresh session")
		return
	}

	log.Info("refresh session successfully created",
		slog.String("user_id", userSession.UserID),
	)

	ctx.JSON(http.StatusOK, response.TokenPair{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	})
}
