-- +goose Up
addresses (
  id UUID PRIMARY KEY,
  country VARCHAR,
  city VARCHAR,
  street VARCHAR,
  building VARCHAR,
  latitude DECIMAL,
  longitude DECIMAL
)

-- +goose Down
DROP TABLE IF EXISTS addresses;