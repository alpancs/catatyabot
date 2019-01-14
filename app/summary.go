package app

import (
	"fmt"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func commandSummary(msg *telegram.Message) (bool, error) {
	if msg.Command() != "rangkuman" {
		return false, nil
	}

	_, err := sendMessage(url.Values{
		"chat_id": {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":    {"rangkuman kosong bos, kan cuma pura-pura"},
	})
	return true, err
}
