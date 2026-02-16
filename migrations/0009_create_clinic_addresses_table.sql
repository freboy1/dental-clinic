-- +goose Up
clinic_addresses (
  id UUID PRIMARY KEY,
  clinic_id UUID REFERENCES clinics(id),
  address_id UUID REFERENCES addresses(id),
  is_main BOOLEAN
)


-- +goose Down
DROP TABLE IF EXISTS clinic_addresses;