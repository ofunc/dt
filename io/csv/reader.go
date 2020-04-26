package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/ofunc/dt"
	helper "github.com/ofunc/dt/io"
	"golang.org/x/text/transform"
)

// Reader is the CSV reader.
type Reader struct {
	drop             int
	tail             int
	comma            rune
	comment          rune
	lazyQuotes       bool
	trimLeadingSpace bool
	suffix           string
	transformer      transform.Transformer
}

// NewReader creates a new reader.
func NewReader() *Reader {
	return new(Reader)
}

// Drop is the drop option.
func (a *Reader) Drop(o int) *Reader {
	if o < 0 {
		panic("dt/io/csv: invalid drop: " + strconv.Itoa(o))
	}
	a.drop = o
	return a
}

// Tail is the tail option.
func (a *Reader) Tail(o int) *Reader {
	if o < 0 {
		panic("dt/io/csv: invalid tail: " + strconv.Itoa(o))
	}
	a.tail = o
	return a
}

// Comma is the comma option.
func (a *Reader) Comma(o rune) *Reader {
	a.comma = o
	return a
}

// Comment is the comment option.
func (a *Reader) Comment(o rune) *Reader {
	a.comment = o
	return a
}

// LazyQuotes is the lazy quotes option.
func (a *Reader) LazyQuotes(o bool) *Reader {
	a.lazyQuotes = o
	return a
}

// TrimLeadingSpace is the trim leading space option.
func (a *Reader) TrimLeadingSpace(o bool) *Reader {
	a.trimLeadingSpace = o
	return a
}

// Suffix is the suffix quotes option.
func (a *Reader) Suffix(o string) *Reader {
	a.suffix = o
	return a
}

// Transformer is the transformer quotes option.
func (a *Reader) Transformer(o transform.Transformer) *Reader {
	a.transformer = o
	return a
}

// ReadFile reads a frame from the file.
func (a *Reader) ReadFile(name string) (*dt.Frame, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return a.Read(f)
}

// Read reads a frame from the io.Reader.
func (a *Reader) Read(r io.Reader) (*dt.Frame, error) {
	var err error
	r, err = a.reader(r)
	if err != nil {
		return nil, err
	}
	cr := csv.NewReader(r)
	if a.comma != 0 {
		cr.Comma = a.comma
	}
	cr.Comment = a.comment
	cr.LazyQuotes = a.lazyQuotes
	cr.TrimLeadingSpace = a.trimLeadingSpace

	rs, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	rs = rs[a.drop:]
	rs = cutEmpty(rs)
	rs = rs[:len(rs)-a.tail]
	if len(rs) < 1 {
		return nil, errors.New("dt/io/csv.Reader: empty data")
	}

	frame := dt.NewFrame(helper.Keys(rs[0], a.suffix)...)
	lists := frame.Lists()
	for _, r := range rs[1:] {
		for i, l := range lists {
			lists[i] = append(l, value(r, i))
		}
	}
	return frame, nil
}

func (a *Reader) reader(r io.Reader) (io.Reader, error) {
	if a.transformer != nil {
		r = transform.NewReader(r, a.transformer)
	}
	br := bufio.NewReader(r)
	xs, err := br.Peek(3)
	if err != nil {
		return nil, err
	}
	if bytes.Equal(xs, []byte{0xef, 0xbb, 0xbf}) {
		if _, err := br.Discard(3); err != nil {
			return nil, err
		}
	}
	return br, nil
}

func value(r []string, i int) dt.Value {
	if i >= len(r) {
		return nil
	}
	return helper.Value(r[i])
}

func cutEmpty(rs [][]string) [][]string {
	i := len(rs) - 1
	for ; i >= 0; i-- {
		if !isEmpty(rs[i]) {
			break
		}
	}
	return rs[:i+1]
}

func isEmpty(r []string) bool {
	for _, v := range r {
		if strings.TrimSpace(v) != "" {
			return false
		}
	}
	return true
}
