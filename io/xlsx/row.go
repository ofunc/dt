package xlsx

import (
	"strings"
)

// Row is a row.
type Row struct {
	Ref   string  `xml:"r,attr"`
	Cells []*Cell `xml:"c"`
}

// CellIter returns the cell iter.
func (a *Row) CellIter() *CellIter {
	return &CellIter{
		i:     -1,
		j:     -1,
		c:     -1,
		cells: a.Cells,
	}
}

// IsEmpty checks if a is empty.
func (a *Row) IsEmpty() bool {
	if a == nil {
		return true
	}
	for _, cell := range a.Cells {
		if strings.TrimSpace(cell.Value) != "" {
			return false
		}
	}
	return true
}
