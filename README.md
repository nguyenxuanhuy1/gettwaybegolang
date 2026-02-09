# API Gateway

Production-ready API Gateway for SaaS with API key authentication, rate limiting, and coin-based billing.

## Features

- ✅ **API Key Authentication** - SHA256 hashed keys with Bearer token format
- ✅ **Plan Management** - User subscription validation with expiration
- ✅ **Rate Limiting** - Redis-based per-user rate limiting
- ✅ **Coin Billing** - Transactional coin deduction with idempotency
- ✅ **Usage Tracking** - Comprehensive API usage statistics
- ✅ **Clean Architecture** - Handler → Service → Repository layers

## Project Structure

```
F:\BE\GATEWAY\
├── cmd/server/main.go              # Application entry point
├── config/                          # Configuration
│   ├── config.go                   # Environment config
│   ├── db.go                       # PostgreSQL connection
│   └── redis.go                    # Redis connection
├── internal/
│   ├── domain/models.go            # Domain models
│   ├── repository/                 # Data access layer
│   │   ├── postgres/               # PostgreSQL repositories
│   │   │   ├── user.go
│   │   │   ├── apikey.go
│   │   │   ├── product.go
│   │   │   ├── transaction.go
│   │   │   └── apilog.go
│   │   └── redis/
│   │       └── ratelimit.go        # Rate limiting
│   ├── service/                    # Business logic
│   │   ├── auth.go                 # Authentication
│   │   ├── plan.go                 # Plan validation
│   │   ├── ratelimit.go            # Rate limiting
│   │   ├── billing.go              # Coin billing
│   │   └── usage.go                # Usage stats
│   ├── middleware/
│   │   └── auth.go                 # API key middleware
│   ├── handler/                    # HTTP handlers
│   │   ├── rotate_key.go           # POST /api/rotate-key
│   │   ├── usage.go                # GET /api/usage
│   │   └── request.go              # POST /api/request
│   └── router/router.go            # Route setup
├── pkg/utils/crypto.go             # Crypto utilities
└── migrations/                     # Database migrations
```

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis

### Installation

1. **Clone and install dependencies**
   ```bash
   cd F:\BE\GATEWAY
   go mod download
   ```

2. **Configure environment**
   
   The `.env` file is already configured with:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=thanghuy
   DB_NAME=getway
   REDIS_URL=redis://localhost:6379/0
   ```

3. **Run migrations**
   
   Migrations are in `migrations/000001_init_archery.up.sql`
   ```bash
   # Apply migrations using your preferred tool (migrate, goose, etc.)
   ```

4. **Start the server**
   ```bash
   go run cmd/server/main.go
   ```

   Server runs on `http://localhost:8081`

## Quick Start - Get Your API Key

### Option 1: Direct API Key Creation (Recommended for Testing)

```bash
# Create user and get API key
curl -X POST http://localhost:8081/auth/create-key \
  -H "Content-Type: application/json" \
  -d '{
    "google_id": "test_user_123",
    "email": "test@example.com",
    "username": "Test User"
  }'
```

Response will include your API key:
```json
{
  "api_key": "sk_abc123...",
  "user_id": 1,
  "email": "test@example.com"
}
```

### Option 2: Google OAuth Flow

1. Navigate to: `http://localhost:8081/auth/google/login`
2. Login with Google
3. After redirect, call `/auth/create-key` with your Google ID

### Setup Test Data

Run this SQL to setup plan and coins:

```sql
-- See scripts/setup_test_data.sql for full script
INSERT INTO products (code, name, rate_limit, active)
VALUES ('free', 'Free Plan', 100, true);

UPDATE users SET coin = 1000 WHERE id = 1;

INSERT INTO user_products (user_id, product_code, active, expired_at)
VALUES (1, 'free', true, NOW() + INTERVAL '30 days');
```

Or run: `psql -U postgres -d getway -f scripts/setup_test_data.sql`

## API Endpoints

### Public Endpoints

#### Health Check
```bash
GET /health
```

#### Google OAuth
```bash
GET /auth/google/login      # Start OAuth flow
GET /auth/google/callback   # OAuth callback
```

#### Create API Key
```bash
POST /auth/create-key
Content-Type: application/json

{
  "google_id": "123456789",
  "email": "user@gmail.com",
  "username": "User Name"
}
```

**Response:**
```json
{
  "message": "API key created successfully",
  "api_key": "sk_abc123...",
  "user_id": 1,
  "email": "user@gmail.com"
}
```

### Protected Endpoints (Require API Key)

All protected endpoints require:
```
Authorization: Bearer sk_<your_api_key>
```

#### 1. Rotate API Key
```bash
POST /api/rotate-key
```

**Response:**
```json
{
  "message": "API key rotated successfully",
  "api_key": "sk_new_key_here",
  "warning": "Save this key securely. You won't be able to see it again."
}
```

#### 2. Get Usage Statistics
```bash
GET /api/usage
```

**Response:**
```json
{
  "user_id": 1,
  "coin_balance": 1000,
  "product_code": "premium",
  "product_name": "Premium Plan",
  "rate_limit": 100,
  "request_count": 42,
  "total_cost": 84,
  "plan_expired_at": "2026-12-31 23:59:59"
}
```

#### 3. Gateway Request (Demo)
```bash
POST /api/request
```

**Response:**
```json
{
  "request_id": "uuid-here",
  "cost": 2,
  "upstream": {
    "status": "success",
    "message": "Request processed successfully",
    "data": {
      "request_id": "uuid-here",
      "timestamp": "2026-02-09T15:45:00Z",
      "mock": true
    }
  }
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid/missing API key
- `403 Forbidden` - Account locked or plan expired
- `429 Too Many Requests` - Rate limit exceeded
- `402 Payment Required` - Insufficient coins

## Request Flow

1. **Authentication** - Verify API key (SHA256 hash lookup)
2. **User Check** - Ensure user is not locked
3. **Plan Validation** - Check plan expiration
4. **Rate Limiting** - Check Redis counter (per minute)
5. **Coin Deduction** - Transactional billing with idempotency
6. **Logging** - Record API request
7. **Upstream** - Forward to upstream service (mocked)

## Database Schema

Key tables:
- `users` - User accounts with coin balance
- `api_keys` - API key credentials (hashed)
- `products` - Subscription products with rate limits
- `user_products` - Active subscriptions (source of truth for expiration)
- `product_prices` - Pricing per unit (request/upload/gb)
- `coin_transactions` - Transaction log with request_id for idempotency
- `api_logs` - API request logs

## Key Design Decisions

### Idempotency
- `coin_transactions.request_id` prevents double charging
- Each request generates a unique UUID
- Duplicate request_id returns 0 cost

### Rate Limiting
- Redis key pattern: `rl:{user_id}:{minute}`
- Auto-expires after 60 seconds
- Per-user, per-minute counter

### Transaction Safety
- `SELECT ... FOR UPDATE` locks user row
- Atomic coin deduction + transaction log
- Rollback on any error

### API Key Security
- SHA256 hashing for storage
- `sk_` prefix for identification
- Rotation revokes old key atomically

## Development

### Run with hot reload
```bash
# If using air
air
```

### Test endpoints
```bash
# Get usage (replace with your API key)
curl -H "Authorization: Bearer sk_your_key" http://localhost:8081/api/usage

# Make a request
curl -X POST -H "Authorization: Bearer sk_your_key" http://localhost:8081/api/request

# Rotate key
curl -X POST -H "Authorization: Bearer sk_your_key" http://localhost:8081/api/rotate-key
```

## Production Considerations

- [ ] Add proper logging (structured logging with zerolog/zap)
- [ ] Add metrics (Prometheus)
- [ ] Add distributed tracing (OpenTelemetry)
- [ ] Add circuit breakers for upstream calls
- [ ] Add connection pooling configuration
- [ ] Add graceful shutdown
- [ ] Add health checks for dependencies
- [ ] Add rate limit headers (X-RateLimit-*)
- [ ] Add API versioning
- [ ] Add request validation
- [ ] Add comprehensive error handling
- [ ] Add integration tests
