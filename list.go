package dt

// List is the list data structure.
type List []Value

// Number converts the list to float list.
func (a List) Number() List {
	for i, v := range a {
		a[i] = Number(v.Number())
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
