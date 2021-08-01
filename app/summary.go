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
	chanThisMonth := make(chan Price, 1)
	chanPastMonth := make(chan Price, 1)
	chanAllTime := make(chan Price, 1)

	go sumInterval(chatID, TextToday, chanToday, chanError)
	go sumInterval(chatID, TextYesterday, chanYesterday, chanError)
	go sumInterval(chatID, TextThisWeek, chanThisWeek, chanError)
	go sumInterval(chatID, TextThisMonth, chanThisMonth, chanError)
	go sumInterval(chatID, TextPastMonth, chanPastMonth, chanError)
	go sumInterval(chatID, TextAllTime, chanAllTime, chanError)

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
		TextToday, <-chanToday,
		TextYesterday, <-chanYesterday,
		TextThisWeek, <-chanThisWeek,
		TextThisMonth, <-chanThisMonth,
		TextPastMonth, <-chanPastMonth,
		TextAllTime, <-chanAllTime,
	), nil
}

func sumInterval(chatID int64, interval string, chanSum chan Price, chanError chan error) {
	s, err := querySum(chatID, interval)
	chanSum <- s
	chanError <- err
}

func querySum(chatID int64, interval string) (Price, error) {
	start, end := buildIntervalSQL(interval)
	query := "SELECT COALESCE(SUM(price), 0) FROM items WHERE chat_id = $1 AND $2 <= created_at AND created_at < $3;"
	var sum Price
	err := db.QueryRow(query, chatID, start, end).Scan(&sum)
	return sum, err
}
