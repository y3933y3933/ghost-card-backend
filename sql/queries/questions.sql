-- name: GetQuestionByID :one
SELECT * FROM questions
WHERE id = $1;

-- name: GetUnusedQuestion :one
SELECT * FROM questions
WHERE level = $1
  AND id NOT IN (
    SELECT question_id FROM rounds WHERE game_id = $2
  )
ORDER BY RANDOM()
LIMIT 1;

-- name: CreateQuestion :one
INSERT INTO questions (
  level, content
) VALUES (
  $1, $2
)
RETURNING *;

-- name: ListQuestionsByLevel :many
SELECT * FROM questions
WHERE level = $1
ORDER BY id;

