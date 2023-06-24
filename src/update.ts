import { parseItemMatch } from "./create";
import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export async function replyForItemUpdate(send: SendTextFn, edit: EditTextFn, chatId: number, replyToMessage: Message, itemMatch: RegExpMatchArray, db: D1Database) {
    const { name, price, hashtags } = parseItemMatch(itemMatch);
    let statements = [
        db.prepare("UPDATE items SET name = ?3, price = ?4 WHERE chat_id = ?1 AND message_id = ?2 RETURNING *;")
            .bind(chatId, replyToMessage.message_id, name, price),
        db.prepare("DELETE FROM hashtags WHERE chat_id = ?1 AND message_id = ?2;")
            .bind(chatId, replyToMessage.message_id),
    ];
    // TODO: fix broken references from table items
    const insertHashtagStmt = db.prepare("INSERT INTO hashtags (chat_id, message_id, hashtag) VALUES (?1, ?2, ?3);");
    for (const hashtag of hashtags) statements.push(insertHashtagStmt.bind(chatId, replyToMessage.message_id, hashtag));

    try {
        const updatedItems = (await db.batch(statements))[0].results;
        if (updatedItems?.length) await edit(replyToMessage.message_id, `~${replyToMessage.text}~\n*${escapeUserInput(name)}* *${thousandSeparated(price)}*`);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        await send(`ada masalah pas lagi ubah catatan 😵`);
    }
}
