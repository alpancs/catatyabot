CREATE TABLE IF NOT EXISTS items (
    chat_id INTEGER,
    message_id INTEGER,
    name TEXT,
    price REAL,
    created_at INTEGER,
    updated_at INTEGER,
    PRIMARY KEY (chat_id, message_id)
);
