package telegram

var (
	confusingMsgs = []string{
		"ngomong apa to bos? 🤔",
		"mbuh bos, gak ngerti 😒",
		"aku orak paham boooss 😔",
	}
)

func confusing(u Update) (*Response, error) {
	res := DefaultResponse
	res.ChatID = u.Message.Chat.ID
	res.Text = sample(confusingMsgs)
	res.ReplyToMessageId = u.Message.MessageID
	return &res, nil
}
