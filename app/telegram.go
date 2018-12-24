package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var telegramAPI = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("TELEGRAM_BOT_TOKEN"))

func TelegramHandler(w http.ResponseWriter, r *http.Request) {
	update, err := parseUpdate(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Println(err)
		return
	}

	err = respond(update)
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

func respond(update *telegram.Update) error {
	return nil
}
