-- name: CreateTransfer :one
INSERT INTO tranfers (to_account_id, from_account_id, amount)
  VALUES (@to_account_id, @from_account_id, @amount)
RETURNING *;

-- name: GetTransfer :one
SELECT * FROM tranfers
WHERE id = @id LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM tranfers
WHERE
  from_account_id = @from_account_id OR
  to_account_id = @to_account_id
ORDER BY id
LIMIT @size_limit
OFFSET @size_offset;
