-- name: GetChirp :one
SELECT * FROM chirps
WHERE ID = $1;

