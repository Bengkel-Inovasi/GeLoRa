PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS gelora_sessions (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    node_id INTEGER,

    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES gelora_users(id)
        ON DELETE SET NULL,

    FOREIGN KEY (node_id)
        REFERENCES gelora_nodes(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_gelora_sessions_user_id ON gelora_sessions(user_id);

CREATE INDEX idx_gelora_sessions_node_id ON gelora_sessions(node_id);

CREATE INDEX idx_gelora_sessions_started_at ON gelora_sessions(started_at);

CREATE INDEX idx_gelora_sessions_ended_at ON gelora_sessions(ended_at);