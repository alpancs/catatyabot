import { Env } from "./env"

export interface Update {
    update_id: number;
}

export async function getUpdateResponse(update: Update, env: Env) {
    return new Response(null);
}
