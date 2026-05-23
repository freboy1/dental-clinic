-- +goose Up
CREATE TABLE ai_booking_state (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  doctor_id UUID,
  service_id UUID,
  clinic_address_id UUID,
  date DATE,
  time TIME,
  step TEXT,
  updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS ai_booking_state;
