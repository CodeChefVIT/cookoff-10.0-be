package utils

import (
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
)

// InterfaceToNumeric converts an interface{} to pgtype.Numeric
// Supports string, float64, int, int64, and float32 types
func InterfaceToNumeric(val interface{}) (pgtype.Numeric, error) {
	numeric := pgtype.Numeric{}

	switch v := val.(type) {
	case string:
		// Try to parse as float first
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			err := numeric.Scan(f)
			return numeric, err
		}
		// If not a valid float, try as int
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			err := numeric.Scan(i)
			return numeric, err
		}
		return numeric, fmt.Errorf("could not convert string to numeric: %v", v)

	case float64:
		err := numeric.Scan(v)
		return numeric, err

	case float32:
		err := numeric.Scan(float64(v))
		return numeric, err

	case int:
		err := numeric.Scan(int64(v))
		return numeric, err

	case int64:
		err := numeric.Scan(v)
		return numeric, err

	case uint:
		err := numeric.Scan(int64(v))
		return numeric, err

	case uint64:
		// Handle potential overflow
		if v > 1<<63-1 {
			return numeric, fmt.Errorf("value too large for int64: %v", v)
		}
		err := numeric.Scan(int64(v))
		return numeric, err

	case int32:
		err := numeric.Scan(int64(v))
		return numeric, err

	case uint32:
		err := numeric.Scan(int64(v))
		return numeric, err

	case int16:
		err := numeric.Scan(int64(v))
		return numeric, err

	case uint16:
		err := numeric.Scan(int64(v))
		return numeric, err

	case int8:
		err := numeric.Scan(int64(v))
		return numeric, err

	case uint8:
		err := numeric.Scan(int64(v))
		return numeric, err

	default:
		return numeric, fmt.Errorf("unsupported type for numeric conversion: %T", v)
	}
}
