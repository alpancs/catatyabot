package app

import (
	"fmt"
	"net/url"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func update(msg *telegram.Message) (bool, error) {
	if msg.ReplyToMessage == nil {
		return false, nil
	}

	priceText := patternPrice.FindString(msg.Text)
	item := strings.TrimSpace(msg.Text[:len(msg.Text)-len(priceText)])
	if item == "" || priceText == "" {
		return false, nil
	}
	price := ParsePrice(priceText)

	result, err := db.Exec("UPDATE items SET name = $3, price = $4 WHERE chat_id = $1 AND message_id = $2;", msg.Chat.ID, msg.ReplyToMessage.MessageID, item, price)
	if err != nil {
		return true, err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return true, err
	}

	if rowsAffected == 0 {
		return false, nil
	}

	err = editMessage(url.Values{
		"chat_id":    {fmt.Sprintf("%d", msg.Chat.ID)},
		"message_id": {fmt.Sprintf("%d", msg.ReplyToMessage.MessageID)},
		"text":       {fmt.Sprintf(SaveTemplate, item, price)},
		"parse_mode": {"Markdown"},
	})
	if err != nil {
		return true, err
	}
	_, err = sendMessage(url.Values{
		"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":                {"sudah diubah nih bos 👆"},
		"reply_to_message_id": {fmt.Sprintf("%d", msg.ReplyToMessage.MessageID)},
	})
	return true, err
}
