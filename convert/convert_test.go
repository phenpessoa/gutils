package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToInt64(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input any
		want  int64
	}{
		{"int", int(10), 10},
		{"int8", int8(10), 10},
		{"int16", int16(10), 10},
		{"int32", int32(10), 10},
		{"int64", int64(10), 10},
		{"float32", float32(10), 10},
		{"float64", float64(10), 10},
		{"uint", uint(10), 10},
		{"uint8", uint8(10), 10},
		{"uint16", uint16(10), 10},
		{"uint32", uint32(10), 10},
		{"uint64", uint64(10), 10},
		{"string-10", "10", 10},
		{"string-a", "a", 0},
		{"time-nanosecond", time.Nanosecond, 1},
		{"time-month-5", time.Month(5), 5},
		{"time-weekday-5", time.Weekday(5), 5},
		{"true", true, 1},
		{"false", false, 0},
		{"empty-[]byte", []byte{}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, ToInt64(tc.input))
		})
	}
}

func TestStringToByteSliceUnsafe(t *testing.T) {
	str := "foo, bar, baz, qux, quux"
	assert.Equal(t, []byte(str), StringToByteSliceUnsafe(str))
}

func TestByteSliceToStringUnsafe(t *testing.T) {
	bs := []byte("foo, bar, baz, qux, quux")
	assert.Equal(t, "foo, bar, baz, qux, quux", ByteSliceToStringUnsafe(bs))
}
