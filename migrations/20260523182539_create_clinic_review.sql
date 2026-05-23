-- +goose Up
CREATE TABLE clinic_reviews (
    id UUID PRIMARY KEY,

    appointment_id UUID REFERENCES appointments(id),

    clinic_id UUID REFERENCES clinics(id),
    user_id UUID REFERENCES users(id),

    rating INT CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,

    created_at TIMESTAMP DEFAULT NOW()
);
-- +goose Down
DROP TABLE IF EXISTS clinic_reviews;
