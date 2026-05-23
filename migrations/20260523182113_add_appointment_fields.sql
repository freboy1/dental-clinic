-- +goose Up
ALTER TABLE appointments
    ADD COLUMN IF NOT EXISTS review_email_sent BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS is_reviewed  BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE appointments
DROP COLUMN IF EXISTS review_email_sent, is_reviewed;