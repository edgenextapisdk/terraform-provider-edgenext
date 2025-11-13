package data

import (
	"fmt"
	"strconv"
)

// convertTotalToString converts total from interface{} (can be string or number) to string
func convertTotalToString(total interface{}) string {
	if total == nil {
		return "0"
	}

	// If it's already a string, return it
	if str, ok := total.(string); ok {
		return str
	}

	// If it's a number, convert to string
	if num, ok := total.(float64); ok {
		return strconv.FormatFloat(num, 'f', -1, 64)
	}
	if num, ok := total.(int); ok {
		return strconv.Itoa(num)
	}
	if num, ok := total.(int64); ok {
		return strconv.FormatInt(num, 10)
	}

	// Fallback: convert to string using fmt
	return fmt.Sprintf("%v", total)
}
