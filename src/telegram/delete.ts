import asyncPool from "tiny-async-pool";

import { escapeUserInput } from "./send";

export const noItemToDelete = `mau hapus pesan yang mana?
_\\*cara hapus suatu catatan: balas pesan bot yang ada tanda ✅nya pakai perintah /hapus_`;

export async function replyForItemDeletion(db: D1Database, chatId: number, replyToMessage: Message, actions: TelegramActions): Promise<void> {
    try {
        const deletedItem = await db.prepare("DELETE FROM items WHERE chat_id = ?1 AND message_id = ?2 RETURNING *")
            .bind(chatId, replyToMessage.message_id).first<Item | null>();
        if (deletedItem) await actions.edit(replyToMessage.message_id, `~${escapeUserInput(replyToMessage.text!)}~`);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause?.message });
        await actions.send(`ada masalah pas lagi hapus catatan 😵`);
    }
}

export async function deleteMessages(db: D1Database, chatId: number, messageText: string, actions: TelegramActions): Promise<void> {
    const clause = messageText.replace(/^\s*\/?delete(@catatyabot)?\s*/i, "");
    let items: Item[] = [];
    try {
        const { results } = await db.prepare(`SELECT chat_id, message_id FROM items WHERE chat_id = ?1 AND ${clause}`).bind(chatId).all<Item>();
        items = results;
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause?.message });
        return;
    }
    const itemToDeleteLimit = 50;
    const toDeleteItems = items.slice(0, itemToDeleteLimit);
    for await (const v of asyncPool(2, toDeleteItems, async (item: Item) => {
        await actions.delete(item.message_id);
        await db.prepare("DELETE FROM items WHERE chat_id = ?1 AND message_id = ?2").bind(chatId, item.message_id).run();
    }));
    if (toDeleteItems.length < items.length) return Promise.reject("The request needs to be repeated due to subrequest limit.");
}
