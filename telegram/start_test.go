package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	u := Update{}
	u.Message.Text = "/start"
	response, err := Respond(u)

	assert.NoError(t, err, "should not error")
	assert.NotNil(t, response, "response should not nil")

	assert.Equal(t, startMsg, response.Text, "response text not match")
	assert.Equal(t, 0, response.ReplyToMessageId, "response should not reply")
}
