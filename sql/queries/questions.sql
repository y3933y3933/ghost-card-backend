-- name: GetRandomQuestionByLevel :one
SELECT id, content FROM questions
WHERE level = $1
ORDER BY RANDOM()
LIMIT 1;


-- name: GetQuestionByID :one
SELECT content FROM questions
WHERE id = $1;