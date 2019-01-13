package main

import (
	"strconv"
	"strings"
)

type Price int64

func ParsePrice(text string) Price {
	num, _ := strconv.ParseFloat(strings.Replace(patternNumber.FindString(text), ",", ".", 1), 64)
	switch {
	case patternThousand.MatchString(text):
		return Price(num * 1000)
	case patternMillion.MatchString(text):
		return Price(num * 1000 * 1000)
	default:
		return Price(num)
	}
}

func (p Price) String() string {
	return separate(strconv.FormatInt(int64(p), 10))
}
func separate(s string) string {
	if len(s) <= 3 {
		return s
	}
	return separate(s[:len(s)-3]) + "." + s[len(s)-3:]
}
