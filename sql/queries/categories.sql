-- name: GetCategoriesIncome :many
SELECT c.id, c.name, c.type, c.is_system
FROM categories c
LEFT JOIN transactions t
ON t.category_id = c.id
WHERE c.type = 'income'
AND c.is_system = true
OR c.type = 'income'
AND c.user_id = $1
GROUP BY c.id, c.name, c.type, c.is_system
ORDER BY COUNT(t.id) DESC;

-- name: GetCategoriesExpense :many
SELECT c.id, c.name, c.type, c.is_system
FROM categories c
LEFT JOIN transactions t
ON t.category_id = c.id
WHERE c.type = 'expense'
AND c.is_system = true
OR c.type = 'expense'
AND c.user_id = $1
GROUP BY c.id, c.name, c.type, c.is_system
ORDER BY COUNT(t.id) DESC;

-- name: GetCategories :many
SELECT c.id, c.name, c.type, c.is_system
FROM categories c
LEFT JOIN transactions t
ON t.category_id = c.id
WHERE c.is_system = true
OR c.user_id = $1
GROUP BY c.id, c.name, c.type, c.is_system
ORDER BY COUNT(t.id) DESC;

-- name: CreateCategory :one
INSERT INTO categories (name, type, user_id)
VALUES ($1, $2, $3)
RETURNING id, name, type;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1
AND is_system = 'false';

-- name: GetCategoryID :one
SELECT id
FROM categories
WHERE name = $1;