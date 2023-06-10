export const askForNewItemsText = "apa saja yang mau dicatat?";

export async function askForNewItems(ask: SendTextFn) {
    return ask(askForNewItemsText);
}
