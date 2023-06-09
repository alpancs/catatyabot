import { sendHelpMessage } from "./help";

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env);
    return new Response();
}

async function respondMessage(message: Message, env: Env) {
    const send = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text);
    const reply = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, message.message_id);
    if (message.text === "/start" || message.text === "/bantuan") return sendHelpMessage(send);
    // dummy
    if (message.text === "/semua") return respondListAll(reply, message.chat.id, env.DB);
    console.info(JSON.stringify({ status: "ignored", message }));
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

let problematicChars = ['=', '.', '-', '#', '(', ')'];
async function sendMessage(botToken: string, chatId: number, text: string, replyToMessageId?: number) {
    for (const problem of problematicChars) {
        text = text.replaceAll(problem, `\\${problem}`);
    }
    return sendCleanMessage(botToken, chatId, text, replyToMessageId);
}

async function sendCleanMessage(botToken: string, chatId: number, text: string, replyToMessageId?: number) {
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
    if (response.status >= 400) {
        const responseText = await response.text();
        console.error(responseText);
        if (response.status < 500) {
            const needToEscapePattern = /.*Character '(.)' is reserved and must be escaped.*/;
            if (needToEscapePattern.test(responseText)) {
                const problem = responseText.replace(needToEscapePattern, "$1");
                problematicChars.push(problem);
                return sendCleanMessage(botToken, chatId, text.replaceAll(problem, `\\${problem}`), replyToMessageId);
            }
        }
    }
    return response;
}
