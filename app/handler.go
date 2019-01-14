package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
	if right, err := commandDelete(u.Message); right {
		return err
	}
	if right, err := commandList(u.Message); right {
		return err
	}
	if right, err := commandSummary(u.Message); right {
		return err
	}

	if right, err := insert(u.Message); right {
		return err
	}
	if right, err := list(u.Message); right {
		return err
	}
	if right, err := update(u.Message); right {
		return err
	}

	return handleElse(u.Message)
}

func handleElse(msg *telegram.Message) error {
	_, err := sendMessage(url.Values{
		"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":                {"ngapain bos? 🙄"},
		"reply_to_message_id": {fmt.Sprintf("%d", msg.MessageID)},
	})
	return err
}
