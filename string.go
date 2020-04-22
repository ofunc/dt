package dt

import (
	"math"
	"strconv"
	"strings"
)

// String is a string value.
type String string

// Float returns as a float64 value.
func (a String) Float() float64 {
	if v, err := strconv.ParseFloat(strings.TrimSpace(string(a)), 64); err == nil {
		return v
	}
	return math.NaN()
}

// String returns as a string value.
func (a String) String() string {
	return string(a)
}
