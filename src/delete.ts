import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export const noItemToDelete = `mau hapus pesan yang mana?
_\\*cara hapus suatu catatan: balas pesan bot yang ada tanda ✅nya pakai perintah /hapus_`;

export async function replyForItemDeletion(send: SendTextFn, edit: EditTextFn, chatId: number, replyToMessageId: number, db: D1Database) {
    try {
        const deletedItem = await db.prepare("DELETE FROM items WHERE chat_id = ?1 AND message_id = ?2 RETURNING *")
            .bind(chatId, replyToMessageId).first<Item | null>();
        if (deletedItem)
            return edit(replyToMessageId, `~*${escapeUserInput(deletedItem.name)}* *${thousandSeparated(deletedItem.price)}* dicatat ✅~`);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        return send(`ada masalah pas lagi hapus catatan 😵`);
    }
}
