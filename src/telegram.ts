import { askToCreateItems, createItemsQuestion, replyForItemsCreation } from "./create";
import { readItemsQuestion, replyForReadItems } from "./list";
import { sendHelpMessage } from "./help";
import { editMessage, sendMessage } from "./send";

export async function getUpdateResponse(update: Update, env: Env) {
    if (update.message) await respondMessage(update.message, env)
    else console.info(JSON.stringify({ status: "ignored", reason: "the update does not contain a message", update }));
    return new Response();
}

async function respondMessage(message: Message, env: Env) {
    const send = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text);
    const reply = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, message.message_id);
    const ask = (text: string) => sendMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, text, message.message_id, true);
    const edit = (messageId: number, text: string) => editMessage(env.TELEGRAM_BOT_TOKEN, message.chat.id, messageId, text);
    if (message.text === "/start" || message.text === "/bantuan") return sendHelpMessage(send);
    if (message.text === "/catat") return askToCreateItems(ask);
    if (message.text === "/lihat") return ask(readItemsQuestion);
    if (message.reply_to_message?.text === createItemsQuestion && message.text)
        return replyForItemsCreation(reply, edit, message.text, env.DB);
    if (message.reply_to_message?.text === readItemsQuestion && message.text)
        return replyForReadItems(reply, message.chat.id, message.text, env.DB);
    console.info(JSON.stringify({ status: "ignored", reason: "the message does not match any cases", message }));
}
