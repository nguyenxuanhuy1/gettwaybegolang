package service

import (
	"context"
	"database/sql"
	"errors"
	"gateway/internal/repository/postgres"
	"gateway/pkg/utils"
)

type AuthService struct {
	apiKeyRepo *postgres.APIKeyRepository
	userRepo   *postgres.UserRepository
}

func NewAuthService(apiKeyRepo *postgres.APIKeyRepository, userRepo *postgres.UserRepository) *AuthService {
	return &AuthService{
		apiKeyRepo: apiKeyRepo,
		userRepo:   userRepo,
	}
}

// VerifyAPIKey verifies an API key and returns the user ID
func (s *AuthService) VerifyAPIKey(ctx context.Context, rawKey string) (int, error) {
	// Hash the key
	keyHash := utils.HashAPIKey(rawKey)

	// Verify key exists and not revoked
	apiKey, err := s.apiKeyRepo.GetByHash(ctx, keyHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("invalid API key")
		}
		return 0, err
	}

	// Update last used timestamp (async, don't block on error)
	go s.apiKeyRepo.UpdateLastUsed(context.Background(), apiKey.UserID)

	return apiKey.UserID, nil
}
