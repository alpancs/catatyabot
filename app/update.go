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
		return true, nil
	}
	price := ParsePrice(priceText)

	return true, editMessage(url.Values{
		"chat_id":    {fmt.Sprintf("%d", msg.Chat.ID)},
		"message_id": {fmt.Sprintf("%d", msg.ReplyToMessage.MessageID)},
		"text":       {fmt.Sprintf("%s dengan harga %s ((pura-pura)) dicatat ya bos 👌", item, price)},
	})
}
