-- +goose Up
CREATE TYPE transaction_type AS ENUM ('income', 'expense');

CREATE TABLE categories(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL UNIQUE,
    type transaction_type NOT NULL
);

INSERT INTO categories (name, type) VALUES
('rent', 'expense'),
('groceries', 'expense'),
('utilities', 'expense'),
('paycheck', 'income'),
('dividends', 'income');

-- +goose Down
DROP TABLE categories;
DROP TYPE transaction_type;