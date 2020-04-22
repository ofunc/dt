package dt

// Join is the join option.
type Join struct {
	lframe *Frame
	rframe *Frame
	lkeys  []string
	rkeys  []string
}

// Keys sets the left keys and right keys.
func (a *Join) Keys(keys ...string) *Join {
	a.lkeys = keys
	a.rkeys = keys
	return a
}

// LKeys sets the left keys.
func (a *Join) LKeys(keys ...string) *Join {
	a.lkeys = keys
	return a
}

// RKeys sets the right keys.
func (a *Join) RKeys(keys ...string) *Join {
	a.rkeys = keys
	return a
}

// Do does the join operation.
func (a *Join) Do(lprefix, rprefix string) *Frame {
	if len(a.lkeys) != len(a.rkeys) {
		panic("dt.Join: number of the left keys not equals to the right keys")
	}
	if len(a.lkeys) == 0 {
		panic("dt.Join: keys can not be empty")
	}

	m := len(a.lframe.lists)
	keys := make([]string, m+len(a.rframe.lists))
	for key, i := range a.lframe.index {
		keys[i] = lprefix + key
	}
	for key, i := range a.rframe.index {
		keys[i+m] = rprefix + key
	}
	frame := NewFrame(keys...)
	copy(frame.lists, a.lframe.lists)

	n := a.lframe.Len()
	for i := range a.rframe.lists {
		frame.lists[i+m] = make(List, n)
	}

	idx := a.index()
	for iter := a.lframe.Iter(); iter.Next(); {
		r := iter.Record().(record)
		if i, ok := idx[makeKey(r, a.lkeys)]; ok {
			for j, l := range a.rframe.lists {
				frame.lists[m+j][r.index] = l[i]
			}
		}
	}
	return frame
}

func (a *Join) index() map[string]int {
	frame := a.rframe
	idx := make(map[string]int, frame.Len())
	for iter := frame.Iter(); iter.Next(); {
		r := iter.Record().(record)
		idx[makeKey(r, a.rkeys)] = r.index
	}
	return idx
}
