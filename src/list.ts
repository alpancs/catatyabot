import { escapeUserInput } from "./send";

export const readItemsQuestion = "mau lihat catatan dari berapa hari yang lalu?";
const answerPattern = /^\s*(\d+)\s*(hari)?\s*(y(an)?g)?\s*(lalu)?\s*(ya|aja|\.*)?\s*$/;

export async function replyForReadItems(reply: SendTextFn, ask: SendTextFn, chatId: number, text: string, db: D1Database) {
    text = text.toLowerCase();
    const match = text.match(answerPattern);
    if (!(match || text.startsWith("dari awal") || text.startsWith("semua"))) {
        return ask(readItemsQuestion);
    }

    let query = `SELECT chat_id, message_id, name, price, datetime(created_at, '+7 hours') created_at FROM items WHERE chat_id = ${chatId}`;
    if (match) query += ` AND created_at >= datetime('now', '-${match[1]} days')`;
    try {
        return replyWithItems(reply, (await db.prepare(query).all<Item>()).results);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        return reply("maaf lagi ada masalah nih, gak bisa lihat daftar catatan 🙏");
    }
}

async function replyWithItems(reply: SendTextFn, items?: Item[]) {
    if (!items?.length) return reply("_catatan masih kosong_");
    let text = "*=== DAFTAR CATATAN ===*";
    let lastCreationDate = "0000-00-00";
    for (const item of items) {
        if (!item.created_at.startsWith(lastCreationDate)) {
            lastCreationDate = item.created_at.substring(0, 10);
            text += `\n\n__${escapeUserInput(lastCreationDate)}__`;
        }
        text += escapeUserInput(`\n${item.created_at.substring(11, 16)} ${item.name} ${thousandSeparated(item.price)}`);
    }
    return reply(text);
}

function thousandSeparated(n: number): string {
    if (n < 0) return `-${thousandSeparated(-n)}`;
    if (n < 1000) return n.toString();
    return `${thousandSeparated(Math.floor(n / 1000))}.`
        + (n % 1000).toString().padStart(3, "0");
}
