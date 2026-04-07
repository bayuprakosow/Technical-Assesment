package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zentara/technical_assesment/internal/config"
)

// BasicAuth memerlukan BASIC_AUTH_USER dan BASIC_AUTH_PASSWORD.
// Jika salah satu kosong, rute internal mengembalikan 503 (lihat PROJECT_PLAN.md).
func BasicAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg.BasicAuthUser == "" || cfg.BasicAuthPassword == "" {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "internal routes require BASIC_AUTH_USER and BASIC_AUTH_PASSWORD",
			})
			return
		}
		user, pass, ok := c.Request.BasicAuth()
		if !ok || user != cfg.BasicAuthUser || pass != cfg.BasicAuthPassword {
			c.Header("WWW-Authenticate", `Basic realm="internal"`)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Next()
	}
}
