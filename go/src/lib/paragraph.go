package lib

import (
	"fmt"
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
