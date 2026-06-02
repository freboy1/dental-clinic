-- +goose Up
CREATE TABLE clinic_address_gallery (
    id UUID PRIMARY KEY,

    clinic_address_id UUID REFERENCES clinic_addresses(id),

    image_url TEXT NOT NULL,

    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS clinic_address_gallery;
