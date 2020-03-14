package xlsx

import (
	"fmt"

	"github.com/ofunc/dt"
)

// Saver is the xlsx saver.
type Saver struct {
	Template string
	File     string
	Sheet    string
}

// Save saves frame to a xlsx file.
func (a Saver) Save(frame *dt.Frame) (err error) {
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
			case dt.Int:
				cell.Type = "n"
			case dt.Float:
				cell.Type = "n"
			case dt.Bool:
				cell.Type = "b"
			case dt.String:
				cell.Type = "inlineStr"
			default:
				cell.Type = "e"
			}
			row.Cells = append(row.Cells, cell)
		}
		rows = append(rows, row)
	}

	workbook, err := OpenFile(a.Template)
	if err != nil {
		return err
	}
	sheet := workbook.Sheet(a.Sheet)
	sheet.Data().Rows = rows
	if err := sheet.Update(); err != nil {
		return err
	}
	return workbook.SaveFile(a.File)
}
