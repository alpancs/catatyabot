package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	// tempe 2 ribu
	// sayur asem 1k
	patternCatat = regexp.MustCompile(`^.+ \d+(,\d+)?( *(ribu|rb|k|juta|jt))?$`)
	patternPrice = regexp.MustCompile(`\d+(,\d+)?( *(ribu|rb|k|juta|jt))?$`)
)

func commandCatat(msg *telegram.Message) (bool, error) {
	if msg.Command() == "catat" {
		resp, err := http.PostForm(sendMessageURL, url.Values{
			"chat_id":      {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":         {"apa saja yang ingin dicatat, bos?"},
			"reply_markup": {`{"force_reply": true}`},
		})

		if err == nil && resp.StatusCode >= 300 {
			fmt.Println("response code", resp.StatusCode)
			if resp.Body != nil {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("response body", body)
			}
		}
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
	if patternCatat.MatchString(text) {
		price := patternPrice.FindString(text)
		item := strings.TrimSpace(text[:len(text)-len(price)])

		resp, err := http.PostForm(sendMessageURL, url.Values{
			"chat_id": {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":    {fmt.Sprintf("%s dengan harga %s akan dicatat ya bos 👌", item, price)},
		})

		if err == nil && resp.StatusCode >= 300 {
			fmt.Println("response code", resp.StatusCode)
			if resp.Body != nil {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				fmt.Println("response body", body)
			}
		}
		return err
	}
	return nil
}
