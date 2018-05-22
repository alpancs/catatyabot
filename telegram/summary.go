package telegram

var (
	summaryMsg = "*== RANGKUMAN TOTAL BELANJA ==*"
)

func summary(u Update) (*Response, error) {
	res := DefaultResponse
	res.ChatID = u.Message.Chat.ID
	res.Text = summaryMsg
	return &res, nil
}
