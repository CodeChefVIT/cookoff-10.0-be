-- +goose Up
ALTER TABLE submissions ADD CONSTRAINT fk_submissions_question FOREIGN KEY(question_id) REFERENCES questions(id) ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE testcases ADD CONSTRAINT fk_testcases FOREIGN KEY(question_id) REFERENCES questions(id) ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE submissions ADD CONSTRAINT fk_submissions_user FOREIGN KEY(user_id) REFERENCES users(id) ON UPDATE NO ACTION ON DELETE CASCADE;

ALTER TABLE submission_results ADD CONSTRAINT fk_submission_results FOREIGN KEY(submission_id) REFERENCES submissions(id) ON UPDATE NO ACTION ON DELETE CASCADE;

-- +goose Down
ALTER TABLE submissions DROP CONSTRAINT fk_submissions_question;
ALTER TABLE testcases DROP CONSTRAINT fk_testcases;
ALTER TABLE submissions DROP CONSTRAINT fk_submissions_user;
ALTER TABLE submission_results DROP CONSTRAINT fk_submission_results;
