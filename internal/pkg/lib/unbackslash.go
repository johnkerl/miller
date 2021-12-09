// ================================================================
// See cst.BuildStringLiteralNode for more context.
// ================================================================

package lib

import (
	"bytes"
)

var unbackslashReplacements = map[byte]string{
	'a': "\a",
	'b': "\b",
	'f': "\f",
	'n': "\n",
	'r': "\r",
	't': "\t",
	'v': "\v",
	// At the Miller-user level this means "\\" becomes a single backslash
	// character.  It looks less clear here since here we are accommodating Go
	// conventions for backslashing conventions as well.
	'\\': "\\",
	// Similarly, "\'" becomes "'"
	'\'': "'",
	'"':  "\"",
	'?':  "?",
}

// UnbackslashStringLiteral replaces "\t" with TAB, etc. for DSL expresions
// like '$foo = "a\tb"'.  See also
// https://en.wikipedia.org/wiki/Escape_sequences_in_C
// (predates the port of Miller from C to Go).
//
// Note that a CST-build pre-pass intentionally excludes regex literals (2nd
// argument to sub/gsub/regextract/etc) from being modified here.
//
// Note "\0" .. "\9" are used for regex captures within the DSL CST builder
// and are not touched here. (See also lib/regex.go.)
func UnbackslashStringLiteral(input string) string {
	var buffer bytes.Buffer

	n := len(input)

	for i := 0; i < n; /* increment in loop */ {
		if input[i] != '\\' {
			buffer.WriteByte(input[i])
			i++
			continue
		}

		if i == n-1 {
			buffer.WriteByte(input[i])
			i++
			continue
		}

		next := input[i+1]
		replacement, ok := unbackslashReplacements[next]
		if ok {
			buffer.WriteString(replacement)
			i += 2
		} else if ok, code := isBackslashOctal(input[i:]); ok {
			buffer.WriteByte(byte(code))
			i += 4
		} else if ok, code := isBackslashHex(input[i:]); ok {
			buffer.WriteByte(byte(code))
			i += 4
		} else {
			buffer.WriteByte('\\')
			buffer.WriteByte(next)
			i += 2
		}
	}

	return buffer.String()
}

// UnhexStringLiteral is like UnbackslashStringLiteral but only unhexes things
// like "\x1f". This is for IFS and IPS setup; see the cli package.
func UnhexStringLiteral(input string) string {
	var buffer bytes.Buffer

	n := len(input)

	for i := 0; i < n; /* increment in loop */ {
		if input[i] != '\\' {
			buffer.WriteByte(input[i])
			i++
			continue
		}

		if i == n-1 {
			buffer.WriteByte(input[i])
			i++
			continue
		}

		next := input[i+1]
		if ok, code := isBackslashHex(input[i:]); ok {
			buffer.WriteByte(byte(code))
			i += 4
		} else {
			buffer.WriteByte('\\')
			buffer.WriteByte(next)
			i += 2
		}
	}

	return buffer.String()
}

// If the string starts with backslash followed by three octal digits, convert
// the next 3 characters from octal. E.g. "\123" becomes 83 (in decimal).
func isBackslashOctal(input string) (bool, int) {
	if len(input) < 4 {
		return false, 0
	}

	if input[0] != '\\' {
		return false, 0
	}

	ok, digit := isOctalDigit(input[1])
	if !ok {
		return false, 0
	}
	code := int(digit)

	ok, digit = isOctalDigit(input[2])
	if !ok {
		return false, 0
	}
	code = 8*code + int(digit)

	ok, digit = isOctalDigit(input[3])
	if !ok {
		return false, 0
	}
	code = 8*code + int(digit)

	return true, code
}

func isOctalDigit(b byte) (bool, byte) {
	if '0' <= b && b <= '7' {
		return true, b - '0'
	}
	return false, 0
}

// If the string starts with leading \x, convert the next 2 characters from hex.
// E.g.  "\xff" becomes 255 (in decimal).
func isBackslashHex(input string) (bool, int) {
	if len(input) < 4 {
		return false, 0
	}

	if input[0] != '\\' {
		return false, 0
	}

	if input[1] != 'x' && input[1] != 'X' {
		return false, 0
	}

	ok, nybble := isHexDigit(input[2])
	if !ok {
		return false, 0
	}
	code := 16 * int(nybble)

	ok, nybble = isHexDigit(input[3])
	if !ok {
		return false, 0
	}
	code += int(nybble)

	return true, code
}

func isHexDigit(b byte) (bool, byte) {
	if '0' <= b && b <= '9' {
		return true, b - '0'
	}
	if 'a' <= b && b <= 'f' {
		return true, b - 'a' + 10
	}
	if 'A' <= b && b <= 'F' {
		return true, b - 'A' + 10
	}
	return false, 0
}
