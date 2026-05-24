-- +goose Up
CREATE TABLE CREATE TABLE address_inventory (
    id UUID PRIMARY KEY,

    clinic_address_id UUID REFERENCES clinic_addresses(id),

    product_id UUID REFERENCES products(id),

    quantity NUMERIC,

    updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS address_inventory;
