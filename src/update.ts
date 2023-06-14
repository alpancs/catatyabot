import { parseItemMatch } from "./create";
import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export async function replyForItemUpdate(send: SendTextFn, edit: EditTextFn, chatId: number, replyToMessageId: number, itemMatch: RegExpMatchArray, db: D1Database) {
    try {
        const existingItem = await db.prepare("SELECT * FROM items WHERE chat_id = ?1 AND message_id = ?2")
            .bind(chatId, replyToMessageId).first<Item | null>();
        if (existingItem) {
            const { name, price } = parseItemMatch(itemMatch);
            await db.prepare("UPDATE items SET name = ?3, price = ?4 WHERE chat_id = ?1 AND message_id = ?2")
                .bind(chatId, replyToMessageId, name, price).run();
            await edit(replyToMessageId, `~*${escapeUserInput(existingItem.name)}* *${thousandSeparated(existingItem.price)}*~ *${escapeUserInput(name)}* *${thousandSeparated(price)}* dicatat ✅`);
        }
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        await send(`ada masalah pas lagi ubah catatan 😵`);
    }
}
