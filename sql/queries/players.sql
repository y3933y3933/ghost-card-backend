-- name: CreatePlayer :one
INSERT INTO players (
  game_id, nickname, is_host
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetPlayerByID :one
SELECT * FROM players
WHERE id = $1;

-- name: ListPlayersByGame :many
SELECT * FROM players
WHERE game_id = $1
ORDER BY joined_at ASC;

-- name: DeletePlayerByID :exec
DELETE FROM players
WHERE id = $1;

-- name: DeletePlayersByGame :exec
DELETE FROM players
WHERE game_id = $1;

-- name: CountPlayersInGame :one 
SELECT COUNT(*) FROM players 
WHERE game_id = $1;