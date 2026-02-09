package handler

import (
	"context"
	"gateway/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UsageHandler struct {
	usageService *service.UsageService
}

func NewUsageHandler(usageService *service.UsageService) *UsageHandler {
	return &UsageHandler{
		usageService: usageService,
	}
}

// Handle retrieves usage statistics for the authenticated user
func (h *UsageHandler) Handle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Get usage stats
	stats, err := h.usageService.GetUsageStats(ctx, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get usage stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
