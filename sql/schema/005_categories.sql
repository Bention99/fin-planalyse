-- +goose Up
ALTER TABLE categories
ADD COLUMN user_id UUID NULL,
ADD COLUMN is_system BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE categories
DROP COLUMN user_id
DROP COLUMN is_system;