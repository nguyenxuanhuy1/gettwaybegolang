create table users (
   id         bigserial primary key,
   google_id  text not null unique,
   email      text not null unique,
   username   text,
   avatar     text,
   role       text not null default 'user' check ( role in ( 'user',
                                                       'admin' ) ),
   locked     boolean not null default false,
   created_at timestamptz not null default now()
);

create index idx_users_email on
   users (
      email
   );
create index idx_users_google_id on
   users (
      google_id
   );


-- API KEYS
create table api_keys (
   id           bigserial primary key,
   user_id      bigint not null
      references users ( id )
         on delete cascade,
   name         text not null,
   key_hash     text not null unique,
   revoked      boolean not null default false,
   last_used_at timestamptz,
   created_at   timestamptz not null default now()
);

create index idx_api_keys_user on
   api_keys (
      user_id
   );
create index idx_api_keys_hash on
   api_keys (
      key_hash
   );


-- PLANS
create table plans (
   id                 bigserial primary key,
   code               text not null unique,
   name               text not null,
   price              bigint not null default 0,
   rate_limit_per_sec integer not null,
   monthly_quota      bigint,
   active             boolean not null default true,
   created_at         timestamptz not null default now()
);


-- USER SUBSCRIPTIONS
create table user_subscriptions (
   id         bigserial primary key,
   user_id    bigint not null
      references users ( id )
         on delete cascade,
   plan_id    bigint not null
      references plans ( id ),
   status     text not null check ( status in ( 'active',
                                            'expired',
                                            'cancelled',
                                            'trial' ) ),
   started_at timestamptz not null default now(),
   expired_at timestamptz,
   created_at timestamptz not null default now()
);

create index idx_user_subscriptions_user on
   user_subscriptions (
      user_id
   );

create index idx_user_subscriptions_status on
   user_subscriptions (
      status
   );


-- WALLET (1 user = 1 ví)
create table wallets (
   user_id    bigint primary key
      references users ( id )
         on delete cascade,
   balance    bigint not null default 0,
   updated_at timestamptz not null default now()
);


-- WALLET TRANSACTIONS
create table wallet_transactions (
   id         bigserial primary key,
   user_id    bigint not null
      references users ( id )
         on delete cascade,
   amount     bigint not null,
   type       text not null check ( type in ( 'topup',
                                        'deduct',
                                        'refund' ) ),
   reason     text,
   request_id text unique,
   created_at timestamptz not null default now()
);

create index idx_wallet_tx_user on
   wallet_transactions (
      user_id
   );


-- USAGE LOGS
create table usage_logs (
   id         bigserial primary key,
   user_id    bigint
      references users ( id )
         on delete set null,
   api_key_id bigint
      references api_keys ( id )
         on delete set null,
   endpoint   text,
   cost       bigint not null default 0,
   request_id text,
   created_at timestamptz not null default now()
);

create index idx_usage_user on
   usage_logs (
      user_id
   );
create index idx_usage_created on
   usage_logs (
      created_at
   );


-- SAMPLE PLANS (OPTIONAL)
insert into plans (
   code,
   name,
   price,
   rate_limit_per_sec,
   monthly_quota
) values ( 'free',
           'Free Plan',
           0,
           3,
           10000 ),( 'basic',
                     'Basic Plan',
                     10000,
                     10,
                     100000 ),( 'pro',
                                'Pro Plan',
                                30000,
                                20,
                                500000 );