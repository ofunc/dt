package dt

import (
	"strconv"
)

// Bool is a int value.
type Bool bool

// Int returns as a int value.
func (a Bool) Int() int {
	if a {
		return 1
	}
	return 0
}

// Float returns as a float64 value.
func (a Bool) Float() float64 {
	if a {
		return 1
	}
	return 0
}

// Bool returns as a bool value.
func (a Bool) Bool() bool {
	return bool(a)
}

// String returns as a string value.
func (a Bool) String() string {
	return strconv.FormatBool(bool(a))
}
