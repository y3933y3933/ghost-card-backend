-- name: CreateGame :one
INSERT INTO games (code, level, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGameByCode :one
SELECT * FROM games
WHERE code = $1;


-- name: UpdateGameStatus :exec
UPDATE games SET status = $2 WHERE id = $1;

