package xlsx

import (
	"archive/zip"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/ofunc/dt"
)

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

var regCalcID = regexp.MustCompile(`<\s*calcPr\s+calcId\s*=\s*"\d*"`)
var regDigits = regexp.MustCompile(`^\d+$`)

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

// Time converts a dt.Value to time.Time.
func Time(v dt.Value) time.Time {
	return time.Unix(int64((v.Number()-25569)*86400), 0).UTC()
}

// Value converts a time.Time to dt.Value.
func Value(v time.Time) dt.Value {
	return dt.Number(v.Unix())/86400 + 25569
}

func readZipFile(f *zip.File) ([]byte, error) {
	r, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return ioutil.ReadAll(r)
}

func saveAs(tar, src string) error {
	ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	defer ole.CoUninitialize()

	unknown, err := oleutil.CreateObject("Excel.Application")
	if err != nil {
		return err
	}
	app, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return err
	}
	defer app.Release()
	defer app.CallMethod("Quit")

	if _, err := app.PutProperty("DisplayAlerts", false); err != nil {
		return err
	}
	if _, err := app.PutProperty("Visible", false); err != nil {
		return err
	}

	olevar, err := app.GetProperty("Workbooks")
	if err != nil {
		return err
	}
	workbooks := olevar.ToIDispatch()

	src, err = filepath.Abs(src)
	if err != nil {
		return err
	}
	olevar, err = workbooks.CallMethod("Open", src)
	if err != nil {
		return err
	}
	file := olevar.ToIDispatch()

	if _, err := file.CallMethod("SaveAs", tar, 51); err != nil {
		return err
	}
	if _, err := file.CallMethod("Close"); err != nil {
		return err
	}
	return nil
}
