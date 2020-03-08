package xlsx

// Cell is a cell.
type Cell struct {
	Ref      string `xml:"r,attr"`
	Type     string `xml:"t,attr"`
	RawValue string `xml:"v"`
	workbook *Workbook
}
