package xlsx

import (
	"bytes"
	"encoding/xml"
	"errors"
)

// Sheet is a sheet.
type Sheet struct {
	ID        string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr"`
	Name      string `xml:"name,attr"`
	target    string
	workbook  *Workbook
	sheetdata *Data
}

func (a *Sheet) data() *Data {
	if a.sheetdata != nil {
		return a.sheetdata
	}
	a.sheetdata = new(Data)
	if data, ok := a.workbook.files["xl/"+a.target]; ok {
		if err := xml.NewDecoder(bytes.NewBuffer(data)).Decode(a.sheetdata); err != nil {
			panic(err)
		}
	} else {
		panic(errors.New("dt/io/xlsx: invalid xlsx file"))
	}

	rows := a.sheetdata.Rows
	for i := len(rows) - 1; i >= 0; i-- {
		if row := rows[i]; row != nil && !row.isEmpty() {
			a.sheetdata.Rows = rows[:i+1]
			break
		}
	}
	return a.sheetdata
}

func (a *Sheet) update() error {
	buf := bytes.NewBuffer([]byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`))
	if err := xml.NewEncoder(buf).Encode(a.sheetdata); err != nil {
		return err
	}
	a.workbook.files["xl/"+a.target] = buf.Bytes()
	return nil
}
