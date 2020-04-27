package io

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ofunc/dt"
)

var regDigits = regexp.MustCompile(`^\d+$`)

// Keys rename duplicate keys with suffix.
func Keys(keys []string, suffix string) []string {
	m := make(map[string]int, len(keys))
	for j, key := range keys {
		k := m[key]
		keys[j] = key + strings.Repeat(suffix, k)
		m[key] = k + 1
	}
	return keys
}

// Value returns the dt.Value by the value.
func Value(value string) dt.Value {
	x := strings.TrimSpace(value)
	if len(x) < 16 || !regDigits.MatchString(x) {
		if v, err := strconv.ParseFloat(x, 64); err == nil {
			return dt.Number(v)
		}
	}
	return dt.String(value)
}
