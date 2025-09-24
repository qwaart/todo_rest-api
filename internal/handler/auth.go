package handler

import (
	"net/http"
	"strconv"

	"rest_api/internal/auth"

	"github.com/gin-gonic/gin"
	"log/slog"
)

type AuthHandler struct {
	service *auth.Service
	log *slog.Logger
}

func NewAuthorization(service *auth.Service, log *slog.Logger) *AuthHandler {
	return &AuthHandler{service: service, log: log}
}

// GET /admin/keys
func (h *AuthHandler) ListKeys(c *gin.Context) {
	keys, err := h.service.ListKeys()
	if err != nil {
		h.log.Error("failed to list keys", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list keys"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"keys": keys})
}

// POST /register
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Owner string `json:"owner" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("invalid register request", slog.String("err", err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := h.service.RegisterAPIKey(req.Owner)
	if err != nil {
		h.log.Error("failed to register api key", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list keys"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"api_key": key})
}

// POST /admin/permission
func (h *AuthHandler) CreatePermission(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.CreatePermission(req.Name); err != nil {
		h.log.Error("failed to create permission", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create permission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "permission created"})
}

// POST /admin/key/:id/permission
func (h *AuthHandler) GrantPermission(c *gin.Context) {
	keyID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key id"})
		return
	}

	var req struct {
		PermissionID int64 `json:"permission" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.GrantPermission(keyID, req.PermissionID); err != nil {
		h.log.Error("failed to grand permission", slog.String("err", err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not grant permission"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "permission granted"})
}

// list permission
func (h *AuthHandler) ListPermissions(c *gin.Context) {
    perms, err := h.service.ListPermissions()
    if err != nil {
        h.log.Error("failed to list permissions", slog.String("err", err.Error()))
        c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list permissions"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"permissions": perms})
}