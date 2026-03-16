-- +goose Up
CREATE TABLE appointments (
  id UUID PRIMARY KEY,
  doctor_id UUID REFERENCES doctors(id),
  clinic_address_id UUID REFERENCES clinic_addresses(id),
  service_id UUID REFERENCES services(id),
  user_id UUID REFERENCES users(id),

  start_time TIMESTAMP,
  end_time TIMESTAMP,

  status VARCHAR,
  created_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS appointments;