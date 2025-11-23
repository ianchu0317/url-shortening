CREATE TABLE IF NOT EXISTS shortened (
    id SERIAL PRIMARY KEY,                    -- auto-increment
    url TEXT NOT NULL,                        -- original url
    shortened_url VARCHAR(10) UNIQUE NOT NULL, -- shortened url (ej: "abc123")
    created_at TIMESTAMP DEFAULT NOW(),       -- auto timestamp
    updated_at TIMESTAMP DEFAULT NOW(),
    accessed INT DEFAULT 0                    -- access counter
);
