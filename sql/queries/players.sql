-- name: ListPlayersByGameCode :many
SELECT p.id, p.nickname, p.is_host, p.joined_at
FROM players p
JOIN games g ON p.game_id = g.id
WHERE g.code = $1
ORDER BY p.joined_at;


-- name: CountPlayersInGame :one
SELECT COUNT(*) FROM players WHERE game_id = $1;

-- name: CreatePlayer :one
INSERT INTO players(game_id, nickname, is_host)
VALUES ($1, $2, $3)
RETURNING id, nickname, is_host, joined_at;