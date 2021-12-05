package app

import (
	"regexp"
	"strconv"
	"strings"
)

type Price int64

var (
	patternPrice    = regexp.MustCompile(` -?\d+([,.]\d+)?( *(ribu|rb|k|juta|jt))?$`)
	patternNumber   = regexp.MustCompile(`-?\d+([,.]\d+)?`)
	patternThousand = regexp.MustCompile(`ribu|rb|k`)
	patternMillion  = regexp.MustCompile(`juta|jt`)
)

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
	if (s[0] == '-' && len(s) <= 4) || len(s) <= 3 {
		return s
	}
	return separate(s[:len(s)-3]) + "." + s[len(s)-3:]
}
