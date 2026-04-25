-- +goose Up
ALTER TABLE doctors
    ADD COLUMN user_id UUID,
ADD CONSTRAINT fk_user
    FOREIGN KEY (user_id) REFERENCES users(id);

-- +goose Down
ALTER TABLE doctors
DROP CONSTRAINT fk_user,
    DROP COLUMN user_id;