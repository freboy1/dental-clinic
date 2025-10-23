CREATE TABLE verification_tokens (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '24 hours'
);
