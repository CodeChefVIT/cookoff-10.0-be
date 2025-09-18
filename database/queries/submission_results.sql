-- name: CreateSubmissionResult :exec
INSERT INTO submission_results (id, testcase_id, submission_id, runtime, memory, points_awarded, status, description)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetStatsForFinalSubEntry :many
SELECT 
    runtime, 
    memory,   
    status
FROM submission_results
WHERE submission_id = $1;