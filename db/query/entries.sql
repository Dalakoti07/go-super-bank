-- name: CreateEntries :one
INSERT INTO entries(
    account_id,
    amount
) values (
             $1, $2
         ) RETURNING *;

-- name: GetEntries :one
SELECT * from entries
WHERE id = $1 LIMIT 1;

-- name: ListEntries :many
SELECT * from entries
ORDER BY id
    LIMIT $1
OFFSET $2;
