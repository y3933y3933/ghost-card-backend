-- name: CreateRound :one
INSERT INTO rounds (
  game_id, question_id, current_player_id, is_joker, status
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetRoundByID :one
SELECT * FROM rounds
WHERE id = $1;

-- name: ListRoundsByGame :many
SELECT * FROM rounds
WHERE game_id = $1
ORDER BY created_at ASC;

-- name: RevealRound :exec
UPDATE rounds
SET is_joker = $2,
    status = 'revealed'
WHERE id = $1;

-- name: DeleteRoundsByGame :exec
DELETE FROM rounds
WHERE game_id = $1;


-- name: GetLatestRoundByGameID :one
SELECT *
FROM rounds
WHERE game_id = $1
ORDER BY created_at DESC
LIMIT 1;