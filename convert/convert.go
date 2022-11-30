package convert

import (
	"fmt"
	"strconv"
	"time"
)

// ToInt converts from to an int.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToInt(from any) (int, error) {
	return toInteger[int](from, false)
}

// ToInt8 converts from to an int8.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToInt8(from any) (int8, error) {
	return toInteger[int8](from, false)
}

// ToInt16 converts from to an int16.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToInt16(from any) (int16, error) {
	return toInteger[int16](from, false)
}

// ToInt32 converts from to an int32.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToInt32(from any) (int32, error) {
	return toInteger[int32](from, false)
}

// ToInt64 converts from to an int64.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToInt64(from any) (int64, error) {
	return toInteger[int64](from, false)
}

// ToUint converts from to an uint.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToUint(from any) (uint, error) {
	return toInteger[uint](from, true)
}

// ToUint8 converts from to an uint8.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToUint8(from any) (uint8, error) {
	return toInteger[uint8](from, true)
}

// ToUint16 converts from to an uint16.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToUint16(from any) (uint16, error) {
	return toInteger[uint16](from, true)
}

// ToUint32 converts from to an uint32.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToUint32(from any) (uint32, error) {
	return toInteger[uint32](from, true)
}

// ToUint64 converts from to an uint64.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToUint64(from any) (uint64, error) {
	return toInteger[uint64](from, true)
}

// ToUintptr converts from to an uintptr.
//
// If from is an integer or a float, a time.Duration,
// a time.Month or a time.Weekday a direct
// type convertion will be made.
//
// If from is a bool, it will return 1 for
// true and 0 for false.
//
// If from is time.Time, it will return
// the Unix time.
//
// If from is a string, strconv will be used.
//
// No other types are allowed and will result
// in an error.
func ToUintptr(from any) (uintptr, error) {
	return toInteger[uintptr](from, true)
}

func toInteger[To strictInteger](from any, unsigned bool) (To, error) {
	switch t := from.(type) {
	case int:
		return To(t), nil
	case int8:
		return To(t), nil
	case int16:
		return To(t), nil
	case int32:
		return To(t), nil
	case int64:
		return To(t), nil
	case float32:
		return To(t), nil
	case float64:
		return To(t), nil
	case uint:
		return To(t), nil
	case uint8:
		return To(t), nil
	case uint16:
		return To(t), nil
	case uint32:
		return To(t), nil
	case uint64:
		return To(t), nil
	case uintptr:
		return To(t), nil
	case time.Duration:
		return To(t), nil
	case time.Month:
		return To(t), nil
	case time.Weekday:
		return To(t), nil
	case time.Time:
		return To(t.Unix()), nil
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	case string:
		if unsigned {
			x, err := strconv.ParseUint(t, 10, 64)
			return To(x), err
		}
		x, err := strconv.ParseInt(t, 10, 64)
		return To(x), err
	default:
		var zero To
		return 0, fmt.Errorf("can not convert type '%T' to '%T'", t, zero)
	}
}

type strictInteger interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr
}
