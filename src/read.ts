import { escapeUserInput } from "./send";

export const readItemsQuestion = "mau lihat catatan dari berapa hari yang lalu?";
const answerPattern = /^\s*(\d+)\s*(hari|hr|pekan|minggu|bulan|bln|tahun|th|thn)?\s*(y(an)?g)?\s*(lalu|terakhir)?\s*\.*\s*$/;
const months = ["Januari", "Februari", "Maret", "April", "Mei", "Juni",
    "Juli", "Agustus", "September", "Oktober", "November", "Desember"];

export async function replyForItemsReading(send: SendTextFn, ask: SendTextFn, chatId: number, text: string, db: D1Database) {
    text = text.toLowerCase();
    const match = text.match(answerPattern);
    if (!(match || text.startsWith("dari awal") || text.startsWith("semua"))) return ask(readItemsQuestion);

    let query = `SELECT chat_id, message_id, name, price, datetime(created_at, '+7 hours') created_at FROM items WHERE chat_id = ${chatId}`;
    let title = "*=== SEMUA CATATAN ===*";
    if (match) {
        const days = parseDays(match);
        query += ` AND created_at >= datetime('now', '-${days} days')`;
        title = `*=== CATATAN DARI ${days} HARI YANG LALU ===*`;
    }
    try {
        return replyWithItems(send, title, (await db.prepare(query).all<Item>()).results);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        return send("maaf lagi ada masalah nih, gak bisa lihat daftar catatan 😵");
    }
}

async function replyWithItems(send: SendTextFn, title: string, items?: Item[]) {
    if (!items?.length) return send(`${title}\n\n_masih kosong_`);
    let text = title;
    let lastCreationDate = "0000-00-00";
    let count = 0;
    let total = 0;
    let grandTotal = 0;
    for (const { name, price, created_at } of items) {
        if (!created_at.startsWith(lastCreationDate)) {
            if (count > 1) text += `\n_total: ${thousandSeparated(total)}_`;
            lastCreationDate = created_at.substring(0, 10);
            text += `\n\n*__${idDateFormat(lastCreationDate)}__*`;
            count = 0;
            total = 0;
        }
        text += `\n_${created_at.substring(11, 16)}_ ${escapeUserInput(name)} ${thousandSeparated(price)}`;
        count += 1;
        total += price;
        grandTotal += price;
    }
    if (count > 1) text += `\n_total: ${thousandSeparated(total)}_`;
    text += `\n\n*_grand total: ${thousandSeparated(grandTotal)}_*`;
    return send(text);
}

function idDateFormat(date: string) {
    return `${parseInt(date.substring(8, 10))} ${months[parseInt(date.substring(5, 7)) - 1]} ${date.substring(0, 4)}`
}

function parseDays(match: RegExpMatchArray) {
    let days = parseInt(match[1]);
    const unit = match[2];
    if (unit === "pekan" || unit === "minggu") days *= 7;
    else if (unit === "bulan" || unit === "bln") days *= 30;
    else if (unit === "tahun" || unit === "thn" || unit === "th") days *= 365;
    return days;
}

export function thousandSeparated(n: number): string {
    if (n < 0) return `-${thousandSeparated(-n)}`;
    if (n < 1000) return n.toString();
    return `${thousandSeparated(Math.floor(n / 1000))}.`
        + (n % 1000).toString().padStart(3, "0");
}
