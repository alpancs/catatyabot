import { editMessage, sendMessage } from "./send";
import { createItemsQuestion, replyForItemsCreation, itemPattern } from "./create";
import { readItemsQuestion, replyForItemsReading } from "./read";
import { replyForItemUpdate } from "./update";
import { noItemToDelete, replyForItemDeletion } from "./delete";
import { helpMessage } from "./help";
import { migrateItems } from "./migration";

let ignoredMessageCounts: { [key: number]: number } = {};

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env);
    else console.info({ status: "ignored", reason: "the update did not contain a message", update });
    return new Response();
}

async function respondMessage(message: Message, env: Env) {
    const send = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text);
    const ask = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, true);
    const edit = (messageId: number, text: string) => editMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, messageId, text);

    if (message.text?.match(/^\s*\/?(start|bantuan)(@catatyabot)?\s*$/i)) return send(helpMessage);
    if (message.text?.match(/^\s*\/?catat(@catatyabot)?\s*$/i)) return ask(createItemsQuestion);
    if (message.text?.match(/^\s*\/?lihat(@catatyabot)?\s*$/i)) return ask(readItemsQuestion);
    if (message.text?.match(/^\s*\/?hapus(@catatyabot)?\s*$/i)) return message.reply_to_message ?
        replyForItemDeletion(send, edit, message.chat.id, message.reply_to_message, env.DB) : send(noItemToDelete);

    if (message.reply_to_message?.text === createItemsQuestion && message.text)
        return replyForItemsCreation(send, edit, message.text, env.DB);
    if (message.reply_to_message?.text === readItemsQuestion && message.text)
        return replyForItemsReading(send, ask, message.chat.id, message.text, env.DB);
    const itemMatch = message.text?.match(itemPattern);
    if (message.reply_to_message && itemMatch)
        return replyForItemUpdate(send, edit, message.chat.id, message.reply_to_message, itemMatch, env.DB);

    if (message.migrate_from_chat_id) return migrateItems(send, message.migrate_from_chat_id, env.DB);

    console.info({ status: "ignored", reason: "the message did not match any cases", message });
    ignoredMessageCounts[message.chat.id] = ((ignoredMessageCounts[message.chat.id] ?? 0) + 1) % 3;
    if (ignoredMessageCounts[message.chat.id] === 0)
        await send("kalau bingung bisa pencet /bantuan atau tanya langsung ke @alpancs 💁‍♂️");
}
