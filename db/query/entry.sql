-- name: CreateEntry :one
INSERT INTO entries (account_id,amount)
  VALUES(@account_id, @amount)
RETURNING *;

-- name: GetEntry :one
SELECT * FROM entries
WHERE id = @id LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
WHERE account_id = @account_id
ORDER BY id
LIMIT @size_limit
OFFSET @size_offset;

