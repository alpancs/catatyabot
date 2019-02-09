package app

import (
	"fmt"
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

	_, err = sendMessage(msg.Chat.ID, formatItemsDebug("rincian 30 catatan terakhir", items), 0)
	return true, err
}

func formatItemsDebug(title string, items []Item) string {
	text := fmt.Sprintf("*=== %s ===*\n", strings.ToUpper(title))
	sum := Price(0)
	for i, item := range items {
		if i == 0 || item.CreatedAt.Day() != items[i-1].CreatedAt.Day() {
			text += fmt.Sprintf("\n_%d %s_\n", item.CreatedAt.Day(), monthNames[item.CreatedAt.Month()-time.January])
		}
		text += fmt.Sprintf("- %s %s %s\n", item.Name, item.Price, item.CreatedAt.Format(time.RFC3339))
		sum += item.Price
	}
	return fmt.Sprintf("%s\n*TOTAL: %s*", text, sum)
}
