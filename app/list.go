package app

import (
	"encoding/json"
	"fmt"
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
	monthNames = []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
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

	_, err := sendMessageCustom(msg.Chat.ID, ListText, 0, replyMarkupList)
	return true, err
}

func list(msg *telegram.Message) (bool, error) {
	start, end := buildIntervalSQL(strings.ToLower(msg.Text))
	if start == "" || end == "" {
		return false, nil
	}

	query := fmt.Sprintf("SELECT name, price, created_at FROM items WHERE chat_id = $1 AND (%s) <= created_at AND created_at < (%s) ORDER BY created_at;", start, end)
	items, err := execQuerySelect(query, msg.Chat.ID)
	if err != nil {
		return true, err
	}

	_, err = sendMessage(msg.Chat.ID, formatItems("catatan "+msg.Text, items), 0)
	return true, err
}

func buildIntervalSQL(interval string) (string, string) {
	timeToday := "DATE_TRUNC('DAY', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')"
	timeTomorrow := timeToday + " + INTERVAL '1 DAY'"
	timeYesterday := timeToday + " - INTERVAL '1 DAY'"
	timeBeginOfWeek := fmt.Sprintf("%s - INTERVAL '%d DAY'", timeToday, time.Now().In(time.FixedZone("Asia/Jakarta", 7*60*60)).Weekday())
	timeBeginOfPastWeek := timeBeginOfWeek + " - INTERVAL '7 DAYS'"
	timeBeginOfMonth := "DATE_TRUNC('MONTH', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Jakarta')"
	timeBeginOfPastMonth := timeBeginOfMonth + " - INTERVAL '1 MONTH'"
	switch interval {
	case Today:
		return timeToday, timeTomorrow
	case Yesterday:
		return timeYesterday, timeToday
	case ThisWeek:
		return timeBeginOfWeek, timeTomorrow
	case PastWeek:
		return timeBeginOfPastWeek, timeBeginOfWeek
	case ThisMonth:
		return timeBeginOfMonth, timeTomorrow
	case PastMonth:
		return timeBeginOfPastMonth, timeBeginOfMonth
	}
	return "", ""
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
