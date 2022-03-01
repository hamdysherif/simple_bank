-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (
  username, full_name, email, hashed_password
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;
