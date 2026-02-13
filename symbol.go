package trade

// SymbolType classifies a trading symbol as fiat or cryptocurrency.
type SymbolType int8

const (
	SymbolFiat   SymbolType = iota // Fiat currency (USD, EUR, GBP, ...)
	SymbolCrypto                   // Cryptocurrency (BTC, ETH, SOL, ...)
)

// String returns a human-readable label for the symbol type.
func (s SymbolType) String() string {
	switch s {
	case SymbolFiat:
		return "fiat"
	case SymbolCrypto:
		return "crypto"
	default:
		return "unknown"
	}
}

// Symbol represents a market trading symbol (e.g. BTC, USDT, EUR).
//
// Symbols form a tree structure through ParentID, enabling hierarchical
// relationships between base assets and their derivatives:
//
//	USD  (id:1,  parent:0, fiat)
//	├── USDT  (id:2, parent:1, crypto)   — stablecoin pegged to USD
//	├── USDC  (id:3, parent:1, crypto)   — stablecoin pegged to USD
//	├── EUR   (id:4, parent:1, fiat)
//	└── GBP   (id:5, parent:1, fiat)
//	BTC  (id:14, parent:0, crypto)
//	├── BTCFT  (id:22, parent:14, crypto) — BTC futures token
//	└── BTCM24 (id:23, parent:14, crypto) — BTC June 2024 future
//
// This tree allows grouping derivatives under their base asset and resolving
// instrument relationships across exchanges.
type Symbol struct {
	ID          uint       `json:"id"`                    // Unique identifier
	Type        SymbolType `json:"type"`                  // Fiat or crypto
	ParentID    uint       `json:"parentId,omitempty"`    // Parent symbol ID (0 = root)
	Code        string     `json:"code"`                  // Ticker code (e.g. "BTC", "USD")
	Name        string     `json:"name"`                  // Full name (e.g. "Bitcoin")
	Description string     `json:"description,omitempty"` // Optional description
	Website     string     `json:"website,omitempty"`     // Project/issuer website
}

// IsRoot returns true if this symbol has no parent (top-level asset).
func (s Symbol) IsRoot() bool {
	return s.ParentID == 0
}

// IsFiat returns true if this is a fiat currency symbol.
func (s Symbol) IsFiat() bool {
	return s.Type == SymbolFiat
}

// IsCrypto returns true if this is a cryptocurrency symbol.
func (s Symbol) IsCrypto() bool {
	return s.Type == SymbolCrypto
}

// Symbols is a collection of Symbol values with lookup helpers.
type Symbols []Symbol

// GetByCode returns the first symbol matching the given code, or nil.
func (ss Symbols) GetByCode(code string) *Symbol {
	for i := range ss {
		if ss[i].Code == code {
			return &ss[i]
		}
	}
	return nil
}

// GetByID returns the symbol with the given ID, or nil.
func (ss Symbols) GetByID(id uint) *Symbol {
	for i := range ss {
		if ss[i].ID == id {
			return &ss[i]
		}
	}
	return nil
}

// Children returns all symbols whose ParentID matches the given ID.
func (ss Symbols) Children(parentID uint) Symbols {
	var result Symbols
	for _, s := range ss {
		if s.ParentID == parentID {
			result = append(result, s)
		}
	}
	return result
}

// Roots returns all top-level symbols (ParentID == 0).
func (ss Symbols) Roots() Symbols {
	return ss.Children(0)
}

// Fiats returns all fiat currency symbols.
func (ss Symbols) Fiats() Symbols {
	var result Symbols
	for _, s := range ss {
		if s.IsFiat() {
			result = append(result, s)
		}
	}
	return result
}

// Cryptos returns all cryptocurrency symbols.
func (ss Symbols) Cryptos() Symbols {
	var result Symbols
	for _, s := range ss {
		if s.IsCrypto() {
			result = append(result, s)
		}
	}
	return result
}
