-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
OFFSET $1 LIMIT $2;

-- name: CreateAccount :one
INSERT INTO accounts (
  owner, balance, currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateAccount :one
UPDATE accounts SET balance = $2
WHERE id = $1
RETURNING *;

-- name: EnoughAccountBalance :one
SELECT (balance >= $1) FROM accounts WHERE id = $2;

-- name: AddAccountBalance :one
UPDATE accounts SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;