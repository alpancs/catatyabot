import { escapeUserInput } from "./send";

export const createItemsQuestion = "apa saja yang mau dicatat?";

export async function askToCreateItems(ask: SendTextFn) {
    return ask(createItemsQuestion);
}

export async function replyForItemsCreation(reply: SendTextFn, text: string) {
    return reply(escapeUserInput(text));
}
