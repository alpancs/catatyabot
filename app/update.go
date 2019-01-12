package main

import (
	"fmt"
	"net/url"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func update(msg *telegram.Message) (bool, error) {
	priceText := patternPrice.FindString(msg.Text)
	item := strings.TrimSpace(msg.Text[:len(msg.Text)-len(priceText)])
	if item == "" || priceText == "" {
		return true, nil
	}
	price := parsePrice(priceText)

	return true, editMessage(msg, url.Values{
		"chat_id":    {fmt.Sprintf("%d", msg.Chat.ID)},
		"message_id": {fmt.Sprintf("%d", msg.ReplyToMessage.MessageID)},
		"text":       {fmt.Sprintf("%s dengan harga %s ((pura-pura)) dicatat ya bos 👌", item, price)},
	})
}
