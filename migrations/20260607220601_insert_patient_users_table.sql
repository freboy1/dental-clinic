-- +goose Up
INSERT INTO users (id, email, password, name, surname, role, gender, age, push_consent, is_verified, created_at)
VALUES ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'patient@gmail.com', '$2a$10$zHyHY4RI07o5SNhDNPrRCOSLgE3cGBUA8/WincB2EmIInkPodu8x.', 'Alish', '', 'patient', 'Male', 34, false, true, '2026-01-08 01:06:24.765831');


-- +goose Down
DELETE FROM users
WHERE id = 'ffffffff-ffff-ffff-ffff-ffffffffffff';