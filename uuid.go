package trade

import (
	"encoding/json"

	"github.com/google/uuid"
)

// UUID wraps uuid.UUID with custom JSON unmarshaling for string UUIDs.
type UUID struct {
	Value *uuid.UUID
}

// UnmarshalJSON parses a JSON string into a uuid.UUID.
func (u *UUID) UnmarshalJSON(data []byte) (err error) {
	if len(data) < 7 {
		return
	}
	var aux string
	if err = json.Unmarshal(data, &aux); err != nil {
		return
	}
	id, err := uuid.Parse(aux)
	if err != nil {
		return
	}
	u.Value = &id
	return
}

// String returns the UUID string representation, or empty string if nil.
func (u UUID) String() string {
	if u.Value == nil {
		return ""
	}
	return u.Value.String()
}
