package trade

import "time"

// TimeAndSale represents an atomic trade transaction captured from an exchange.
// This is the fundamental unit of market data — each instance records a single
// trade event with its price, volume, aggressor side, and timing.
type TimeAndSale struct {
	InternalID UUID // Internal transaction ID

	ID                 string // Exchange transaction ID / Trade ID
	ExchangeID         int64  // Exchange identifier
	DataFeedProviderID int64  // Data feed provider identifier

	TradeSequence     int64 // Sequence number for same-timestamp trades
	TradeOpenInterest int64 // Open interest at time of trade

	Ticker string    // Market trading symbol (e.g. "BTCUSDT", "ES", "NQ")
	Time   time.Time // Trade timestamp
	Sale
}

// Sale holds the core trade data: price, direction, and volume.
type Sale struct {
	Price         float64       `json:"price"`         // Trade price
	AggressorSide AggressorSide `json:"aggressorSide"` // Buy or sell initiated
	Volume        int           `json:"volume"`        // Quantity of contracts/shares
}

// DataFeedProvider identifies a data feed source.
type DataFeedProvider struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	URN  string `json:"urn"`
}

// Exchange identifies a trading exchange.
type Exchange struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// OrderBookEntry represents a single level in an order book snapshot.
type OrderBookEntry struct {
	Ticker     string    `json:"ticker"`
	ExchangeID int64     `json:"exchangeId"`
	Time       time.Time `json:"time"`
	BestBid    Sale      `json:"bestBid"`
	BestAsk    Sale      `json:"bestAsk"`
}

// OrderBook is a collection of order book snapshots.
type OrderBook []OrderBookEntry
