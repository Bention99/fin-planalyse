-- name: GetCategories :many
SELECT * FROM categories
ORDER BY type;

-- name: CreateCategory :one
INSERT INTO categories (name, type)
VALUES ($1, $2)
RETURNING id, name, type;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;