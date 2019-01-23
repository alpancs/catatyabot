package app

import (
	"strings"

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
		_, err = sendMessage(msg.Chat.ID, "hapus apa to bos? 🙄", msg.MessageID)
	} else {
		editMessage(msg.Chat.ID, msg.ReplyToMessage.MessageID, strings.Replace(msg.ReplyToMessage.Text, " #catatan", "", -1))
		_, err = sendMessage(msg.Chat.ID, "ini sudah dihapus ya bos 🚮", msg.ReplyToMessage.MessageID)
	}
	return true, err
}

func helpDelete(msg *telegram.Message) error {
	_, err := sendMessage(msg.Chat.ID, "tolong reply ke catatan yang pengen dihapus ya bos", msg.MessageID)
	return err
}
