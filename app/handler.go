package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
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
