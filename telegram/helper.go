package telegram

import (
	"math/rand"
)

func sample(set []string) string {
	return set[rand.Int()%len(set)]
}
