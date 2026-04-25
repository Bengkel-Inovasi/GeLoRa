PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS gelora_records (
    id INTEGER PRIMARY KEY,
    session_id INTEGER,
    time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    heart_rate REAL,
    body_temperature REAL,
    latitude REAL,
    longitude REAL,

    FOREIGN KEY (session_id)
        REFERENCES gelora_sessions(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_gelora_records_session_id ON gelora_records(session_id);

CREATE INDEX idx_gelora_records_time ON gelora_records(time);