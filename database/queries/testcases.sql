-- name: GetTestcases :many
SELECT
    id,
    expected_output,
    memory,
    input,
    hidden,
    runtime,
    question_id
FROM testcases
WHERE question_id = $1;