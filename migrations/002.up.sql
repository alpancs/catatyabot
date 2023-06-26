CREATE TABLE IF NOT EXISTS hashtags (
    chat_id INTEGER NOT NULL,
    message_id INTEGER NOT NULL,
    hashtag TEXT NOT NULL,
    FOREIGN KEY (chat_id, message_id) REFERENCES items (chat_id, message_id)
        ON UPDATE CASCADE ON DELETE CASCADE
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_hashtags_1
    ON hashtags (chat_id, message_id, lower(hashtag));
