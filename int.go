package dt

import (
	"strconv"
)

// Int is a int value.
type Int int

// Int returns as a int value.
func (a Int) Int() int {
	return int(a)
}

// Float returns as a float64 value.
func (a Int) Float() float64 {
	return float64(a)
}

// Bool returns as a bool value.
func (a Int) Bool() bool {
	return a != 0
}

// String returns as a string value.
func (a Int) String() string {
	return strconv.Itoa(int(a))
}
