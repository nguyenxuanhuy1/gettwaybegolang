create table users (
   id         serial primary key,
   username   text not null,
   role       text not null default 'user',
   avatar     text,
   locked     boolean default false,
   coin       integer not null default 0,
   google_id  text not null unique,
   email      text not null unique,
   created_at timestamptz default now()
);

create index idx_users_google_id on
   users (
      google_id
   );
create index idx_users_email on
   users (
      email
   );

create table products (
   id         serial primary key,
   code       text unique not null,
   name       text not null,
   rate_limit integer,
   active     boolean default true
);

create table user_products (
   user_id      integer
      references users ( id ),
   product_code text
      references products ( code ),
   active       boolean default true,
   started_at   timestamptz default now(),
   expired_at   timestamptz,
   primary key ( user_id,
                 product_code )
);

create index idx_user_products_active on
   user_products (
      user_id,
      active
   );

create table product_prices (
   id           serial primary key,
   product_code text
      references products ( code ),
   unit         text not null check ( unit in ( 'request',
                                        'upload',
                                        'gb' ) ),
   price        integer not null check ( price >= 0 ),
   active       boolean default true,
   created_at   timestamptz default now()
);

create table coin_transactions (
   id         serial primary key,
   user_id    integer not null
      references users ( id ),
   amount     integer not null check ( amount <> 0 ),
   type       text not null check ( type in ( 'topup',
                                        'deduct' ) ),
   reason     text,
   request_id text unique,
   created_at timestamptz default now()
);

create index idx_coin_tx_user on
   coin_transactions (
      user_id
   );

create table api_keys (
   user_id      integer primary key
      references users ( id ),
   key_hash     text not null unique,
   last_used_at timestamptz,
   revoked      boolean default false,
   created_at   timestamptz default now()
);

create index idx_api_keys_hash on
   api_keys (
      key_hash
   );

create table api_logs (
   id           bigserial primary key,
   user_id      integer
      references users ( id ),
   product_code text,
   endpoint     text,
   cost         integer,
   request_id   text,
   created_at   timestamptz default now()
);