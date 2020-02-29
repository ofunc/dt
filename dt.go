package dt

import (
	"math"
	"reflect"
)

var tvalue = reflect.TypeOf(Value(nil))

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

func makeKey(t reflect.Type, r Record, keys []string) interface{} {
	k := reflect.New(t).Elem()
	for i, key := range keys {
		k.Index(i).Set(reflect.ValueOf(r.Value(key)))
	}
	return k.Interface()
}
