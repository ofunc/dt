package xlsx

// Data is the sheet data.
type Data struct {
	Rows     []*Row `xml:"row"`
	workbook *Workbook
}
