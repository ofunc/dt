package dt

import (
	"math"
	"strings"
)

// IsNA checks if a is NA.
func IsNA(a Value) bool {
	switch v := a.(type) {
	case nil:
		return true
	case Number:
		return math.IsNaN(float64(v))
	default:
		return false
	}
}

func makeKey(r Record, keys []string) string {
	ks := make([]string, len(keys))
	for i, key := range keys {
		ks[i] = string(r.String(key))
	}
	return strings.Join(ks, "\r\t\n")
}
