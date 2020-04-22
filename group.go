package dt

// Group is a group data structure.
type Group struct {
	frame *Frame
	data  map[string]([]int)
	keys  []string
	names []string
	funcs [](func(List) Value)
}

// Apply applies the aggregate function to group a.
func (a *Group) Apply(key string, name string, f func(List) Value) *Group {
	a.keys = append(a.keys, key)
	a.names = append(a.names, name)
	a.funcs = append(a.funcs, f)
	return a
}

// Do does the group.
func (a *Group) Do() *Frame {
	frame := NewFrame(a.names...)
	for _, is := range a.data {
		for j, key := range a.keys {
			list := a.frame.Get(key)
			l := make(List, len(is))
			for k, i := range is {
				l[k] = list[i]
			}
			frame.lists[j] = append(frame.lists[j], a.funcs[j](l))
		}
	}
	return frame
}
