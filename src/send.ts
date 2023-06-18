const allEscapeeChars = ['_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'];
const userInputEscapeeChars = ['*', '_', '~'];
const nonUserInputEscapeeChars = allEscapeeChars.filter(c => !userInputEscapeeChars.includes(c));

export function escapeUserInput(text: string) {
    for (const c of userInputEscapeeChars) text = text.replaceAll(c, `\\${c}`);
    return text;
}

function escapeNonUserInput(text: string) {
    for (const c of nonUserInputEscapeeChars) text = text.replaceAll(c, `\\${c}`);
    return text;
}

export async function sendMessage(botToken: string, chatId: number, text: string, forceReply?: boolean): Promise<Response> {
    return sendCleanMessage(botToken, chatId, escapeNonUserInput(text), forceReply);
}

async function sendCleanMessage(botToken: string, chatId: number, text: string, forceReply?: boolean): Promise<Response> {
    const { head, tail } = splitOnTelegramLimit(text);
    const response = await fetch(`https://api.telegram.org/bot${botToken}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: chatId,
            text: head,
            parse_mode: "MarkdownV2",
            reply_markup: forceReply ? { force_reply: true } : undefined,
        }),
    });
    if (!response.ok) throw new Error(await response.text());
    return tail ? sendCleanMessage(botToken, chatId, tail, forceReply) : response;
}

export async function editMessage(botToken: string, chatId: number, messageId: number, text: string): Promise<Response> {
    const response = await fetch(`https://api.telegram.org/bot${botToken}/editMessageText`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: chatId,
            message_id: messageId,
            text: escapeNonUserInput(text),
            parse_mode: "MarkdownV2",
        }),
    });
    if (!response.ok) throw new Error(await response.text());
    return response;
}

function splitOnTelegramLimit(text: string) {
    const limit = 4096;
    if (text.length <= limit) return { head: text, tail: "" };
    for (let i = limit; i >= 0; --i)
        if (text[i] === "\n")
            return { head: text.substring(0, i), tail: text.substring(i + 1) };
    return { head: text.substring(0, limit), tail: text.substring(limit) };
}
