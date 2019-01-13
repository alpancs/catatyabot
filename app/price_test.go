package main

import "testing"

func assertEqual(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected %v, actual %v", expected, actual)
	}
}

func TestParsePrice0(t *testing.T) {
	assertEqual(t, Price(0), ParsePrice("0"))
}
func TestParsePrice25ribu(t *testing.T) {
	assertEqual(t, Price(25000), ParsePrice("25rb"))
}
func TestParsePrice1200(t *testing.T) {
	assertEqual(t, Price(1200), ParsePrice("1,2 k"))
}
func TestParsePrice5000000(t *testing.T) {
	assertEqual(t, Price(5000000), ParsePrice("5 juta"))
}
func TestParsePriceNotValid(t *testing.T) {
	assertEqual(t, Price(0), ParsePrice("not valid"))
}

func TestPriceString0(t *testing.T) {
	assertEqual(t, "0", Price(0).String())
}
func TestPriceString12(t *testing.T) {
	assertEqual(t, "12", Price(12).String())
}
func TestPriceString2500(t *testing.T) {
	assertEqual(t, "2.500", Price(2500).String())
}
func TestPriceString1234567(t *testing.T) {
	assertEqual(t, "1.234.567", Price(1234567).String())
}

func BenchmarkPriceString1234567(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Price(1234567).String()
	}
	b.ReportAllocs()
}
