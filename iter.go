package dt

// Iter is a iter of frame.
type Iter struct {
	index int
	frame *Frame
}

// Next check if has next record.
func (a *Iter) Next() bool {
	a.index++
	return a.index < a.frame.Len()
}

// Record returns the current record.
func (a *Iter) Record() Record {
	return record{
		index: a.index,
		frame: a.frame,
	}
}
