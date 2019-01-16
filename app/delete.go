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

	result, err := db.Exec("DELETE FROM items WHERE chat_id = $1 AND message_id = $2;", msg.Chat.ID, msg.ReplyToMessage.MessageID)
	if err != nil {
		return true, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return true, err
	}

	if rowsAffected == 0 {
		_, err = sendMessage(url.Values{
			"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":                {"hapus apa to bos? 🙄"},
			"reply_to_message_id": {fmt.Sprintf("%d", msg.MessageID)},
		})
	} else {
		_, err = sendMessage(url.Values{
			"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":                {"ini sudah dihapus ya bos 🚮"},
			"reply_to_message_id": {fmt.Sprintf("%d", msg.ReplyToMessage.MessageID)},
		})
	}
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
