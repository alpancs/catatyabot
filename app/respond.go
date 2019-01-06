package main

import (
	"fmt"
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
	if right, err := catat(update.Message); right {
		return err
	}

	return nil
}
