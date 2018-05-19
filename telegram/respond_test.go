package telegram

import (
	"testing"
)

func TestRespondStart(t *testing.T) {
	u := Update{}
	u.Message.Text = "/start"
	response, err := Respond(u)
	notErrOrNil(t, response, err)
	if response.Text != StartMsg {
		t.Error("response text not match")
	}
	if response.ReplyToMessageId != 0 {
		t.Error("response should not reply")
	}
}

func TestRespondConfusing(t *testing.T) {
	u := Update{}
	u.Message.Text = "halo bot"
	u.Message.MessageID = 1
	response, err := Respond(u)
	notErrOrNil(t, response, err)
	if !contains(ConfusingMsgs, response.Text) {
		t.Error("response text not match")
	}
	if response.ReplyToMessageId == 0 {
		t.Error("response should reply")
	}
}

func notErrOrNil(t *testing.T, response *Response, err error) {
	if err != nil {
		t.Error(err)
	}
	if response == nil {
		t.Error("response should not nil")
	}
}

func contains(set []string, s string) bool {
	for _, e := range set {
		if e == s {
			return true
		}
	}
	return false
}
