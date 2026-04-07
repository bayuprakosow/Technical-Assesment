package middleware

import (
	"github.com/gin-gonic/gin"
)

func UserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get(CtxUserIDKey)
	if !ok {
		return 0, false
	}
	uid, ok := v.(int64)
	if !ok {
		return 0, false
	}
	return uid, true
}
