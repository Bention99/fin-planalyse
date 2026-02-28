-- +goose Up

ALTER TABLE categories
DROP CONSTRAINT IF EXISTS categories_name_key;

CREATE UNIQUE INDEX IF NOT EXISTS categories_user_id_lower_name_unique
ON categories (user_id, lower(name))
WHERE user_id IS NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS categories_system_lower_name_unique
ON categories (lower(name))
WHERE is_system = true;

-- +goose Down
DROP INDEX IF EXISTS categories_system_name_unique;
DROP INDEX IF EXISTS categories_user_id_name_unique;

ALTER TABLE categories
ADD CONSTRAINT categories_name_key UNIQUE (name);