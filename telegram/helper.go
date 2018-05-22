package telegram

import (
	"math/rand"
)

func sample(set []string) string {
	if len(set) == 0 {
		return ""
	}
	return set[rand.Int()%len(set)]
}
