-- +goose Up
CREATE TABLE products (
   id UUID PRIMARY KEY,
   name TEXT,
   unit TEXT,
   created_at TIMESTAMP
);
-- +goose Down
DROP TABLE IF EXISTS products;
