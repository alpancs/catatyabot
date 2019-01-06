package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	apiURL         = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("BOT_TOKEN"))
	sendMessageURL = apiURL + "sendMessage"
)

func respondUpdate(update *telegram.Update) error {
	if update.Message == nil {
		return nil
	}

	if right, err := commandCatat(update.Message); right {
		return err
	}

	return nil
}

func commandCatat(msg *telegram.Message) (bool, error) {
	if msg.Command() == "catat" {
		resp, err := http.PostForm(sendMessageURL, url.Values{
			"chat_id":      {fmt.Sprintf("%d", msg.Chat.ID)},
			"text":         {"apa saja yang ingin dicatat?"},
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
