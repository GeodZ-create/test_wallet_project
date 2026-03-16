-- +goose Up
CREATE TABLE IF NOT EXISTS wallets(
    id SERIAL PRIMARY KEY,
    wallet_id UUID UNIQUE NOT NULL,
    balance BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS wallets;