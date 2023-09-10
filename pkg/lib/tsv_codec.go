package lib

import (
	"bytes"
)

// * https://en.wikipedia.org/wiki/Tab-separated_values
// * https://www.iana.org/assignments/media-types/text/tab-separated-values
//   \n for newline,
//   \r for carriage return,
//   \t for tab,
//   \\ for backslash.

// TSVDecodeField is for the TSV record-reader.
func TSVDecodeField(input string) string {
	var buffer bytes.Buffer
	n := len(input)
	for i := 0; i < n; /* increment in loop */ {
		c := input[i]
		if c == '\\' && i < n-1 {
			d := input[i+1]
			if d == '\\' {
				buffer.WriteByte('\\')
				i += 2
			} else if d == 'n' {
				buffer.WriteByte('\n')
				i += 2
			} else if d == 'r' {
				buffer.WriteByte('\r')
				i += 2
			} else if d == 't' {
				buffer.WriteByte('\t')
				i += 2
			} else {
				buffer.WriteByte(c)
				i++
			}
		} else {
			buffer.WriteByte(c)
			i++
		}
	}
	return buffer.String()
}

// TSVEncodeField is for the TSV record-writer.
func TSVEncodeField(input string) string {
	var buffer bytes.Buffer
	for _, r := range input {
		if r == '\\' {
			buffer.WriteByte('\\')
			buffer.WriteByte('\\')
		} else if r == '\n' {
			buffer.WriteByte('\\')
			buffer.WriteByte('n')
		} else if r == '\r' {
			buffer.WriteByte('\\')
			buffer.WriteByte('r')
		} else if r == '\t' {
			buffer.WriteByte('\\')
			buffer.WriteByte('t')
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}
