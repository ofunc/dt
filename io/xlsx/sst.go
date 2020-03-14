package xlsx

import (
	"strconv"
)

// SST is the shared strings table.
type SST struct {
	Values []string `xml:"si>t"`
}

// Value returns the shared string value.
func (a SST) Value(v string) string {
	if i, err := strconv.Atoi(v); err == nil {
		if i < len(a.Values) {
			return a.Values[i]
		}
	}
	return v
}
