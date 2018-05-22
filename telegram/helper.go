package telegram

import (
	"math/rand"
	"testing"
)

func sample(set []string) string {
	return set[rand.Int()%len(set)]
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
