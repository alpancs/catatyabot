package app

import (
	"encoding/json"
	"fmt"
	"net/url"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	ListText  = "pengen lihat daftar catatan yang mana bos? 👀"
	Today     = "hari ini"
	Yesterday = "kemarin"
	ThisWeek  = "pekan ini"
	PastWeek  = "pekan lalu"
	ThisMonth = "bulan ini"
	PastMonth = "bulan lalu"
)

var replyMarkupList = buildReplyMarkupList()

func buildReplyMarkupList() string {
	raw, err := json.Marshal(telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: Today}, {Text: Yesterday}},
			{{Text: ThisWeek}, {Text: PastWeek}},
			{{Text: ThisMonth}, {Text: PastMonth}},
		},
		ResizeKeyboard:  true,
		OneTimeKeyboard: true,
		Selective:       true,
	})
	if err != nil {
		panic(err)
	}
	return string(raw)
}

func commandList(msg *telegram.Message) (bool, error) {
	if msg.Command() != "lihat" {
		return false, nil
	}

	_, err := sendMessage(url.Values{
		"chat_id":             {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":                {ListText},
		"reply_to_message_id": {fmt.Sprintf("%d", msg.MessageID)},
		"reply_markup":        {replyMarkupList},
	})
	return true, err
}

func list(msg *telegram.Message) (bool, error) {
	if msg.ReplyToMessage == nil || msg.ReplyToMessage.Text != ListText {
		return false, nil
	}

	_, err := sendMessage(url.Values{
		"chat_id": {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":    {"masih kosong bos, kan cuma pura-pura nyatat"},
	})
	return true, err
}
