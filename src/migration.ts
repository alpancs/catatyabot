import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export async function migrateItems(send: SendTextFn, fromChatId: number, db: D1Database) {
    const itemsPromise = getItems(fromChatId, db);
    await send("Wah _group_ ini barusan banget di-_upgrade_ jadi _supergroup_ dan kode _group_-nya jadi berubah. Ijin catat ulang ya 🙏");
    for (const item of await itemsPromise) await migrateItem(send, item, db);
}

async function migrateItem(send: SendTextFn, item: Item, db: D1Database) {
    const message = `*${escapeUserInput(item.name)}* *${thousandSeparated(item.price)}* dicatat ulang ✅`;
    const { result } = await (await send(message)).json<{ result: Message }>();
    await db.prepare("UPDATE items SET chat_id = ?3, message_id = ?4 WHERE chat_id = ?1 AND message_id = ?2")
        .bind(item.chat_id, item.message_id, result.chat.id, result.message_id).run();
}

async function getItems(chatId: number, db: D1Database) {
    const { results } = await db.prepare("SELECT * FROM items WHERE chat_id = ?").bind(chatId).all<Item>();
    return results ?? [];
}
