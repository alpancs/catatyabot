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
	editMessageURL = apiURL + "editMessageText"
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

func respondUpdate(u *telegram.Update) error {
	if u.Message == nil {
		return nil
	}

	if right, err := commandInsert(u.Message); right {
		return err
	}
	if right, err := bulkInsert(u.Message); right {
		return err
	}
	if right, err := update(u.Message); right {
		return err
	}

	return handleElse(u.Message)
}

func handleElse(msg *telegram.Message) error {
	_, err := sendMessage(url.Values{
		"chat_id": {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":    {"ngapain bos? 🙄"},
	})
	return err
}

func sendMessage(data url.Values) (*telegram.Message, error) {
	resp, err := http.PostForm(sendMessageURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		resp.Write(os.Stdout)
		return nil, nil
	}

	var respMsg telegram.Message
	err = json.NewDecoder(resp.Body).Decode(&respMsg)
	return &respMsg, err
}

func editMessage(data url.Values) error {
	resp, err := http.PostForm(editMessageURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		resp.Write(os.Stdout)
	}
	return nil
}
