-- +goose Up
ALTER TABLE submissions
ALTER COLUMN user_id SET NOT NULL;

-- +goose Down
ALTER TABLE submissions
ALTER COLUMN user_id DROP NOT NULL;
