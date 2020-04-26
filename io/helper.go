package io

import "strings"

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
