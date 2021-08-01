package app

import (
	"errors"
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func update(msg *telegram.Message) (bool, error) {
	if msg.ReplyToMessage == nil {
		return false, nil
	}

	item, priceText, err := findOneItem(msg.Text)
	if err != nil {
		_, err = sendMessage(msg.Chat.ID, "kalau untuk _update_ catatan .. satu-satu ya, tidak bisa beberapa catatan sekaligus 😬", 0)
		return true, err
	}
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

	err = editMessage(msg.Chat.ID, msg.ReplyToMessage.MessageID, fmt.Sprintf(SaveTemplate, item, price))
	if err != nil {
		return true, err
	}
	_, err = sendMessage(msg.Chat.ID, "sudah diubah nih bos 👆", msg.ReplyToMessage.MessageID)
	return true, err
}

func findOneItem(text string) (item, priceText string, err error) {
	count := 0
	for _, line := range strings.Split(text, "\n") {
		priceText = patternPrice.FindString(line)
		item = strings.TrimSpace(line[:len(line)-len(priceText)])
		if item != "" && priceText != "" {
			count++
		}
	}
	if count > 1 {
		err = errors.New("found multiple items")
	}
	return
}
