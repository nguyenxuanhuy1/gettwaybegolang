package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AppConfig struct {
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPass                  string
	DBName                  string
	GoogleClientID          string
	GoogleClientSecret      string
	GoogleRedirectURL       string
	FrontendAuthRedirectURL string
	KongAdminURL            string
	InternalSecret          string
}

var Config AppConfig
var GoogleOAuthConfig *oauth2.Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	Config = AppConfig{
		DBHost:                  os.Getenv("DB_HOST"),
		DBPort:                  os.Getenv("DB_PORT"),
		DBUser:                  os.Getenv("DB_USER"),
		DBPass:                  os.Getenv("DB_PASSWORD"),
		DBName:                  os.Getenv("DB_NAME"),
		GoogleClientID:          os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:      os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:       os.Getenv("GOOGLE_REDIRECT_URL"),
		FrontendAuthRedirectURL: os.Getenv("FRONTEND_AUTH_REDIRECT_URL"),
		KongAdminURL:            os.Getenv("KONG_ADMIN_URL"),
		InternalSecret:          os.Getenv("INTERNAL_SECRET"),
	}

	InitGoogleOAuth()
}

func InitGoogleOAuth() {
	GoogleOAuthConfig = &oauth2.Config{
		ClientID:     Config.GoogleClientID,
		ClientSecret: Config.GoogleClientSecret,
		RedirectURL:  Config.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
