export interface Env {
	DB: D1Database;
}

export default {
	async fetch(request: Request, env: Env) {
		const { results } = await env.DB.prepare("SELECT * FROM items").all();
		return Response.json(results);
	},
};
