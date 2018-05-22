package telegram

import (
	"testing"
)

func TestStart(t *testing.T) {
	u := Update{}
	u.Message.Text = "/start"
	response, err := Respond(u)
	notErrOrNil(t, response, err)
	if response.Text != startMsg {
		t.Error("response text not match")
	}
	if response.ReplyToMessageId != 0 {
		t.Error("response should not reply")
	}
}
