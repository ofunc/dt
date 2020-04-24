package dt

import "math"

// First returns the first of list l.
func First(l List) Value {
	if len(l) > 0 {
		return l[0]
	}
	return nil
}

// Last returns the last of list l.
func Last(l List) Value {
	if n := len(l); n > 0 {
		return l[n-1]
	}
	return nil
}

// Count returns the count of list l.
func Count(l List) Value {
	return Number(len(l))
}

// Sum returns the sum of list l.
func Sum(l List) Value {
	s := 0.0
	for _, v := range l {
		if v == nil {
			return Number(math.NaN())
		}
		s += v.Number()
	}
	return Number(s)
}

// Mean returns the mean of list l.
func Mean(l List) Value {
	return Number(Sum(l).Number()) / Number(len(l))
}

// Var returns the var of list l.
func Var(l List) Value {
	x, y, n := 0.0, 0.0, float64(len(l))
	for _, v := range l {
		if v == nil {
			return Number(math.NaN())
		}
		z := v.Number()
		x += z * z
		y += z
	}
	y /= n
	return Number(x/n - y*y)
}

// Std returns the std of list l.
func Std(l List) Value {
	if v := Var(l).Number(); v > 0 {
		return Number(math.Sqrt(v))
	}
	return Number(0)
}

// Min returns the min of list l.
func Min(l List) Value {
	m := math.Inf(1)
	for _, v := range l {
		if v != nil {
			if x := v.Number(); x < m {
				m = x
			}
		}
	}
	return Number(m)
}

// Max returns the max of list l.
func Max(l List) Value {
	m := math.Inf(-1)
	for _, v := range l {
		if v != nil {
			if x := v.Number(); x > m {
				m = x
			}
		}
	}
	return Number(m)
}
