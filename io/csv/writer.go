package csv

import (
	"encoding/csv"
	"io"
	"os"

	"github.com/ofunc/dt"
	"golang.org/x/text/transform"
)

// Writer is the CSV writer.
type Writer struct {
	Comma       rune
	UseCRLF     bool
	Transformer transform.Transformer
}

// WriteFile write the frame to the file.
func (a Writer) WriteFile(frame *dt.Frame, name string) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := a.Write(frame, f); err != nil {
		return err
	}
	return f.Sync()
}

// WriteFile write the frame to the io.Writer.
func (a Writer) Write(frame *dt.Frame, w io.Writer) error {
	if a.Transformer != nil {
		w = transform.NewWriter(w, a.Transformer)
	}
	cw := csv.NewWriter(w)
	if a.Comma != 0 {
		cw.Comma = a.Comma
	}
	cw.UseCRLF = a.UseCRLF

	if err := cw.Write(frame.Keys()); err != nil {
		return err
	}
	n, lists := frame.Len(), frame.Lists()
	r := make([]string, len(lists))
	for i := 0; i < n; i++ {
		for j, l := range lists {
			if v := l[i]; v == nil {
				r[j] = ""
			} else {
				r[j] = v.String()
			}
		}
		if err := cw.Write(r); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}
