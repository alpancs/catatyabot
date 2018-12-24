package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var apiURL = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("BOT_TOKEN"))

func respondUpdate(update *telegram.Update) error {
	_, err := http.PostForm(apiURL+"sendMessage", url.Values{
		"chat_id": {fmt.Sprintf("%d", update.Message.Chat.ID)},
		"text":    {fmt.Sprintf("%#v", update)},
	})
	return err
}
