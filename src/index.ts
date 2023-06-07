import { Env } from "./env";
import { Update, getUpdateResponse } from "./telegram";

export default {
	async fetch(request: Request, env: Env) {
		const { pathname } = new URL(request.url);
		if (request.method == "POST" && pathname == "/webhook/telegram") {
			return this.handleRequestWebhookTelegram(request, env);
		}
		return new Response(null, { status: 404 });
	},

	async handleRequestWebhookTelegram(request: Request, env: Env) {
		if (!request.headers.has("X-Telegram-Bot-Api-Secret-Token")) {
			return new Response(null, { status: 401 });
		}
		if (request.headers.get("X-Telegram-Bot-Api-Secret-Token") != env.TELEGRAM_BOT_TOKEN) {
			return new Response(null, { status: 403 });
		}
		let update: Update;
		try {
			update = await request.json();
		} catch (error: any) {
			return new Response(error, { status: 422 })
		}
		return getUpdateResponse(update, env);
	},
};
