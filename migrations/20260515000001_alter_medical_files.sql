-- +goose Up
ALTER TABLE medical_files
    ADD COLUMN IF NOT EXISTS medical_record_id UUID REFERENCES medical_records(id) ON DELETE CASCADE,
    ADD COLUMN IF NOT EXISTS file_name TEXT,
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT NOW();

-- +goose Down
ALTER TABLE medical_files
DROP COLUMN IF EXISTS medical_record_id,
    DROP COLUMN IF EXISTS file_name,
    DROP COLUMN IF EXISTS created_at;