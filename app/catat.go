package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	patternPrice    = regexp.MustCompile(`\d+(,\d+)?( *(ribu|rb|k|juta|jt))?$`)
	patternNumber   = regexp.MustCompile(`\d+(,\d+)?`)
	patternThousand = regexp.MustCompile(`ribu|rb|k`)
	patternMillion  = regexp.MustCompile(`juta|jt`)
)

func commandCatat(msg *telegram.Message) (bool, error) {
	if msg.Command() == "catat" {
		_, err := sendMessage(msg, url.Values{
			"chat_id":      {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":         {"apa saja yang ingin dicatat, bos?"},
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
	if item != "" && priceText != "" {
		price := parsePrice(priceText)
		_, err := sendMessage(msg, url.Values{
			"chat_id": {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":    {fmt.Sprintf("%s dengan harga %d ((pura-pura)) dicatat ya bos 👌", item, price)},
		})
		return err
	}
	return nil
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
