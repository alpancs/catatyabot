import { Env } from "./env.d";
import { Item } from "./item.d";
import { Update, Message } from "./telegram.d";

type SendTextFn = (text: string) => Promise<Response>;

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env);
    return new Response();
}

async function respondMessage(message: Message, env: Env) {
    if (message.text === "/semua") return respondListAll(message, env);
    console.info(JSON.stringify({ status: "ignored", message }));
}

async function respondListAll(message: Message, env: Env) {
    const { results } = await env.DB.prepare("SELECT * FROM items WHERE chat_id = ?").bind(message.chat.id).all<Item>();
    const send = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text);
    if (results?.length) return sendItemList(send, "*=== DAFTAR SEMUANYA ===*", results);
    return send("_catatan masih kosong_");
}

async function sendItemList(send: SendTextFn, title: string, items: Item[]) {
    const text = `${title}\n\n` + items.map(i => `[${i.created_at.slice(0, 16)}] ${i.name}: ${i.price}`).join("\n");
    return send(text);
}

async function sendMessage(botToken: string, chatId: number, text: string, replyToMessageId?: number) {
    text = text.replace(/(=|-|\.)/g, "\\$1");
    const response = await fetch(`https://api.telegram.org/bot${botToken}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: chatId,
            reply_to_message_id: replyToMessageId,
            text: text,
            parse_mode: "MarkdownV2",
        }),
    });
    if (response.status >= 400) console.error(await response.clone().text());
    return response;
}
