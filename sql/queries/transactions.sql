-- name: GetTransactions :many
SELECT
  t.id,
  t.user_id,
  t.category_id,
  t.date,
  t.amount,
  t.is_optional,
  t.created_at,
  t.updated_at,
  c.name  AS category_name,
  c.type  AS category_type
FROM transactions t
JOIN categories c
  ON t.category_id = c.id
WHERE t.user_id = $1
ORDER BY t.date DESC;

-- name: CreateTransaction :one
WITH inserted AS (
    INSERT INTO transactions (
        id,
        user_id,
        category_id,
        date,
        amount,
        is_optional
    )
    VALUES (
        gen_random_uuid(),
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT i.*, c.name, c.type
FROM inserted i
JOIN categories c ON i.category_id = c.id;

-- name: DeleteTransaction :execrows
DELETE FROM transactions
WHERE id = $1 AND user_id = $2;