-- +goose Up
CREATE TABLE inventory_transactions (
    id UUID PRIMARY KEY,

    clinic_address_id UUID REFERENCES clinic_addresses(id),

    product_id UUID REFERENCES products(id),

    quantity NUMERIC,

    transaction_type TEXT,

    appointment_id UUID NULL,

    created_at TIMESTAMP
);
-- +goose Down
DROP TABLE IF EXISTS inventory_transactions;
