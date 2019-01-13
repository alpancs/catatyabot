package main

import (
	"fmt"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func commandList(msg *telegram.Message) (bool, error) {
	if msg.Command() != "lihat" {
		return false, nil
	}

	_, err := sendMessage(url.Values{
		"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":                {"pengen lihat daftar catatan yang mana bos? 👀"},
		"reply_to_message_id": {fmt.Sprintf("%d", msg.MessageID)},
		"reply_markup": {`{
			"keyboard": [
				[{"text":"hari ini"},{"text":"kemarin"}],
				[{"text":"pekan ini"},{"text":"pekan lalu"}],
				[{"text":"bulan ini"},{"text":"bulan lalu"}]
			],
			"resize_keyboard": true,
			"one_time_keyboard": true,
			"selective": true
		}`},
	})
	return true, err
}
