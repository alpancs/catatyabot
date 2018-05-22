package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSample(t *testing.T) {
	assert.Panics(t, func() { sample([]string{}) }, "sampling empty slice should panic")
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
