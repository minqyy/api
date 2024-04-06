package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minqyy/api/internal/app/handler"
	"github.com/minqyy/api/internal/app/middleware"
	"github.com/minqyy/api/internal/config"
	"github.com/minqyy/api/internal/service"
	"github.com/minqyy/api/pkg/requestid"
	"github.com/minqyy/api/pkg/str"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log/slog"
)

type Router struct {
	config  *config.Config
	log     *slog.Logger
	handler *handler.Handler
	mw      *middleware.Middleware
	service *service.Service
}

// New creates a new instance of router.Router
func New(cfg *config.Config, log *slog.Logger, s *service.Service) *Router {
	return &Router{
		config:  cfg,
		log:     log,
		handler: handler.New(cfg, log),
		mw:      middleware.New(cfg, log),
		service: s,
	}
}

// InitRoutes initializes all routes and logs them
func (r *Router) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(requestid.Set)
	router.Use(r.mw.RequestLog)

	api := router.Group("/api")
	{
		api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		auth := api.Group("/auth")
		{
			auth.POST("/sign-up", r.handler.Register)
		}
	}

	r.logRoutes(router.Routes())

	return router
}

// logRoutes logs all routes of router.Router
func (r *Router) logRoutes(routes gin.RoutesInfo) {
	for _, route := range routes {
		var method, path string
		if r.config.Env == config.EnvLocal {
			method = str.CompleteStringToLength(route.Method, 9, ' ')
			path = str.CompleteStringToLength(route.Path, 25, ' ')
		} else {
			method = route.Method
			path = route.Path
		}

		routeLog := fmt.Sprintf("%s %s --> %s", method, path, route.Handler)

		r.log.Debug(routeLog)
	}
}
