CREATE TABLE IF NOT EXISTS tags (
    chat_id INTEGER NOT NULL,
    message_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    normalized_name TEXT NOT NULL,
    PRIMARY KEY (chat_id, message_id, normalized_name),
    FOREIGN KEY (chat_id, message_id) REFERENCES items(chat_id, message_id)
);
