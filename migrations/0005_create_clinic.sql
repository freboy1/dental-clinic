CREATE TABLE clinics (
    id UUID PRIMARY KEY,
    name VARCHAR,
    description TEXT,
    phone VARCHAR,
    email VARCHAR,
    website VARCHAR,
    rating FLOAT DEFAULT 0,
    reviews_count INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP
);
