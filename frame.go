package dt

import (
	"bytes"
	"errors"
	"fmt"
	"sort"
)

// Frame is the frame data structure.
type Frame struct {
	index map[string]int
	lists []List
}

// NewFrame creates a new frame.
func NewFrame(keys ...string) *Frame {
	n := len(keys)
	index := make(map[string]int, n)
	for j, key := range keys {
		if _, ok := index[key]; ok {
			panic("dt: key already exists: " + key)
		}
		index[key] = j
	}
	return &Frame{
		index: index,
		lists: make([]List, n),
	}
}

// Empty returns a empty frame like frame a.
func (a *Frame) Empty() *Frame {
	index := make(map[string]int, len(a.lists))
	for key, j := range a.index {
		index[key] = j
	}
	return &Frame{
		index: index,
		lists: make([]List, len(a.lists)),
	}
}

// Copy makes a copy of frame a.
func (a *Frame) Copy(deep bool) *Frame {
	b := a.Empty()
	copy(b.lists, a.lists)
	if deep {
		for j, l := range b.lists {
			t := make(List, len(l))
			copy(t, l)
			b.lists[j] = t
		}
	}
	return b
}

// Keys returns the keys of frame a.
func (a *Frame) Keys() []string {
	keys := make([]string, len(a.lists))
	for key, j := range a.index {
		keys[j] = key
	}
	return keys
}

// Lists returns the lists of frame a.
func (a *Frame) Lists() []List {
	return a.lists
}

// Len returns the length of frame a.
func (a *Frame) Len() int {
	if len(a.lists) == 0 {
		return 0
	}
	return len(a.lists[0])
}

// Has checks if the frame has the keys.
func (a *Frame) Has(keys ...string) bool {
	for _, key := range keys {
		if _, ok := a.index[key]; !ok {
			return false
		}
	}
	return true
}

// Check checks if the frame has the keys.
func (a *Frame) Check(keys ...string) error {
	for _, key := range keys {
		if _, ok := a.index[key]; !ok {
			return errors.New("dt: key not found: " + key)
		}
	}
	return nil
}

// Get gets the list by key.
func (a *Frame) Get(key string) List {
	if j, ok := a.index[key]; ok {
		return a.lists[j]
	}
	panic("dt: key not found: " + key)
}

// Set sets the list by key.
func (a *Frame) Set(key string, list List) *Frame {
	a.check(list)
	if j, ok := a.index[key]; ok {
		a.lists[j] = list
		return a
	}
	a.index[key] = len(a.lists)
	a.lists = append(a.lists, list)
	return a
}

// Add adds the list with key.
func (a *Frame) Add(key string, list List) *Frame {
	if _, ok := a.index[key]; ok {
		panic("dt: key already exists: " + key)
	}
	a.check(list)
	a.index[key] = len(a.lists)
	a.lists = append(a.lists, list)
	return a
}

// Del deletes the list by keys.
func (a *Frame) Del(keys ...string) *Frame {
	for _, key := range keys {
		a.del(key)
	}
	return a
}

// Rename renames the key.
func (a *Frame) Rename(old, new string) *Frame {
	if j, ok := a.index[old]; ok {
		if _, ok := a.index[new]; ok {
			panic("dt: key already exists: " + new)
		}
		delete(a.index, old)
		a.index[new] = j
	} else {
		panic("dt: key not found: " + old)
	}
	return a
}

// Pick picks some lists and returns a new frame,
func (a *Frame) Pick(key string, keys ...string) *Frame {
	b := NewFrame()
	b.Set(key, a.Get(key))
	for _, key := range keys {
		b.Set(key, a.Get(key))
	}
	return b
}

// Iter returns a iter of frame a.
func (a *Frame) Iter() *Iter {
	return &Iter{
		frame: a,
		index: -1,
	}
}

// Slice gets the slice of frame a.
func (a *Frame) Slice(i, j int) *Frame {
	n := a.Len()
	if i < 0 {
		i += n
	}
	if j < 0 {
		j += n
	}
	b := a.Copy(false)
	for k, list := range b.lists {
		b.lists[k] = list[i:j]
	}
	return b
}

// Concat concats frame a with b.
func (a *Frame) Concat(b *Frame) *Frame {
	for key, j := range a.index {
		a.lists[j] = append(a.lists[j], b.Get(key)...)
	}
	return a
}

// Append appends x to frames a.
func (a *Frame) Append(rs ...Record) *Frame {
	for key, j := range a.index {
		for _, r := range rs {
			a.lists[j] = append(a.lists[j], r.Value(key))
		}
	}
	return a
}

// Sort sorts frame a by function f.
func (a *Frame) Sort(f func(Record, Record) bool) *Frame {
	sort.Sort(sorter{
		frame: a,
		cmp:   f,
	})
	return a
}

// Map maps frame a to list by function f.
func (a *Frame) Map(f func(Record) Value) List {
	list := make(List, 0, a.Len())
	for iter := a.Iter(); iter.Next(); {
		list = append(list, f(iter.Record()))
	}
	return list
}

// MapTo maps frame a to the key list.
func (a *Frame) MapTo(key string, f func(Record) Value) *Frame {
	list := a.Map(f)
	a.Set(key, list)
	return a
}

// Filter filters the frame with function f.
func (a *Frame) Filter(f func(Record) bool) *Frame {
	b := a.Empty()
	for iter := a.Iter(); iter.Next(); {
		r := iter.Record().(record)
		if f(r) {
			for j, l := range b.lists {
				b.lists[j] = append(l, a.lists[j][r.index])
			}
		}
	}
	return b
}

// DropNA drops NA value.
func (a *Frame) DropNA(keys ...string) *Frame {
	if len(keys) == 0 {
		keys = a.Keys()
	}
	return a.Filter(func(r Record) bool {
		for _, key := range keys {
			if IsNA(r.Value(key)) {
				return false
			}
		}
		return true
	})
}

// FillNA fills NA value with v.
func (a *Frame) FillNA(value Value, keys ...string) *Frame {
	if len(keys) == 0 {
		keys = a.Keys()
	}
	for _, key := range keys {
		a.Get(key).FillNA(value)
	}
	return a
}

// Join joins frame a and b.
func (a *Frame) Join(b *Frame, key string, keys ...string) *Join {
	return &Join{
		lframe: a,
		rframe: b,
		rkeys:  append(keys, key),
	}
}

// GroupBy groups records by keys.
func (a *Frame) GroupBy(key string, keys ...string) *Group {
	m := len(keys)
	keys = append(keys, key)
	keys[0], keys[m] = keys[m], keys[0]

	lists := make([]List, len(keys))
	for j, key := range keys {
		lists[j] = a.Get(key)
	}
	data := make(map[string]([]int))
	for i, n := 0, a.Len(); i < n; i++ {
		k := makeKey(i, lists)
		data[k] = append(data[k], i)
	}
	g := &Group{
		frame: a,
		data:  data,
	}
	for _, key := range keys {
		g.Apply(key, key, First)
	}
	return g
}

// String shows frame a as string.
func (a *Frame) String() string {
	buf := new(bytes.Buffer)
	m := len(a.lists) - 1
	for j, key := range a.Keys() {
		fmt.Fprint(buf, key)
		if j == m {
			fmt.Fprintln(buf)
		} else {
			fmt.Fprint(buf, "\t")
		}
	}
	n := a.Len()
	for i := 0; i < n; i++ {
		for j, l := range a.lists {
			fmt.Fprint(buf, l[i])
			if j == m {
				fmt.Fprintln(buf)
			} else {
				fmt.Fprint(buf, "\t")
			}
		}
	}
	return buf.String()
}

func (a *Frame) check(list List) {
	if n, m := a.Len(), len(list); n != m {
		panic(fmt.Sprintf("dt: invalid list length, expected %v, got %v", n, m))
	}
}

func (a *Frame) del(key string) {
	if j, ok := a.index[key]; ok {
		delete(a.index, key)
		copy(a.lists[j:], a.lists[j+1:])
		a.lists = a.lists[:len(a.lists)-1]
		for key, k := range a.index {
			if k > j {
				a.index[key] = k - 1
			}
		}
	}
}
