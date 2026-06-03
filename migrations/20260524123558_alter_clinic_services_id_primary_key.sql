-- +goose Up
UPDATE clinic_services SET id = gen_random_uuid() WHERE id IS NULL;
ALTER TABLE clinic_services ALTER COLUMN id SET NOT NULL;
ALTER TABLE clinic_services ADD PRIMARY KEY (id);

-- +goose Down
ALTER TABLE clinic_services DROP CONSTRAINT IF EXISTS clinic_services_pkey;