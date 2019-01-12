package main

import "fmt"

type Price int64

func (p Price) String() string {
	return separate(fmt.Sprintf("%d", p))
}
func separate(s string) string {
	if len(s) <= 3 {
		return s
	}
	return fmt.Sprintf("%s.%s", separate(s[:len(s)-3]), s[len(s)-3:])
}
