package xlsx

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"github.com/ofunc/dt"
)

// Reader is the xlsx reader.
type Reader struct {
	Drop  int
	Head  int
	Tail  int
	Sep   string
	Sheet string
}

// ReadFile reads a frame from the file.
func (a Reader) ReadFile(name string) (*dt.Frame, error) {
	workbook, err := OpenFile(name)
	if err != nil {
		return nil, err
	}
	return a.ReadWorkbook(workbook)
}

// Read reads a frame from the io.Reader.
func (a Reader) Read(r io.Reader) (*dt.Frame, error) {
	workbook, err := OpenReader(r)
	if err != nil {
		return nil, err
	}
	return a.ReadWorkbook(workbook)
}

// ReadZip reads a frame from the zip.Reader.
func (a Reader) ReadZip(zr *zip.Reader) (*dt.Frame, error) {
	workbook, err := OpenZip(zr)
	if err != nil {
		return nil, err
	}
	return a.ReadWorkbook(workbook)
}

// ReadWorkbook reads a frame from the workbook.
func (a Reader) ReadWorkbook(workbook *Workbook) (frame *dt.Frame, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	rowiter := workbook.Sheet(a.Sheet).Data().RowIter()
	for i := 0; i < a.Drop; i++ {
		if !rowiter.Next() {
			break
		}
	}

	head := 1
	if a.Head > 0 {
		head = a.Head
	}
	hs := make([][]string, head)
	for i := range hs {
		if !rowiter.Next() {
			break
		}
		row := rowiter.Row()
		if row != nil {
			for celliter := row.CellIter(); celliter.Next(); {
				if v := workbook.Value(celliter.Cell()); v != nil {
					hs[i] = append(hs[i], v.String())
				} else {
					hs[i] = append(hs[i], "")
				}
			}
		}
	}

	keys := a.makeKeys(cleanHeads(hs))
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
				v = workbook.Value(celliter.Cell())
			}
			lists[i] = append(list, v)
		}
	}

	n := frame.Len() - a.Tail
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
				keys[j] = keys[j] + a.Sep + x
			}
		}
	}
	return keys
}

func cleanHeads(hs [][]string) [][]string {
	for i, h := range hs {
		for j := len(h) - 1; j >= 0; j-- {
			if strings.TrimSpace(h[j]) != "" {
				hs[i] = h[:j+1]
				break
			}
		}
	}
	for i := len(hs) - 1; i >= 0; i-- {
		if len(hs[i]) > 0 {
			hs = hs[:i+1]
			break
		}
	}
	return hs
}
