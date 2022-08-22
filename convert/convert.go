package convert

import (
	"strconv"
	"time"
)

// ToInt converts from to an int64.
func ToInt64(from any) int64 {
	switch t := from.(type) {
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case int64:
		return t
	case float32:
		return int64(t)
	case float64:
		return int64(t)
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		return int64(t)
	case string:
		parsed, _ := strconv.ParseInt(t, 10, 64)
		return parsed
	case time.Duration:
		return int64(t)
	case time.Month:
		return int64(t)
	case time.Weekday:
		return int64(t)
	case bool:
		if t {
			return 1
		}
		return 0
	default:
		return 0
	}
}
