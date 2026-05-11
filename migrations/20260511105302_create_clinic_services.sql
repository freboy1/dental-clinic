-- +goose Up
CREATE TABLE clinic_services (
   id UUID,
   clinic_id UUID,
   service_id UUID,
   price DECIMAL,
   duration_minutes INT,
   is_active BOOLEAN
);

-- +goose Down
DROP TABLE IF EXISTS clinic_services;