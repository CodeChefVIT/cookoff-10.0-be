-- +goose Up
ALTER TABLE submission_results DROP CONSTRAINT fk_submission_results;

-- +goose Down
ALTER TABLE submission_results ADD CONSTRAINT fk_submission_results FOREIGN KEY(submission_id) REFERENCES submissions(id) ON UPDATE NO ACTION ON DELETE CASCADE;
