-- +goose Up
CREATE TABLE doctor_services (
     id UUID,
     doctor_id UUID,
     clinic_service_id UUID
);

-- +goose Down
DROP TABLE IF EXISTS doctor_services;