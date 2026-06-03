-- +goose Up
ALTER TABLE services ADD COLUMN IF NOT EXISTS name_en VARCHAR(255);
ALTER TABLE services ADD COLUMN IF NOT EXISTS name_kaz VARCHAR(255);

-- Backfill: existing name goes to name_en by default
UPDATE services SET name_en = name WHERE name_en IS NULL;

-- +goose Down
ALTER TABLE services DROP COLUMN IF EXISTS name_en;
ALTER TABLE services DROP COLUMN IF EXISTS name_kaz;