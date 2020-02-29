package dt

// List is the list data structure.
type List []Value

// Filter filters the list with function f.
func (a List) Filter(f func(Value) bool) List {
	var b List
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
