package service

import (
	"context"
	"database/sql"
	"errors"
	"gateway/internal/domain"
	"gateway/internal/repository/postgres"
	"time"
)

type PlanService struct {
	productRepo *postgres.ProductRepository
}

func NewPlanService(productRepo *postgres.ProductRepository) *PlanService {
	return &PlanService{
		productRepo: productRepo,
	}
}

// ValidateUserPlan checks if user has an active, non-expired plan
func (s *PlanService) ValidateUserPlan(ctx context.Context, userID int) (*domain.UserPlanInfo, error) {
	plan, err := s.productRepo.GetActiveUserProduct(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("no active plan found")
		}
		return nil, err
	}

	// Check if plan is expired
	if plan.ExpiredAt != nil && plan.ExpiredAt.Before(time.Now()) {
		return nil, errors.New("plan has expired")
	}

	return plan, nil
}
