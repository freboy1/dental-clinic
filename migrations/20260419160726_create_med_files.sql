-- +goose Up
CREATE TABLE medical_files (
  id UUID,
  file_url TEXT
);

-- +goose Down
DROP TABLE IF EXISTS medical_files;