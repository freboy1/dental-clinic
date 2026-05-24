-- +goose Up
ALTER TABLE clinic_services
    ADD PRIMARY KEY (id);
-- +goose Down
ALTER TABLE clinic_services
DROP CONSTRAINT IF EXISTS clinic_services_pkey;
