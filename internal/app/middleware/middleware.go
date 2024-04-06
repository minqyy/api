package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/pkg/requestid"
	"log/slog"
	"strings"
)

type Middleware struct {
	config *config.Config
	log    *slog.Logger
}

// New returns a pointer to a new instance of Middleware
func New(cfg *config.Config, log *slog.Logger) *Middleware {
	return &Middleware{
		config: cfg,
		log:    log,
	}
}

// RequestLog logs every request with parameters: method, path, client_ip, remote_addr, user_agent, status and duration
func (m *Middleware) RequestLog(ctx *gin.Context) {
	if strings.HasPrefix(ctx.Request.URL.Path, "/api/docs/") { // ignore logging swagger documentation
		return
	}

	m.log.Info("Request handled",
		slog.String("request_id", requestid.Get(ctx)),
		slog.String("method", ctx.Request.Method),
		slog.String("path", ctx.Request.URL.Path),
		slog.String("client_ip", ctx.ClientIP()),
	)
}
