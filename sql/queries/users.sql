-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password)
VALUES (gen_random_uuid(), $1, $2)
RETURNING id, email, hashed_password;

-- name: GetUserByEmail :one
SELECT id, email, hashed_password
FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT id, email
FROM users
WHERE id = $1;