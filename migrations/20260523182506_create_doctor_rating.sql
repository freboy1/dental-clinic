-- +goose Up
CREATE TABLE doctor_ratings (
    id UUID PRIMARY KEY,

    appointment_id UUID REFERENCES appointments(id),

    doctor_id UUID REFERENCES doctors(id),
    user_id UUID REFERENCES users(id),

    rating INT CHECK (rating >= 1 AND rating <= 5),

    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS doctor_ratings;
