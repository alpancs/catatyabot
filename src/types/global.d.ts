export { };

declare global {
    interface Env {
        DB: D1Database;
        TELEGRAM_BOT_TOKEN: string;
        TELEGRAM_WEBHOOK_SECRET_TOKEN: string;
    }

    type SendTextFn = (text: string) => Promise<Response>;
    type EditTextFn = (messageId: number, text: string) => Promise<Response>;

    interface Item {
        chat_id: number;
        message_id: number;
        name: string;
        price: number;
        created_at: string;
    }

    interface Update {
        message?: Message;
    }

    interface Message {
        message_id: number;
        chat: Chat;
        reply_to_message?: Message;
        text?: string;
        migrate_from_chat_id?: number;
    }

    interface Chat {
        id: number;
    }
}
