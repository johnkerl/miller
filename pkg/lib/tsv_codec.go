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
			switch d {
			case '\\':
				buffer.WriteByte('\\')
				i += 2
			case 'n':
				buffer.WriteByte('\n')
				i += 2
			case 'r':
				buffer.WriteByte('\r')
				i += 2
			case 't':
				buffer.WriteByte('\t')
				i += 2
			default:
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
		switch r {
		case '\\':
			buffer.WriteByte('\\')
			buffer.WriteByte('\\')
		case '\n':
			buffer.WriteByte('\\')
			buffer.WriteByte('n')
		case '\r':
			buffer.WriteByte('\\')
			buffer.WriteByte('r')
		case '\t':
			buffer.WriteByte('\\')
			buffer.WriteByte('t')
		default:
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}
