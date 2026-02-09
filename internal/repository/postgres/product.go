package postgres

import (
	"context"
	"database/sql"
	"gateway/internal/domain"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// GetActiveUserProduct retrieves user's active product with details
func (r *ProductRepository) GetActiveUserProduct(ctx context.Context, userID int) (*domain.UserPlanInfo, error) {
	query := `
		SELECT 
			up.user_id,
			up.product_code,
			p.name,
			p.rate_limit,
			up.expired_at
		FROM user_products up
		JOIN products p ON up.product_code = p.code
		WHERE up.user_id = $1 AND up.active = true
		LIMIT 1
	`

	var plan domain.UserPlanInfo
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&plan.UserID, &plan.ProductCode, &plan.ProductName,
		&plan.RateLimit, &plan.ExpiredAt,
	)

	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// GetRequestPrice retrieves the price per request for a product
func (r *ProductRepository) GetRequestPrice(ctx context.Context, productCode string) (int, error) {
	query := `
		SELECT price
		FROM product_prices
		WHERE product_code = $1 AND unit = 'request' AND active = true
		LIMIT 1
	`

	var price int
	err := r.db.QueryRowContext(ctx, query, productCode).Scan(&price)
	if err != nil {
		return 0, err
	}
	return price, nil
}
