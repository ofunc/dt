package dt

type sorter struct {
	frame *Frame
	cmp   func(Record, Record) bool
}

// Len is the number of elements in the collection.
func (a sorter) Len() int {
	return a.frame.Len()
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (a sorter) Less(i, j int) bool {
	return a.cmp(record{
		frame: a.frame,
		index: i,
	}, record{
		frame: a.frame,
		index: j,
	})
}

// Swap swaps the elements with indexes i and j.
func (a sorter) Swap(i, j int) {
	for _, list := range a.frame.lists {
		list[i], list[j] = list[j], list[i]
	}
}
