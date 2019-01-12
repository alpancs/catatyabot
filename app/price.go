package main

import "strconv"

type Price int64

func (p Price) String() string {
	return separate(strconv.FormatInt(int64(p), 10))
}
func separate(s string) string {
	if len(s) <= 3 {
		return s
	}
	return separate(s[:len(s)-3]) + "." + s[len(s)-3:]
}
