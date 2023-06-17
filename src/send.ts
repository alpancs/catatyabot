const allEscapeeChars = ['_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'];
const userInputEscapeeChars = ['*', '_', '~'];
const nonUserInputEscapeeChars = allEscapeeChars.filter(c => !userInputEscapeeChars.includes(c));
const sendTextLimit = 4096;

export function escapeUserInput(text: string) {
    for (const c of userInputEscapeeChars) text = text.replaceAll(c, `\\${c}`);
    return text;
}

function escapeNonUserInput(text: string) {
    for (const c of nonUserInputEscapeeChars) text = text.replaceAll(c, `\\${c}`);
    return text;
}

export async function sendMessage(botToken: string, chatId: number, text: string, forceReply?: boolean) {
    text = escapeNonUserInput(text);
    while (text) {
        let i = sendTextLimit;
        if (text.length > sendTextLimit) while (text[i] !== "\n") i--;
        await sendPartialMessage(botToken, chatId, text.substring(0, i), forceReply);
        text = text.substring(i + 1);
    }
}

async function sendPartialMessage(botToken: string, chatId: number, text: string, forceReply?: boolean) {
    const response = await fetch(`https://api.telegram.org/bot${botToken}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: chatId,
            text: text,
            parse_mode: "MarkdownV2",
            reply_markup: forceReply ? { force_reply: true, selective: true } : undefined,
        }),
    });
    if (!response.ok) throw new Error(await response.text());
    return response;
}

export async function editMessage(botToken: string, chatId: number, messageId: number, text: string) {
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
