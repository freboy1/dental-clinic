-- +goose Up
ALTER TABLE medical_files
    ADD COLUMN IF NOT EXISTS mime_type TEXT;

-- +goose Down
ALTER TABLE medical_files
DROP COLUMN IF EXISTS mime_type;