-- +goose Up
CREATE TABLE doctor_time_slots (
  id UUID PRIMARY KEY,

  doctor_id UUID REFERENCES doctors(id),
  clinic_address_id UUID REFERENCES clinic_addresses(id),

  slot_start TIMESTAMP,
  slot_end TIMESTAMP,

  status VARCHAR,

  created_at TIMESTAMP
);


-- +goose Down
DROP TABLE IF EXISTS doctor_time_slots;
