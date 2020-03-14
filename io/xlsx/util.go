package xlsx

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/ofunc/dt"
)

// RowRef returns the row ref by index.
func RowRef(i int) string {
	if i < 0 {
		panic(fmt.Errorf("dt/io/xlsx: invalid row index: %v", i))
	}
	return strconv.Itoa(i + 1)
}

// RowIndex returns the row index by ref.
func RowIndex(r string) int {
	i, err := strconv.Atoi(r)
	if err != nil {
		panic(err)
	}
	if i < 1 {
		panic(fmt.Errorf("dt/io/xlsx: invalid row ref: %v", r))
	}
	return i - 1
}

// ColRef returns the col ref by index.
func ColRef(i int) string {
	if i < 0 {
		panic(fmt.Errorf("dt/io/xlsx: invalid col index: %v", i))
	}
	if i < 26 {
		return string('A' + i)
	}
	return ColRef(i/26-1) + ColRef(i%26)
}

// ColIndex returns the col index by ref.
func ColIndex(r string) int {
	n := len(r)
	if n < 1 {
		panic(errors.New("dt/io/xlsx: empty col ref"))
	}
	if n == 1 {
		c := r[0]
		if c < 'A' || c > 'Z' {
			panic(fmt.Errorf("dt/io/xlsx: invalid col ref: %v", r))
		}
		return int(c - 'A')
	}
	return ColIndex(r[n-1:n]) + 26*(ColIndex(r[:n-1])+1)
}

// CellRef returns the cell ref by index.
func CellRef(i, j int) string {
	return ColRef(j) + RowRef(i)
}

// CellIndex returns the cell index by ref.
func CellIndex(r string) (int, int) {
	var k int
	n := len(r)
	for k = 0; k < n; k++ {
		if r[k] < 'A' || r[k] > 'Z' {
			break
		}
	}
	return RowIndex(r[k:]), ColIndex(r[:k])
}

func readZipFile(f *zip.File) ([]byte, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func cellValue(workbook *Workbook, cell *Cell) dt.Value {
	if cell == nil {
		return nil
	}
	switch cell.Type {
	case "b":
		return dt.Bool(cell.Value != "0")
	case "inlineStr":
		return dt.String(cell.Value)
	case "s":
		return dt.String(cell.Value) // TODO
	default:
		s := strings.TrimSpace(cell.Value)
		if v, err := strconv.Atoi(s); err == nil {
			return dt.Int(v)
		}
		if len(s) <= 15 {
			if v, err := strconv.ParseFloat(s, 64); err == nil {
				return dt.Float(v)
			}
		}
		return dt.String(cell.Value)
	}
}

func emptyRow(row *Row) bool {
	if row == nil {
		return true
	}
	for _, cell := range row.Cells {
		if strings.TrimSpace(cell.Value) != "" {
			return false // TODO 注意sharestring ""
		}
	}
	return true
}
