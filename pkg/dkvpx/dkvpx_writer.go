// Package dkvpx (see also dkvpx_reader.go) handles DKVPX records: comma-delimited
// key=value pairs with selective quoting. Keys and values are double-quoted only
// when they contain comma, equals, newline, or double-quote.
package dkvpx

import (
	"strings"
)

// needsQuoting reports whether the string must be quoted (contains comma,
// equals, newline, or double-quote).
func needsQuoting(s string) bool {
	return strings.ContainsAny(s, ",\n\r=\"")
}

// FormatField returns the string formatted for DKVPX output: quoted and escaped
// if it contains comma, equals, newline, or quote; otherwise unchanged.
// Useful for callers that need to apply colorization or other wrapping.
func FormatField(s string) string {
	if !needsQuoting(s) {
		return s
	}
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"':
			b.WriteString(`""`)
		default:
			b.WriteByte(s[i])
		}
	}
	b.WriteByte('"')
	return b.String()
}
