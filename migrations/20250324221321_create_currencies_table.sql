-- +goose Up
CREATE TABLE IF NOT EXISTS currencies (
    code TEXT PRIMARY KEY,
    rate FLOAT DEFAULT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true,
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS currencies;