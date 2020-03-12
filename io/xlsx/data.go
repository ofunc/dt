package xlsx

// Data is the sheet data.
type Data struct {
	Rows []*Row `xml:"sheetData>row"`
}

// RowIter returns the row iter.
func (a *Data) RowIter() *RowIter {
	return &RowIter{
		i:    -1,
		j:    -1,
		r:    -1,
		rows: a.Rows,
	}
}
