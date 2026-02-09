package handler

import (
	"context"
	"database/sql"
	"gateway/internal/repository/postgres"
	"gateway/pkg/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateKeyHandler struct {
	userRepo   *postgres.UserRepository
	apiKeyRepo *postgres.APIKeyRepository
}

func NewCreateKeyHandler(userRepo *postgres.UserRepository, apiKeyRepo *postgres.APIKeyRepository) *CreateKeyHandler {
	return &CreateKeyHandler{
		userRepo:   userRepo,
		apiKeyRepo: apiKeyRepo,
	}
}

// Handle creates API key for logged-in user
func (h *CreateKeyHandler) Handle(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Get user info from cookies (set by Google OAuth callback)
	googleID, err := c.Cookie("google_id")
	if err != nil || googleID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not logged in"})
		return
	}

	email, _ := c.Cookie("user_email")
	username, _ := c.Cookie("user_name")

	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user email"})
		return
	}

	// Get or create user
	userID, err := h.getOrCreateUser(ctx, googleID, email, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get/create user"})
		return
	}

	// Generate new API key
	newKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate API key"})
		return
	}

	// Hash and store
	keyHash := utils.HashAPIKey(newKey)
	err = h.apiKeyRepo.RotateKey(ctx, userID, keyHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create API key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "API key created successfully",
		"api_key": newKey,
		"user_id": userID,
		"email":   email,
		"warning": "Save this key securely. You won't be able to see it again.",
	})
}

func (h *CreateKeyHandler) getOrCreateUser(ctx context.Context, googleID, email, username string) (int, error) {
	// Try to find existing user by google_id
	query := `SELECT id FROM users WHERE google_id = $1`
	var userID int
	err := h.userRepo.GetDB().QueryRowContext(ctx, query, googleID).Scan(&userID)

	if err == nil {
		// User exists
		return userID, nil
	}

	if err != sql.ErrNoRows {
		// Real error
		return 0, err
	}

	// User doesn't exist, create new one
	if username == "" {
		username = email
	}

	insertQuery := `
		INSERT INTO users (username, email, google_id, role, coin, locked, created_at)
		VALUES ($1, $2, $3, 'user', 0, false, NOW())
		RETURNING id
	`

	err = h.userRepo.GetDB().QueryRowContext(ctx, insertQuery, username, email, googleID).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
