export interface Env {
	DB: D1Database;
	TELEGRAM_TOKEN: string;
}

export default {
	async fetch(request: Request, env: Env) {
		const { pathname } = new URL(request.url);
		if (request.method == "POST" && pathname == "/webhook/telegram") {
			return handleRequestWebhookTelegram(request, env);
		}
		return new Response(null, { status: 404 });
	},
};

async function handleRequestWebhookTelegram(request: Request, env: Env) {
	if (!request.headers.has("X-Telegram-Bot-Api-Secret-Token")) {
		return new Response(null, { status: 401 });
	}
	if (request.headers.get("X-Telegram-Bot-Api-Secret-Token") != env.TELEGRAM_TOKEN) {
		return new Response(null, { status: 403 });
	}
	return Response.json({ message: "OK" });
}
