package csv

import (
	"github.com/ofunc/dt"
	util "github.com/ofunc/dt/io"
)

// Value gets the value by index i and j.
func Value(rs [][]string, i, j int) dt.Value {
	if i < len(rs) && i >= 0 {
		if r := rs[i]; j < len(r) && j >= 0 {
			return util.Value(r[j])
		}
	}
	return nil
}
