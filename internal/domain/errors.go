package domain

import "errors"

// Sentinel errors used across the application
var (
	ErrInsufficientBalance = errors.New("insufficient wallet balance")
	ErrNotFound            = errors.New("not found")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrDuplicate           = errors.New("duplicate entry")
)
