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
	TextList      = "pengen lihat daftar catatan yang mana bos? 👀"
	TextToday     = "hari ini"
	TextYesterday = "kemarin"
	TextThisWeek  = "pekan ini"
	TextThisMonth = "bulan ini"
	TextPastMonth = "bulan lalu"
	TextAllTime   = "semuanya"
)

var (
	tzWIB      = time.FixedZone("WIB", 7*60*60)
	monthNames = []string{"Januari", "Februari", "Maret", "April", "Mei", "Juni",
		"Juli", "Agustus", "September", "Oktober", "November", "Desember"}
	replyMarkupList = buildReplyMarkupList()
	removeMarkup    = `{"remove_keyboard":true}`
)

func buildReplyMarkupList() string {
	raw, err := json.Marshal(telegram.ReplyKeyboardMarkup{
		Keyboard: [][]telegram.KeyboardButton{
			{{Text: TextToday}, {Text: TextYesterday}},
			{{Text: TextThisWeek}, {Text: TextThisMonth}},
			{{Text: TextPastMonth}, {Text: TextAllTime}},
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

	_, err := sendMessageCustom(msg.Chat.ID, TextList, 0, replyMarkupList)
	return true, err
}

func list(msg *telegram.Message) (bool, error) {
	if msg.ReplyToMessage.Text != TextList {
		return false, nil
	}

	start, end := buildIntervalSQL(strings.ToLower(msg.Text))
	if start.IsZero() && end.IsZero() {
		return false, nil
	}

	query := "SELECT name, price, created_at FROM items WHERE chat_id = $1 AND $2 <= created_at AND created_at < $3 ORDER BY created_at;"
	items, err := execQuerySelect(query, msg.Chat.ID, start, end)
	if err != nil {
		return true, err
	}

	_, err = sendMessageCustom(msg.Chat.ID, formatItems("catatan "+msg.Text, items), 0, removeMarkup)
	return true, err
}

func buildIntervalSQL(interval string) (time.Time, time.Time) {
	now := time.Now().In(tzWIB)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, tzWIB).In(time.UTC)
	tomorrow := today.AddDate(0, 0, 1)
	beginOfWeek := today.AddDate(0, 0, int(time.Sunday-now.Weekday()))
	beginOfMonth := today.AddDate(0, 0, 1-now.Day())
	switch interval {
	case TextToday:
		return today, tomorrow
	case TextYesterday:
		return today.AddDate(0, 0, -1), today
	case TextThisWeek:
		return beginOfWeek, tomorrow
	case TextThisMonth:
		return beginOfMonth, tomorrow
	case TextPastMonth:
		return beginOfMonth.AddDate(0, -1, 0), beginOfMonth
	case TextAllTime:
		return time.Time{}, tomorrow
	default:
		return time.Time{}, time.Time{}
	}
}

func execQuerySelect(query string, queryParams ...interface{}) ([]Item, error) {
	rows, err := db.Query(query, queryParams...)
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
		item.CreatedAt = item.CreatedAt.In(tzWIB)
		items = append(items, item)
	}
	return items, nil
}

func formatItems(title string, items []Item) string {
	text := fmt.Sprintf("*=== %s ===*\n", strings.ToUpper(title))
	sum := Price(0)
	for i, item := range items {
		if i == 0 || !(item.CreatedAt.Day() == items[i-1].CreatedAt.Day() && item.CreatedAt.Month() == items[i-1].CreatedAt.Month() && item.CreatedAt.Year() == items[i-1].CreatedAt.Year()) {
			text += fmt.Sprintf("\n%d %s %d\n", item.CreatedAt.Day(), monthNames[item.CreatedAt.Month()-time.January], item.CreatedAt.Year())
		}
		text += fmt.Sprintf("%02d:%02d - %s %s\n", item.CreatedAt.Hour(), item.CreatedAt.Minute(), item.Name, item.Price)
		sum += item.Price
	}
	return fmt.Sprintf("%s\n*TOTAL: %s*", text, sum)
}
