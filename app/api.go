package main

import (
	"encoding/json"
	"fmt"
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
