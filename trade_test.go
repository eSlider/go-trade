package trade

import (
	"encoding/json"
	"testing"
	"time"
)

func TestOrderUnmarshalJSON(t *testing.T) {
	t.Run("valid data", func(t *testing.T) {
		data := []byte(`{
			"order_id": "b68d69564a79dea4776afa33d1d2fcab",
			"customer_id": "41",
			"order_status": "shipped",
			"order_approved_at": "2018-02-28 10:40:35",
			"order_delivered_customer_date": "2018-03-05 16:10:13",
			"order_delivered_carrier_date": "2018-03-05 16:10:13",
			"order_estimated_delivery_date": "2018-03-23 00:00:00",
			"order_purchase_timestamp": "2018-02-28 08:57:03"
		}`)

		order := &Order{}
		err := order.UnmarshalJSON(data)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if order.OrderID.String() != "b68d6956-4a79-dea4-776a-fa33d1d2fcab" {
			t.Errorf("OrderID = %s, want b68d6956-4a79-dea4-776a-fa33d1d2fcab", order.OrderID.String())
		}

		if order.CustomerID != 41 {
			t.Errorf("CustomerID = %d, want 41", order.CustomerID)
		}

		if order.Status != "shipped" {
			t.Errorf("Status = %s, want shipped", order.Status)
		}

		if order.ApprovedAt == nil {
			t.Fatal("ApprovedAt is nil")
		}
	})

	t.Run("invalid customer_id", func(t *testing.T) {
		data := []byte(`{
			"order_id": "b68d69564a79dea4776afa33d1d2fcab",
			"customer_id": "invalid",
			"order_status": "shipped",
			"order_approved_at": "2018-02-28 10:40:35",
			"order_delivered_customer_date": "2018-03-05 16:10:13",
			"order_delivered_carrier_date": "2018-03-05 16:10:13",
			"order_estimated_delivery_date": "2018-03-23 00:00:00",
			"order_purchase_timestamp": "2018-02-28 08:57:03"
		}`)

		order := &Order{}
		err := order.UnmarshalJSON(data)
		if err == nil {
			t.Error("Expected error for invalid customer_id, got nil")
		}
	})
}

func TestCandleMethods(t *testing.T) {
	c := Candle{
		Open:      100.0,
		High:      110.0,
		Low:       95.0,
		Close:     108.0,
		Ask:       108.5,
		Bid:       107.5,
		TimeOpen:  time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
		TimeClose: time.Date(2025, 1, 1, 10, 5, 0, 0, time.UTC),
	}

	if c.IsEmpty() {
		t.Error("IsEmpty() = true, want false")
	}
	if !c.IsBullish() {
		t.Error("IsBullish() = false, want true (close > open)")
	}
	if c.IsBearish() {
		t.Error("IsBearish() = true, want false")
	}
	if c.Range() != 15.0 {
		t.Errorf("Range() = %f, want 15.0", c.Range())
	}
	if c.Spread() != 1.0 {
		t.Errorf("Spread() = %f, want 1.0", c.Spread())
	}
	if c.Duration() != 5*time.Minute {
		t.Errorf("Duration() = %v, want 5m", c.Duration())
	}
}

func TestCandleBearish(t *testing.T) {
	c := Candle{Open: 100.0, Close: 90.0}
	if !c.IsBearish() {
		t.Error("IsBearish() = false, want true")
	}
	if c.IsBullish() {
		t.Error("IsBullish() = true, want false")
	}
}

func TestAggressorSideString(t *testing.T) {
	tests := []struct {
		side AggressorSide
		want string
	}{
		{AggressorNone, "none"},
		{AggressorBuy, "buy"},
		{AggressorSell, "sell"},
		{AggressorUnknown, "unknown"},
	}
	for _, tt := range tests {
		if got := tt.side.String(); got != tt.want {
			t.Errorf("AggressorSide(%d).String() = %s, want %s", tt.side, got, tt.want)
		}
	}
}

func TestDateTimeUnmarshalJSON(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		dt := &DateTime{}
		err := dt.UnmarshalJSON([]byte(`"2025-06-15 14:30:00"`))
		if err != nil {
			t.Fatalf("UnmarshalJSON error = %v", err)
		}
		if dt.Time == nil {
			t.Fatal("Time is nil")
		}
		if dt.Time.Year() != 2025 {
			t.Errorf("Year = %d, want 2025", dt.Time.Year())
		}
	})

	t.Run("empty", func(t *testing.T) {
		dt := &DateTime{}
		err := dt.UnmarshalJSON([]byte(`""`))
		if err != nil {
			t.Fatalf("UnmarshalJSON error = %v", err)
		}
		if dt.Time != nil {
			t.Error("Time should be nil for empty input")
		}
	})
}

func TestUUIDUnmarshalJSON(t *testing.T) {
	u := &UUID{}
	data := []byte(`"550e8400-e29b-41d4-a716-446655440000"`)
	if err := u.UnmarshalJSON(data); err != nil {
		t.Fatalf("UnmarshalJSON error = %v", err)
	}
	if u.Value == nil {
		t.Fatal("Value is nil")
	}
	if u.String() != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("String() = %s", u.String())
	}
}

func TestUUIDEmptyString(t *testing.T) {
	u := UUID{}
	if u.String() != "" {
		t.Errorf("String() = %q, want empty", u.String())
	}
}

func TestMarketString(t *testing.T) {
	m := Market{FromSymbol: "BTC", ToSymbol: "USDT"}
	if m.String() != "BTC/USDT" {
		t.Errorf("String() = %s, want BTC/USDT", m.String())
	}
}

func TestTemperaturesRoundTrip(t *testing.T) {
	original := Temperatures{
		AskBid:         AskBid{Ask: 100, Bid: 50},
		TimeDurationMS: 5000,
		TradesCount:    42,
	}

	data, err := original.Marshal()
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}

	parsed, err := UnmarshalTemperatures(data)
	if err != nil {
		t.Fatalf("UnmarshalTemperatures error = %v", err)
	}

	if parsed.TradesCount != 42 {
		t.Errorf("TradesCount = %d, want 42", parsed.TradesCount)
	}
}

func TestPriceClustersJSON(t *testing.T) {
	pc := PriceClusters{Ask: 1.5, Bid: 2.5, Duration: 100, Trades: 10}
	data, err := json.Marshal(pc)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}

	var parsed PriceClusters
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if parsed.Ask != 1.5 {
		t.Errorf("Ask = %f, want 1.5", parsed.Ask)
	}
}
