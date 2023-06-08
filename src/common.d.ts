export interface Env {
    DB: D1Database;
    TELEGRAM_BOT_TOKEN: string;
}

export interface Item {
    chat_id: number;
    message_id: number;
    name: string;
    price: number;
    created_at: string;
}
