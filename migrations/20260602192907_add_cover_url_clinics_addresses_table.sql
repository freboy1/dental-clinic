-- +goose Up
ALTER TABLE clinic_addresses
    ADD COLUMN IF NOT EXISTS cover_image_url  TEXT;

-- +goose Down
ALTER TABLE clinic_addresses
DROP COLUMN IF EXISTS cover_image_url ;