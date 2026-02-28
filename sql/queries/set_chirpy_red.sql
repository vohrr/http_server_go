-- name: SetChirpyRed :exec
UPDATE users
SET is_chirpy_red = $1
WHERE ID = $2;
