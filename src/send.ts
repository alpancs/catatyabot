const allEscapeeChars = ['_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'];
const userInputEscapeeChars = ['*', '_', '~'];
let nonUserInputEscapeeChars = allEscapeeChars.filter(c => !userInputEscapeeChars.includes(c));

export function escapeUserInput(text: string) {
    for (const c of userInputEscapeeChars) text = text.replaceAll(c, `\\${c}`);
    return text;
}

function escapeNonUserInput(text: string) {
    for (const c of nonUserInputEscapeeChars) text = text.replaceAll(c, `\\${c}`);
    return text;
}

export async function sendMessage(botToken: string, chatId: number, text: string, replyToMessageId?: number, forceReply?: boolean) {
    return sendCleanMessage(botToken, chatId, escapeNonUserInput(text), replyToMessageId, forceReply);
}

async function sendCleanMessage(botToken: string, chatId: number, text: string, replyToMessageId?: number, forceReply?: boolean) {
    const response = await fetch(`https://api.telegram.org/bot${botToken}/sendMessage`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: chatId,
            reply_to_message_id: replyToMessageId,
            text: text,
            parse_mode: "MarkdownV2",
            reply_markup: {
                force_reply: forceReply,
                selective: true,
            },
        }),
    });
    if (response.status >= 400) {
        const responseText = await response.text();
        const match = responseText.match(/.*Character '(?<problem>.)' is reserved and must be escaped.*/);
        if (match) {
            console.warn(responseText);
            const problem = match.groups?.problem!;
            nonUserInputEscapeeChars.push(problem);
            return sendCleanMessage(botToken, chatId, text.replaceAll(problem, `\\${problem}`), replyToMessageId, forceReply);
        }
        throw new Error(responseText);
    }
    return response;
}

export async function editMessage(botToken: string, chatId: number, messageId: number, text: string) {
    return fetch(`https://api.telegram.org/bot${botToken}/editMessageText`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            chat_id: chatId,
            message_id: messageId,
            text: escapeNonUserInput(text),
            parse_mode: "MarkdownV2",
        }),
    });
}
