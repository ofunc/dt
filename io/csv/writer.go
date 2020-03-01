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
func NewWriter() Writer {
	return Writer{}
}

// Comma sets the comma.
func (a Writer) Comma(v rune) Writer {
	a.comma = v
	return a
}

// UseCRLF sets the use CRLF.
func (a Writer) UseCRLF(v bool) Writer {
	a.useCRLF = v
	return a
}

// Transformer sets the transformer.
func (a Writer) Transformer(v transform.Transformer) Writer {
	a.transformer = v
	return a
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
	n := frame.Len()
	lists := frame.Lists()
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
