package AuthMiddleware

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodGet {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader != "Bearer."+token {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Next()
	}
}