package dt

import (
	"math"
	"reflect"
)

var na = Value(nil)
var tvalue = reflect.TypeOf(&na).Elem()

// IsNA checks if a is NA.
func IsNA(a Value) bool {
	switch v := a.(type) {
	case nil:
		return true
	case Float:
		return math.IsNaN(float64(v))
	default:
		return false
	}
}

// First returns the first of list l.
func First(l List) Value {
	return l.First()
}

// Last returns the last of list l.
func Last(l List) Value {
	return l.Last()
}

// Count returns the count of list l.
func Count(l List) Value {
	return l.Count()
}

// Sum returns the sum of list l.
func Sum(l List) Value {
	return l.Sum()
}

// Mean returns the mean of list l.
func Mean(l List) Value {
	return l.Mean()
}

// Var returns the var of list l.
func Var(l List) Value {
	return l.Var()
}

// Std returns the std of list l.
func Std(l List) Value {
	return l.Std()
}

// Min returns the min of list l.
func Min(l List) Value {
	return l.Min()
}

// Max returns the max of list l.
func Max(l List) Value {
	return l.Max()
}

func makeKey(t reflect.Type, r Record, keys []string) interface{} {
	k := reflect.New(t).Elem()
	for i, key := range keys {
		k.Index(i).Set(reflect.ValueOf(r.Value(key)))
	}
	return k.Interface()
}
