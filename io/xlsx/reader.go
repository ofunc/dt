package xlsx

import (
	"errors"
	"io"
	"os"
	"strings"

	"github.com/ofunc/dt"
	"github.com/plandem/xlsx"
	"github.com/plandem/xlsx/types"
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
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return a.Read(f)
}

// Read reads a frame from the io.Reader.
func (a Reader) Read(r io.Reader) (*dt.Frame, error) {
	xl, err := xlsx.Open(r)
	if err != nil {
		return nil, err
	}
	defer xl.Close()
	var sheet xlsx.Sheet
	if a.sheet == "" {
		sheet = xl.Sheet(0)
	} else {
		sheet = xl.SheetByName(a.sheet)
	}
	defer sheet.Close()

	nc, nr := sheet.Dimension()
	nr = cutEmpty(sheet, nc, nr) - a.tail
	frame := makeFrame(sheet, nc, nr, a.drop, a.head, a.sep)
	lists := frame.Lists()
	ir := a.drop + a.head
	for i := ir; i < nr; i++ {
		for j := 0; j < nc; j++ {
			lists[j] = append(lists[j], value(sheet.Cell(j, i)))
		}
	}
	return frame, nil
}

func makeFrame(s xlsx.Sheet, nc, nr int, ir int, h int, sep string) *dt.Frame {
	frame := dt.NewFrame()
	for j := 0; j < nc; j++ {
		frame.Set(makeKey(s, nr, ir, j, h, sep), nil)
	}
	return frame
}

func makeKey(s xlsx.Sheet, nr int, ir, ic int, h int, sep string) string {
	key := ""
	n := ir + h
	if n > nr {
		n = nr
	}
	for i := ir; i < n; i++ {
		for j := ic; j >= 0; j-- {
			if v := s.Cell(j, i).String(); v != "" {
				if key == "" {
					key = v
				} else if key != v {
					key = key + sep + v
				}
				break
			}
		}
	}
	return key
}

func cutEmpty(s xlsx.Sheet, nc, nr int) int {
	for i := nr - 1; i >= 0; i-- {
		if !isEmpty(s, i, nc) {
			return i + 1
		}
	}
	return 0
}

func isEmpty(s xlsx.Sheet, i int, nc int) bool {
	for j := 0; j < nc; j++ {
		if strings.TrimSpace(s.Cell(j, i).Value()) != "" {
			return false
		}
	}
	return true
}

func value(c *xlsx.Cell) dt.Value {
	if c == nil {
		return nil
	}
	switch c.Type() {
	case types.CellTypeSharedString, types.CellTypeInlineString:
		return dt.String(c.Value())
	case types.CellTypeBool:
		if v, err := c.Bool(); err == nil {
			return dt.Bool(v)
		}
		return dt.Bool(false)
	default:
		if v, err := c.Int(); err == nil {
			return dt.Int(v)
		}
		if v, err := c.Float(); err == nil {
			return dt.Float(v)
		}
		return dt.String(c.String())
	}
}
