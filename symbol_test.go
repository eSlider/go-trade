package trade

import (
	"encoding/json"
	"testing"
)

func sampleSymbols() Symbols {
	return Symbols{
		{ID: 1, Type: SymbolFiat, ParentID: 0, Code: "USD", Name: "US Dollar"},
		{ID: 2, Type: SymbolCrypto, ParentID: 1, Code: "USDT", Name: "Tether"},
		{ID: 3, Type: SymbolCrypto, ParentID: 1, Code: "USDC", Name: "USD Coin"},
		{ID: 4, Type: SymbolFiat, ParentID: 1, Code: "EUR", Name: "Euro"},
		{ID: 5, Type: SymbolFiat, ParentID: 1, Code: "GBP", Name: "British Pound"},
		{ID: 14, Type: SymbolCrypto, ParentID: 0, Code: "BTC", Name: "Bitcoin", Website: "https://bitcoin.org"},
		{ID: 22, Type: SymbolCrypto, ParentID: 14, Code: "BTCFT", Name: "BTC Futures Token"},
		{ID: 23, Type: SymbolCrypto, ParentID: 14, Code: "BTCM24", Name: "BTC June 2024 Future"},
	}
}

func TestSymbolTypeString(t *testing.T) {
	tests := []struct {
		st   SymbolType
		want string
	}{
		{SymbolFiat, "fiat"},
		{SymbolCrypto, "crypto"},
		{SymbolType(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.st.String(); got != tt.want {
			t.Errorf("SymbolType(%d).String() = %s, want %s", tt.st, got, tt.want)
		}
	}
}

func TestSymbolIsRoot(t *testing.T) {
	ss := sampleSymbols()
	usd := ss.GetByCode("USD")
	if !usd.IsRoot() {
		t.Error("USD should be root")
	}
	usdt := ss.GetByCode("USDT")
	if usdt.IsRoot() {
		t.Error("USDT should not be root")
	}
}

func TestSymbolIsFiatCrypto(t *testing.T) {
	ss := sampleSymbols()
	usd := ss.GetByCode("USD")
	if !usd.IsFiat() || usd.IsCrypto() {
		t.Error("USD should be fiat")
	}
	btc := ss.GetByCode("BTC")
	if !btc.IsCrypto() || btc.IsFiat() {
		t.Error("BTC should be crypto")
	}
}

func TestSymbolsGetByCode(t *testing.T) {
	ss := sampleSymbols()
	if s := ss.GetByCode("BTC"); s == nil || s.Name != "Bitcoin" {
		t.Errorf("GetByCode(BTC) = %v, want Bitcoin", s)
	}
	if s := ss.GetByCode("NONEXISTENT"); s != nil {
		t.Errorf("GetByCode(NONEXISTENT) = %v, want nil", s)
	}
}

func TestSymbolsGetByID(t *testing.T) {
	ss := sampleSymbols()
	if s := ss.GetByID(14); s == nil || s.Code != "BTC" {
		t.Errorf("GetByID(14) = %v, want BTC", s)
	}
	if s := ss.GetByID(999); s != nil {
		t.Errorf("GetByID(999) = %v, want nil", s)
	}
}

func TestSymbolsChildren(t *testing.T) {
	ss := sampleSymbols()

	// USD children: USDT, USDC, EUR, GBP
	usdChildren := ss.Children(1)
	if len(usdChildren) != 4 {
		t.Errorf("USD children = %d, want 4", len(usdChildren))
	}

	// BTC children: BTCFT, BTCM24
	btcChildren := ss.Children(14)
	if len(btcChildren) != 2 {
		t.Errorf("BTC children = %d, want 2", len(btcChildren))
	}

	// No children
	leafChildren := ss.Children(22)
	if len(leafChildren) != 0 {
		t.Errorf("BTCFT children = %d, want 0", len(leafChildren))
	}
}

func TestSymbolsRoots(t *testing.T) {
	ss := sampleSymbols()
	roots := ss.Roots()
	if len(roots) != 2 {
		t.Fatalf("Roots() = %d, want 2 (USD, BTC)", len(roots))
	}
	codes := map[string]bool{}
	for _, r := range roots {
		codes[r.Code] = true
	}
	if !codes["USD"] || !codes["BTC"] {
		t.Errorf("Roots should be USD and BTC, got %v", codes)
	}
}

func TestSymbolsFiats(t *testing.T) {
	ss := sampleSymbols()
	fiats := ss.Fiats()
	if len(fiats) != 3 {
		t.Errorf("Fiats() = %d, want 3 (USD, EUR, GBP)", len(fiats))
	}
	for _, f := range fiats {
		if !f.IsFiat() {
			t.Errorf("Fiats() contains non-fiat: %s", f.Code)
		}
	}
}

func TestSymbolsCryptos(t *testing.T) {
	ss := sampleSymbols()
	cryptos := ss.Cryptos()
	if len(cryptos) != 5 {
		t.Errorf("Cryptos() = %d, want 5", len(cryptos))
	}
	for _, c := range cryptos {
		if !c.IsCrypto() {
			t.Errorf("Cryptos() contains non-crypto: %s", c.Code)
		}
	}
}

func TestSymbolJSON(t *testing.T) {
	s := Symbol{
		ID: 14, Type: SymbolCrypto, Code: "BTC",
		Name: "Bitcoin", Website: "https://bitcoin.org",
	}
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}

	var parsed Symbol
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}
	if parsed.Code != "BTC" || parsed.Type != SymbolCrypto {
		t.Errorf("Round-trip failed: %+v", parsed)
	}
}
