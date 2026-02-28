-- name: GetCategoriesIncome :many
SELECT tc.id, tc.name, tc.type
FROM (
SELECT categories.*, COUNT(transactions.id) AS transaction_count
FROM categories 
LEFT JOIN transactions 
ON transactions.category_id = categories.id
WHERE categories.type = 'income'
GROUP BY categories.id) tc
ORDER BY transaction_count DESC;

-- name: GetCategoriesExpense :many
SELECT tc.id, tc.name, tc.type
FROM (
SELECT categories.*, COUNT(transactions.id) AS transaction_count
FROM categories 
LEFT JOIN transactions 
ON transactions.category_id = categories.id
WHERE categories.type = 'expense'
GROUP BY categories.id) tc
ORDER BY transaction_count DESC;

-- name: GetCategories :many
SELECT tc.id, tc.name, tc.type
FROM (
SELECT categories.*, COUNT(transactions.id) AS transaction_count
FROM categories 
LEFT JOIN transactions 
ON transactions.category_id = categories.id
GROUP BY categories.id) tc
ORDER BY transaction_count DESC;

-- name: CreateCategory :one
INSERT INTO categories (name, type)
VALUES ($1, $2)
RETURNING id, name, type;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

-- name: GetCategoryID :one
SELECT id
FROM categories
WHERE name = $1;