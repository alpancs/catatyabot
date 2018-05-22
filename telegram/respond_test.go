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
