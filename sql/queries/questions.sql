-- name: CreateQuestion :one
INSERT INTO questions (
  level, content
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetQuestionByID :one
SELECT * FROM questions
WHERE id = $1;

-- name: ListQuestions :many
SELECT * FROM questions
ORDER BY created_at DESC;

-- name: ListQuestionsByLevel :many
SELECT * FROM questions
WHERE level = $1
ORDER BY created_at DESC;

-- name: GetRandomQuestionByLevel :one
SELECT *
FROM questions
WHERE level = $1
ORDER BY RANDOM()
LIMIT 1;

-- name: UpdateQuestion :exec
UPDATE questions
SET content = $2,
    level = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteQuestion :exec
DELETE FROM questions
WHERE id = $1;
