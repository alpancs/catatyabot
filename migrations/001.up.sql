CREATE TABLE IF NOT EXISTS items (
    chat_id INTEGER NOT NULL,
    message_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    price REAL NOT NULL,
    created_at TEXT NOT NULL,
    PRIMARY KEY (chat_id, message_id)
);
