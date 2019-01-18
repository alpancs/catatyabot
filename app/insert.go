package app

import (
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	NewItemsText = "apa saja yang pengen dicatat, bos?"
	SaveTemplate = "*%s %s* dicatat ya bos 👌 #catatan"
)

func commandInsert(msg *telegram.Message) (bool, error) {
	if msg.Command() != "catat" {
		return false, nil
	}

	_, err := sendMessageCustom(msg.Chat.ID, NewItemsText, msg.MessageID, `{"force_reply":true,"selective":true}`)
	return true, err
}

func insert(msg *telegram.Message) (bool, error) {
	if msg.ReplyToMessage == nil || msg.ReplyToMessage.Text != NewItemsText {
		return false, nil
	}

	lines := strings.Split(msg.Text, "\n")
	errs := make(chan error, len(lines))
	for _, text := range lines {
		go insertSpecificLine(msg, strings.TrimSpace(text), errs)
	}
	for range lines {
		err := <-errs
		if err != nil {
			return true, err
		}
	}
	return true, nil
}

func insertSpecificLine(msg *telegram.Message, text string, errs chan error) {
	priceText := patternPrice.FindString(text)
	item := strings.TrimSpace(text[:len(text)-len(priceText)])
	if item == "" || priceText == "" {
		errs <- nil
		return
	}

	price := ParsePrice(priceText)
	resp, err := sendMessage(msg.Chat.ID, fmt.Sprintf(SaveTemplate, item, price), 0)
	if err != nil {
		errs <- err
		return
	}

	_, err = db.Exec("INSERT INTO items VALUES ($1, $2, $3, $4);", resp.Chat.ID, resp.MessageID, item, price)
	errs <- err
	if err != nil {
		revertReport(msg, resp, item, price)
	}
}

func revertReport(req, resp *telegram.Message, item string, price Price) {
	deleteMessage(resp.Chat.ID, resp.MessageID)
	sendMessage(req.Chat.ID, fmt.Sprintf("%s %s gagal dicatat bos 😔 #gagalmaningsonson", item, price), req.MessageID)
}
