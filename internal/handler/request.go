package handler

import (
	"context"
	"fmt"
	"gateway/internal/domain"
	"gateway/internal/repository/postgres"
	"gateway/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RequestHandler struct {
	rateLimitService *service.RateLimitService
	billingService   *service.BillingService
	apiLogRepo       *postgres.APILogRepository
}

func NewRequestHandler(
	rateLimitService *service.RateLimitService,
	billingService *service.BillingService,
	apiLogRepo *postgres.APILogRepository,
) *RequestHandler {
	return &RequestHandler{
		rateLimitService: rateLimitService,
		billingService:   billingService,
		apiLogRepo:       apiLogRepo,
	}
}

// Handle processes an API gateway request
func (h *RequestHandler) Handle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()

	// Get user ID and plan from context
	userID, _ := c.Get("user_id")
	planInterface, _ := c.Get("plan")
	plan := planInterface.(*domain.UserPlanInfo)

	// Generate unique request ID
	requestID := uuid.New().String()

	// 1. Check rate limit
	rateLimit := 0
	if plan.RateLimit != nil {
		rateLimit = *plan.RateLimit
	}

	if rateLimit > 0 {
		err := h.rateLimitService.CheckRateLimit(ctx, userID.(int), rateLimit)
		if err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":      "rate limit exceeded",
				"rate_limit": rateLimit,
			})
			return
		}
	}

	// 2. Deduct coins
	cost, err := h.billingService.DeductCoins(ctx, userID.(int), plan.ProductCode, requestID)
	if err != nil {
		if err.Error() == "insufficient coins" {
			c.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient coins"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "billing error"})
		return
	}

	// 3. Log API request
	endpoint := c.Request.URL.Path
	err = h.apiLogRepo.Create(ctx, userID.(int), plan.ProductCode, endpoint, cost, requestID)
	if err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to log API request: %v\n", err)
	}

	// 4. Forward to upstream (mocked)
	upstreamResponse := map[string]interface{}{
		"status":  "success",
		"message": "Request processed successfully",
		"data": map[string]interface{}{
			"request_id": requestID,
			"timestamp":  time.Now().Format(time.RFC3339),
			"mock":       true,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"request_id": requestID,
		"cost":       cost,
		"upstream":   upstreamResponse,
	})
}
