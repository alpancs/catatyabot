package app

import (
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func commandSummary(msg *telegram.Message) (bool, error) {
	if msg.Command() != "rangkum" {
		return false, nil
	}

	summary, err := buildSummary(msg.Chat.ID)
	if err != nil {
		return true, err
	}

	_, err = sendMessage(msg.Chat.ID, summary, 0)
	return true, err
}

func buildSummary(chatID int64) (string, error) {
	chanError := make(chan error, 6)
	chanToday := make(chan Price, 1)
	chanYesterday := make(chan Price, 1)
	chanThisWeek := make(chan Price, 1)
	chanPastWeek := make(chan Price, 1)
	chanThisMonth := make(chan Price, 1)
	chanPastMonth := make(chan Price, 1)

	go sumInterval(chatID, Today, chanToday, chanError)
	go sumInterval(chatID, Yesterday, chanYesterday, chanError)
	go sumInterval(chatID, ThisWeek, chanThisWeek, chanError)
	go sumInterval(chatID, PastWeek, chanPastWeek, chanError)
	go sumInterval(chatID, ThisMonth, chanThisMonth, chanError)
	go sumInterval(chatID, PastMonth, chanPastMonth, chanError)

	for i := 0; i < 6; i++ {
		if err := <-chanError; err != nil {
			return "", err
		}
	}

	return fmt.Sprintf(`*=== RANGKUMAN ===*

- %s: %s
- %s: %s

- %s: %s
- %s: %s

- %s: %s
- %s: %s`,
		Today, <-chanToday,
		Yesterday, <-chanYesterday,
		ThisWeek, <-chanThisWeek,
		PastWeek, <-chanPastWeek,
		ThisMonth, <-chanThisMonth,
		PastMonth, <-chanPastMonth,
	), nil
}

func sumInterval(chatID int64, interval string, chanSum chan Price, chanError chan error) {
	s, err := execQuerySum(buildQuerySum(interval), chatID)
	chanSum <- s
	chanError <- err
}

func buildQuerySum(interval string) string {
	query := "SELECT COALESCE(SUM(price), 0) FROM items WHERE chat_id = $1 AND (%s) <= created_at AND created_at < (%s);"
	return fillQuery(query, interval)
}

func execQuerySum(query string, chatID int64) (Price, error) {
	var sum Price
	err := db.QueryRow(query, chatID).Scan(&sum)
	return sum, err
}
