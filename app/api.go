package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	apiURL           = fmt.Sprintf("https://api.telegram.org/bot%s/", os.Getenv("BOT_TOKEN"))
	sendMessageURL   = apiURL + "sendMessage"
	editMessageURL   = apiURL + "editMessageText"
	deleteMessageURL = apiURL + "deleteMessage"
)

func sendMessage(chatID int64, text string, replyToMessageID int) (*telegram.Message, error) {
	return sendMessageCustom(chatID, text, replyToMessageID, "")
}
func sendMessageCustom(chatID int64, text string, replyToMessageID int, replyMarkup string) (*telegram.Message, error) {
	resp, err := http.PostForm(sendMessageURL, url.Values{
		"chat_id":             {strconv.FormatInt(chatID, 10)},
		"text":                {text},
		"parse_mode":          {"MarkdownV2"},
		"reply_to_message_id": {strconv.Itoa(replyToMessageID)},
		"reply_markup":        {replyMarkup},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		resp.Write(os.Stdout)
		return nil, nil
	}

	var respAPI telegram.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&respAPI)
	if err != nil {
		return nil, err
	}
	if !respAPI.Ok {
		return nil, fmt.Errorf("error_code: %d, description: %s", respAPI.ErrorCode, respAPI.Description)
	}
	var respMsg telegram.Message
	err = json.Unmarshal(respAPI.Result, &respMsg)
	return &respMsg, err
}

func updateMessage(url string, data url.Values) error {
	resp, err := http.PostForm(url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		resp.Write(os.Stdout)
	}
	return nil
}

func editMessage(chatID int64, messageID int, text string) error {
	return updateMessage(editMessageURL, url.Values{
		"chat_id":    {strconv.FormatInt(chatID, 10)},
		"message_id": {strconv.Itoa(messageID)},
		"text":       {text},
		"parse_mode": {"MarkdownV2"},
	})
}

func deleteMessage(chatID int64, messageID int) error {
	return updateMessage(deleteMessageURL, url.Values{
		"chat_id":    {strconv.FormatInt(chatID, 10)},
		"message_id": {strconv.Itoa(messageID)},
	})
}
