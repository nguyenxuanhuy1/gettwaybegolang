package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"gateway/config"

	"github.com/gin-gonic/gin"
)

type GoogleAuthHandler struct{}

func NewGoogleAuthHandler() *GoogleAuthHandler {
	return &GoogleAuthHandler{}
}

// HandleLogin initiates Google OAuth flow
func (h *GoogleAuthHandler) HandleLogin(c *gin.Context) {
	// Generate state token for CSRF protection
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Store state in session/cookie (simplified - use secure session in production)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	// Redirect to Google OAuth
	url := config.GoogleOAuthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// HandleCallback handles Google OAuth callback
func (h *GoogleAuthHandler) HandleCallback(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// Verify state token
	state := c.Query("state")
	savedState, err := c.Cookie("oauth_state")
	if err != nil || state != savedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state token"})
		return
	}

	// Exchange code for token
	code := c.Query("code")
	token, err := config.GoogleOAuthConfig.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	// Get user info from Google
	client := config.GoogleOAuthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	if err := c.BindJSON(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user info"})
		return
	}

	// Create session token (simplified - use JWT or Redis session in production)
	// For now, store google_id as session token
	sessionToken := userInfo.ID // In production: generate JWT with user info

	// Set session cookie
	c.SetCookie("session_token", sessionToken, 3600*24*7, "/", "", false, true) // 7 days

	// Store user info in cookie for frontend (temporary)
	c.SetCookie("user_email", userInfo.Email, 3600*24*7, "/", "", false, false)
	c.SetCookie("user_name", userInfo.Name, 3600*24*7, "/", "", false, false)
	c.SetCookie("google_id", userInfo.ID, 3600*24*7, "/", "", false, false)

	// Redirect to frontend
	redirectURL := config.Config.FrontendAuthRedirectURL + 
		"?logged_in=true&email=" + userInfo.Email + 
		"&name=" + userInfo.Name

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
