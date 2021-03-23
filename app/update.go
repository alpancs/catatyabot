package app

import (
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	UpdateTemplate = "ini sudah dicoret ya bos 👆\ndiganti sama *%s %s*"
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

	err = editMessage(msg.Chat.ID, msg.ReplyToMessage.MessageID, cross(msg.ReplyToMessage.Text))
	if err != nil {
		return true, err
	}
	_, err = sendMessage(msg.Chat.ID, fmt.Sprintf(UpdateTemplate, item, price), msg.ReplyToMessage.MessageID)
	return true, err
}

func cross(text string) string {
	text = strings.Replace(text, "*", "~", 2)

	restrictedChars := []byte{'_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!'}
	for _, c := range restrictedChars {
		if c != '~' {
			s := string(c)
			text = strings.ReplaceAll(text, s, `\`+s)
		}
	}
	return text
}
