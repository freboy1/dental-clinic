-- +goose Up
ALTER TABLE doctors
ADD COLUMN is_deleted Int NOT NULL DEFAULT 0;


-- +goose Down
ALTER TABLE doctors
DROP COLUMN is_deleted;
