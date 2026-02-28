-- +goose Up
INSERT INTO addresses (id, country, city, street, building, latitude, longitude) VALUES
  ('11111111-1111-1111-1111-111111111111', 'USA', 'New York', '5th Avenue', '711', 40.775036, -73.965088),
  ('22222222-2222-2222-2222-222222222222', 'France', 'Paris', 'Champs-Élysées', '50', 48.869796, 2.307266),
  ('33333333-3333-3333-3333-333333333333', 'Japan', 'Tokyo', 'Chuo Dori', '3-2-1', 35.671780, 139.765052);

-- +goose Down
DELETE FROM addresses
WHERE id IN (
  '11111111-1111-1111-1111-111111111111',
  '22222222-2222-2222-2222-222222222222',
  '33333333-3333-3333-3333-333333333333'
);