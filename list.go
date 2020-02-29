package dt

// List is the list data structure.
type List []Value

// FillNA fills NA value with value.
func (a List) FillNA(value Value) List {
	for i, v := range a {
		if IsNA(v) {
			a[i] = value
		}
	}
	return a
}
