-- +goose Up
CREATE TABLE medical_records (
  id UUID PRIMARY KEY,
  appointment_id UUID ,
  doctor_id UUID,
  patient_id UUID,

  diagnosis TEXT,
  notes TEXT,
  is_checked BOOLEAN,

  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS medical_records;