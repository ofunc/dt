package dt

import (
	"math"
	"strconv"
)

// Float is a float64 value.
type Float float64

// Int returns as a int value.
func (a Float) Int() int {
	return int(math.Round(float64(a)))
}

// Float returns as a float64 value.
func (a Float) Float() float64 {
	return float64(a)
}

// Bool returns as a bool value.
func (a Float) Bool() bool {
	return a != 0 && !math.IsNaN(float64(a))
}

// String returns as a string value.
func (a Float) String() string {
	return strconv.FormatFloat(float64(a), 'g', -1, 64)
}
