-- name: GetUserByEmail :one
SELECT * FROM users
WHERE EMAIL = $1;
