package domain

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	GoogleID  string    `json:"google_id"`
	Role      string    `json:"role"`
	Avatar    *string   `json:"avatar"`
	Coin      int       `json:"coin"`
	Locked    bool      `json:"locked"`
	CreatedAt time.Time `json:"created_at"`
}

// APIKey represents an API key for authentication
type APIKey struct {
	UserID     int        `json:"user_id"`
	KeyHash    string     `json:"key_hash"`
	Revoked    bool       `json:"revoked"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Product represents a subscription product
type Product struct {
	ID        int    `json:"id"`
	Code      string `json:"code"`
	Name      string `json:"name"`
	RateLimit *int   `json:"rate_limit"`
	Active    bool   `json:"active"`
}

// UserProduct represents a user's active subscription
type UserProduct struct {
	UserID      int        `json:"user_id"`
	ProductCode string     `json:"product_code"`
	Active      bool       `json:"active"`
	StartedAt   time.Time  `json:"started_at"`
	ExpiredAt   *time.Time `json:"expired_at"`
}

// ProductPrice represents pricing for a product
type ProductPrice struct {
	ID          int       `json:"id"`
	ProductCode string    `json:"product_code"`
	Unit        string    `json:"unit"` // 'request', 'upload', 'gb'
	Price       int       `json:"price"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
}

// CoinTransaction represents a coin transaction
type CoinTransaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Amount    int       `json:"amount"` // negative for deduct, positive for topup
	Type      string    `json:"type"`   // 'topup' or 'deduct'
	Reason    *string   `json:"reason"`
	RequestID *string   `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}

// APILog represents an API request log
type APILog struct {
	ID          int64      `json:"id"`
	UserID      int        `json:"user_id"`
	ProductCode *string    `json:"product_code"`
	Endpoint    *string    `json:"endpoint"`
	Cost        *int       `json:"cost"`
	RequestID   *string    `json:"request_id"`
	CreatedAt   time.Time  `json:"created_at"`
}

// UserPlanInfo combines user product with product details
type UserPlanInfo struct {
	UserID      int
	ProductCode string
	ProductName string
	RateLimit   *int
	ExpiredAt   *time.Time
}
