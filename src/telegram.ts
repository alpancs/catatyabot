import { editMessage, sendMessage } from "./send";
import { createItemsQuestion, replyForItemsCreation, itemPattern } from "./create";
import { readItemsQuestion, replyForItemsReading } from "./read";
import { replyForItemUpdate } from "./update";
import { noItemToDelete, replyForItemDeletion } from "./delete";
import { helpMessage } from "./help";

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env)
    else console.info(JSON.stringify({ status: "ignored", reason: "the update does not contain a message", update }));
    return new Response();
}

async function respondMessage(message: Message, env: Env) {
    const send = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text);
    const ask = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, true);
    const edit = (messageId: number, text: string) => editMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, messageId, text);
    if (message.text === "/start" || message.text === "/bantuan") return send(helpMessage);
    if (message.text === "/catat") return ask(createItemsQuestion);
    if (message.text === "/lihat") return ask(readItemsQuestion);
    if (message.text === "/hapus" && !message.reply_to_message) return send(noItemToDelete);
    if (message.text === "/hapus" && message.reply_to_message)
        return replyForItemDeletion(send, edit, message.chat.id, message.reply_to_message.message_id, env.DB);
    if (message.reply_to_message?.text === createItemsQuestion && message.text)
        return replyForItemsCreation(send, edit, message.text, env.DB);
    if (message.reply_to_message?.text === readItemsQuestion && message.text)
        return replyForItemsReading(send, ask, message.chat.id, message.text, env.DB);
    const itemMatch = message.text?.match(itemPattern);
    if (message.reply_to_message && itemMatch)
        return replyForItemUpdate(send, edit, message.chat.id, message.reply_to_message.message_id, itemMatch, env.DB);
    console.info(JSON.stringify({ status: "ignored", reason: "the message does not match any cases", message }));
}
