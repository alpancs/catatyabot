package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var apiURL = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("BOT_TOKEN"))

func respondUpdate(update *telegram.Update) error {
	resp, err := json.MarshalIndent(update, "", " ")
	if err != nil {
		return err
	}

	_, err = http.PostForm(apiURL+"sendMessage", url.Values{
		"chat_id":    {fmt.Sprintf("%d", update.Message.Chat.ID)},
		"text":       {"```" + string(resp) + "```"},
		"parse_mode": "Markdown",
	})

	return err
}
