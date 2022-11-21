package convert

import (
	"testing"
	"time"
)

func TestToInt64(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   any
		want    int64
		wantErr bool
	}{
		{"int", int(10), 10, false},
		{"int8", int8(10), 10, false},
		{"int16", int16(10), 10, false},
		{"int32", int32(10), 10, false},
		{"int64", int64(10), 10, false},
		{"float32", float32(10), 10, false},
		{"float64", float64(10), 10, false},
		{"uint", uint(10), 10, false},
		{"uint8", uint8(10), 10, false},
		{"uint16", uint16(10), 10, false},
		{"uint32", uint32(10), 10, false},
		{"uint64", uint64(10), 10, false},
		{"string-10", "10", 10, false},
		{"string-a", "a", 0, true},
		{"time-nanosecond", time.Nanosecond, 1, false},
		{"time-month-5", time.Month(5), 5, false},
		{"time-weekday-5", time.Weekday(5), 5, false},
		{"true", true, 1, false},
		{"false", false, 0, false},
		{"empty-[]byte", []byte{}, 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ToInt64(tc.input)
			if ((err != nil) != tc.wantErr) || tc.want != got {
				t.Errorf("\ntest '%s' failed to convert\nwant: %v\ngot: %v\nwantErr: %v\nerr: %v", tc.name, tc.want, got, tc.wantErr, err)
			}
		})
	}
}
