package dt

// Join is the join option.
type Join struct {
	lframe *Frame
	rframe *Frame
	lkeys  []string
	rkeys  []string
}

// On sets the left keys.
func (a *Join) On(key string, keys ...string) *Join {
	a.lkeys = append(keys, key)
	return a
}

// Do does the join operation.
func (a *Join) Do(prefix string) *Frame {
	if len(a.lkeys) == 0 {
		a.lkeys = a.rkeys
	}
	if len(a.lkeys) != len(a.rkeys) {
		panic("dt.Join: number of the left keys not equals to the right keys")
	}

	m := len(a.lframe.lists)
	rframe := a.rframe.Copy(false).Del(a.rkeys...)
	keys := make([]string, m+len(rframe.lists))
	for key, j := range rframe.index {
		keys[j+m] = prefix + key
	}
	frame := NewFrame(keys...)
	copy(frame.lists, a.lframe.lists)

	n := a.lframe.Len()
	for j := range rframe.lists {
		frame.lists[j+m] = make(List, n)
	}

	idx := a.index()
	for iter := a.lframe.Iter(); iter.Next(); {
		r := iter.Record().(record)
		if i, ok := idx[makeKey(r, a.lkeys)]; ok {
			for j, l := range rframe.lists {
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
