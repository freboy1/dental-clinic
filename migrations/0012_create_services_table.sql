-- +goose Up
CREATE TABLE services (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          name VARCHAR(255) NOT NULL,
                          description TEXT,
                          price NUMERIC(10, 2) NOT NULL DEFAULT 0,
                          duration INT NOT NULL DEFAULT 30,
                          clinic_id UUID REFERENCES clinics(id) ON DELETE CASCADE,
                          is_active BOOLEAN NOT NULL DEFAULT TRUE,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS services;