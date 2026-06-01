-- +goose Up
CREATE TABLE clinic_admins (
   id UUID PRIMARY KEY,
   clinic_id UUID REFERENCES clinics(id),
   user_id UUID REFERENCES users(id),
   created_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS clinic_admins;
