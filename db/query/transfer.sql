-- name: GetTransfer :one
select * from transfers WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
select * from transfers ORDER BY created_at DESC OFFSET $1 LIMIT $2;

-- name: CreateTransfer :one
INSERT INTO transfers(from_account_id, to_account_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListTransfersByFrom :many
select * from transfers
WHERE from_account_id = $1
ORDER BY created_at DESC OFFSET $2 LIMIT $3;
