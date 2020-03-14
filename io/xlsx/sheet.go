package xlsx

import (
	"bytes"
	"encoding/xml"
)

// Sheet is a sheet.
type Sheet struct {
	ID       string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr"`
	Name     string `xml:"name,attr"`
	target   string
	workbook *Workbook
	data     *Data
}

// Data returns the sheet data.
func (a *Sheet) Data() *Data {
	if a.data != nil {
		return a.data
	}
	a.data = new(Data)
	if err := xml.NewDecoder(bytes.NewBuffer(a.workbook.files["xl/"+a.target])).Decode(a.data); err != nil {
		panic(err)
	}
	rows := a.data.Rows
	i := len(rows) - 1
	for ; i >= 0; i-- {
		if row := rows[i]; row != nil && !row.IsEmpty() {
			break
		}
	}
	a.data.Rows = rows[:i+1]
	return a.data
}

// Update updates the sheet data.
func (a *Sheet) Update() error {
	buf := bytes.NewBuffer([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`))
	if err := xml.NewEncoder(buf).Encode(a.data); err != nil {
		return err
	}
	a.workbook.files["xl/"+a.target] = buf.Bytes()
	return nil
}
