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

func commandCatat(msg *telegram.Message) (bool, error) {
	if msg.Command() == "catat" {
		_, err := sendMessage(msg, url.Values{
			"chat_id":      {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":         {NewNoteText},
			"reply_markup": {`{"force_reply": true}`},
		})
		return true, err
	}
	return false, nil
}

func catat(msg *telegram.Message) (bool, error) {
	for _, text := range strings.Split(msg.Text, "\n") {
		if err := catatText(msg, strings.TrimSpace(text)); err != nil {
			return true, err
		}
	}
	return true, nil
}

func catatText(msg *telegram.Message, text string) error {
	priceText := patternPrice.FindString(text)
	item := strings.TrimSpace(text[:len(text)-len(priceText)])
	if item == "" || priceText == "" {
		return nil
	}

	price := parsePrice(priceText)
	if msg.ReplyToMessage.Text == NewNoteText {
		resp, err := sendMessage(msg, url.Values{
			"chat_id": {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":    {fmt.Sprintf("%s dengan harga %d ((pura-pura)) dicatat ya bos 👌", item, price)},
		})
		fmt.Println(resp)
		return err
	} else {
		return editMessage(msg, url.Values{
			"chat_id":    {fmt.Sprintf("%d", msg.Chat.ID)},
			"message_id": {fmt.Sprintf("%d", msg.ReplyToMessage.MessageID)},
			"text":       {fmt.Sprintf("%s dengan harga %d ((pura-pura)) dicatat ya bos 👌", item, price)},
		})
	}
}

func parsePrice(text string) int64 {
	num, _ := strconv.ParseFloat(strings.Replace(patternNumber.FindString(text), ",", ".", 1), 64)
	switch {
	case patternThousand.MatchString(text):
		return int64(num * 1000)
	case patternMillion.MatchString(text):
		return int64(num * 1000 * 1000)
	default:
		return int64(num)
	}
}
