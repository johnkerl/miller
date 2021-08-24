package lib

import (
	"fmt"
	"os"
)

// For online help contexts like printing all the built-in DSL functions, or
// the list of all verbs.
func PrintWordsAsParagraph(words []string, o *os.File) {
	separator := " "
	maxlen := 80

	separatorlen := len(separator)
	linelen := 0
	j := 0

	for _, word := range words {
		wordlen := len(word)
		linelen += separatorlen + wordlen
		if linelen >= maxlen {
			fmt.Fprintf(o, "\n")
			linelen = separatorlen + wordlen
			j = 0
		}
		if j > 0 {
			fmt.Fprint(o, separator)
		}
		fmt.Fprint(o, word)
		j++
	}

	fmt.Printf("\n")
}
