// Package dkvpx (see also dkvpx_reader.go) handles DKVPX records: comma-delimited
// key=value pairs with selective quoting. Keys and values are double-quoted only
// when they contain comma, equals, newline, or double-quote.
package dkvpx

import (
	"strings"
)

// needsQuoting reports whether the string must be quoted (contains the
// pair separator, the key-value separator, newline, or double-quote).
func needsQuoting(s, ofs, ops string) bool {
	return strings.ContainsAny(s, "\n\r\"") || strings.Contains(s, ofs) || strings.Contains(s, ops)
}

// FormatField returns the string formatted for DKVPX output with the default
// separators (comma and equals): quoted and escaped if it contains a
// separator, newline, or quote; otherwise unchanged.
func FormatField(s string) string {
	return FormatFieldWithSeparators(s, ",", "=")
}

// FormatFieldWithSeparators returns the string formatted for DKVPX output:
// quoted and escaped if it contains the given pair separator (OFS), the given
// key-value separator (OPS), newline, or quote; otherwise unchanged.
// Useful for callers that need to apply colorization or other wrapping.
func FormatFieldWithSeparators(s, ofs, ops string) string {
	if !needsQuoting(s, ofs, ops) {
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
