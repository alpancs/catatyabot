import { escapeUserInput } from "./send";
import { answerPattern as itemPattern } from "./create";

export const noItemToDelete = "mau hapus pesan yang mana? balas pesan bot 👆 yang berisi catatan pakai perintah /hapus";

export async function replyForItemDeletion(send: SendTextFn, edit: EditTextFn, chatId: number, replyToMessageId: number, replyToMessageText: string, db: D1Database) {
    try {
        await db.prepare("DELETE FROM items WHERE chat_id = ?1 AND message_id = ?2")
            .bind(chatId, replyToMessageId).run();
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        return send(`"${escapeUserInput(replyToMessageText)}" ga bisa dihapus 😵`);
    }

    try {
        await edit(replyToMessageId, `~${escapeUserInput(replyToMessageText)}~`);
    } catch (error: any) {
        console.error({ replyToMessageId, replyToMessageText, error });
    }
    return send(`${escapeUserInput(replyToMessageText)} sudah dihapus 🚮`);
}
