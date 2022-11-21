package convert

import (
	"fmt"
	"strconv"
	"time"
)

// ToInt64 converts from to an int64.
func ToInt64(from any) (int64, error) {
	switch t := from.(type) {
	case int:
		return int64(t), nil
	case int8:
		return int64(t), nil
	case int16:
		return int64(t), nil
	case int32:
		return int64(t), nil
	case int64:
		return t, nil
	case float32:
		return int64(t), nil
	case float64:
		return int64(t), nil
	case uint:
		return int64(t), nil
	case uint8:
		return int64(t), nil
	case uint16:
		return int64(t), nil
	case uint32:
		return int64(t), nil
	case uint64:
		return int64(t), nil
	case string:
		return strconv.ParseInt(t, 10, 64)
	case time.Duration:
		return int64(t), nil
	case time.Month:
		return int64(t), nil
	case time.Weekday:
		return int64(t), nil
	case bool:
		if t {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("convert: can not convert type %T to int64", t)
	}
}
