package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	update, err := parseUpdate(r.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	go respondUpdateAsync(update)
}

func respondUpdateAsync(update *telegram.Update) {
	err := respondUpdate(update)
	if err != nil {
		fmt.Println(err)
	}
}

func parseUpdate(body io.Reader) (*telegram.Update, error) {
	var update telegram.Update
	err := json.NewDecoder(body).Decode(&update)
	return &update, err
}

func respondUpdate(u *telegram.Update) error {
	msg := u.Message
	if msg == nil {
		return nil
	}

	if right, err := commandInsert(msg); right {
		return err
	}
	if right, err := commandDelete(msg); right {
		return err
	}
	if right, err := commandList(msg); right {
		return err
	}
	if right, err := commandSummary(msg); right {
		return err
	}
	if right, err := commandDebug(msg); right {
		return err
	}

	if right, err := insert(msg); right {
		return err
	}
	if right, err := list(msg); right {
		return err
	}
	if right, err := update(msg); right {
		return err
	}

	if right, err := help(msg); right {
		return err
	}

	return nil
}
