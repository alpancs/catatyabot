import { escapeUserInput } from "./send";

export const readItemsQuestion = "mau lihat catatan dari berapa hari yang lalu?";
export const peekItemsQuestion = "mau intip catatan dari berapa hari yang lalu?";
const answerPattern = /^[\s-]*(?:(?<hashtagLeft>#\w+)\s+)?(?<answer>dari\s+awal|semua|semuanya|(?<coef>\d+\.?\d*)\s*(?<unit>hari|hr|pekan|minggu|bulan|bln|tahun|th|thn)?(?:\s+(?:yang\s+lalu|yg\s+lalu|terakhir))?)[\s.]*(?:\s+(?<hashtagRight>#\w+))?\s*$/i;
const months = ["Januari", "Februari", "Maret", "April", "Mei", "Juni",
    "Juli", "Agustus", "September", "Oktober", "November", "Desember"];

export async function replyForItemsReading(db: D1Database, chatId: number, text: string, showItems: boolean, actions: TelegramActions): Promise<void> {
    const match = text.match(answerPattern);
    if (!match) { await actions.ask(readItemsQuestion); return; }

    const { days, hashtag } = parseDaysMatch(match);
    const hashtagOnTitle = hashtag ? ` ${hashtag}` : "";
    let title = days === undefined ? `*=== SEMUA CATATAN${hashtagOnTitle} ===*` : `*=== CATATAN${hashtagOnTitle} DARI ${days} HARI YANG LALU ===*`;

    let query = `
        SELECT chat_id, message_id, name, price, datetime(created_at, '+7 hours') created_at
        FROM items
        ${hashtag ? "JOIN hashtags USING (chat_id, message_id)" : ""}
        WHERE chat_id = ?
        ${days === undefined ? "" : "AND created_at >= datetime('now', ?)"}
        ${hashtag === undefined ? "" : "AND lower(hashtag) = lower(?)"}`;
    let values: any[] = [chatId];
    if (days !== undefined) values.push(`-${days} days`);
    if (hashtag !== undefined) values.push(hashtag);

    try {
        await replyWithItems(title, (await db.prepare(query).bind(...values).all<Item>()).results ?? [], showItems, actions);
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause?.message });
        await actions.send("maaf lagi ada masalah nih, gak bisa lihat daftar catatan 😵");
    }
}

async function replyWithItems(title: string, items: Item[], showItems: boolean, actions: TelegramActions) {
    if (items.length === 0) return actions.send(`${title}\n\n_masih kosong_`);
    let text = title;
    let lastCreationDate = "0000-00-00";
    let count = 0;
    let total = 0;
    let grandTotal = 0;
    for (const { name, price, created_at } of items) {
        if (!created_at.startsWith(lastCreationDate)) {
            if (count > 1 || !showItems) text += `\n_total: ${thousandSeparated(total)}_`;
            lastCreationDate = created_at.substring(0, 10);
            text += `\n\n*__${idDateFormat(lastCreationDate)}__*`;
            count = 0;
            total = 0;
        }
        if (showItems) text += `\n_${created_at.substring(11, 16)}_ ${escapeUserInput(name)} ${thousandSeparated(price)}`;
        count += 1;
        total += price;
        grandTotal += price;
    }
    if (count > 1) text += `\n_total: ${thousandSeparated(total)}_`;
    text += `\n\n*_grand total: ${thousandSeparated(grandTotal)}_*`;
    return actions.send(text);
}

function idDateFormat(date: string) {
    return `${parseInt(date.substring(8, 10))} ${months[parseInt(date.substring(5, 7)) - 1]} ${date.substring(0, 4)}`
}

function parseDaysMatch(match: RegExpMatchArray) {
    const groups: { [key: string]: string | undefined } = match.groups!;
    const hashtag = groups.hashtagLeft ?? groups.hashtagRight;
    const answer = groups.answer?.toLowerCase();
    if (answer === "dari awal" || answer === "semua" || answer === "semuanya")
        return { days: undefined, hashtag };

    let days = parseFloat(groups.coef!);
    if (groups.unit === "pekan" || groups.unit === "minggu") days *= 7;
    else if (groups.unit === "bulan" || groups.unit === "bln") days *= 30;
    else if (groups.unit === "tahun" || groups.unit === "thn" || groups.unit === "th") days *= 365;
    return { days, hashtag };
}

export function thousandSeparated(n: number): string {
    if (n < 0) return `-${thousandSeparated(-n)}`;
    if (n < 1000) return Math.round(n).toString();
    return `${thousandSeparated(Math.floor(n / 1000))}.`
        + (Math.round(n) % 1000).toString().padStart(3, "0");
}
