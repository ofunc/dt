package xlsx

// Data is the sheet data.
type Data struct {
	Rows     []*Row `xml:"sheetData>row"`
	workbook *Workbook
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

// // Value returns the value at index i and j.
// func (a *Data) Value(i, j int) dt.Value {
// 	r := a.Row(i)
// 	if r == nil {
// 		return nil
// 	}
// 	c := r.Cell(j)
// 	if c == nil {
// 		return nil
// 	}
// 	switch c.Type {
// 	case "b":
// 		return dt.Bool(c.Value != "0")
// 	case "s":
// 		// TODO shared string
// 		return dt.String("")
// 	default:
// 		return dt.String(c.Value) // TODO int, float, string 注意精度溢出
// 	}
// }
