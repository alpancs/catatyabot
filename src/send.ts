const allEscapeeChars = ['_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'];
const userInputEscapeeChars = ['*', '_', '~'];
let generalEscapeeChars = allEscapeeChars.filter(c => !userInputEscapeeChars.includes(c));

export async function sendMessage(botToken: string, chatId: number, text: string, replyToMessageId?: number, forceReply?: boolean) {
    for (const c of generalEscapeeChars) text = text.replaceAll(c, `\\${c}`);
    return sendCleanMessage(botToken, chatId, text, replyToMessageId, forceReply);
}

export function escapeUserInput(text: string) {
    for (const c of userInputEscapeeChars) text = text.replaceAll(c, `\\${c}`);
    return text;
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
        const needToEscapePattern = /.*Character '(.)' is reserved and must be escaped.*/;
        if (needToEscapePattern.test(responseText)) {
            console.warn(responseText);
            const problem = responseText.replace(needToEscapePattern, "$1");
            generalEscapeeChars.push(problem);
            return sendCleanMessage(botToken, chatId, text.replaceAll(problem, `\\${problem}`), replyToMessageId, forceReply);
        }
        throw responseText;
    }
}
