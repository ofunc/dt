package xlsx

// Sheet is a sheet.
type Sheet struct {
	ID       string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id"`
	Name     string `xml:"name"`
	target   string
	workbook *Workbook
	data     *Data
}
