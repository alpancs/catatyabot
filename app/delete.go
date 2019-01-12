package main

import (
	"fmt"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func commandDelete(msg *telegram.Message) (bool, error) {
	if msg.Command() != "hapus" || msg.ReplyToMessage == nil {
		return false, nil
	}

	_, err := sendMessage(url.Values{
		"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":                {"ini ((pura-pura)) sudah dihapus ya bos 👌"},
		"reply_to_message_id": {fmt.Sprintf("%d", msg.ReplyToMessage.MessageID)},
	})
	return true, err
}
