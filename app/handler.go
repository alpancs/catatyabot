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

var (
	apiURL         = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("BOT_TOKEN"))
	sendMessageURL = apiURL + "sendMessage"
)

func handler(w http.ResponseWriter, r *http.Request) {
	update, err := parseUpdate(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	err = respondUpdate(update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println(err)
	}
}

func parseUpdate(body io.ReadCloser) (*telegram.Update, error) {
	defer body.Close()
	var update telegram.Update
	err := json.NewDecoder(body).Decode(&update)
	return &update, err
}

func respondUpdate(update *telegram.Update) error {
	if update.Message == nil {
		return nil
	}

	if right, err := commandCatat(update.Message); right {
		return err
	}
	if right, err := catat(update.Message); right {
		return err
	}

	return nil
}

func sendMessage(msg *telegram.Message, data url.Values) (*telegram.Message, error) {
	resp, err := http.PostForm(sendMessageURL, url.Values{
		"chat_id":      {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":         {"apa saja yang ingin dicatat, bos?"},
		"reply_markup": {`{"force_reply": true}`},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		resp.Write(os.Stdout)
		return nil, err
	}

	var respMsg telegram.Message
	err = json.NewDecoder(resp.Body).Decode(&respMsg)
	return &respMsg, err
}
