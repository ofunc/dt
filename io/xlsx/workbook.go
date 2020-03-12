package xlsx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"io"
	"io/ioutil"
)

// Workbook is a workbook.
type Workbook struct {
	XMLBase
	XMLName xml.Name
	CalcID  string   `xml:"calcPr>calcId,attr"`
	Sheets  []*Sheet `xml:"sheets>sheet"`
	files   map[string]([]byte)
	rels    Rels
}

// TODO sharedstrings，归属 workbook，打开时一次性读取
// 保存其他所有文件，原样写入

// OpenFile opens the workbook from a file.
func OpenFile(name string) (*Workbook, error) {
	zr, err := zip.OpenReader(name)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	return OpenZipReader(&zr.Reader)
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
	return OpenZipReader(zr)
}

// OpenZipReader opens the workbook from a zip reader.
func OpenZipReader(zr *zip.Reader) (*Workbook, error) {
	files := make(map[string]([]byte))
	for _, f := range zr.File {
		data, err := readZipFile(f)
		if err != nil {
			return nil, err
		}
		files[f.Name] = data
	}
	workbook := new(Workbook)
	if err := xml.NewDecoder(bytes.NewReader(files["xl/workbook.xml"])).Decode(workbook); err != nil {
		return nil, err
	}
	if err := xml.NewDecoder(bytes.NewReader(files["xl/_rels/workbook.xml.rels"])).Decode(&workbook.rels); err != nil {
		return nil, err
	}
	workbook.CalcID = ""
	workbook.files = files
	return workbook, nil
}

// SaveFile save the workbook to a file.
func (a *Workbook) SaveFile(name string) error {
	return nil // TODO
}

// SaveWriter save the workbook to a writer.
func (a *Workbook) SaveWriter(w io.Writer) error {
	return nil // TODO
}

// SaveZipWriter save the workbook to a zip writer.
func (a *Workbook) SaveZipWriter(zw *zip.Writer) error {
	return nil // TODO
}

// Sheet returns the sheet by name.
func (a *Workbook) Sheet(name string) *Sheet {
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