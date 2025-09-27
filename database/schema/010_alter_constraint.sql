-- +goose Up
ALTER TABLE testcases ALTER COLUMN input DROP NOT NULL;

-- +goose Down
ALTER TABLE testcases ALTER COLUMN input SET NOT NULL;
