import { editMessage, sendMessage } from "./send";
import { createItemsQuestion, replyForItemsCreation, itemPattern } from "./create";
import { readItemsQuestion, replyForItemsReading } from "./read";
import { replyForItemUpdate } from "./update";
import { noItemToDelete, replyForItemDeletion } from "./delete";
import { helpMessage } from "./help";

let ignoredMessageCounts: { [key: number]: number } = {};

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env)
    else console.info({ status: "ignored", reason: "the update does not contain a message", update });
    return new Response();
}

async function respondMessage(message: Message, env: Env) {
    const send = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text);
    const ask = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, true);
    const edit = (messageId: number, text: string) => editMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, messageId, text);
    if (message.text?.match(/^\s*\/?(start|bantuan)(@catatyabot)?\s*$/i)) return send(helpMessage);
    if (message.text?.match(/^\s*\/?catat(@catatyabot)?\s*$/i)) return ask(createItemsQuestion);
    if (message.text?.match(/^\s*\/?lihat(@catatyabot)?\s*$/i)) return ask(readItemsQuestion);
    if (message.text?.includes("hapus") && !message.reply_to_message) return send(noItemToDelete);
    if (message.text?.includes("hapus") && message.reply_to_message)
        return replyForItemDeletion(send, edit, message.chat.id, message.reply_to_message, env.DB);
    if (message.reply_to_message?.text === createItemsQuestion && message.text)
        return replyForItemsCreation(send, edit, message.text, env.DB);
    if (message.reply_to_message?.text === readItemsQuestion && message.text)
        return replyForItemsReading(send, ask, message.chat.id, message.text, env.DB);
    const itemMatch = message.text?.match(itemPattern);
    if (message.reply_to_message && itemMatch)
        return replyForItemUpdate(send, edit, message.chat.id, message.reply_to_message, itemMatch, env.DB);
    console.info({ status: "ignored", reason: "the message does not match any cases", message });
    ignoredMessageCounts[message.chat.id] = (ignoredMessageCounts[message.chat.id] ?? 0) + 1;
    if (ignoredMessageCounts[message.chat.id] > 3) {
        ignoredMessageCounts[message.chat.id] = 0;
        return send("kalau bingung bisa pencet /bantuan atau tanya langsung ke @alpancs 💁‍♂️");
    }
}
