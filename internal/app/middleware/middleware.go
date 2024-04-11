package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/pkg/requestid"
	"log/slog"
	"strings"
	"time"
)

type Middleware struct {
	config *config.Config
	log    *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Middleware {
	return &Middleware{
		config: cfg,
		log:    log,
	}
}

func (m *Middleware) RequestLog(ctx *gin.Context) {
	if strings.HasPrefix(ctx.Request.URL.Path, "/api/docs/") { // ignore logging swagger documentation
		return
	}

	start := time.Now()

	ctx.Next()

	m.log.Info("request completed",
		slog.String("request_id", requestid.Get(ctx)),
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.Request.URL.Path),
		slog.String("client_ip", ctx.ClientIP()),
		slog.String("duration", fmt.Sprintf("%v", time.Now().Sub(start))),
		slog.Int("body_size", ctx.Writer.Size()),
		slog.Int("status", ctx.Writer.Status()),
	)
}
