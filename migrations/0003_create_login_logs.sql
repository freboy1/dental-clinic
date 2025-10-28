-- +goose Up

CREATE TABLE login_logs (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    ip_address VARCHAR(45) NOT NULL,
    login_time TIMESTAMP DEFAULT NOW(),
    success BOOLEAN DEFAULT TRUE
);

-- +goose Down
DROP TABLE IF EXISTS login_logs;