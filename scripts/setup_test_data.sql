-- Quick Setup Script for Testing
-- Run this after creating your first user via /auth/create-key

-- 1. Create a free product plan
INSERT INTO products (code, name, rate_limit, active)
VALUES ('free', 'Free Plan', 100, true)
ON CONFLICT (code) DO NOTHING;

-- 2. Set pricing (2 coins per request)
INSERT INTO product_prices (product_code, unit, price, active)
VALUES ('free', 'request', 2, true);

-- 3. Top-up 1000 coins for user ID 1 (change ID as needed)
UPDATE users SET coin = 1000 WHERE id = 1;

-- 4. Assign plan to user (expires in 30 days)
INSERT INTO user_products (user_id, product_code, active, started_at, expired_at)
VALUES (1, 'free', true, NOW(), NOW() + INTERVAL '30 days')
ON CONFLICT (user_id, product_code) DO UPDATE SET
    active = true,
    started_at = NOW(),
    expired_at = NOW() + INTERVAL '30 days';

-- Verify setup
SELECT 
    u.id,
    u.email,
    u.coin,
    up.product_code,
    p.name as plan_name,
    p.rate_limit,
    up.expired_at,
    pp.price as cost_per_request
FROM users u
LEFT JOIN user_products up ON u.id = up.user_id AND up.active = true
LEFT JOIN products p ON up.product_code = p.code
LEFT JOIN product_prices pp ON p.code = pp.product_code AND pp.unit = 'request'
WHERE u.id = 1;
