package xlsx

import (
	"errors"
	"strings"
)

// CellIter is a cell iter.
type CellIter struct {
	i     int
	j     int
	c     int
	cells []*Cell
}

// Next checks if has next cell.
func (a *CellIter) Next() bool {
	a.i++
	if a.i <= a.c {
		return true
	}
	a.j++
	if a.j >= len(a.cells) {
		return false
	}
	_, a.c = CellIndex(strings.ToUpper(a.cells[a.j].Ref))
	if a.i > a.c {
		panic(errors.New("dt/io/xlsx: invalid xlsx file"))
	}
	return true
}

// Cell returns the current cell.
func (a *CellIter) Cell() *Cell {
	if a.i < a.c {
		return nil
	}
	return a.cells[a.j]
}
