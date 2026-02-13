package currency

import (
	"testing"
)

func TestNew(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	if len(p.Currencies) == 0 {
		t.Fatal("New() returned empty currencies")
	}
	if len(p.Units) == 0 {
		t.Fatal("New() returned empty units")
	}
}

func TestGetCurrency(t *testing.T) {
	p, err := New()
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	btc := p.Currencies.Get("BTC")
	if btc == nil {
		t.Fatal("Get(BTC) returned nil")
	}
	if btc.Description != "Bitcoin" {
		t.Errorf("BTC description = %q, want Bitcoin", btc.Description)
	}
	if !btc.IsCrypto() {
		t.Error("BTC.IsCrypto() = false, want true")
	}

	usd := p.Currencies.Get("USD")
	if usd == nil {
		t.Fatal("Get(USD) returned nil")
	}
	if !usd.IsFiat() {
		t.Error("USD.IsFiat() = false, want true")
	}

	unknown := p.Currencies.Get("ZZZZZ")
	if unknown != nil {
		t.Error("Get(ZZZZZ) should return nil")
	}
}

func TestCryptos(t *testing.T) {
	p, _ := New()
	cryptos := p.Currencies.Cryptos()
	if len(cryptos) == 0 {
		t.Fatal("Cryptos() returned empty list")
	}
	for _, c := range cryptos {
		if c.Kind != Crypto {
			t.Errorf("Cryptos() contains non-crypto: %s (%s)", c.Code, c.Kind)
		}
	}
}

func TestFiats(t *testing.T) {
	p, _ := New()
	fiats := p.Currencies.Fiats()
	if len(fiats) == 0 {
		t.Fatal("Fiats() returned empty list")
	}
	for _, c := range fiats {
		if c.Kind != Fiat {
			t.Errorf("Fiats() contains non-fiat: %s (%s)", c.Code, c.Kind)
		}
	}
}
