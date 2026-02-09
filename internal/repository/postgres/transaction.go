package postgres

import (
	"context"
	"database/sql"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// CheckRequestIDExists checks if a request_id already exists (idempotency)
func (r *TransactionRepository) CheckRequestIDExists(ctx context.Context, requestID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM coin_transactions WHERE request_id = $1)`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, requestID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// CreateTransaction creates a new coin transaction
func (r *TransactionRepository) CreateTransaction(ctx context.Context, tx *sql.Tx, userID int, amount int, txType string, reason string, requestID string) error {
	query := `
		INSERT INTO coin_transactions (user_id, amount, type, reason, request_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := tx.ExecContext(ctx, query, userID, amount, txType, reason, requestID, time.Now())
	return err
}
