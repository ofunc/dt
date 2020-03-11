package xlsx

// RowIter is a row iter.
type RowIter struct {
	i    int
	j    int
	r    int
	rows []*Row
}

// Next checks if has next row.
func (a *RowIter) Next() bool {
	a.i++
	if a.i <= a.r {
		return true
	}
	a.j++
	if a.j >= len(a.rows) {
		return false
	}
	a.r = RowIndex(a.rows[a.j].Ref)
	return true
}

// Row returns the current row.
func (a *RowIter) Row() *Row {
	if a.i < a.r {
		return nil
	}
	return a.rows[a.j]
}
