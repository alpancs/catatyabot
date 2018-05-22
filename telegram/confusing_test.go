package telegram

import (
	"testing"
)

func TestConfusing(t *testing.T) {
	u := Update{}
	u.Message.Text = "halo bot"
	u.Message.MessageID = 1
	response, err := Respond(u)
	notErrOrNil(t, response, err)
	if !contains(confusingMsgs, response.Text) {
		t.Error("response text not match")
	}
	if response.ReplyToMessageId == 0 {
		t.Error("response should reply")
	}
}
