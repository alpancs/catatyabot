package telegram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSample(t *testing.T) {
	assert.NotPanics(t, func() { sample([]string{}) }, "sampling empty slice should not panic")
}
