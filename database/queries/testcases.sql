-- name: CreateTestCase :one
INSERT INTO testcases (
    id,
    expected_output,
    memory,
    input,
    hidden,
    runtime,
    question_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetTestCase :one
SELECT
    id,
    expected_output,
    memory,
    input,
    hidden,
    runtime,
    question_id
FROM testcases
WHERE id = $1;

-- name: GetTestCasesByQuestion :many
SELECT
    id,
    memory,
    input,
    hidden,
    runtime,
    question_id
FROM testcases
WHERE question_id = $1;

-- name: GetPublicTestCasesByQuestion :many
SELECT
    id,
    memory,
    input,
    hidden,
    runtime,
    question_id
FROM testcases
WHERE question_id = $1
AND hidden = false;

-- name: GetAllTestCases :many
SELECT
    id,
    expected_output,
    memory,
    input,
    hidden,
    runtime,
    question_id
FROM testcases;

-- name: UpdateTestCase :one
UPDATE testcases
SET
    expected_output = $2,
    memory = $3,
    input = $4,
    hidden = $5,
    runtime = $6,
    question_id = $7
WHERE id = $1
RETURNING *;

-- name: DeleteTestCase :exec
DELETE FROM testcases
WHERE id = $1;

-- name: GetAllTestCasesByQuestion :many
SELECT
    id,
    memory,
    expected_output,
    input,
    hidden,
    runtime,
    question_id
FROM testcases
WHERE question_id = $1;
