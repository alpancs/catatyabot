package app

import (
	"fmt"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	NewItemsText = `apa saja yang pengen dicatat, bos?

\*contoh:
_sayur kangkung 2 ribu_
_lombok 1/2 kg 3,5k_`
	SaveTemplate = "*%s %s* dicatat ya bos 👌 #catatan"
)

func commandInsert(msg *telegram.Message) (bool, error) {
	if msg.Command() != "catat" {
		return false, nil
	}

	_, err := sendMessageCustom(msg.Chat.ID, NewItemsText, 0, `{"force_reply":true}`)
	return true, err
}

func insert(msg *telegram.Message) (bool, error) {
	if msg.ReplyToMessage == nil || msg.ReplyToMessage.Text != NewItemsText {
		return false, nil
	}

	lines := strings.Split(msg.Text, "\n")
	numWorkers := len(lines)
	if numWorkers > 5 {
		numWorkers = 5
	}
	return true, insertByWorkers(numWorkers, msg, lines)
}

func insertByWorkers(numWorkers int, msg *telegram.Message, lines []string) error {
	chanText := make(chan string)
	chanError := make(chan error, len(lines))
	for i := 0; i < numWorkers; i++ {
		go worker(msg, chanText, chanError)
	}
	for _, text := range lines {
		chanText <- text
	}
	close(chanText)

	for range lines {
		err := <-chanError
		if err != nil {
			return err
		}
	}
	return nil
}

func worker(msg *telegram.Message, chanText <-chan string, chanError chan error) {
	for text := range chanText {
		chanError <- insertSpecificLine(msg, strings.TrimSpace(text))
	}
}

func insertSpecificLine(msg *telegram.Message, text string) error {
	priceText := patternPrice.FindString(text)
	item := strings.TrimSpace(text[:len(text)-len(priceText)])
	if item == "" || priceText == "" {
		return nil
	}

	price := ParsePrice(priceText)
	resp, err := sendMessage(msg.Chat.ID, fmt.Sprintf(SaveTemplate, item, price), 0)
	if err != nil {
		return err
	}

	_, err = db.Exec("INSERT INTO items VALUES ($1, $2, $3, $4);", resp.Chat.ID, resp.MessageID, item, price)
	if err != nil {
		revertReport(msg, resp, item, price)
	}
	return err
}

func revertReport(req, resp *telegram.Message, item string, price Price) {
	deleteMessage(resp.Chat.ID, resp.MessageID)
	sendMessage(req.Chat.ID, fmt.Sprintf("%s %s gagal dicatat bos 😔 #gagalmaningsonson", item, price), req.MessageID)
}
