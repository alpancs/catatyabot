import { escapeUserInput } from "./send";

export const readItemsQuestion = "mau lihat catatan dari berapa hari yang lalu?";
const answerPattern = /^\s*(\d+)\s*(hari)?\s*(y(an)?g)?\s*(lalu)?\s*(ya|aja|\.*)?\s*$/;
const months = ["Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"];

export async function replyForReadItems(reply: SendTextFn, ask: SendTextFn, chatId: number, text: string, db: D1Database) {
    text = text.toLowerCase();
    const match = text.match(answerPattern);
    if (!(match || text.startsWith("dari awal") || text.startsWith("semua"))) {
        return ask(readItemsQuestion);
    }

    let query = `SELECT chat_id, message_id, name, price, datetime(created_at, '+7 hours') created_at FROM items WHERE chat_id = ${chatId}`;
    if (match) query += ` AND created_at >= datetime('now', '-${match[1]} days')`;
    let title = match ? `*=== CATATAN DARI ${match[1]} HARI YANG LALU ===*` : "*=== SEMUA CATATAN ===*";
    try {
        return replyWithItems(reply, title, (await db.prepare(query).all<Item>()).results);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        return reply("maaf lagi ada masalah nih, gak bisa lihat daftar catatan 🙏");
    }
}

async function replyWithItems(reply: SendTextFn, title: string, items?: Item[]) {
    if (!items?.length) return reply("_catatan masih kosong_");
    let text = title;
    let lastCreationDate = "0000-00-00";
    let total = 0;
    let grandTotal = 0;
    for (const { name, price, created_at } of items) {
        if (!created_at.startsWith(lastCreationDate)) {
            if (lastCreationDate !== "0000-00-00") {
                text += `\n_total: ${thousandSeparated(total)}_`;
                total = 0;
            }
            lastCreationDate = created_at.substring(0, 10);
            text += `\n\n*__${idFormatted(lastCreationDate)}__*`;
        }
        total += price;
        grandTotal += price;
        text += `\n${created_at.substring(11, 16)} ${escapeUserInput(name)} ${thousandSeparated(price)}`;
    }
    text += `\n_total: ${thousandSeparated(total)}_\n\n*_grand total: ${thousandSeparated(grandTotal)}_*`;
    return reply(text);
}

function idFormatted(date: string) {
    return `${parseInt(date.substring(8, 10))} ${months[parseInt(date.substring(5, 7)) - 1]} ${date.substring(0, 4)}`
}

function thousandSeparated(n: number): string {
    if (n < 0) return `-${thousandSeparated(-n)}`;
    if (n < 1000) return n.toString();
    return `${thousandSeparated(Math.floor(n / 1000))}.`
        + (n % 1000).toString().padStart(3, "0");
}
