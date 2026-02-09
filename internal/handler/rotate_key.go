package handler

import (
	"context"
	"gateway/internal/repository/postgres"
	"gateway/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RotateKeyHandler struct {
	apiKeyRepo *postgres.APIKeyRepository
}

func NewRotateKeyHandler(apiKeyRepo *postgres.APIKeyRepository) *RotateKeyHandler {
	return &RotateKeyHandler{
		apiKeyRepo: apiKeyRepo,
	}
}

// Handle rotates the user's API key
func (h *RotateKeyHandler) Handle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Generate new API key
	newKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate API key"})
		return
	}

	// Hash the new key
	newKeyHash := utils.HashAPIKey(newKey)

	// Rotate key in database
	err = h.apiKeyRepo.RotateKey(ctx, userID.(int), newKeyHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to rotate key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key rotated successfully",
		"api_key": newKey,
		"warning": "Save this key securely. You won't be able to see it again.",
	})
}
