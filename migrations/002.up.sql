CREATE TABLE IF NOT EXISTS hashtags (
    chat_id INTEGER NOT NULL,
    message_id INTEGER NOT NULL,
    hashtag TEXT NOT NULL,
    PRIMARY KEY (chat_id, message_id, hashtag),
    FOREIGN KEY (chat_id, message_id) REFERENCES items(chat_id, message_id)
        ON UPDATE CASCADE ON DELETE CASCADE
);
