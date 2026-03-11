-- +goose Up
CREATE TABLE doctor_working_hours (
  id UUID PRIMARY KEY,
  doctor_id UUID REFERENCES doctors(id),
  clinic_address_id UUID REFERENCES clinic_addresses(id),

  day_of_week INT, -- 1-7
  start_time TIME,
  end_time TIME,
);


-- +goose Down
DROP TABLE IF EXISTS doctor_working_hours;
