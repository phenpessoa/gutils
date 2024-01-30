package jsonx

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/phenpessoa/gutils/unsafex"
)

// Int64 is a type that represents an int64 that can will marshaled to JSON as
// a string and can be unmarshaled from JSON as either a JSON number or a JSON
// string.
type Int64 int64

// MarshalJSON implements the json.Marshaler interface for Int64.
// It converts the Int64 value to a JSON string representation.
func (i64 Int64) MarshalJSON() ([]byte, error) {
	str := strconv.FormatInt(int64(i64), 10)
	return unsafex.ByteSlice(`"` + str + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface for Int64 type.
// It unmarshals a JSON value into an Int64 value. The JSON value can either be
// a string or a number.
//
// An empty string is considered valid and will make Int64 be zero.
func (i64 *Int64) UnmarshalJSON(b []byte) error {
	str := unsafex.String(b)
	str = strings.ReplaceAll(str, `"`, "")

	if str == "" {
		*i64 = 0
		return nil
	}

	parsed, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return fmt.Errorf("can not parse %s into Int64: %w", str, err)
	}

	*i64 = Int64(parsed)
	return nil
}
