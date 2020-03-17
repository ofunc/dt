package xlsx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/ofunc/dt"
)

// Workbook is a workbook.
type Workbook struct {
	Sheets []*Sheet `xml:"sheets>sheet"`
	files  map[string]([]byte)
	rels   Rels
	sst    SSTable
}

// OpenFile opens the workbook from a file.
func OpenFile(name string) (*Workbook, error) {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return OpenZip(&zr.Reader)
}

// OpenReader opens the workbook from a reader.
func OpenReader(r io.Reader) (*Workbook, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	zr, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return nil, err
	}
	return OpenZip(zr)
}

// OpenZip opens the workbook from a zip reader.
func OpenZip(zr *zip.Reader) (*Workbook, error) {
	files := make(map[string]([]byte))
	for _, f := range zr.File {
		data, err := readZipFile(f)
		if err != nil {
			return nil, err
		}
		files[f.Name] = data
	}

	workbook := new(Workbook)
	if data, ok := files["xl/workbook.xml"]; ok {
		if err := xml.NewDecoder(bytes.NewReader(data)).Decode(workbook); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("dt/io/xlsx: invalid xlsx file")
	}

	if data, ok := files["xl/_rels/workbook.xml.rels"]; ok {
		if err := xml.NewDecoder(bytes.NewReader(data)).Decode(&workbook.rels); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("dt/io/xlsx: invalid xlsx file")
	}

	if data, ok := files["xl/sharedStrings.xml"]; ok {
		if err := xml.NewDecoder(bytes.NewReader(data)).Decode(&workbook.sst); err != nil {
			return nil, err
		}
	}
	workbook.files = files
	return workbook, nil
}

// WriteFile writes the workbook to a file.
func (a *Workbook) WriteFile(name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := a.Write(f); err != nil {
		return err
	}
	return f.Sync()
}

// Write writes the workbook to a writer.
func (a *Workbook) Write(w io.Writer) error {
	return a.WriteZip(zip.NewWriter(w))
}

// WriteZip writes the workbook to a zip writer.
func (a *Workbook) WriteZip(zw *zip.Writer) error {
	a.files["xl/workbook.xml"] = regCalcID.ReplaceAll(a.files["xl/workbook.xml"], []byte(`<calcPr calcId=""`))
	for name, body := range a.files {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		if _, err := w.Write(body); err != nil {
			return err
		}
	}
	return zw.Close()
}

// Sheet returns the sheet by name.
func (a *Workbook) Sheet(name string) *Sheet {
	if name == "" {
		return a.Sheet(a.Sheets[0].Name)
	}
	for _, sheet := range a.Sheets {
		if sheet.Name == name {
			if sheet.workbook == nil {
				sheet.workbook = a
				for _, rel := range a.rels.Rels {
					if rel.ID == sheet.ID {
						sheet.target = rel.Target
						break
					}
				}
			}
			return sheet
		}
	}
	return nil
}

// Value returns the cell value.
func (a *Workbook) Value(cell *Cell) dt.Value {
	if cell == nil {
		return nil
	}
	switch cell.Type {
	case "e":
		return nil
	case "b":
		return dt.Bool(cell.Value != "0")
	case "s":
		return dt.String(a.sst.Value(cell.Value))
	case "inlineStr":
		return dt.String(cell.Value)
	default:
		x := strings.TrimSpace(cell.Value)
		if len(x) < 16 {
			if v, err := strconv.ParseFloat(x, 64); err == nil {
				return dt.Float(v)
			}
		}
		return dt.String(cell.Value)
	}
}
