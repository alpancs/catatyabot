package main

import "testing"

func TestPriceString0(t *testing.T) {
	p := Price(0)
	if p.String() != "0" {
		t.Error("expected 0, got", p)
	}
}

func TestPriceString12(t *testing.T) {
	p := Price(12)
	if p.String() != "12" {
		t.Error("expected 12, got", p)
	}
}

func TestPriceString2500(t *testing.T) {
	p := Price(2500)
	if p.String() != "2.500" {
		t.Error("expected 2500, got", p)
	}
}

func TestPriceString1234567(t *testing.T) {
	p := Price(1234567)
	if p.String() != "1.234.567" {
		t.Error("expected 1234567, got", p)
	}
}
