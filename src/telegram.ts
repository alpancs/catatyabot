import { Env } from "./env.d";
import { Update, Message } from "./telegram.d";

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env);
    return new Response(null);
}

async function respondMessage(message: Message, env: Env) {
    if (message.text === "/semua") {
        return respondListAll(message, env);
    }
    console.debug(JSON.stringify({ status: "ignored", message }));
}

async function respondListAll(message: Message, env: Env) {
    const { results } = await env.DB.prepare("SELECT * FROM items WHERE chat_id = ?").bind(message.chat.id).all();
    return sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, "```json\n" + JSON.stringify(results) + "\n```");
}

async function sendMessage(botToken: string, chatId: number, text: string) {
    return fetch(`https://api.telegram.org/bot${botToken}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ chat_id: chatId, text: text, parse_mode: "MarkdownV2" }),
    });
}
