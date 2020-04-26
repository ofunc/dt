package xlsx

// RowIter is a row iter.
type RowIter struct {
	i    int
	j    int
	r    int
	rows []*Row
}

func (a *RowIter) next() bool {
	a.i++
	if a.i <= a.r {
		return true
	}
	a.j++
	if a.j >= len(a.rows) {
		return false
	}
	a.r = RowIndex(a.rows[a.j].Ref)
	if a.i > a.r {
		panic("dt/io/xlsx: invalid xlsx file")
	}
	return true
}

func (a *RowIter) row() *Row {
	if a.i < a.r {
		return nil
	}
	return a.rows[a.j]
}
