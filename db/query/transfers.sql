-- name: CreateTransfer :one
INSERT INTO transfers(
                      from_account_id,
                      to_account_id,
                      amount
) values (
             $1, $2, $3
         ) RETURNING *;

-- name: GetTransfer :one
SELECT * from transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfers :many
SELECT * from transfers
ORDER BY id
    LIMIT $1
OFFSET $2;

