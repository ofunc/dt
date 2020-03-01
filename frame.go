package dt

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

// Frame is the frame data structure.
type Frame struct {
	index map[string]int
	lists []List
}

// NewFrame creates a new frame.
func NewFrame() *Frame {
	return &Frame{
		index: make(map[string]int),
	}
}

// Like returns a empty frame like frame a.
func (a *Frame) Like() *Frame {
	index := make(map[string]int, len(a.lists))
	for key, i := range a.index {
		index[key] = i
	}
	return &Frame{
		index: index,
		lists: make([]List, len(a.lists)),
	}
}

// Copy makes a copy of frame a.
func (a *Frame) Copy(deep bool) *Frame {
	index := make(map[string]int, len(a.lists))
	lists := make([]List, len(a.lists))
	copy(lists, a.lists)
	if deep {
		for i, l := range lists {
			t := make(List, len(l))
			copy(t, l)
			lists[i] = t
		}
	}
	for key, i := range a.index {
		index[key] = i
	}
	return &Frame{
		index: index,
		lists: lists,
	}
}

// Keys returns the keys of frame a.
func (a *Frame) Keys() []string {
	keys := make([]string, len(a.lists))
	for key, i := range a.index {
		keys[i] = key
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

// Get gets the list by key.
func (a *Frame) Get(key string) List {
	if i, ok := a.index[key]; ok {
		return a.lists[i]
	}
	return make(List, a.Len())
}

// Set sets the list by key.
func (a *Frame) Set(key string, list List) {
	a.check(list)
	if i, ok := a.index[key]; ok {
		a.lists[i] = list
	}
	a.index[key] = len(a.lists)
	a.lists = append(a.lists, list)
}

// Add adds the list with key.
func (a *Frame) Add(key string, list List) {
	if _, ok := a.index[key]; ok {
		a.Add(key+" ", list)
		return
	}
	a.check(list)
	a.index[key] = len(a.lists)
	a.lists = append(a.lists, list)
}

// Del deletes the list by key.
func (a *Frame) Del(key string) List {
	if i, ok := a.index[key]; ok {
		list := a.lists[i]
		delete(a.index, key)
		copy(a.lists[i:], a.lists[i+1:])
		a.lists = a.lists[:len(a.lists)-1]
		return list
	}
	return nil
}

// Pick picks some lists and returns a new frame,
func (a *Frame) Pick(keys ...string) *Frame {
	if len(keys) == 0 {
		return a
	}
	b := NewFrame()
	for _, key := range keys {
		b.Set(key, a.Get(key))
	}
	return b
}

// Iter returns a iter of frame a.
func (a *Frame) Iter() *Iter {
	return &Iter{
		index: -1,
		frame: a,
	}
}

// Slice gets the slice of frame a.
func (a *Frame) Slice(i, j int) *Frame {
	b := a.Copy(false)
	for i, list := range b.lists {
		b.lists[i] = list[i:j]
	}
	return b
}

// Concat concats frame a with b.
func (a *Frame) Concat(b *Frame) *Frame {
	for key, i := range a.index {
		a.lists[i] = append(a.lists[i], b.Get(key)...)
	}
	return a
}

// Append appends x to frames a.
func (a *Frame) Append(rs ...Record) *Frame {
	for key, i := range a.index {
		for _, r := range rs {
			a.lists[i] = append(a.lists[i], r.Value(key))
		}
	}
	return a
}

// Sort sorts frame a by function f.
func (a *Frame) Sort(f func(Record, Record) bool) *Frame {
	sort.Sort(sorter{
		cmp:   f,
		frame: a,
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

// Filter filters the frame with function f.
func (a *Frame) Filter(f func(Record) bool) *Frame {
	b := a.Like()
	for iter := a.Iter(); iter.Next(); {
		r := iter.Record().(record)
		if f(r) {
			for i, l := range b.lists {
				b.lists[i] = append(l, a.lists[i][r.index])
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
func (a *Frame) Join(b *Frame) Join {
	return Join{
		lframe: a,
		rframe: b,
	}
}

// Group groups records by keys.
func (a *Frame) Group(keys ...string) Group {
	data := make(map[interface{}]([]int))
	typ := reflect.ArrayOf(len(keys), tvalue)
	for iter := a.Iter(); iter.Next(); {
		r := iter.Record().(record)
		k := makeKey(typ, r, keys)
		data[k] = append(data[k], r.index)
	}
	return Group{
		frame: a,
		data:  data,
	}
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
	if len(a.lists) == 0 {
		return
	}
	if n, m := len(a.lists[0]), len(list); n != m {
		panic(fmt.Errorf("dt: invalid list length, expected %v, got %v", n, m))
	}
}
