-- name: CreateUser :one
INSERT INTO users (id, email, reg_no, password, role, round_qualified, score, name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;

-- name: GetUserByUsername :one
SELECT *
FROM users
WHERE name = $1;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT *
FROM users;

-- name: UpgradeUserToRound :exec
UPDATE users
SET round_qualified = round_qualified + 1
WHERE id = $1;

-- name: BanUser :exec
UPDATE users SET is_banned = TRUE
WHERE id = $1;

-- name: UnbanUser :exec
UPDATE users
SET is_banned = FALSE
WHERE id = $1;

-- name: GetLeaderboard :many
select id, name, score from users
order by score;

-- name: UpdateProfile :exec
UPDATE users SET reg_no = $1, password = $2, name = $3
WHERE id = $4;

-- name: GetSubmissionByUser :many
SELECT id, question_id, testcases_passed, testcases_failed, runtime,
       submission_time, source_code, language_id, description, memory,
       user_id, status
FROM submissions
WHERE user_id = $1;

-- name: GetUsersWithCursor :many
SELECT id, email, reg_no, password, role, round_qualified, score, name, is_banned
FROM users
WHERE ($1::uuid IS NULL OR id > $1)
ORDER BY id ASC
LIMIT $2;

-- name: GetUserRound :one
SELECT round_qualified
FROM users
WHERE id = $1;