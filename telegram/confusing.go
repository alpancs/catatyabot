package telegram

var (
	ConfusingMsgs = []string{
		"ngomong apa to bos? 🤔",
		"mbuh bos, gak ngerti 😒",
		"aku orak paham boooss 😔",
	}
)

func responseConfusing(u Update) (*Response, error) {
	res := DefaultResponse
	res.ChatID = u.Message.Chat.ID
	res.Text = sample(ConfusingMsgs)
	res.ReplyToMessageId = u.Message.MessageID
	return &res, nil
}
