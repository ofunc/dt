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

func makeKey(i int, lists []List) string {
	ks := make([]string, len(lists))
	for j, list := range lists {
		if v := list[i]; v != nil {
			ks[j] = v.String()
		}
	}
	return strings.Join(ks, "\r\t\n")
}
