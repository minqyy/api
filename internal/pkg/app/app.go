package app

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/minqyy/api/internal/app/router"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/lib/log/prettyslog"
	"github.com/minqyy/api/internal/lib/log/sl"
	"github.com/minqyy/api/internal/service"
	"github.com/minqyy/api/internal/service/repository"
	"github.com/minqyy/api/internal/service/repository/postgres"
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
	return &App{
		config: cfg,
		log:    initLogger(cfg.Env),
	}
}

// Run runs entire application and services
func (a *App) Run() {
	gin.SetMode(gin.ReleaseMode)

	a.log.Info("Configuring server...", slog.String("env", a.config.Env))

	postgresDB, err := postgres.New(a.config.Postgres)
	if err != nil {
		a.log.Error("Could not connect to postgres database")
		os.Exit(1)
	}

	repo := repository.New(postgresDB, a.config)
	srv := service.New(repo)
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
				a.log.Error("Could not start the server", sl.Err(err))
			}
		}
	}()

	a.log.Info("Server started", slog.String("address", a.config.Server.Address))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	a.log.Info("Server is shutting down...")

	err = server.Shutdown(context.Background())
	if err != nil {
		a.log.Error("Error occurred on server shutting down", sl.Err(err))
	}

	a.log.Info("Server stopped")

	err = postgresDB.Close()
	if err != nil {
		a.log.Error("Could not close postgres connection")
	}

	a.log.Info("Postgres connection closed")

	// TODO: Close all db connections
}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case config.EnvLocal:
		log = initPrettyLogger()
	case config.EnvDevelopment:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProduction:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func initPrettyLogger() *slog.Logger {
	opts := prettyslog.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	prettyHandler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(prettyHandler)
}
