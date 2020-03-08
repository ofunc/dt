package xlsx

// Row is a row.
type Row struct {
	Ref      int     `xml:"r,attr"`
	Cells    []*Cell `xml:"c"`
	workbook *Workbook
}
