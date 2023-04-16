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

-- name: ListTransfers :many
SELECT *
FROM transfer
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateTransfer :one
UPDATE transfer
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteTransfer :exec
DELETE
FROM transfer
WHERE id = $1;