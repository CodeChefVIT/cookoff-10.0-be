-- name: CreateSubmission :exec
INSERT INTO submissions (
    id,
    question_id,
    language_id,
    source_code,
    testcases_passed,
    testcases_failed,
    runtime,
    memory,
    status,
    submission_time,
    description,
    user_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);



-- name: GetSubmissionByID :one
SELECT
    id,
    question_id,
    language_id,
    source_code,
    testcases_passed,
    testcases_failed,
    runtime,
    memory,
    status,
    submission_time,
    description,
    user_id
FROM submissions
WHERE id = $1;
