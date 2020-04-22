package dt

import (
	"math"
)

// List is the list data structure.
type List []Value

// Float converts the list to float list.
func (a List) Float() List {
	for i, v := range a {
		a[i] = Float(v.Float())
	}
	return a
}

// String converts the list to string list.
func (a List) String() List {
	for i, v := range a {
		a[i] = String(v.String())
	}
	return a
}

// Map maps the list by function f.
func (a List) Map(f func(Value) Value) List {
	b := make(List, len(a))
	for i, v := range a {
		b[i] = f(v)
	}
	return b
}

// Filter filters the list with function f.
func (a List) Filter(f func(Value) bool) List {
	b := make(List, 0, len(a))
	for _, v := range a {
		if f(v) {
			b = append(b, v)
		}
	}
	return b
}

// FillNA fills NA value with value.
func (a List) FillNA(value Value) List {
	for i, v := range a {
		if IsNA(v) {
			a[i] = value
		}
	}
	return a
}

// First returns the first of list a.
func (a List) First() Value {
	if len(a) > 0 {
		return a[0]
	}
	return nil
}

// Last returns the last of list a.
func (a List) Last() Value {
	if n := len(a); n > 0 {
		return a[n-1]
	}
	return nil
}

// Count returns the count of list a.
func (a List) Count() Value {
	return Float(len(a))
}

// Sum returns the sum of list a.
func (a List) Sum() Value {
	s := 0.0
	for _, v := range a {
		if IsNA(v) {
			return Float(math.NaN())
		}
		s += v.Float()
	}
	return Float(s)
}

// Mean returns the mean of list a.
func (a List) Mean() Value {
	return a.Sum().(Float) / Float(len(a))
}

// Var returns the var of list a.
func (a List) Var() Value {
	x, y, n := 0.0, 0.0, float64(len(a))
	for _, v := range a {
		if IsNA(v) {
			return Float(math.NaN())
		}
		z := v.Float()
		x += z * z
		y += z
	}
	y /= n
	return Float(x/n - y*y)
}

// Std returns the std of list a.
func (a List) Std() Value {
	v := a.Var().Float()
	if v <= 0 {
		return Float(0)
	}
	return Float(math.Sqrt(v))
}

// Min returns the min of list a.
func (a List) Min() Value {
	m := math.Inf(1)
	for _, v := range a {
		if v != nil {
			if x := v.Float(); x < m {
				m = x
			}
		}
	}
	return Float(m)
}

// Max returns the max of list a.
func (a List) Max() Value {
	m := math.Inf(-1)
	for _, v := range a {
		if v != nil {
			if x := v.Float(); x > m {
				m = x
			}
		}
	}
	return Float(m)
}
