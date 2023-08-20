import { thousandSeparated } from "./read";
import { escapeUserInput } from "./send";

export const createItemsQuestion = "apa saja yang mau dicatat?";
export const itemPattern = /^(?<name>.+)\s+(?:(?<withUnit>(?<priceFloat>[+-]?\d+[,.]?\d*)\s*(?<unit>ribu|rb|k|juta|jt))|(?<priceInt>-?\d+(?:[,.]\d*)*))(?<rawHashtags>(?:\s+#\w+)*)\s*$/i;

export async function replyForItemsCreation(db: D1Database, text: string, actions: TelegramActions): Promise<void> {
    const prices = (await Promise.all(text.split("\n").map(l => replyForItemCreation(db, l.match(itemPattern), actions)))).filter(p => p);
    if (prices.length > 1) await actions.send(`Totalnya barusan: *${thousandSeparated(prices.reduce((p, c) => p + c))}*`);
}

async function replyForItemCreation(db: D1Database, match: RegExpMatchArray | null, actions: TelegramActions): Promise<number> {
    if (!match) return 0;

    const { name, price, hashtags } = parseItemMatch(match);
    const message = `*${escapeUserInput(name)}* *${thousandSeparated(price)}* dicatat ✅`;
    const { result } = await (await actions.send(message)).json<{ result: Message }>();

    let statements = [
        db.prepare("INSERT INTO items (chat_id, message_id, name, price, created_at) VALUES (?1, ?2, ?3, ?4, datetime('now'));")
            .bind(result.chat.id, result.message_id, name, price)
    ];
    const insertHashtagStmt = db.prepare("INSERT OR IGNORE INTO hashtags (chat_id, message_id, hashtag) VALUES (?1, ?2, ?3);");
    for (const hashtag of hashtags) statements.push(insertHashtagStmt.bind(result.chat.id, result.message_id, hashtag));
    try {
        await db.batch(statements);
        return price;
    } catch (error: any) {
        console.error({ message: error.message, cause: error.cause?.message });
        await actions.edit(result.message_id, `*${escapeUserInput(name)}* gagal dicatat 😵`);
        return 0;
    }
}

export function parseItemMatch(match: RegExpMatchArray) {
    const groups = match.groups!;
    const name = groups.name.trim() + groups.rawHashtags;
    const hashtags = (name.match(/(?:^|\s+)#\w+/ig) ?? []).map(s => s.trim());
    if (groups.withUnit) {
        let price = parseFloat(groups.priceFloat.replace(",", "."));
        const unit = groups.unit.toLowerCase();
        if (unit === 'ribu' || unit === 'rb' || unit === 'k') price *= 1000;
        else if (unit === 'juta' || unit === 'jt') price *= 1000000;
        return { name, price, hashtags };
    }
    return { name, price: parseInt(groups.priceInt.replace(/[,.]/g, "")), hashtags };
}
