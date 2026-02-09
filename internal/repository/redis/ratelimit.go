package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimitRepository struct {
	client *redis.Client
}

func NewRateLimitRepository(client *redis.Client) *RateLimitRepository {
	return &RateLimitRepository{client: client}
}

// IncrementAndCheck increments the rate limit counter and checks if limit is exceeded
// Returns current count and error if limit exceeded
func (r *RateLimitRepository) IncrementAndCheck(ctx context.Context, userID int, limit int) (int64, error) {
	now := time.Now()
	minute := now.Format("2006-01-02-15:04")
	key := fmt.Sprintf("rl:%d:%s", userID, minute)

	// Increment counter
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Set expiration on first increment
	if count == 1 {
		r.client.Expire(ctx, key, 60*time.Second)
	}

	return count, nil
}
