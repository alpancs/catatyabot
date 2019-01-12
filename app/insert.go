package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	NewNoteText = "apa saja yang ingin dicatat, bos?"
)

var (
	patternPrice    = regexp.MustCompile(` \d+(,\d+)?( *(ribu|rb|k|juta|jt))?$`)
	patternNumber   = regexp.MustCompile(`\d+(,\d+)?`)
	patternThousand = regexp.MustCompile(`ribu|rb|k`)
	patternMillion  = regexp.MustCompile(`juta|jt`)
)

func commandInsert(msg *telegram.Message) (bool, error) {
	if msg.Command() != "catat" {
		return false, nil
	}

	_, err := sendMessage(url.Values{
		"chat_id":      {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":         {NewNoteText},
		"reply_markup": {`{"force_reply": true}`},
	})
	return true, err
}

func insert(msg *telegram.Message) (bool, error) {
	if msg.ReplyToMessage == nil || msg.ReplyToMessage.Text != NewNoteText {
		return false, nil
	}

	for _, text := range strings.Split(msg.Text, "\n") {
		if err := insertSpecificLine(msg, strings.TrimSpace(text)); err != nil {
			return true, err
		}
	}
	return true, nil
}

func insertSpecificLine(msg *telegram.Message, text string) error {
	priceText := patternPrice.FindString(text)
	item := strings.TrimSpace(text[:len(text)-len(priceText)])
	if item == "" || priceText == "" {
		return nil
	}
	price := parsePrice(priceText)

	resp, err := sendMessage(url.Values{
		"chat_id": {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":    {fmt.Sprintf("%s dengan harga %s ((pura-pura)) dicatat ya bos 👌", item, price)},
	})
	fmt.Printf("%+v\n", resp)
	return err
}

func parsePrice(text string) Price {
	num, _ := strconv.ParseFloat(strings.Replace(patternNumber.FindString(text), ",", ".", 1), 64)
	switch {
	case patternThousand.MatchString(text):
		return Price(num * 1000)
	case patternMillion.MatchString(text):
		return Price(num * 1000 * 1000)
	default:
		return Price(num)
	}
}
