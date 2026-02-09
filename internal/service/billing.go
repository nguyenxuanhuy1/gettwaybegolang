package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gateway/internal/repository/postgres"
)

type BillingService struct {
	db          *sql.DB
	userRepo    *postgres.UserRepository
	productRepo *postgres.ProductRepository
	txRepo      *postgres.TransactionRepository
}

func NewBillingService(
	db *sql.DB,
	userRepo *postgres.UserRepository,
	productRepo *postgres.ProductRepository,
	txRepo *postgres.TransactionRepository,
) *BillingService {
	return &BillingService{
		db:          db,
		userRepo:    userRepo,
		productRepo: productRepo,
		txRepo:      txRepo,
	}
}

// DeductCoins deducts coins for a request (idempotent via request_id)
func (s *BillingService) DeductCoins(ctx context.Context, userID int, productCode string, requestID string) (int, error) {
	// Check if request_id already exists (idempotency)
	exists, err := s.txRepo.CheckRequestIDExists(ctx, requestID)
	if err != nil {
		return 0, err
	}
	if exists {
		// Already processed, return 0 cost
		return 0, nil
	}

	// Get price for this product
	price, err := s.productRepo.GetRequestPrice(ctx, productCode)
	if err != nil {
		if err == sql.ErrNoRows {
			// No price configured, allow free requests
			return 0, nil
		}
		return 0, err
	}

	// Start transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// Lock user row and get current balance
	user, err := s.userRepo.GetByID(ctx, tx, userID)
	if err != nil {
		return 0, err
	}

	// Check sufficient balance
	if user.Coin < price {
		return 0, errors.New("insufficient coins")
	}

	// Update balance
	newBalance := user.Coin - price
	err = s.userRepo.UpdateCoinBalance(ctx, tx, userID, newBalance)
	if err != nil {
		return 0, err
	}

	// Record transaction
	reason := fmt.Sprintf("API request to %s", productCode)
	err = s.txRepo.CreateTransaction(ctx, tx, userID, -price, "deduct", reason, requestID)
	if err != nil {
		return 0, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return price, nil
}
