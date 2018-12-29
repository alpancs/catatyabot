package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var apiURL = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("BOT_TOKEN"))

func respondUpdate(update *telegram.Update) error {
	replyText, err := json.MarshalIndent(update, "", " ")
	if err != nil {
		return err
	}

	if update.Message == nil {
		update.Message = update.EditedMessage
	}
	if update.Message == nil {
		return nil
	}

	resp, err := http.PostForm(apiURL+"sendMessage", url.Values{
		"chat_id":    {fmt.Sprintf("%d", update.Message.Chat.ID)},
		"text":       {"```\n" + string(replyText) + "\n```"},
		"parse_mode": {"Markdown"},
	})

	if resp.StatusCode >= 300 {
		fmt.Println("response code:", resp.StatusCode)
		fmt.Println("response body:")
		io.Copy(os.Stdout, resp.Body)
	}

	return err
}
