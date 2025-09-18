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

-- name: UpdateSubmission :exec
UPDATE submissions
SET 
    runtime = $1, 
    memory = $2, 
    status = $3,
    testcases_passed = $4,
    testcases_failed = $5
WHERE id = $6;

-- name: UpdateScore :exec
WITH best_submissions AS (
    SELECT 
        s.user_id AS user_id,
        s.question_id,
        MAX((s.testcases_passed) * q.points / (s.testcases_passed + s.testcases_failed)::numeric) AS best_score
    FROM submissions s
    INNER JOIN questions q ON s.question_id = q.id
    INNER JOIN users u on s.user_id = u.id 
    WHERE s.user_id = (select user_id from submissions where id = $1) AND q.round = u.round_qualified
    GROUP BY s.user_id, s.question_id
)
UPDATE users
SET score = (
    SELECT SUM(best_score)
    FROM best_submissions
)
WHERE users.id = (select user_id from submissions s where s.id = $1);