package lib

import (
	"bytes"
	"fmt"
	"strings"
)

// For online help contexts like printing all the built-in DSL functions, or
// the list of all verbs.
func PrintWordsAsParagraph(words []string) {
	separator := " "
	maxlen := 80

	separatorlen := len(separator)
	linelen := 0
	j := 0

	for _, word := range words {
		wordlen := len(word)
		linelen += separatorlen + wordlen
		if linelen >= maxlen {
			fmt.Printf("\n")
			linelen = separatorlen + wordlen
			j = 0
		}
		if j > 0 {
			fmt.Print(separator)
		}
		fmt.Print(word)
		j++
	}

	fmt.Printf("\n")
}

// For online help contexts like printing all the built-in DSL functions, or
// the list of all verbs. Max width is nominally 80.
func FormatAsParagraph(text string, maxWidth int) []string {
	lines := make([]string, 0)
	words := strings.Fields(text)

	separator := " "
	separatorlen := len(separator)
	linelen := 0
	j := 0

	var buffer bytes.Buffer
	for _, word := range words {
		wordlen := len(word)
		linelen += separatorlen + wordlen
		if linelen >= maxWidth {
			line := buffer.String()
			lines = append(lines, line)
			buffer.Reset()
			linelen = separatorlen + wordlen
			j = 0
		}
		if j > 0 {
			buffer.WriteString(separator)
		}
		buffer.WriteString(word)
		j++
	}
	line := buffer.String()
	if line != "" {
		lines = append(lines, line)
	}

	return lines
}
