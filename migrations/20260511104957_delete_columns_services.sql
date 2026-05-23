-- +goose Up
ALTER TABLE services
DROP COLUMN price,
DROP COLUMN duration,
DROP COLUMN clinic_id,
DROP COLUMN is_active,
DROP COLUMN created_at;
-- ADD CONSTRAINT fk_user
--     FOREIGN KEY (user_id) REFERENCES users(id);

-- +goose Down
-- ALTER TABLE doctors
-- DROP CONSTRAINT fk_user,
--     DROP COLUMN user_id;