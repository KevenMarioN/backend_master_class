-- name: CreateAccount :one
INSERT INTO accounts (owner, balance, currency)
	VALUES (@owner,@balance,@currency)
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = @id LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = @id LIMIT 1 FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = @owner
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = @balance
WHERE id = @id
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = balance + @amount
WHERE id = @id
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = @id;
