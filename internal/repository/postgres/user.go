package postgres

import (
	"context"
	"database/sql"
	"gateway/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves a user by ID with FOR UPDATE lock (for transactions)
func (r *UserRepository) GetByID(ctx context.Context, tx *sql.Tx, userID int) (*domain.User, error) {
	query := `
		SELECT id, username, email, google_id, role, avatar, coin, locked, created_at
		FROM users
		WHERE id = $1
		FOR UPDATE
	`

	var user domain.User
	var err error

	if tx != nil {
		err = tx.QueryRowContext(ctx, query, userID).Scan(
			&user.ID, &user.Username, &user.Email, &user.GoogleID,
			&user.Role, &user.Avatar, &user.Coin, &user.Locked, &user.CreatedAt,
		)
	} else {
		err = r.db.QueryRowContext(ctx, query, userID).Scan(
			&user.ID, &user.Username, &user.Email, &user.GoogleID,
			&user.Role, &user.Avatar, &user.Coin, &user.Locked, &user.CreatedAt,
		)
	}

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByIDSimple retrieves a user without locking
func (r *UserRepository) GetByIDSimple(ctx context.Context, userID int) (*domain.User, error) {
	query := `
		SELECT id, username, email, google_id, role, avatar, coin, locked, created_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.GoogleID,
		&user.Role, &user.Avatar, &user.Coin, &user.Locked, &user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateCoinBalance updates user's coin balance
func (r *UserRepository) UpdateCoinBalance(ctx context.Context, tx *sql.Tx, userID int, newBalance int) error {
	query := `UPDATE users SET coin = $1 WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, newBalance, userID)
	return err
}

// GetDB returns the database connection (for special queries)
func (r *UserRepository) GetDB() *sql.DB {
	return r.db
}
