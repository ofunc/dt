package xlsx

import (
	"errors"
	"fmt"
	"io"

	"github.com/ofunc/dt"
)

// Reader is the xlsx reader.
type Reader struct {
	drop  int
	head  int
	tail  int
	sep   string
	sheet string
}

// NewReader creates a new reader.
func NewReader(sep string) Reader {
	return Reader{
		head: 1,
		sep:  sep,
	}
}

// Drop drops the first n records.
func (a Reader) Drop(n int) Reader {
	if n < 0 {
		panic(errors.New("dt/io/xlsx.Reader.Drop: invalid arguments"))
	}
	a.drop = n
	return a
}

// Head sets the head lines.
func (a Reader) Head(n int) Reader {
	if n < 1 {
		panic(errors.New("dt/io/xlsx.Reader.Head: invalid arguments"))
	}
	a.head = n
	return a
}

// Tail sets the tail lines.
func (a Reader) Tail(n int) Reader {
	if n < 0 {
		panic(errors.New("dt/io/xlsx.Reader.Tail: invalid arguments"))
	}
	a.tail = n
	return a
}

// Sheet sets the sheet.
func (a Reader) Sheet(s string) Reader {
	a.sheet = s
	return a
}

// ReadFile reads a frame from the file.
func (a Reader) ReadFile(name string) (*dt.Frame, error) {
	workbook, err := OpenFile(name)
	if err != nil {
		return nil, err
	}
	return a.read(workbook)
}

// Read reads a frame from the io.Reader.
func (a Reader) Read(r io.Reader) (*dt.Frame, error) {
	workbook, err := OpenReader(r)
	if err != nil {
		return nil, err
	}
	return a.read(workbook)
}

func (a Reader) read(workbook *Workbook) (frame *dt.Frame, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	rowiter := workbook.Sheet(a.sheet).Data().RowIter()
	for i := 0; i < a.drop; i++ {
		if !rowiter.Next() {
			break
		}
	}

	hs := make([][]string, a.head)
	for i := range hs {
		if !rowiter.Next() {
			break
		}
		row := rowiter.Row()
		if row != nil {
			for celliter := row.CellIter(); celliter.Next(); {
				cell := celliter.Cell()
				if cell == nil {
					hs[i] = append(hs[i], "")
				} else if cell.Type == "s" {
					hs[i] = append(hs[i], workbook.sst.Value(cell.Value))
				} else {
					hs[i] = append(hs[i], cell.Value)
				}
			}
		}
	}

	keys := a.makeKeys(hs)
	frame = dt.NewFrame(keys...)
	lists := frame.Lists()
	for rowiter.Next() {
		row := rowiter.Row()
		if row == nil {
			continue
		}
		celliter := row.CellIter()
		for i, list := range lists {
			var v dt.Value
			if celliter.Next() {
				if cell := celliter.Cell(); cell != nil {
					v = cellValue(workbook, cell)
				}
			}
			lists[i] = append(list, v)
		}
	}

	n := frame.Len() - a.tail
	if n < 0 {
		n = 0
	}
	for i, list := range lists {
		lists[i] = list[:n]
	}
	return
}

func (a Reader) makeKeys(hs [][]string) []string {
	n := len(hs[len(hs)-1])
	keys := make([]string, n)
	xs := make([]string, len(hs))
	for j := 0; j < n; j++ {
		ok := false
		for i, h := range hs {
			if j < len(h) && h[j] != "" {
				xs[i] = h[j]
				ok = true
			} else if ok {
				xs[i] = ""
			}
		}
		for _, x := range xs {
			if keys[j] == "" {
				keys[j] = x
			} else if x != "" {
				keys[j] = keys[j] + a.sep + x
			}
		}
	}
	return keys
}
