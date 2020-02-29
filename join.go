package dt

import (
	"errors"
	"reflect"
)

// Join is the join option.
type Join struct {
	lframe *Frame
	rframe *Frame
	lkeys  []string
	rkeys  []string
}

// Keys sets the left keys and right keys.
func (a Join) Keys(keys ...string) Join {
	a.lkeys = keys
	a.rkeys = keys
	return a
}

// LKeys sets the left keys.
func (a Join) LKeys(keys ...string) Join {
	a.lkeys = keys
	return a
}

// RKeys sets the right keys.
func (a Join) RKeys(keys ...string) Join {
	a.rkeys = keys
	return a
}

// Do does the join operation.
func (a Join) Do(lprefix, rprefix string) *Frame {
	if len(a.lkeys) != len(a.rkeys) {
		panic(errors.New("dt.Join: number of the left keys not equals to the right keys"))
	}
	if len(a.lkeys) == 0 {
		panic(errors.New("dt.Join: keys can not be empty"))
	}

	n, m := len(a.lframe.lists), a.lframe.Len()
	frame := NewFrame(n + len(a.rframe.lists))
	copy(frame.lists, a.lframe.lists)
	for key, i := range a.lframe.index {
		frame.index[lprefix+key] = i
	}
	for key, i := range a.rframe.index {
		frame.index[rprefix+key] = n + i
		frame.lists[n+i] = make(List, m)
	}

	idx := a.index(a.rframe)
	typ := reflect.ArrayOf(len(a.lkeys), tvalue)
	for iter := a.lframe.Iter(); iter.Next(); {
		r := iter.Record().(record)
		if i, ok := idx[makeKey(typ, r, a.lkeys)]; ok {
			for j, l := range a.rframe.lists {
				frame.lists[n+j][r.index] = l[i]
			}
		}
	}
	return frame
}

func (a Join) index(frame *Frame) map[interface{}]int {
	idx := make(map[interface{}]int, frame.Len())
	typ := reflect.ArrayOf(len(a.rkeys), tvalue)
	for iter := frame.Iter(); iter.Next(); {
		r := iter.Record().(record)
		idx[makeKey(typ, r, a.rkeys)] = r.index
	}
	return idx
}
