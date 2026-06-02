-- +goose Up
ALTER TABLE clinics
    ADD COLUMN IF NOT EXISTS logo_url TEXT;

-- +goose Down
ALTER TABLE clinics
DROP COLUMN IF EXISTS logo_url;