package app

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Item struct {
	Name      string
	Price     Price
	CreatedAt time.Time
}

const (
	ListText  = "pengen lihat daftar catatan yang mana bos? 👀"
	Today     = "hari ini"
	Yesterday = "kemarin"
	ThisWeek  = "pekan ini"
	PastWeek  = "pekan lalu"
	ThisMonth = "bulan ini"
	PastMonth = "bulan lalu"
)

var (
	monthNames      = []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni", "Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	replyMarkupList = buildReplyMarkupList()
)

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

	query := buildQuerySelect(msg.Text)
	if query == "" {
		return false, nil
	}

	items, err := execQuerySelect(query, msg.Chat.ID)
	if err != nil {
		return true, err
	}

	_, err = sendMessage(url.Values{
		"chat_id":      {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":         {formatItems("catatan "+msg.Text, items)},
		"parse_mode":   {"Markdown"},
		"reply_markup": {`{"remove_keyboard": true}`},
	})
	return true, err
}

func buildQuerySelect(interval string) string {
	query := "SELECT name, price, created_at FROM items WHERE chat_id = $1 AND (%s) <= created_at AND created_at < (%s) ORDER BY created_at;"

	today := "DATE_TRUNC('DAY', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')"
	tomorrow := today + " + INTERVAL '1 DAY'"
	yesterday := today + " - INTERVAL '1 DAY'"

	beginOfWeek := fmt.Sprintf("%s - INTERVAL '%d DAY'", today, time.Now().In(time.FixedZone("Asia/Jakarta", 7*60*60)).Weekday())
	beginOfPastWeek := beginOfWeek + " - INTERVAL '7 DAYS'"

	beginOfMonth := "DATE_TRUNC('MONTH', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')"
	beginOfPastMonth := beginOfMonth + " - INTERVAL '1 MONTH'"

	switch interval {
	case Today:
		return fmt.Sprintf(query, today, tomorrow)
	case Yesterday:
		return fmt.Sprintf(query, yesterday, today)
	case ThisWeek:
		return fmt.Sprintf(query, beginOfWeek, tomorrow)
	case PastWeek:
		return fmt.Sprintf(query, beginOfPastWeek, beginOfWeek)
	case ThisMonth:
		return fmt.Sprintf(query, beginOfMonth, tomorrow)
	case PastMonth:
		return fmt.Sprintf(query, beginOfPastMonth, beginOfMonth)
	default:
		return ""
	}
}

func execQuerySelect(query string, chatID int64) ([]Item, error) {
	rows, err := db.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	var items []Item
	for rows.Next() {
		var item Item
		err = rows.Scan(&item.Name, &item.Price, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func formatItems(title string, items []Item) string {
	text := fmt.Sprintf("*=== %s ===*\n", strings.ToUpper(title))
	sum := Price(0)
	lastDay := 0
	for _, item := range items {
		if day := item.CreatedAt.Day(); day != lastDay {
			text += fmt.Sprintf("\n_%d %s_\n", day, monthNames[item.CreatedAt.Month()-1])
			lastDay = day
		}
		text += fmt.Sprintf("- %s %s\n", item.Name, item.Price)
		sum += item.Price
	}
	return fmt.Sprintf("%s\n*TOTAL: %s*", text, sum)
}
