package app

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func commandDebug(msg *telegram.Message) (bool, error) {
	if msg.Command() != "rinci" {
		return false, nil
	}

	query := "SELECT name, price, created_at FROM items WHERE chat_id = $1 ORDER BY created_at DESC LIMIT 30;"

	items, err := execQuerySelect(query, msg.Chat.ID)
	if err != nil {
		return true, err
	}

	_, err = sendMessage(url.Values{
		"chat_id":    {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":       {formatItemsDebug("rincian 30 catatan terakhir", items)},
		"parse_mode": {"Markdown"},
	})
	return true, err
}

func formatItemsDebug(title string, items []Item) string {
	text := fmt.Sprintf("*=== %s ===*\n", strings.ToUpper(title))
	sum := Price(0)
	lastDay := 0
	for _, item := range items {
		if day := item.CreatedAt.Day(); day != lastDay {
			text += fmt.Sprintf("\n_%d %s_\n", day, monthNames[item.CreatedAt.Month()-1])
			lastDay = day
		}
		text += fmt.Sprintf("- %s %s %s\n", item.Name, item.Price, item.CreatedAt.Format(time.RFC3339))
		sum += item.Price
	}
	return fmt.Sprintf("%s\n*TOTAL: %s*", text, sum)
}
