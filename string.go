package dt

import (
	"math"
	"strconv"
	"strings"
)

// String is a string value.
type String string

// Int returns as a int value.
func (a String) Int() int {
	if v, err := strconv.ParseInt(strings.TrimSpace(string(a)), 10, 64); err == nil {
		return int(v)
	}
	if v := math.Round(a.Float()); !math.IsNaN(v) {
		return int(v)
	}
	return 0
}

// Float returns as a float64 value.
func (a String) Float() float64 {
	if v, err := strconv.ParseFloat(strings.TrimSpace(string(a)), 64); err == nil {
		return v
	}
	return math.NaN()
}

// Bool returns as a bool value.
func (a String) Bool() bool {
	return a != ""
}

// String returns as a string value.
func (a String) String() string {
	return string(a)
}
