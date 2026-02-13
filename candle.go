package trade

import (
	"encoding/json"
	"time"
)

// AskBidPrice represents aggregated trade details within a candlestick.
type AskBidPrice struct {
	Price       float64   `json:"price"`       // Trade price
	AskBid      AskBid    `json:"askBid"`      // Ask/bid counts
	Duration    time.Time `json:"duration"`     // Time delta between trades
	TradesCount int32     `json:"tradesCount"` // Number of events in the deal
}

// Candle represents an OHLC candlestick with extended market microstructure data.
type Candle struct {
	Empty bool `json:"empty,omitempty"`

	// Time
	TimeOpen  time.Time `json:"timeOpen,omitempty"`
	TimeClose time.Time `json:"timeClose,omitempty"`

	// Price
	Open  float64 `json:"open,omitempty"`
	High  float64 `json:"high,omitempty"`
	Low   float64 `json:"low,omitempty"`
	Close float64 `json:"close,omitempty"`
	Ask   float64 `json:"ask,omitempty"`
	Bid   float64 `json:"bid,omitempty"`

	TradesCount   float64                  `json:"tradesCount,omitempty"`
	PriceClusters map[string]PriceClusters `json:"priceClusters,omitempty"`

	DeltaLevels  *CandleDeltaLevels `json:"deltaLevels,omitempty"`
	VolumeLevels *CandleDeltaLevels `json:"volumeLevels,omitempty"`

	PriceSnake []float64 `json:"priceSnake,omitempty"`
}

// IsEmpty returns true if the candle contains no data.
func (c *Candle) IsEmpty() bool {
	return c.Empty
}

// Spread returns the ask-bid spread.
func (c *Candle) Spread() float64 {
	return c.Ask - c.Bid
}

// Range returns the high-low price range.
func (c *Candle) Range() float64 {
	return c.High - c.Low
}

// Duration returns the candle's time span.
func (c *Candle) Duration() time.Duration {
	return c.TimeClose.Sub(c.TimeOpen)
}

// IsBullish returns true if the close is above the open.
func (c *Candle) IsBullish() bool {
	return c.Close > c.Open
}

// IsBearish returns true if the close is below the open.
func (c *Candle) IsBearish() bool {
	return c.Close < c.Open
}

// Temperatures aggregates ask/bid counts with timing data.
type Temperatures struct {
	AskBid         AskBid `json:"askBid"`
	TimeDurationMS int64  `json:"timeDurationMs"`
	TradesCount    int64  `json:"tradesCount"`
}

// UnmarshalTemperatures parses JSON into a Temperatures value.
func UnmarshalTemperatures(data []byte) (Temperatures, error) {
	var r Temperatures
	err := json.Unmarshal(data, &r)
	return r, err
}

// Marshal serializes Temperatures to JSON.
func (r *Temperatures) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// AskBid holds ask and bid counts.
type AskBid struct {
	Ask int64 `json:"ask"`
	Bid int64 `json:"bid"`
}
