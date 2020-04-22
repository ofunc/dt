package dt

import (
	"math"
)

// Record is the record interface.
type Record interface {
	Value(key string) Value
	Float(key string) float64
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

// Float returns the float64 value by key.
func (a record) Float(key string) float64 {
	if v := a.Value(key); v != nil {
		return v.Float()
	}
	return math.NaN()
}

// String returns the string value by key.
func (a record) String(key string) string {
	if v := a.Value(key); v != nil {
		return v.String()
	}
	return ""
}
