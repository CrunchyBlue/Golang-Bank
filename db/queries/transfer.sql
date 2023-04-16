-- name: CreateTransfer :one
INSERT INTO transfer (source_account_id,
                      destination_account_id,
                      amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetTransfer :one
SELECT *
FROM transfer
WHERE id = $1;

-- name: GetTransfers :many
SELECT *
FROM transfer
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: GetOutboundTransfersForAccount :many
SELECT *
FROM transfer
WHERE source_account_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: GetInboundTransfersForAccount :many
SELECT *
FROM transfer
WHERE destination_account_id = $1
ORDER BY id
LIMIT $2 OFFSET $3;

-- name: UpdateTransfer :one
UPDATE transfer
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransfer :exec
DELETE
FROM transfer
WHERE id = $1;