package xlsx

import (
	"bytes"
	"encoding/xml"
)

// Sheet is a sheet.
type Sheet struct {
	ID       string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id"`
	Name     string `xml:"name"`
	target   string
	workbook *Workbook
	data     *Data
}

// Data returns the sheet data.
func (a *Sheet) Data() (*Data, error) {
	if a.data != nil {
		return a.data, nil
	}
	a.data = new(Data)
	if err := xml.NewDecoder(bytes.NewBuffer(a.workbook.files["xl/"+a.target])).Decode(a.data); err != nil {
		return nil, err
	}
	// TODO trim empty rows
	a.data.workbook = a.workbook
	return a.data, nil
}

// Update updates the sheet data.
func (a *Sheet) Update() error {
	return nil // TODO
}
