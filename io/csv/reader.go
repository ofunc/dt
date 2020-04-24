package csv

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/ofunc/dt"
	"golang.org/x/text/transform"
)

var regDigits = regexp.MustCompile(`^\d+$`)

// Reader is the CSV reader.
type Reader struct {
	Drop             int
	Tail             int
	Comma            rune
	Comment          rune
	LazyQuotes       bool
	TrimLeadingSpace bool
	Transformer      transform.Transformer
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
	var err error
	r, err = a.reader(r)
	if err != nil {
		return nil, err
	}
	cr := csv.NewReader(r)
	if a.Comma != 0 {
		cr.Comma = a.Comma
	}
	cr.Comment = a.Comment
	cr.LazyQuotes = a.LazyQuotes
	cr.TrimLeadingSpace = a.TrimLeadingSpace

	rs, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}
	rs = rs[a.Drop:]
	rs = cutEmpty(rs)
	rs = rs[:len(rs)-a.Tail]
	if len(rs) < 1 {
		return nil, errors.New("dt/io/csv.Reader: empty data")
	}

	frame := dt.NewFrame(rs[0]...)
	lists := frame.Lists()
	for _, r := range rs[1:] {
		for i, l := range lists {
			lists[i] = append(l, value(r, i))
		}
	}
	return frame, nil
}

func (a Reader) reader(r io.Reader) (io.Reader, error) {
	if a.Transformer != nil {
		r = transform.NewReader(r, a.Transformer)
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
	x := strings.TrimSpace(r[i])
	if len(x) < 16 || !regDigits.MatchString(x) {
		if v, err := strconv.ParseFloat(x, 64); err == nil {
			return dt.Number(v)
		}
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
