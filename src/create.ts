import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export const createItemsQuestion = "apa saja yang mau dicatat?";
export const itemPattern = /^\s*(.+)\s+(-?\d+[,.]?\d*)\s*(ribu|rb|k|juta|jt)?\s*$/i;

export async function replyForItemsCreation(send: SendTextFn, edit: EditTextFn, text: string, db: D1Database) {
    for (const match of text.split("\n").map(l => l.match(itemPattern))) {
        if (match) await replyForItemCreation(send, edit, match, db);
    }
}

async function replyForItemCreation(send: SendTextFn, edit: EditTextFn, match: RegExpMatchArray, db: D1Database) {
    const { name, price } = parseItemMatch(match)
    const replyResponse = await send(`*${escapeUserInput(name)}* *${thousandSeparated(price)}* dicatat ✅`);
    const { result } = await replyResponse.json<{ result: Message }>();

    try {
        await db.prepare("INSERT INTO items (chat_id, message_id, name, price, created_at) VALUES (?1, ?2, ?3, ?4, datetime('now'));")
            .bind(result.chat.id, result.message_id, name, price).run();
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        await edit(result.message_id, `*${escapeUserInput(name)}* gagal dicatat 😵`);
    }
}

export function parseItemMatch(match: RegExpMatchArray) {
    let price = parseFloat(match[2].replace(",", "."));
    const unit = match[3]?.toLowerCase();
    if (unit === 'ribu' || unit === 'rb' || unit === 'k') price *= 1000;
    else if (unit === 'juta' || unit === 'jt') price *= 1000000;
    return { name: match[1], price };
}
