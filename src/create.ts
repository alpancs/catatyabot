import { escapeUserInput } from "./send";

export const createItemsQuestion = "apa saja yang mau dicatat?";
const itemPattern = /^\s*(?<name>.+)\s+(?<price>-?\d+)\s*(?<unit>ribu|rb|k|juta|jt)?\s*$/i;

export async function askToCreateItems(ask: SendTextFn) {
    return ask(createItemsQuestion);
}

export async function replyForItemsCreation(reply: SendTextFn, edit: EditTextFn, text: string, db: D1Database) {
    for (const line of text.split("\n")) {
        await replyForItemCreation(reply, edit, line, db);
    }
}

async function replyForItemCreation(reply: SendTextFn, edit: EditTextFn, text: string, db: D1Database) {
    const match = text.match(itemPattern);
    if (!match) return reply(`⚠️ "${escapeUserInput(text)}" tidak dicatat karena tidak ada harganya/tidak ada namanya.`);

    const { name, price } = parse(match.groups!)
    const replyResponse = await reply(`*${escapeUserInput(name)}* *${price}* dicatat ✅`);
    const { result } = await replyResponse.json<{ result: Message }>();

    try {
        const { success, error } = await db
            .prepare("INSERT INTO items (chat_id, message_id, name, price, created_at) VALUES (?, ?, ?, ?, datetime('now'));")
            .bind(result.chat.id, result.message_id, name, price).run();
        if (!success) throw new Error(error);
    }
    catch (error: any) {
        await edit(result.message_id, `*${escapeUserInput(name)}* *${price}* gagal dicatat ❌`);
        throw error;
    }
}

function parse(groups: { [key: string]: string; }) {
    let price = parseInt(groups.price!);
    const unit = groups.unit.toLowerCase();
    if (unit === 'ribu' || unit === 'rb' || unit === 'k') price *= 1000;
    else if (unit === 'juta' || unit === 'jt') price *= 1000000;
    return { name: groups.name!, price };
}
