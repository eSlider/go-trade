package trade

import "time"

// DateTime wraps time.Time with custom JSON unmarshaling
// using the "2006-01-02 15:04:05" format common in exchange APIs.
type DateTime struct {
	*time.Time
}

// UnmarshalJSON parses a date-time string in "2006-01-02 15:04:05" format.
func (t *DateTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) < 7 {
		return
	}
	data = data[1 : len(data)-1]
	var ts time.Time
	ts, err = time.Parse(time.DateTime, string(data))
	t.Time = &ts
	return
}
