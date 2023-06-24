import { parseItemMatch } from "./create";
import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export async function replyForItemUpdate(send: SendTextFn, edit: EditTextFn, chatId: number, replyToMessage: Message, itemMatch: RegExpMatchArray, db: D1Database) {
    const { name, price } = parseItemMatch(itemMatch);
    try {
        const updatedItem = await db.prepare("UPDATE items SET name = ?3, price = ?4 WHERE chat_id = ?1 AND message_id = ?2 RETURNING *")
            .bind(chatId, replyToMessage.message_id, name, price).first<Item | null>();
        if (updatedItem) await edit(replyToMessage.message_id, `~${replyToMessage.text}~\n*${escapeUserInput(name)}* *${thousandSeparated(price)}*`);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        await send(`ada masalah pas lagi ubah catatan 😵`);
    }
}
