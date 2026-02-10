package router

import (
	"database/sql"
	"gateway/internal/handler"
	"gateway/internal/middleware"
	"gateway/internal/repository/postgres"
	"gateway/internal/repository/redis"
	"gateway/internal/service"

	"github.com/gin-gonic/gin"
	redisClient "github.com/redis/go-redis/v9"
)

func SetupRouter(db *sql.DB, redisConn *redisClient.Client) *gin.Engine {
	r := gin.Default()

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	apiKeyRepo := postgres.NewAPIKeyRepository(db)
	productRepo := postgres.NewProductRepository(db)
	txRepo := postgres.NewTransactionRepository(db)
	apiLogRepo := postgres.NewAPILogRepository(db)
	rateLimitRepo := redis.NewRateLimitRepository(redisConn)

	// Initialize services
	authService := service.NewAuthService(apiKeyRepo, userRepo)
	planService := service.NewPlanService(productRepo)
	rateLimitService := service.NewRateLimitService(rateLimitRepo)
	billingService := service.NewBillingService(db, userRepo, productRepo, txRepo)
	usageService := service.NewUsageService(userRepo, productRepo, apiLogRepo)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authService, planService, userRepo)

	// Initialize handlers
	rotateKeyHandler := handler.NewRotateKeyHandler(apiKeyRepo)
	usageHandler := handler.NewUsageHandler(usageService)
	requestHandler := handler.NewRequestHandler(rateLimitService, billingService, apiLogRepo)
	googleAuthHandler := handler.NewGoogleAuthHandler(userRepo)
	createKeyHandler := handler.NewCreateKeyHandler(userRepo, apiKeyRepo)

	// Public routes (health check)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
 
	// Google OAuth routes (public)
	auth := r.Group("/api")
	{
		auth.GET("/auth/google", googleAuthHandler.HandleLogin)
		auth.GET("/auth/google/callback", googleAuthHandler.HandleCallback)
	}

	// API key creation (public - called after Google login from frontend)
	r.POST("/api/auth/create-key", createKeyHandler.Handle)

	// Protected API routes
	api := r.Group("/api")
	api.Use(authMiddleware.RequireAPIKey()) 
	{
		api.POST("/rotate-key", rotateKeyHandler.Handle)
		api.GET("/usage", usageHandler.Handle)
		api.POST("/request", requestHandler.Handle)
	}

	return r
}
