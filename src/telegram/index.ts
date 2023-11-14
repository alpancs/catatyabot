import { sendMessage, editMessage, deleteMessage, responseToSendMessage } from "./send";
import { createItemsQuestion, replyForItemsCreation, itemPattern } from "./create";
import { readItemsQuestion, peekItemsQuestion, replyForItemsReading } from "./read";
import { replyForItemUpdate } from "./update";
import { noItemToDelete, replyForItemDeletion, deleteMessages } from "./delete";
import { helpMessage } from "./help";
import { migrateItems } from "./migration";

let ignoredMessageCounts: { [key: number]: number } = {};

export async function respondTelegramUpdate(update: Update, env: Env): Promise<Response | void> {
    if (update.message) return respondMessage(update.message, env);
    console.info({ status: "ignored", reason: "the update did not contain a message", update });
}

async function respondMessage(message: Message, env: Env): Promise<Response | void> {
    const quickSend = (text: string, forceReply?: boolean) => responseToSendMessage(message.chat.id, text, forceReply);
    const quickAsk = (text: string) => quickSend(text, true);
    const actions = {
        send: (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text),
        ask: (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, true),
        edit: (messageId: number, text: string) => editMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, messageId, text),
        delete: (messageId: number) => deleteMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, messageId),
    };

    if (message.text?.match(/^\s*\/?(start|bantuan)(@catatyabot)?\s*$/i)) return quickSend(helpMessage);
    if (message.text?.match(/^\s*\/?catat(@catatyabot)?\s*$/i)) return quickAsk(createItemsQuestion);
    if (message.text?.match(/^\s*\/?lihat(@catatyabot)?\s*$/i)) return quickAsk(readItemsQuestion);
    if (message.text?.match(/^\s*\/?intip(@catatyabot)?\s*$/i)) return quickAsk(peekItemsQuestion);
    if (message.text?.match(/^\s*\/?hapus(@catatyabot)?\s*$/i)) return message.reply_to_message ?
        replyForItemDeletion(env.DB, message.chat.id, message.reply_to_message, actions) : quickSend(noItemToDelete);
    if (message.text?.match(/^\s*\/?delete(@catatyabot)?\s*/i) && message.from?.username === "alpancs")
        return deleteMessages(env.DB, message.chat.id, message.text, actions);

    if (message.reply_to_message?.text === createItemsQuestion && message.text)
        return replyForItemsCreation(env.DB, message.text, actions);
    if (message.reply_to_message?.text === readItemsQuestion && message.text)
        return replyForItemsReading(env.DB, message.chat.id, message.text, true, actions);
    if (message.reply_to_message?.text === peekItemsQuestion && message.text)
        return replyForItemsReading(env.DB, message.chat.id, message.text, false, actions);
    const itemMatch = message.text?.match(itemPattern);
    if (message.reply_to_message && itemMatch)
        return replyForItemUpdate(env.DB, message.chat.id, message.reply_to_message, itemMatch, actions);

    if (message.migrate_from_chat_id) return migrateItems(env.DB, message.migrate_from_chat_id, actions);

    console.info({ status: "ignored", reason: "the message did not match any cases", message });
    ignoredMessageCounts[message.chat.id] = ((ignoredMessageCounts[message.chat.id] ?? 0) + 1) % 3;
    if (ignoredMessageCounts[message.chat.id] === 0)
        return quickSend("kalau bingung bisa pencet /bantuan atau tanya langsung ke @alpancs 💁‍♂️");
}
