-- name: GetCurrentRoundByGameCode :one
SELECT r.id, r.current_player_id, r.question_id, r.is_joker, r.created_at, g.id AS game_id,
       r.status, g.level,
       q.content AS question_content
FROM rounds r
JOIN games g ON r.game_id = g.id
JOIN questions q ON r.question_id = q.id
WHERE g.code = $1
ORDER BY r.created_at DESC
LIMIT 1;




-- name: CreateRound :one
INSERT INTO rounds (game_id, question_id, current_player_id, status)
VALUES ($1, $2, $3, 'pending')
RETURNING id, question_id, current_player_id, status, created_at;