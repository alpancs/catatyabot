package telegram

import (
	"os"
)

type Update struct {
	Message struct {
		MessageID int `json:"message_id"`
		Chat      struct {
			ID int
		}
		Text string
	}
}

type Response struct {
	ChatID           int    `json:"chat_id"`
	Text             string `json:"text"`
	ParseMode        string `json:"parse_mode,omitempty"`
	ReplyToMessageId int    `json:"reply_to_message_id,omitempty"`
}

var (
	Username        = os.Getenv("TELEGRAM_BOT_USERNAME")
	DefaultResponse = Response{ParseMode: "Markdown"}
)

func Respond(u Update) (*Response, error) {
	switch u.Message.Text {
	case "/start":
		return responseStart(u)
	default:
		return responseConfusing(u)
	}
}
