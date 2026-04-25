CREATE TABLE IF NOT EXISTS gelora_users (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL,
    bio TEXT NOT NULL DEFAULT '',
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE INDEX idx_gelora_users_name_nocase ON gelora_users(name COLLATE NOCASE);

CREATE INDEX idx_gelora_users_role ON gelora_users(role);
