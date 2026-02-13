package trade

// PriceClusters describes volume and time distribution at a price level.
type PriceClusters struct {
	Ask      float64 `json:"ask,omitempty"`           // Volume at the ask
	Bid      float64 `json:"bid,omitempty"`           // Volume at the bid
	Duration float64 `json:"durationMilli,omitempty"` // Time at this level (ms)
	Trades   float64 `json:"trades,omitempty"`        // Number of trades
}
