import { Env, Item } from "./common.d";
import { Update, Message } from "./telegram.d";

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env);
    return new Response(null);
}

async function respondMessage(message: Message, env: Env) {
    if (message.text === "/semua") return respondListAll(message, env);
    console.info(JSON.stringify({ status: "ignored", message }));
}

async function respondListAll(message: Message, env: Env) {
    const { results } = await env.DB.prepare("SELECT * FROM items WHERE chat_id = ?").bind(message.chat.id).all<Item>();
    if (results) return sendItemList(env.TELEGRAM_BOT_TOKEN, message.chat.id, "*=== DAFTAR SEMUANYA ===*", results);
    return sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, "_catatan masih kosong_");
}

async function sendMessage(botToken: string, chatId: number, text: string) {
    text = text.replace(/(=|-|\.)/g, "\\$1");
    const response = await fetch(`https://api.telegram.org/bot${botToken}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ chat_id: chatId, text, parse_mode: "MarkdownV2" }),
    });
    if (response.status >= 400) {
        console.error(await response.text());
    }
}

async function sendItemList(botToken: string, chatId: number, title: string, items: Item[]) {
    const text = `${title}\n\n` + items.map(i => `[${i.created_at.slice(0, 16)}] ${i.name}: ${i.price}`).join("\n");
    return sendMessage(botToken, chatId, text);
}
