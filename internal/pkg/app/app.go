package app

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/minqyy/api/internal/app/router"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/lib/log"
	"github.com/minqyy/api/internal/lib/log/prettyslog"
	"github.com/minqyy/api/internal/lib/log/sl"
	"github.com/minqyy/api/internal/service"
	"github.com/minqyy/api/internal/service/hasher"
	"github.com/minqyy/api/internal/service/repository"
	"github.com/minqyy/api/internal/service/repository/postgres"
	"github.com/minqyy/api/internal/service/repository/redis"
	"github.com/minqyy/api/internal/service/token"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	config *config.Config
	log    *slog.Logger
}

func New(cfg *config.Config) *App {
	_ = os.MkdirAll(".logs", os.ModePerm)

	return &App{
		config: cfg,
		log:    initLogger(cfg.Env),
	}
}

func (a *App) Run() {
	gin.SetMode(gin.ReleaseMode)

	a.log.Info("configuring server...", slog.String("env", a.config.Env))

	postgresDB, err := postgres.New(a.config.Postgres)
	if err != nil {
		a.log.Error("could not connect to postgres database", sl.Err(err))
		os.Exit(1)
	}

	redisDB, err := redis.New(a.config.Redis)
	if err != nil {
		a.log.Error("could not connect to redis database", sl.Err(err))
		os.Exit(1)
	}

	tokenManager := token.New(a.config.Token)
	repo := repository.New(a.config, postgresDB, redisDB)
	hash := hasher.New(a.config.Hasher.Salt)
	srv := service.New(repo, hash, tokenManager)
	r := router.New(a.config, a.log, srv)

	server := &http.Server{
		Addr:         a.config.Server.Address,
		Handler:      r.InitRoutes(),
		ReadTimeout:  a.config.Server.Timeout,
		WriteTimeout: a.config.Server.Timeout,
		IdleTimeout:  a.config.Server.IdleTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				a.log.Error("could not start the server", sl.Err(err))
			}
		}
	}()

	a.log.Info("server started", slog.String("address", a.config.Server.Address))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	a.log.Info("server is shutting down...")

	err = server.Shutdown(context.Background())
	if err != nil {
		a.log.Error("error occurred on server shutting down", sl.Err(err))
	}

	a.log.Info("HTTP server stopped")

	err = postgresDB.Close()
	if err != nil {
		a.log.Error("could not close postgres connection")
	}

	a.log.Info("postgres connection closed")

	err = redisDB.Close()
	if err != nil {
		a.log.Error("could not close redis connection", sl.Err(err))
	}

	a.log.Info("redis connection closed")

	a.log.Info("all services stopped")
}

func initLogger(env string) *slog.Logger {
	var l *slog.Logger

	switch env {
	case config.EnvLocal:
		l = initPrettyLogger(slog.LevelDebug)
	case config.EnvDevelopment:
		l = initCustomLogger(slog.LevelDebug)
	case config.EnvProduction:
		l = initCustomLogger(slog.LevelInfo)
	}

	return l
}

func initCustomLogger(level slog.Level) *slog.Logger {
	opts := log.HandlerOptions{
		SlogOptions: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewHandler(os.Stdout)

	return slog.New(handler)
}

func initPrettyLogger(level slog.Level) *slog.Logger {
	opts := prettyslog.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: level,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
