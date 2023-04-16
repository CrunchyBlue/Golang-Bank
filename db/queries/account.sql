-- name: CreateAccount :one
INSERT INTO account (owner,
                     balance,
                     currency)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetAccount :one
SELECT *
FROM account
WHERE id = $1;

-- name: GetAccountForUpdate :one
SELECT *
FROM account
WHERE id = $1
    FOR NO KEY UPDATE;

-- name: GetAccounts :many
SELECT *
FROM account
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateAccount :one
UPDATE account
SET owner    = $1,
    balance  = $2,
    currency = $3
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: UpdateAccountBalance :one
UPDATE account
SET balance = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE
FROM account
WHERE id = $1;