package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfusing(t *testing.T) {
	u := Update{}
	u.Message.Text = "halo bot"
	u.Message.MessageID = 1
	response, err := Respond(u)

	assert.NoError(t, err, "should not error")
	assert.NotNil(t, response, "response should not nil")

	assert.Contains(t, confusingMsgs, response.Text, "response text not match")
	assert.NotEqual(t, 0, response.ReplyToMessageId, "response should reply")
}
