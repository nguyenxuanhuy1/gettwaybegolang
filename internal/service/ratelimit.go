package service

import (
	"context"
	"errors"
	"gateway/internal/repository/redis"
)

type RateLimitService struct {
	rateLimitRepo *redis.RateLimitRepository
}

func NewRateLimitService(rateLimitRepo *redis.RateLimitRepository) *RateLimitService {
	return &RateLimitService{
		rateLimitRepo: rateLimitRepo,
	}
}

// CheckRateLimit checks if user has exceeded their rate limit
func (s *RateLimitService) CheckRateLimit(ctx context.Context, userID int, limit int) error {
	if limit <= 0 {
		// No rate limit set
		return nil
	}

	count, err := s.rateLimitRepo.IncrementAndCheck(ctx, userID, limit)
	if err != nil {
		return err
	}

	if count > int64(limit) {
		return errors.New("rate limit exceeded")
	}

	return nil
}
