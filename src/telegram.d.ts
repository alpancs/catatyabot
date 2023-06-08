export interface Update {
    update_id: number;
    message?: Message;
}

export interface Message {
    message_id: number;
    chat: Chat;
    reply_to_message?: Message;
    text?: string;
}

export interface Chat {
    id: number;
}
