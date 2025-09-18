-- name: CreateQuestion :exec
INSERT INTO questions (
    id,
    description, 
    title,
    qType,
    isBountyActive,
    input_format,
    points,
    round,
    constraints,
    output_format,
    sample_test_input,
    sample_test_output,
    explanation
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
);

-- name: GetQuestion :one
SELECT * FROM questions
WHERE id = $1;

-- name: GetAllQuestions :many
SELECT * FROM questions ORDER BY id;

-- name: UpdateQuestion :exec
UPDATE questions
SET description = $2,
    qType = $3,
    isBountyActive = $4,
    input_format = $5,
    points = $6,
    round = $7,
    constraints = $8,
    output_format = $9,
    sample_test_input = $10,
    sample_test_output = $11,
    explanation = $12
WHERE id = $1;

-- name: UpdateQuestionBountyActive :exec
UPDATE questions
SET isBountyActive = true
WHERE id = $1;

-- name: UpdateQuestionBountyInactive :exec
UPDATE questions
SET isBountyActive = false
WHERE id = $1;

-- name: DeleteQuestion :exec
DELETE FROM questions WHERE id = $1;

-- name: GetQuestionsByRound :many
SELECT * FROM questions
WHERE round = $1
ORDER BY qType;