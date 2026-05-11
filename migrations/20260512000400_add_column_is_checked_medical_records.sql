-- +goose Up
ALTER TABLE medical_records
    ADD COLUMN is_checked boolean;

-- +goose Down
ALTER TABLE medical_records
    DROP COLUMN is_checked;