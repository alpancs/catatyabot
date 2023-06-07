import { Env } from "./env";

export interface Update {
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

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env);
    return new Response(null);
}

async function respondMessage(message: Message, env: Env) {
    await fetch(`https://api.telegram.org/bot${env.TELEGRAM_BOT_TOKEN}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: message.chat.id,
            text: `your message was "${message.text}", right?`,
        }),
    });
}
