-- +goose Up
CREATE TABLE service_materials (
                                   id UUID PRIMARY KEY,

                                   service_id UUID REFERENCES clinic_services(id),

                                   product_id UUID REFERENCES products(id),

                                   quantity_required NUMERIC
);

-- +goose Down
DROP TABLE IF EXISTS service_materials;
