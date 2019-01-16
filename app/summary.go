package app

import (
	"fmt"
	"net/url"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func commandSummary(msg *telegram.Message) (bool, error) {
	if msg.Command() != "rangkuman" {
		return false, nil
	}

	summary, err := sum(msg.Chat.ID)
	if err != nil {
		return true, err
	}

	_, err = sendMessage(url.Values{
		"chat_id":    {fmt.Sprintf("%d", msg.Chat.ID)},
		"text":       {summary},
		"parse_mode": {"Markdown"},
	})
	return true, err
}

func sum(chatID int64) (string, error) {
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

	return fmt.Sprintf(`*==== RANGKUMAN ====*

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
	query := "SELECT SUM(price) FROM items WHERE chat_id = $1 AND (%s) <= created_at AND created_at < (%s);"

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

func execQuerySum(query string, chatID int64) (Price, error) {
	var sum *Price
	err := db.QueryRow(query, chatID).Scan(&sum)
	if sum == nil {
		return Price(0), err
	}
	return *sum, err
}
