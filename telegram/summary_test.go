package telegram

import (
	"testing"
)

func TestSummary(t *testing.T) {
	u := Update{}
	u.Message.Text = "/rangkuman"
	u.Message.MessageID = 1
	response, err := Respond(u)
	notErrOrNil(t, response, err)
	if response.Text != summaryMsg {
		t.Error("response text not match")
	}
	if response.ReplyToMessageId != 0 {
		t.Error("response should not reply")
	}
}
