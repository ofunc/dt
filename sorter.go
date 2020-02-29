package dt

type sorter struct {
	cmp   func(Record, Record) bool
	frame *Frame
}

// Len is the number of elements in the collection.
func (a sorter) Len() int {
	return a.frame.Len()
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (a sorter) Less(i, j int) bool {
	return a.cmp(record{
		index: i,
		frame: a.frame,
	}, record{
		index: j,
		frame: a.frame,
	})
}

// Swap swaps the elements with indexes i and j.
func (a sorter) Swap(i, j int) {
	for _, list := range a.frame.lists {
		list[i], list[j] = list[j], list[i]
	}
}
