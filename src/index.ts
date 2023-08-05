import { getUpdateResponse } from "./telegram/index";

export default {
	async fetch(request: Request, env: Env) {
		if (request.method === "POST" && new URL(request.url).pathname === "/webhook/telegram")
			return this.handleWebhookTelegram(request, env);
		return new Response(undefined, { status: 404 });
	},

	async handleWebhookTelegram(request: Request, env: Env) {
		if (request.headers.get("X-Telegram-Bot-Api-Secret-Token") !== env.TELEGRAM_WEBHOOK_SECRET_TOKEN)
			return new Response(undefined, { status: 401 });
		try {
			return getUpdateResponse(await request.json(), env);
		} catch (error: any) {
			return new Response(error, { status: 422 })
		}
	},
};
