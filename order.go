package trade

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

// Order represents a trading order with custom JSON unmarshaling
// that handles string-encoded numeric fields common in exchange APIs.
type Order struct {
	OrderID               uuid.UUID  `mapstructure:"order_id" json:"orderId"`
	CustomerID            int64      `mapstructure:"customer_id" json:"customerId"`
	Status                string     `mapstructure:"order_status" json:"status"`
	ApprovedAt            *time.Time `mapstructure:"order_approved_at" json:"approvedAt,omitempty"`
	DeliveredToCustomerAt *time.Time `mapstructure:"order_delivered_customer_date" json:"deliveredToCustomerAt,omitempty"`
	DeliveredToCarrierAt  *time.Time `mapstructure:"order_delivered_carrier_date" json:"deliveredToCarrierAt,omitempty"`
	DeliveryEstimatedAt   *time.Time `mapstructure:"order_estimated_delivery_date" json:"deliveryEstimatedAt,omitempty"`
	PurchaseAt            *time.Time `mapstructure:"order_purchase_timestamp" json:"purchaseAt,omitempty"`
}

// UnmarshalJSON handles string-to-native type conversions for exchange APIs
// that encode numeric and date fields as strings.
func (r *Order) UnmarshalJSON(data []byte) (err error) {
	var aux map[string]interface{}
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: nil,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			switch t.String() {
			case "int64":
				i, err := strconv.ParseInt(data.(string), 10, 64)
				if err != nil {
					return 0, err
				}
				return i, nil
			case "uint8":
				return int64(data.(uint8)), nil
			case "uuid.UUID":
				return uuid.Parse(data.(string))
			case "time.Time":
				s := data.(string)
				if s == "" {
					return time.Time{}, nil
				}
				return time.Parse(time.DateTime, s)
			}
			return data, nil
		}),
		Result: r,
	})
	if err != nil {
		return
	}
	return decoder.Decode(aux)
}
