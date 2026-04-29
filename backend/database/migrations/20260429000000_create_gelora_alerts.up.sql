PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS gelora_alerts (
    id INTEGER PRIMARY KEY,
    user_id INTEGER,
    message TEXT NOT NULL DEFAULT '',
    acknowledged_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (user_id)
        REFERENCES gelora_users(id)
        ON DELETE SET NULL
);

CREATE INDEX idx_gelora_alerts_user_id ON gelora_alerts(user_id);
CREATE INDEX idx_gelora_alerts_created_at ON gelora_alerts(created_at);
CREATE INDEX idx_gelora_alerts_acknowledged_at ON gelora_alerts(acknowledged_at);
