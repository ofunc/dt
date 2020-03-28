package xlsx

import (
	"strconv"
	"strings"
)

// SSTable is the shared strings table.
type SSTable struct {
	Items []SItem `xml:"si"`
}

// SItem is the string item.
type SItem struct {
	Text  string   `xml:"t"`
	Texts []string `xml:"r>t"`
}

func (a SSTable) value(v string) string {
	if i, err := strconv.Atoi(v); err == nil {
		if i < len(a.Items) {
			item := a.Items[i]
			if item.Text != "" {
				return item.Text
			}
			return strings.Join(item.Texts, "")
		}
	}
	return v
}
