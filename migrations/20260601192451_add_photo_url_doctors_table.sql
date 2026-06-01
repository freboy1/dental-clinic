-- +goose Up
ALTER TABLE doctors
    ADD COLUMN photo_url TEXT;

-- +goose Down
ALTER TABLE doctors
DROP COLUMN IF EXISTS photo_url;