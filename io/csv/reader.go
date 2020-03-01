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
	"golang.org/x/text/transform"
)

// Reader is the CSV reader.
type Reader struct {
	drop             int
	head             int
	tail             int
	sep              string
	comma            rune
	comment          rune
	lazyQuotes       bool
	trimLeadingSpace bool
	transformer      transform.Transformer
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
		panic(errors.New("dt/io/csv.Reader.Drop: invalid arguments"))
	}
	a.drop = n
	return a
}

// Head sets the head lines.
func (a Reader) Head(n int) Reader {
	if n < 1 {
		panic(errors.New("dt/io/csv.Reader.Head: invalid arguments"))
	}
	a.head = n
	return a
}

// Tail sets the tail lines.
func (a Reader) Tail(n int) Reader {
	if n < 0 {
		panic(errors.New("dt/io/csv.Reader.Tail: invalid arguments"))
	}
	a.tail = n
	return a
}

// Comma sets the comma.
func (a Reader) Comma(v rune) Reader {
	a.comma = v
	return a
}

// Comment sets the comment.
func (a Reader) Comment(v rune) Reader {
	a.comment = v
	return a
}

// LazyQuotes sets the lazy quotes.
func (a Reader) LazyQuotes(v bool) Reader {
	a.lazyQuotes = v
	return a
}

// TrimLeadingSpace sets the trim leading space.
func (a Reader) TrimLeadingSpace(v bool) Reader {
	a.trimLeadingSpace = v
	return a
}

// Transformer sets the transformer.
func (a Reader) Transformer(v transform.Transformer) Reader {
	a.transformer = v
	return a
}

// ReadFile reads a frame from the file.
func (a Reader) ReadFile(name string) (*dt.Frame, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return a.Read(f)
}

// Read reads a frame from the io.Reader.
func (a Reader) Read(r io.Reader) (*dt.Frame, error) {
	cr := csv.NewReader(a.reader(r))
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

	frame := makeFrame(rs[:a.head], a.sep)
	lists := frame.Lists()
	for _, r := range rs[a.head:] {
		for i, l := range lists {
			lists[i] = append(l, value(r, i))
		}
	}
	return frame, nil
}

func (a Reader) reader(r io.Reader) io.Reader {
	if a.transformer != nil {
		r = transform.NewReader(r, a.transformer)
	}
	br := bufio.NewReader(r)
	xs, err := br.Peek(3)
	if err != nil {
		panic(err)
	}
	if bytes.Equal(xs, []byte{0xef, 0xbb, 0xbf}) {
		if _, err := br.Discard(3); err != nil {
			panic(err)
		}
	}
	return br
}

func value(r []string, i int) dt.Value {
	if i >= len(r) {
		return nil
	}
	x := strings.TrimSpace(r[i])
	if v, err := strconv.ParseInt(x, 10, 64); err == nil {
		return dt.Int(v)
	}
	if v, err := strconv.ParseFloat(x, 64); err == nil {
		return dt.Float(v)
	}
	if x == "true" || x == "TRUE" {
		return dt.Bool(true)
	}
	if x == "false" || x == "FALSE" {
		return dt.Bool(false)
	}
	return dt.String(r[i])
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

func makeFrame(rs [][]string, sep string) *dt.Frame {
	r := rs[len(rs)-1]
	frame := dt.NewFrame()
	for i := range r {
		frame.Set(makeKey(rs, i, sep), nil)
	}
	return frame
}

func makeKey(rs [][]string, i int, sep string) string {
	key := ""
	for _, r := range rs {
		n := len(r)
		for j := i; j >= 0; j-- {
			if j < n && r[j] != "" {
				if key == "" {
					key = r[j]
				} else {
					key = key + sep + r[j]
				}
				break
			}
		}
	}
	return key
}
