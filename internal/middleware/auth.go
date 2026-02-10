package middleware

import (
	"context"
	"gateway/internal/domain"
	"gateway/internal/service"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService *service.AuthService
	planService *service.PlanService
	userRepo    interface {
		GetByIDSimple(ctx context.Context, userID int) (*domain.User, error)
	}
}

func NewAuthMiddleware(
	authService *service.AuthService,
	planService *service.PlanService,
	userRepo interface {
		GetByIDSimple(ctx context.Context, userID int) (*domain.User, error)
	},
) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		planService: planService, 8
		userRepo:    userRepo,
	}
}

// RequireAPIKey middleware verifies API key and validates user plan
func (m *AuthMiddleware) RequireAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// Extract Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Check Bearer format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		apiKey := parts[1]

		// Verify API key
		userID, err := m.authService.VerifyAPIKey(ctx, apiKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}

		// Get user info
		user, err := m.userRepo.GetByIDSimple(ctx, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
			c.Abort()
			return
		}

		// Check if user is locked
		if user.Locked {
			c.JSON(http.StatusForbidden, gin.H{"error": "account is locked"})
			c.Abort()
			return
		}

		// Validate plan
		plan, err := m.planService.ValidateUserPlan(ctx, userID)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Store user info and plan in context
		c.Set("user_id", userID)
		c.Set("user", user)
		c.Set("plan", plan)

		c.Next()
	}
}
