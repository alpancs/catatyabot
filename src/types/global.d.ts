export { };

declare global {
    interface Env {
        DB: D1Database;
        TELEGRAM_BOT_TOKEN: string;
    }

    type SendTextFn = (text: string) => Promise<Response>;

    interface Item {
        chat_id: number;
        message_id: number;
        name: string;
        price: number;
        created_at: string;
    }

    interface Update {
        update_id: number;
        message?: Message;
    }

    interface Message {
        message_id: number;
        chat: Chat;
        reply_to_message?: Message;
        text?: string;
    }

    interface Chat {
        id: number;
    }
}
