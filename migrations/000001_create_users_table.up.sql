
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    phone TEXT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'client',  -- client / doctor / admin
    push_consent BOOLEAN DEFAULT false,
    activated BOOLEAN DEFAULT false,
    activation_token TEXT,
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT now()
    );


CREATE TABLE IF NOT EXISTS login_history (
                                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    ip_address TEXT,
    success BOOLEAN,
    user_agent TEXT,
    attempt_time TIMESTAMP WITH TIME ZONE DEFAULT now()
    );

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
