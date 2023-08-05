export const noItemToDelete = `mau hapus pesan yang mana?
_\\*cara hapus suatu catatan: balas pesan bot yang ada tanda ✅nya pakai perintah /hapus_`;

export async function replyForItemDeletion(db: D1Database, chatId: number, replyToMessage: Message, actions: TelegramActions) {
    try {
        const deletedItem = await db.prepare("DELETE FROM items WHERE chat_id = ?1 AND message_id = ?2 RETURNING *")
            .bind(chatId, replyToMessage.message_id).first<Item | null>();
        if (deletedItem) await actions.edit(replyToMessage.message_id, `~${replyToMessage.text}~`);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        await actions.send(`ada masalah pas lagi hapus catatan 😵`);
    }
}
