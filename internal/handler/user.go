package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"gateway/config"
	"gateway/internal/auth"
	"gateway/internal/repository"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type GoogleUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

/*  GOOGLE LOGIN  */
func GoogleLogin(c *gin.Context) {
	url := config.GoogleOAuthConfig.AuthCodeURL(
		"state",
		oauth2.SetAuthURLParam("prompt", "select_account"),
	)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {

	redirectBase := config.Config.FrontendAuthRedirectURL
	code := c.Query("code")

	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=missing_code")
		return
	}

	token, err := config.GoogleOAuthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=exchange_failed")
		return
	}

	resp, err := http.Get(
		"https://openidconnect.googleapis.com/v1/userinfo?access_token=" + token.AccessToken,
	)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=get_userinfo_failed")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=read_userinfo_failed")
		return
	}

	var googleUser GoogleUserInfo
	if err := json.Unmarshal(body, &googleUser); err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=parse_userinfo_failed")
		return
	}

	if !googleUser.EmailVerified {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=email_not_verified")
		return
	}

	userRepo := repository.NewUserRepository(config.DB)
	user, err := userRepo.FindOrCreateByGoogle(
		googleUser.Sub,
		googleUser.Email,
		googleUser.Name,
		googleUser.Picture,
	)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=create_user_failed")
		return
	}

	if user.Locked {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=account_locked")
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=generate_token_failed")
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect,
			redirectBase+"?error=generate_refresh_token_failed")
		return
	}

	redirectURL := fmt.Sprintf(
		redirectBase+"?accessToken=%s&refreshToken=%s",
		accessToken,
		refreshToken,
	)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

/*  PROFILE  */

func Profile(c *gin.Context) {
	claimsAny, exists := c.Get(auth.ContextUserKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	claims, ok := claimsAny.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userRepo := repository.NewUserRepository(config.DB)
	user, err := userRepo.FindByID(claims.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"username": user.Username,
		"avatar":   user.Avatar,
		// "email":    user.Email,
		"role": user.Role,
	})
}
func RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "refresh_token is required",
		})
		return
	}

	refreshClaims, err := auth.ParseRefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid refresh token",
		})
		return
	}

	userRepo := repository.NewUserRepository(config.DB)
	user, err := userRepo.FindByID(refreshClaims.UserID)
	if err != nil || user.Locked {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not found or locked",
		})
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate access token",
		})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not generate refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
