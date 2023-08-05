import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export async function migrateItems(actions: TelegramActions, fromChatId: number, db: D1Database) {
    const items = await getItems(fromChatId, db);
    if (items.length) {
        await actions.send("Wah group ini barusan banget di-upgrade jadi supergroup dan kode groupnya jadi berubah. Ijin catat ulang ya 🙏");
        for (const item of items) await migrateItem(actions, item, db);
    }
}

async function migrateItem(actions: TelegramActions, item: Item, db: D1Database) {
    const message = `*${escapeUserInput(item.name)}* *${thousandSeparated(item.price)}* dicatat ulang ✅`;
    const { result } = await (await actions.send(message)).json<{ result: Message }>();
    await db.prepare("UPDATE items SET chat_id = ?3, message_id = ?4 WHERE chat_id = ?1 AND message_id = ?2")
        .bind(item.chat_id, item.message_id, result.chat.id, result.message_id).run();
}

async function getItems(chatId: number, db: D1Database) {
    return (await db.prepare("SELECT * FROM items WHERE chat_id = ?").bind(chatId).all<Item>()).results ?? [];
}
