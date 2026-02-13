package trade

// AggressorSide indicates which side initiated a trade.
type AggressorSide int

const (
	AggressorNone    AggressorSide = iota // No side determined
	AggressorSell                         // Seller initiated
	AggressorBuy                          // Buyer initiated
	AggressorUnknown                      // Side could not be determined
)

// String returns a human-readable representation.
func (a AggressorSide) String() string {
	switch a {
	case AggressorSell:
		return "sell"
	case AggressorBuy:
		return "buy"
	case AggressorNone:
		return "none"
	default:
		return "unknown"
	}
}
