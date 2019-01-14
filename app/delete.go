package app

import (
	"fmt"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func commandDelete(msg *telegram.Message) (bool, error) {
	if msg.Command() != "hapus" {
		return false, nil
	}

	if msg.ReplyToMessage == nil {
		return true, helpDelete(msg)
	}

	_, err := sendMessage(url.Values{
		"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":                {"ini ((pura-pura)) sudah dihapus ya bos 👌"},
		"reply_to_message_id": {fmt.Sprintf("%d", msg.ReplyToMessage.MessageID)},
	})
	return true, err
}

func helpDelete(msg *telegram.Message) error {
	_, err := sendMessage(url.Values{
		"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":                {"tolong reply ke catatan yang pengen dihapus ya bos"},
		"reply_to_message_id": {fmt.Sprintf("%d", msg.MessageID)},
	})
	return err
}
