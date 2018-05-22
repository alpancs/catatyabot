package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummary(t *testing.T) {
	u := Update{}
	u.Message.Text = "/rangkuman"
	u.Message.MessageID = 1
	response, err := Respond(u)

	assert.NoError(t, err, "should not error")
	assert.NotNil(t, response, "response should not nil")

	assert.Equal(t, summaryMsg, response.Text, "response text not match")
	assert.Equal(t, 0, response.ReplyToMessageId, "response should not reply")
}
