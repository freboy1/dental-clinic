-- +goose Up
ALTER TABLE clinic_services
    ADD PRIMARY KEY (id);
-- +goose Down
ALTER TABLE clinic_services
DROP PRIMARY KEY;