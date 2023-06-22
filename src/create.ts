import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export const createItemsQuestion = "apa saja yang mau dicatat?";
export const itemPattern = /^(?<name>.+)\s+(?:(?<withUnit>(?<priceFloat>-?\d+[,.]?\d*)\s*(?<unit>ribu|rb|k|juta|jt))|(?<priceInt>-?\d+(?:[,.]\d*)*))$/i;

export async function replyForItemsCreation(send: SendTextFn, edit: EditTextFn, text: string, db: D1Database) {
    for (const match of text.split("\n").map(l => l.match(itemPattern))) {
        if (match) await replyForItemCreation(send, edit, match, db);
    }
}

async function replyForItemCreation(send: SendTextFn, edit: EditTextFn, match: RegExpMatchArray, db: D1Database) {
    const { name, price } = parseItemMatch(match)
    const { result } = await (await send(`*${escapeUserInput(name)}* *${thousandSeparated(price)}* dicatat ✅`)).json<{ result: Message }>();

    try {
        await db.prepare("INSERT INTO items (chat_id, message_id, name, price, created_at) VALUES (?1, ?2, ?3, ?4, datetime('now'));")
            .bind(result.chat.id, result.message_id, name, price).run();
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause.message });
        await edit(result.message_id, `*${escapeUserInput(name)}* gagal dicatat 😵`);
    }
}

export function parseItemMatch(match: RegExpMatchArray) {
    const groups = match.groups!;
    const name = groups.name.trim();
    if (groups.withUnit) {
        let price = parseFloat(groups.priceFloat.replace(",", "."));
        const unit = groups.unit.toLowerCase();
        if (unit === 'ribu' || unit === 'rb' || unit === 'k') price *= 1000;
        else if (unit === 'juta' || unit === 'jt') price *= 1000000;
        return { name, price };
    }
    return { name, price: parseInt(groups.priceInt.replace(/[,.]/g, "")) };
}
