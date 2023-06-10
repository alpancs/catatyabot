let problematicChars = ['_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'];

export async function sendMessage(botToken: string, chatId: number, text: string, replyToMessageId?: number, forceReply?: boolean) {
    for (const problem of problematicChars) {
        text = text.replaceAll(problem, `\\${problem}`);
    }
    return sendCleanMessage(botToken, chatId, text, replyToMessageId, forceReply);
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
        console.error(responseText);
        if (response.status < 500) {
            const needToEscapePattern = /.*Character '(.)' is reserved and must be escaped.*/;
            if (needToEscapePattern.test(responseText)) {
                const problem = responseText.replace(needToEscapePattern, "$1");
                problematicChars.push(problem);
                return sendCleanMessage(botToken, chatId, text.replaceAll(problem, `\\${problem}`), replyToMessageId, forceReply);
            }
        }
    }
}
