CREATE TABLE items (
    chat_id bigint,
    message_id integer,
    name text NOT NULL,
    price bigint NOT NULL,
    created_at timestamp DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta'),
    PRIMARY KEY (chat_id, message_id)
);
