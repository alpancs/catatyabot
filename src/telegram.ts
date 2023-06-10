import { askToCreateItems, createItemsQuestion, replyForItemsCreation } from "./create";
import { sendHelpMessage } from "./help";
import { sendMessage } from "./send";

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env)
    else console.info(JSON.stringify({ status: "ignored", reason: "the update does not contain a message", update }));
    return new Response();
}

async function respondMessage(message: Message, env: Env) {
    const send = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text);
    const reply = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, message.message_id);
    const ask = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, message.message_id, true);
    if (message.text === "/start" || message.text === "/bantuan") return sendHelpMessage(send);
    if (message.text === "/catat") return askToCreateItems(ask);
    if (message.reply_to_message?.text === createItemsQuestion && message.text) return replyForItemsCreation(reply, message.text);

    // dummy
    if (message.text === "/semua") return respondListAll(reply, message.chat.id, env.DB);
    console.info(JSON.stringify({ status: "ignored", reason: "the message does not match any cases", message }));
}

// dummy
async function respondListAll(reply: SendTextFn, chatId: number, db: D1Database) {
    const { results } = await db.prepare("SELECT * FROM items WHERE chat_id = ?").bind(chatId).all<Item>();
    if (results?.length) {
        const title = "*=== DAFTAR SEMUANYA ===*";
        const text = `${title}\n\n` + results.map(i => `[${i.created_at.slice(0, 16)}] ${i.name}: ${i.price}`).join("\n");
        return reply(text);
    }
    return reply("_catatan masih kosong_");
}
