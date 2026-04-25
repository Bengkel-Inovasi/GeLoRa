PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS gelora_nodes (
    id INTEGER PRIMARY KEY,
    mid TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL DEFAULT 'Unnamed Node',
    description TEXT NOT NULL DEFAULT '',
    is_validated INTEGER NOT NULL DEFAULT 0,
    validated_by INTEGER,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,

    FOREIGN KEY (validated_by)
        REFERENCES gelora_users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_gelora_nodes_mid_nocase ON gelora_nodes(mid COLLATE NOCASE);

CREATE INDEX idx_gelora_nodes_name_nocase ON gelora_nodes(name COLLATE NOCASE);

CREATE INDEX idx_gelora_nodes_is_validated ON gelora_nodes(is_validated);

