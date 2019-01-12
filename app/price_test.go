package main

import "testing"

func assertEqual(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("expected %s, actual %s", expected, actual)
	}
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
