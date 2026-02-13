package trade

// InstrumentType classifies a trading instrument.
type InstrumentType string

const (
	InstrumentSpot   InstrumentType = "spot"
	InstrumentFuture InstrumentType = "future"
	InstrumentOption InstrumentType = "option"
	InstrumentFX     InstrumentType = "fx"
)

// Instrument represents a tradable instrument (e.g. BTCUSDT, ES futures, EUR/USD).
type Instrument struct {
	ID                 int            `json:"id"`
	Type               InstrumentType `json:"type"`
	Ticker             string         `json:"ticker"`
	Name               string         `json:"name"`
	Description        string         `json:"description,omitempty"`
	ExchangeID         int            `json:"exchangeId,omitempty"`
	DataFeedProviderID int            `json:"dataFeedProviderId,omitempty"`
}

// Market represents a trading pair with base and quote symbols.
type Market struct {
	ID          int    `json:"id"`
	FromSymbol  string `json:"fromSymbol"` // Base (e.g. BTC)
	ToSymbol    string `json:"toSymbol"`   // Quote (e.g. USDT)
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

// String returns the market as "FROM/TO".
func (m Market) String() string {
	return m.FromSymbol + "/" + m.ToSymbol
}
