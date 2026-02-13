package trade

// CandleDeltaLevels holds min/max delta values and their price levels
// within a candle, used for volume profile and delta analysis.
type CandleDeltaLevels struct {
	MinDeltaValue float64 `json:"minDeltaValue,omitempty"`
	MinDeltaPrice float64 `json:"minDeltaPrice,omitempty"`
	MaxDeltaValue float64 `json:"maxDeltaValue,omitempty"`
	MaxDeltaPrice float64 `json:"maxDeltaPrice,omitempty"`
}
