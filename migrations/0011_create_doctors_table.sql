-- +goose Up
CREATE TABLE doctors (
                         id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         specialization VARCHAR(255) NOT NULL,
                         name VARCHAR(255) NOT NULL,
                         email VARCHAR(255) NOT NULL,
                         experience INT NOT NULL DEFAULT 0,
                         clinic_id UUID REFERENCES clinics(id) ON DELETE SET NULL,
                         bio TEXT,
                         is_available BOOLEAN NOT NULL DEFAULT TRUE,
                         created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS doctors;