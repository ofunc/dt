package xlsx

import (
	"errors"
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

func (a Reader) read(workbook *Workbook) (*dt.Frame, error) {
	// TODO recover panic
	// sheet := workbook.Sheet(a.sheet)
	// data, err := sheet.Data()
	// if err != nil {
	// 	return nil, err
	// }
	// frame := dt.NewFrame()
	// nr := data.Len() - a.tail
	// nc := data.Row(a.drop + a.head - 1).Len()
	// for j := 0; j < nc; j++ {
	// 	frame.Set(a.makeKey(data, j), nil)
	// }
	// lists := frame.Lists()
	// for i := a.drop + a.head; i < nr; i++ {
	// 	for j := 0; j < nc; j++ {
	// 		lists[j] = append(lists[j], data.Value(i, j))
	// 	}
	// }
	return nil, nil
}

// func (a Reader) makeKeys(hs [][]string) []string {
// 	var key string
// 	n := a.drop + a.head
// 	for i := a.drop; i < n; i++ {
// 		for j := c; j >= 0; j-- {
// 			if v := data.Value(i, j); v != nil {
// 				if k := v.String(); k != "" {
// 					if key == "" {
// 						key = k
// 					} else {
// 						key = key + a.sep + k
// 					}
// 					break
// 				}
// 			}
// 		}
// 	}
// 	return key
// }
