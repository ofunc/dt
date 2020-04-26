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
	comma       rune
	useCRLF     bool
	transformer transform.Transformer
}

// NewWriter creates a new writer.
func NewWriter() *Writer {
	return new(Writer)
}

// Comma is the comma option.
func (a *Writer) Comma(o rune) *Writer {
	a.comma = o
	return a
}

// UseCRLF is the use CRLF option.
func (a *Writer) UseCRLF(o bool) *Writer {
	a.useCRLF = o
	return a
}

// Transformer is the transformer option.
func (a *Writer) Transformer(o transform.Transformer) *Writer {
	a.transformer = o
	return a
}

// WriteFile write the frame to the file.
func (a *Writer) WriteFile(frame *dt.Frame, name string) error {
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
func (a *Writer) Write(frame *dt.Frame, w io.Writer) error {
	if a.transformer != nil {
		w = transform.NewWriter(w, a.transformer)
	}
	cw := csv.NewWriter(w)
	if a.comma != 0 {
		cw.Comma = a.comma
	}
	cw.UseCRLF = a.useCRLF

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
