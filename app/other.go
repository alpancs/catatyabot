package app

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func other(msg *telegram.Message) error {
	_, err := sendMessage(msg.Chat.ID, "ngapain bos? 🙄", msg.MessageID)
	return err
}
