package dt

import (
	"math"
)

// Record is the record interface.
type Record interface {
	Value(key string) Value
	Int(key string) int
	Float(key string) float64
	Bool(key string) bool
	String(key string) string
}

// record is a ref record.
type record struct {
	index int
	frame *Frame
}

// Value returns the value by key.
func (a record) Value(key string) Value {
	if i, ok := a.frame.index[key]; ok {
		return a.frame.lists[i][a.index]
	}
	return nil
}

// Int returns the int value by key.
func (a record) Int(key string) int {
	if v := a.Value(key); v != nil {
		return v.Int()
	}
	return 0
}

// Float returns the float64 value by key.
func (a record) Float(key string) float64 {
	if v := a.Value(key); v != nil {
		return v.Float()
	}
	return math.NaN()
}

// Bool returns the bool value by key.
func (a record) Bool(key string) bool {
	if v := a.Value(key); v != nil {
		return v.Bool()
	}
	return false
}

// String returns the string value by key.
func (a record) String(key string) string {
	if v := a.Value(key); v != nil {
		return v.String()
	}
	return ""
}
