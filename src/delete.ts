import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export const noItemToDelete = "mau hapus pesan yang mana? balas pesan bot 👆 yang berisi catatan pakai perintah /hapus";

export async function replyForItemDeletion(send: SendTextFn, edit: EditTextFn, chatId: number, replyToMessageId: number, db: D1Database) {
    try {
        const deletedItem = await db.prepare("DELETE FROM items WHERE chat_id = ?1 AND message_id = ?2 RETURNING *")
            .bind(chatId, replyToMessageId).first<Item | null>();
        if (deletedItem) {
            await edit(replyToMessageId, `~${escapeUserInput(deletedItem.name)} ${thousandSeparated(deletedItem.price)}~`);
            return send(`${escapeUserInput(deletedItem.name)} ${thousandSeparated(deletedItem.price)} sudah dihapus 🚮`);
        }
        return send(`yang mau dihapus tidak ada di catatan 🙄`);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        return send(`ada masalah pas lagi hapus catatan 😵`);
    }
}
