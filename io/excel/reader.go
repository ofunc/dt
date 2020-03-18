package excel

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/ofunc/dt"
	"github.com/ofunc/dt/io/xlsx"
)

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

// Reader is the excel reader.
type Reader xlsx.Reader

// ReadFile reads a frame from the file.
func (a Reader) ReadFile(name string) (*dt.Frame, error) {
	if ext := strings.ToLower(filepath.Ext(name)); ext == ".xlsx" {
		return xlsx.Reader(a).ReadFile(name)
	}

	target := tempFileName()
	if err := saveAs(target, name); err != nil {
		return nil, err
	}
	defer os.Remove(target)
	return xlsx.Reader(a).ReadFile(target)
}

func tempFileName() string {
	return filepath.Join(os.TempDir(), fmt.Sprintf("dt_io_excel_%v.xlsx", rnd.Int()))
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
