package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"

	"gateway/config"
	"gateway/internal/repository/postgres"

	"github.com/gin-gonic/gin"
)

type GoogleAuthHandler struct {
	userRepo *postgres.UserRepository
}

func NewGoogleAuthHandler(userRepo *postgres.UserRepository) *GoogleAuthHandler {
	return &GoogleAuthHandler{
		userRepo: userRepo,
	}
}


func (h *GoogleAuthHandler) HandleLogin(c *gin.Context) {
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	url := config.GoogleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}


func (h *GoogleAuthHandler) HandleCallback(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}

	code := c.Query("code")
	token, err := config.GoogleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "oauth exchange failed"})
		return
	}

	client := config.GoogleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user info"})
		return
	}

	userID, err := h.getOrCreateUser(
		ctx,
		userInfo.ID,
		userInfo.Email,
		userInfo.Name,
		userInfo.Picture,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user creation failed"})
		return
	}

	sessionToken := generateSecureToken()
	c.SetCookie(
		"session_token",
		sessionToken,
		3600*24*7,
		"/",
		"",
		false, 
		true,
	)

	c.Redirect(
		http.StatusTemporaryRedirect,
		config.Config.FrontendAuthRedirectURL,
	)
}


func (h *GoogleAuthHandler) getOrCreateUser(
	ctx context.Context,
	googleID, email, username, avatar string,
) (int, error) {

	var userID int
	err := h.userRepo.GetDB().
		QueryRowContext(ctx, `SELECT id FROM users WHERE google_id = $1`, googleID).
		Scan(&userID)

	if err == nil {
		// Update avatar
		_, _ = h.userRepo.GetDB().
			ExecContext(ctx, `UPDATE users SET avatar = $1 WHERE id = $2`, avatar, userID)
		return userID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	if username == "" {
		username = email
	}

	err = h.userRepo.GetDB().
		QueryRowContext(ctx, `
			INSERT INTO users (username, email, google_id, avatar, role, coin, locked, created_at)
			VALUES ($1, $2, $3, $4, 'user', 0, false, NOW())
			RETURNING id
		`,
			username, email, googleID, avatar,
		).
		Scan(&userID)

	return userID, err
}


func generateSecureToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
