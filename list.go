package dt

// List is the list data structure.
type List []Value

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

// Int converts the list to int list.
func (a List) Int() List {
	for i, v := range a {
		a[i] = Int(v.Int())
	}
	return a
}

// Float converts the list to float list.
func (a List) Float() List {
	for i, v := range a {
		a[i] = Float(v.Float())
	}
	return a
}
