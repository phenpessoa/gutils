package convert

import (
	"testing"
	"time"
)

func TestToSignedInteger(t *testing.T) {
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
		{"uintptr", uintptr(10), 10, false},
		{"string-10", "10", 10, false},
		{"string-a", "a", 0, true},
		{"time-nanosecond", time.Nanosecond, 1, false},
		{"time-month-5", time.Month(5), 5, false},
		{"time-weekday-5", time.Weekday(5), 5, false},
		{"time", time.Date(2022, 11, 30, 18, 41, 15, 10, time.UTC), 1669833675, false},
		{"true", true, 1, false},
		{"false", false, 0, false},
		{"empty-[]byte", []byte{}, 0, true},
		{"negative", -1, -1, false},
		{"negative str", "-1", -1, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := toInteger[int64](tc.input, false)
			if ((err != nil) != tc.wantErr) || tc.want != got {
				t.Errorf("\ntest '%s' failed to convert\nwant: %v\ngot: %v\nwantErr: %v\nerr: %v",
					tc.name, tc.want, got, tc.wantErr, err,
				)
			}
		})
	}
}

func TestToUnsignedInteger(t *testing.T) {
	for _, tc := range []struct {
		name    string
		input   any
		want    uint64
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
		{"uintptr", uintptr(10), 10, false},
		{"string-10", "10", 10, false},
		{"string-a", "a", 0, true},
		{"time-nanosecond", time.Nanosecond, 1, false},
		{"time-month-5", time.Month(5), 5, false},
		{"time-weekday-5", time.Weekday(5), 5, false},
		{"time", time.Date(2022, 11, 30, 18, 41, 15, 10, time.UTC), 1669833675, false},
		{"true", true, 1, false},
		{"false", false, 0, false},
		{"empty-[]byte", []byte{}, 0, true},
		{"negative", -1, 1<<64 - 1, false},
		{"negative str", "-1", 0, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := toInteger[uint64](tc.input, true)
			if ((err != nil) != tc.wantErr) || tc.want != got {
				t.Errorf("\ntest '%s' failed to convert\nwant: %v\ngot: %v\nwantErr: %v\nerr: %v",
					tc.name, tc.want, got, tc.wantErr, err,
				)
			}
		})
	}
}

func TestToIntegerFuncs(t *testing.T) {
	for _, tc := range []struct {
		name string
		f    func() error
	}{
		{
			"ToInt",
			func() error { _, err := ToInt(10); return err },
		},
		{
			"ToInt8",
			func() error { _, err := ToInt8(10); return err },
		},
		{
			"ToInt16",
			func() error { _, err := ToInt16(10); return err },
		},
		{
			"ToInt32",
			func() error { _, err := ToInt32(10); return err },
		},
		{
			"ToInt64",
			func() error { _, err := ToInt64(10); return err },
		},
		{
			"ToUint",
			func() error { _, err := ToUint(10); return err },
		},
		{
			"ToUint",
			func() error { _, err := ToUint(10); return err },
		},
		{
			"ToUint8",
			func() error { _, err := ToUint8(10); return err },
		},
		{
			"ToUint16",
			func() error { _, err := ToUint16(10); return err },
		},
		{
			"ToUint32",
			func() error { _, err := ToUint32(10); return err },
		},
		{
			"ToUint64",
			func() error { _, err := ToUint64(10); return err },
		},
		{
			"ToUintptr",
			func() error { _, err := ToUintptr(10); return err },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.f(); err != nil {
				t.Errorf("\ntest '%s' failed\nerr: %v", tc.name, err)
			}
		})
	}
}
