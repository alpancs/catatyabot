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
    const response = await postToTelegram(botToken, "sendMessage", sendMessagePayload(chatId, head, forceReply));
    return tail ? sendCleanMessage(botToken, chatId, tail, forceReply) : response;
}

export function responseToSendMessage(chatId: number, text: string, forceReply?: boolean): Response {
    return Response.json({ method: "sendMessage", ...sendMessagePayload(chatId, escapeNonUserInput(text), forceReply) });
}

function sendMessagePayload(chatId: number, text: string, forceReply?: boolean): any {
    return {
        chat_id: chatId,
        text: text,
        parse_mode: "MarkdownV2",
        reply_markup: forceReply ? { force_reply: true } : undefined,
    };
}

export async function editMessage(botToken: string, chatId: number, messageId: number, text: string): Promise<Response> {
    return postToTelegram(botToken, "editMessageText", {
        chat_id: chatId,
        message_id: messageId,
        text: escapeNonUserInput(text),
        parse_mode: "MarkdownV2",
    });
}

async function postToTelegram(token: string, methodName: string, body: any): Promise<Response> {
    const response = await fetch(`https://api.telegram.org/bot${token}/${methodName}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
    });
    if (!response.ok) {
        const responseJSON = await response.json<{ description: string }>();
        const description = responseJSON.description;
        const retryDelaySeconds = parseInt(description.match(/Too Many Requests: retry after (\d+)/i)?.[1] ?? "");
        if (retryDelaySeconds) return new Promise(resolve => setTimeout(() => resolve(postToTelegram(token, methodName, body)), retryDelaySeconds));
        else throw new Error(JSON.stringify(responseJSON));
    }
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
