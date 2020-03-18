package xlsx

import (
	"bytes"
	"encoding/xml"
	"errors"
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
	if data, ok := a.workbook.files["xl/"+a.target]; ok {
		if err := xml.NewDecoder(bytes.NewBuffer(data)).Decode(a.data); err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("dt/io/xlsx: invalid xlsx file"))
	}

	rows := a.data.Rows
	for i := len(rows) - 1; i >= 0; i-- {
		if row := rows[i]; row != nil && !row.IsEmpty() {
			a.data.Rows = rows[:i+1]
			break
		}
	}
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
