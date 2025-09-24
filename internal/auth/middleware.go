package auth

import (
	"net/http"
	"strings"
	"log/slog"

	"github.com/gin-gonic/gin"
)


// AuthMiddleware - checks the API key
func AuthMiddleware(storage *Storage, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if !strings.HasPrefix(token, "ApiKey.") {
			log.Warn("missing or invalid Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid Authorization header"})
			c.Abort()
			return
		}

		apiKey := strings.TrimPrefix(token, "ApiKey.")
		ok, err := storage.ValidateKey(apiKey)
		if err != nil {
			log.Error("failed to validate api key", slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			c.Abort()
			return
		}

		if !ok {
			log.Warn("invalid api key")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid api key"})
			c.Abort()
			return
		}

		c.Set("api_key", apiKey) // save for handlers
		log.Debug("api key validate successfully", slog.String("key", apiKey))
		c.Next()
	}
}

// RequirePermission - middleware for checking a specific permissions
func RequirePermission(storage *Storage, permission string, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, exists := c.Get("api_key")
		if !exists {
			log.Warn("API key not found in context")
			c.JSON(http.StatusForbidden, gin.H{"error": "API key not found in context"})
			c.Abort()
			return
		}

		ok, err := storage.HasPermission(apiKey.(string), permission)
		if err != nil {
			log.Error("failed to check permission", slog.String("permission", permission), slog.Any("err", err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			c.Abort()
			return
		}

		if !ok {
			log.Warn("permission denied", slog.String("permission", permission), slog.String("api_key", apiKey.(string)))
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		log.Debug("permission granted", slog.String("permission", permission))
		c.Next()
	}
}