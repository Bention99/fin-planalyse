-- +goose Up
CREATE TABLE transactions(
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES categories(id),
    date DATE NOT NULL,
    amount BIGINT NOT NULL,
    is_optional BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX idx_transactions_user_date
ON transactions(user_id, date);

-- +goose Down
DROP TABLE transactions;
DROP TYPE transaction_type;