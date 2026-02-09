package postgres

import (
	"context"
	"database/sql"
	"time"
)

type APILogRepository struct {
	db *sql.DB
}

func NewAPILogRepository(db *sql.DB) *APILogRepository {
	return &APILogRepository{db: db}
}

// Create inserts a new API log entry
func (r *APILogRepository) Create(ctx context.Context, userID int, productCode, endpoint string, cost int, requestID string) error {
	query := `
		INSERT INTO api_logs (user_id, product_code, endpoint, cost, request_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := r.db.ExecContext(ctx, query, userID, productCode, endpoint, cost, requestID, time.Now())
	return err
}

// GetUsageStats retrieves usage statistics for a user
func (r *APILogRepository) GetUsageStats(ctx context.Context, userID int) (int64, int, error) {
	query := `
		SELECT 
			COUNT(*) as request_count,
			COALESCE(SUM(cost), 0) as total_cost
		FROM api_logs
		WHERE user_id = $1
	`
	
	var requestCount int64
	var totalCost int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&requestCount, &totalCost)
	if err != nil {
		return 0, 0, err
	}
	
	return requestCount, totalCost, nil
}
