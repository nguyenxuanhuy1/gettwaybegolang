package repository

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var (
	ErrUserExists = errors.New("user already exists")
)

type User struct {
	ID       int64
	Username string
	Role     string
	Locked   bool
	GoogleID *string
	Email    *string
	Avatar   *string
}

type UserProfile struct {
	ID       int64   `json:"id"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
	Role     string  `json:"role"`
	Locked   bool    `json:"locked"`
	Email    *string `json:"email"`
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindOrCreateByGoogle(googleID, email, username, avatar string) (*User, error) {
	var u User
	err := r.db.QueryRow(`
		SELECT id, username, role, locked, google_id, email, avatar
		FROM users
		WHERE google_id = $1
	`, googleID).Scan(
		&u.ID,
		&u.Username,
		&u.Role,
		&u.Locked,
		&u.GoogleID,
		&u.Email,
		&u.Avatar,
	)

	if err == nil {
		return &u, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	err = r.db.QueryRow(`
		INSERT INTO users (username, google_id, email, avatar, role)
		VALUES ($1, $2, $3, $4, 'user')
		RETURNING id, username, role, locked, google_id, email, avatar
	`, username, googleID, email, avatar).Scan(
		&u.ID,
		&u.Username,
		&u.Role,
		&u.Locked,
		&u.GoogleID,
		&u.Email,
		&u.Avatar,
	)

	if err == nil {
		return &u, nil
	}

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
		err2 := r.db.QueryRow(`
			UPDATE users
			SET google_id = $1,
			    avatar = COALESCE(NULLIF($2, ''), avatar),
			    username = COALESCE(NULLIF($3, ''), username)
			WHERE email = $4
			RETURNING id, username, role, locked, google_id, email, avatar
		`, googleID, avatar, username, email).Scan(
			&u.ID,
			&u.Username,
			&u.Role,
			&u.Locked,
			&u.GoogleID,
			&u.Email,
			&u.Avatar,
		)

		if err2 == nil {
			return &u, nil
		}
		return nil, err
	}

	return nil, err
}

func (r *UserRepository) FindByID(id int64) (*UserProfile, error) {
	var u UserProfile

	err := r.db.QueryRow(`
		SELECT id, username, avatar, role, locked, email
		FROM users
		WHERE id = $1
	`, id).Scan(
		&u.ID,
		&u.Username,
		&u.Avatar,
		&u.Role,
		&u.Locked,
		&u.Email,
	)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
