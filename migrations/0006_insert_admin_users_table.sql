-- +goose Up
INSERT INTO users (id, email, password, name, surname, role, gender, age, push_consent, is_verified, created_at)
VALUES ('00000000-0000-0000-0000-000000000000', 'alisherdautov22@gmail.com', '$2a$10$4AyQ9e1B/L5THx.bdOmGFOeOiU9vAseAxWs/xPBBQ0rdwmbwjJw1S', 'Alish', '', 'admin', 'Male', 34, false, true, '2026-01-08 01:06:24.765831');


-- +goose Down
DELETE FROM users
WHERE id = '00000000-0000-0000-0000-000000000000';