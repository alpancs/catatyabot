package main

import (
	"fmt"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var apiURL = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("BOT_TOKEN"))

func respondUpdate(update *telegram.Update) error {
	return nil
}
