package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionMiddleware struct{}

func NewSessionMiddleware() *SessionMiddleware {
	return &SessionMiddleware{}
}

// RequireLogin checks if user is logged in via session/JWT
func (m *SessionMiddleware) RequireLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		_ = ctx

		// Check for session token or JWT
		// Option 1: Check cookie
		sessionToken, err := c.Cookie("session_token")
		if err == nil && sessionToken != "" {
			// TODO: Validate session token with Redis/database
			// For now, simple check
			c.Set("session_token", sessionToken)
			c.Next()
			return
		}

		// Option 2: Check Authorization header with JWT
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				jwtToken := parts[1]
				// TODO: Validate JWT and extract user info
				// For now, simple check
				c.Set("jwt_token", jwtToken)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		c.Abort()
	}
}
