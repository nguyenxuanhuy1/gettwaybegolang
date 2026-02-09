package postgres

import (
	"context"
	"database/sql"
	"gateway/internal/domain"
	"time"
)

type APIKeyRepository struct {
	db *sql.DB
}

func NewAPIKeyRepository(db *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// GetByHash retrieves an API key by its hash
func (r *APIKeyRepository) GetByHash(ctx context.Context, keyHash string) (*domain.APIKey, error) {
	query := `
		SELECT user_id, key_hash, revoked, last_used_at, created_at
		FROM api_keys
		WHERE key_hash = $1 AND revoked = false
	`

	var apiKey domain.APIKey
	err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
		&apiKey.UserID, &apiKey.KeyHash, &apiKey.Revoked,
		&apiKey.LastUsedAt, &apiKey.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &apiKey, nil
}

// UpdateLastUsed updates the last_used_at timestamp
func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, userID int) error {
	query := `UPDATE api_keys SET last_used_at = $1 WHERE user_id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

// RotateKey revokes old key and creates new one
func (r *APIKeyRepository) RotateKey(ctx context.Context, userID int, newKeyHash string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Revoke old key
	_, err = tx.ExecContext(ctx, `UPDATE api_keys SET revoked = true WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}

	// Insert new key
	_, err = tx.ExecContext(ctx, `
		INSERT INTO api_keys (user_id, key_hash, revoked, created_at)
		VALUES ($1, $2, false, $3)
		ON CONFLICT (user_id) DO UPDATE SET
			key_hash = $2,
			revoked = false,
			created_at = $3,
			last_used_at = NULL
	`, userID, newKeyHash, time.Now())

	if err != nil {
		return err
	}

	return tx.Commit()
}
