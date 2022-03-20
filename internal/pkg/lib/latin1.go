package lib

import (
	"bytes"
	"fmt"
	"unicode/utf8"
)

func TryLatin1ToUTF8(input string) (string, error) {
	var buffer bytes.Buffer
	for _, b := range []byte(input) {
		// 0x00-0xff map to 0x0000-0xffff
		buffer.WriteRune(rune(b))
	}
	output := buffer.String()
	return output, nil
}

func TryUTF8ToLatin1(input string) (string, error) {
	var buffer bytes.Buffer

	bytes := []byte(input)
	for len(bytes) > 0 {
		r, size := utf8.DecodeRune(bytes)

		if r < 0x0080 {
			buffer.WriteByte(byte(r))
		} else if r >= 0x80 && r <= 0x00ff {
			buffer.WriteByte(byte(r))
		} else {
			return "", fmt.Errorf("character 0x%08x (%v) is not encodable as Latin-1", int(r), r)
		}

		bytes = bytes[size:]
	}
	output := buffer.String()
	return output, nil
}
