-- name: CreateGame :one
INSERT INTO games (code, level, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetGameByCode :one
SELECT * FROM games
WHERE code = $1;

-- -- name: GetGameByID :one
-- SELECT * FROM games
-- WHERE id = $1;

-- -- name: GetGameByCode :one
-- SELECT * FROM games
-- WHERE code = $1;

-- -- name: ListGames :many
-- SELECT * FROM games
-- ORDER BY created_at DESC;

-- -- name: UpdateGameStatus :exec
-- UPDATE games
-- SET status = $2,
--     updated_at = NOW()
-- WHERE id = $1;

-- -- name: DeleteGame :exec
-- DELETE FROM games
-- WHERE id = $1;
