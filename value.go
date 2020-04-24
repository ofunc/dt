package dt

import (
	"math"
	"strconv"
	"strings"
)

// Value is the value interface.
type Value interface {
	Number() float64
	String() string
}

// Number is a float64 value.
type Number float64

// String is a string value.
type String string

// Number returns as a number value.
func (a Number) Number() float64 {
	return float64(a)
}

// String returns as a string value.
func (a String) String() string {
	return string(a)
}

// Number returns as a number value.
func (a String) Number() float64 {
	if v, err := strconv.ParseFloat(strings.TrimSpace(string(a)), 64); err == nil {
		return v
	}
	return math.NaN()
}

// String returns as a string value.
func (a Number) String() string {
	if v := int64(a); Number(v) == a {
		return strconv.FormatInt(v, 10)
	}
	return strconv.FormatFloat(float64(a), 'g', -1, 64)
}
