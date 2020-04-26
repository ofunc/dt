package xlsx

import (
	"fmt"

	"github.com/ofunc/dt"
)

// Writer is the xlsx writer.
type Writer struct {
	template string
	filename string
	sheet    string
}

// NewWriter creates a new writer.
func NewWriter(filename string) *Writer {
	return &Writer{
		template: filename,
		filename: filename,
	}
}

// Sheet is the sheet option.
func (a *Writer) Sheet(o string) *Writer {
	a.sheet = o
	return a
}

// Template is the template option.
func (a *Writer) Template(o string) *Writer {
	a.template = o
	return a
}

// WriteFile writes frame to a xlsx file.
func (a *Writer) WriteFile(frame *dt.Frame) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	var rows []*Row
	row := &Row{
		Ref: RowRef(0),
	}
	j := 0
	for _, key := range frame.Keys() {
		cell := &Cell{
			Ref:   CellRef(0, j),
			Type:  "inlineStr",
			Value: key,
		}
		row.Cells = append(row.Cells, cell)
		j++
	}
	rows = append(rows, row)

	lists := frame.Lists()
	n, m := frame.Len(), len(lists)
	for i := 0; i < n; i++ {
		ref := RowRef(i + 1)
		row := &Row{
			Ref: ref,
		}
		for j := 0; j < m; j++ {
			cell := &Cell{
				Ref: ColRef(j) + ref,
			}
			value := lists[j][i]
			if value != nil {
				cell.Value = value.String()
			}
			switch value.(type) {
			case nil:
				cell.Type = "e"
			case dt.Number:
				cell.Type = "n"
			default:
				cell.Type = "inlineStr"
			}
			row.Cells = append(row.Cells, cell)
		}
		rows = append(rows, row)
	}

	workbook, err := OpenFile(a.template)
	if err != nil {
		return err
	}
	sheet := workbook.sheet(a.sheet)
	sheet.data().Rows = rows
	if err := sheet.update(); err != nil {
		return err
	}
	return workbook.WriteFile(a.filename)
}
