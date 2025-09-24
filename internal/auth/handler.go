package auth

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	storage *Storage
	log     *slog.Logger
}

func NewHandler(storage *Storage, log *slog.Logger) *Handler {
	return &Handler{storage: storage, log: log}
}

// Register — public endpoint for creating key
func (h *Handler) Register(c *gin.Context) {
	var req struct {
		Owner string `json:"owner" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	key, err := h.storage.RegisterKey(req.Owner)
	if err != nil {
		h.log.Error("failed to register API key", slog.String("owner", req.Owner), slog.Any("err", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register key"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"owner":   req.Owner,
		"api_key": key,
	})
}

// ListKeys — only for admin
func (h *Handler) ListKeys(c *gin.Context) {
	keys, err := h.storage.ListAPIKeys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list keys"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// CreatePermission — creation of a new permission
func (h *Handler) CreatePermission(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.storage.CreatePermission(req.Name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create permission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "permission created"})
}

// GrantPermission — binding a right to a key
func (h *Handler) GrantPermission(c *gin.Context) {
	keyID := c.Param("id")
	var req struct {
		Permission string `json:"permission" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.storage.GrantPermission(keyID, req.Permission); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to grant permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "permission granted"})
}