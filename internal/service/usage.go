package service

import (
	"context"
	"gateway/internal/repository/postgres"
)

type UsageService struct {
	userRepo   *postgres.UserRepository
	productRepo *postgres.ProductRepository
	apiLogRepo *postgres.APILogRepository
}

func NewUsageService(
	userRepo *postgres.UserRepository,
	productRepo *postgres.ProductRepository,
	apiLogRepo *postgres.APILogRepository,
) *UsageService {
	return &UsageService{
		userRepo:   userRepo,
		productRepo: productRepo,
		apiLogRepo: apiLogRepo,
	}
}

type UsageStats struct {
	UserID       int     `json:"user_id"`
	CoinBalance  int     `json:"coin_balance"`
	ProductCode  string  `json:"product_code"`
	ProductName  string  `json:"product_name"`
	RateLimit    *int    `json:"rate_limit"`
	RequestCount int64   `json:"request_count"`
	TotalCost    int     `json:"total_cost"`
	PlanExpired  *string `json:"plan_expired_at"`
}

// GetUsageStats retrieves comprehensive usage statistics for a user
func (s *UsageService) GetUsageStats(ctx context.Context, userID int) (*UsageStats, error) {
	// Get user info
	user, err := s.userRepo.GetByIDSimple(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get active plan
	plan, err := s.productRepo.GetActiveUserProduct(ctx, userID)
	if err != nil {
		// User might not have a plan yet
		plan = nil
	}

	// Get usage stats
	requestCount, totalCost, err := s.apiLogRepo.GetUsageStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	stats := &UsageStats{
		UserID:       user.ID,
		CoinBalance:  user.Coin,
		RequestCount: requestCount,
		TotalCost:    totalCost,
	}

	if plan != nil {
		stats.ProductCode = plan.ProductCode
		stats.ProductName = plan.ProductName
		stats.RateLimit = plan.RateLimit
		if plan.ExpiredAt != nil {
			expiredStr := plan.ExpiredAt.Format("2006-01-02 15:04:05")
			stats.PlanExpired = &expiredStr
		}
	}

	return stats, nil
}
