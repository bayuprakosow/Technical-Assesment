package router

import (
	"github.com/gin-gonic/gin"
	"github.com/zentara/technical_assesment/internal/config"
	"github.com/zentara/technical_assesment/internal/docs"
	"github.com/zentara/technical_assesment/internal/handlers"
	"github.com/zentara/technical_assesment/internal/middleware"
)

func New(cfg *config.Config, h *handlers.Handlers) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	docs.Register(r)

	r.GET("/health", h.Health)
	r.GET("/ready", h.Ready)
	r.GET("/api/v1/public/info", h.PublicInfo)

	r.POST("/api/v1/auth/register", h.Register)
	r.POST("/api/v1/auth/login", h.Login)

	authz := r.Group("/api/v1")
	authz.Use(middleware.JWT(cfg))
	{
		authz.GET("/me", h.Me)
		authz.GET("/findings", h.ListFindings)
		authz.POST("/findings", h.CreateFinding)
	}

	internal := r.Group("/internal")
	internal.Use(middleware.BasicAuth(cfg))
	{
		internal.GET("/metrics", h.InternalMetrics)
		internal.POST("/cache/purge", h.InternalPurge)
	}

	return r
}
