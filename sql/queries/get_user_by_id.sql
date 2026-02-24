-- name: GetUserById :one
SELECT * FROM users
WHERE ID = $1;
