package xlsx

import (
	"archive/zip"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ofunc/dt"
	helper "github.com/ofunc/dt/io"
)

// Reader is the xlsx reader.
type Reader struct {
	drop   int
	tail   int
	head   int
	sep    string
	sheet  string
	suffix string
}

// NewReader creates a new reader.
func NewReader() *Reader {
	return &Reader{
		head: 1,
	}
}

// Drop is the drop option.
func (a *Reader) Drop(o int) *Reader {
	if o < 0 {
		panic("dt/io/xlsx: invalid drop: " + strconv.Itoa(o))
	}
	a.drop = o
	return a
}

// Tail is the tail option.
func (a *Reader) Tail(o int) *Reader {
	if o < 0 {
		panic("dt/io/xlsx: invalid tail: " + strconv.Itoa(o))
	}
	a.tail = o
	return a
}

// Head is the head option.
func (a *Reader) Head(o int) *Reader {
	if o < 1 {
		panic("dt/io/xlsx: invalid head: " + strconv.Itoa(o))
	}
	a.drop = o
	return a
}

// Sep is the sep option.
func (a *Reader) Sep(o string) *Reader {
	a.sep = o
	return a
}

// Sheet is the sheet option.
func (a *Reader) Sheet(o string) *Reader {
	a.sheet = o
	return a
}

// Suffix is the suffix option.
func (a *Reader) Suffix(o string) *Reader {
	a.suffix = o
	return a
}

// ReadFile reads a frame from the file.
func (a *Reader) ReadFile(name string) (*dt.Frame, error) {
	workbook, err := OpenFile(name)
	if err != nil {
		return nil, err
	}
	return a.ReadWorkbook(workbook)
}

// Read reads a frame from the io.Reader.
func (a *Reader) Read(r io.Reader) (*dt.Frame, error) {
	workbook, err := OpenReader(r)
	if err != nil {
		return nil, err
	}
	return a.ReadWorkbook(workbook)
}

// ReadZip reads a frame from the zip.Reader.
func (a *Reader) ReadZip(zr *zip.Reader) (*dt.Frame, error) {
	workbook, err := OpenZip(zr)
	if err != nil {
		return nil, err
	}
	return a.ReadWorkbook(workbook)
}

// ReadWorkbook reads a frame from the workbook.
func (a *Reader) ReadWorkbook(workbook *Workbook) (frame *dt.Frame, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("%v", e)
		}
	}()

	rowiter := workbook.sheet(a.sheet).data().rowIter()
	for i := 0; i < a.drop; i++ {
		if !rowiter.next() {
			break
		}
	}

	hs := make([][]string, a.head)
	for i := range hs {
		if !rowiter.next() {
			break
		}
		row := rowiter.row()
		if row != nil {
			for celliter := row.cellIter(); celliter.next(); {
				if v := workbook.value(celliter.cell()); v != nil {
					hs[i] = append(hs[i], v.String())
				} else {
					hs[i] = append(hs[i], "")
				}
			}
		}
	}

	keys := helper.Keys(a.makeKeys(cleanHeads(hs)), a.suffix)
	frame = dt.NewFrame(keys...)
	lists := frame.Lists()
	for rowiter.next() {
		row := rowiter.row()
		if row == nil {
			continue
		}
		celliter := row.cellIter()
		for i, list := range lists {
			var v dt.Value
			if celliter.next() {
				v = workbook.value(celliter.cell())
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

func (a *Reader) makeKeys(hs [][]string) []string {
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
