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

// Bool converts the list to bool list.
func (a List) Bool() List {
	for i, v := range a {
		a[i] = Bool(v.Bool())
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

// Count returns the count of list a.
func (a List) Count() Value {
	return Int(len(a))
}

// Sum returns the sum of list a.
func (a List) Sum() Value {
	var s Value = Int(0)
	for _, v := range a {
		if x, ok := s.(Int); ok {
			if y, ok := v.(Int); ok {
				s = x + y
			}
		}
		s = Float(s.Float() + v.Float())
	}
	return s
}
