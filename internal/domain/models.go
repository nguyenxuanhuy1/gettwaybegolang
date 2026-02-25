package domain

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	GoogleID  string    `json:"google_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Avatar    *string   `json:"avatar"`
	Role      string    `json:"role"` // 'user' | 'admin'
	Locked    bool      `json:"locked"`
	CreatedAt time.Time `json:"created_at"`
}

// Session represents a user login session
type Session struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// APIKey represents an API key for authentication
type APIKey struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Name       string     `json:"name"`
	KeyHash    string     `json:"-"`
	Revoked    bool       `json:"revoked"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Plan represents a subscription plan (aligns with DB table `plans`)
type Plan struct {
	ID               int64     `json:"id"`
	Code             string    `json:"code"`
	Name             string    `json:"name"`
	Price            int64     `json:"price"`
	RateLimitPerSec  int       `json:"rate_limit_per_sec"`
	MonthlyQuota     *int64    `json:"monthly_quota"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
}

// Subscription represents a user's plan subscription (aligns with DB table `user_subscriptions`)
type Subscription struct {
	ID        int64      `json:"id"`
	UserID    int64      `json:"user_id"`
	PlanID    int64      `json:"plan_id"`
	Status    string     `json:"status"` // 'active' | 'expired' | 'cancelled' | 'trial'
	StartedAt time.Time  `json:"started_at"`
	ExpiredAt *time.Time `json:"expired_at"`
	CreatedAt time.Time  `json:"created_at"`
}

// SubscriptionWithPlan combines subscription with its plan details
type SubscriptionWithPlan struct {
	Subscription
	PlanCode        string `json:"plan_code"`
	PlanName        string `json:"plan_name"`
	RateLimitPerSec int    `json:"rate_limit_per_sec"`
}

// Wallet represents a user's wallet
type Wallet struct {
	UserID    int64     `json:"user_id"`
	Balance   int64     `json:"balance"`
	UpdatedAt time.Time `json:"updated_at"`
}

// WalletTransaction represents a wallet transaction
type WalletTransaction struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Amount    int64     `json:"amount"` // positive=topup, negative=deduct
	Type      string    `json:"type"`   // 'topup' | 'deduct' | 'refund'
	Reason    *string   `json:"reason"`
	RequestID *string   `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}

// UsageLog represents an API request log (aligns with DB table `usage_logs`)
type UsageLog struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"user_id"`
	APIKeyID  *int64    `json:"api_key_id"`
	Endpoint  *string   `json:"endpoint"`
	Cost      int64     `json:"cost"`
	RequestID *string   `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}
