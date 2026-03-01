package input

import (
	"fmt"
	"strconv"
	"strings"
)

type fixedWidthSplitter struct {
	// Valid increasing indexes
	indexes []int
}

// NewFixedWidthSplitter creates a new fixed-width field splitter based on the spec.
//
// The spec parameter supports these formats:
//
//	widths:w1,w2,w3...     : Split by explicit widths for each column. Omit last one for variable width ending (e.g., 4,4,5)
//	left-align             : Parse header with left-aligned fields
//	left-align-multi-word  : Parse header with multi-word support
//	right-align            : Parse header with right-aligned fields
//	right-align-multi-word : Parse header with right-aligned fields and multi-word support
//
// For mult-word cases, the column headers can be made up of multiple words
// e.g. "Seq No". Adjacent columns should be seperated by at least two spaces
func NewFixedWidthSplitter(spec, referenceRow string) (*fixedWidthSplitter, error) {
	var indexes []int

	if spec[:7] == "widths:" {
		var err error
		indexes, err = parseWidths(spec[7:])
		if err != nil {
			return nil, err
		}
	} else if spec == "left-align" {
		indexes = parseLeftAlign(referenceRow, false)
	} else if spec == "left-align-multi-word" {
		indexes = parseLeftAlign(referenceRow, true)
	} else if spec == "right-align" {
		indexes = parseRightAlign(referenceRow, false)
	} else if spec == "right-align-multi-word" {
		indexes = parseRightAlign(referenceRow, true)
	} else {
		return nil, fmt.Errorf("Unknown spec: %v", spec)
	}
	return &fixedWidthSplitter{indexes: indexes}, nil

}

func (sp *fixedWidthSplitter) Split(line string) []string {
	return split(line, sp.indexes)
}

func parseWidths(widthsStr string) ([]int, error) {
	var indexes []int
	if len(widthsStr) == 0 {
		return indexes, nil
	}
	widths := strings.Split(widthsStr, ",")
	pos := 0
	for _, w := range widths {
		width, err := strconv.Atoi(strings.TrimSpace(w))
		if err != nil {
			return nil, fmt.Errorf("invalid width: %v, error: %w", w, err)
		}
		if width <= 0 {
			return nil, fmt.Errorf("not a positive width: %v", w)
		}
		pos += width
		indexes = append(indexes, pos)
	}
	return indexes, nil
}

func parseLeftAlign(referenceRow string, allowMultiWord bool) []int {
	var indexes []int
	inWord := true
	firstSpace := false //Used for multi word. True on the first space after a non-space
	for i, c := range referenceRow {
		if c != ' ' {
			if !inWord {
				indexes = append(indexes, i)
			}
			inWord = true
		} else {
			if allowMultiWord {
				if firstSpace {
					inWord = false
					firstSpace = false
				} else {
					firstSpace = true
					// inWord = true // already true
				}
			} else {
				inWord = false
			}
		}
	}
	return indexes
}

func parseRightAlign(referenceRow string, allowMultiWord bool) []int {
	var indexes []int
	inWord := false
	firstSpace := false //Used for multi word. True on the first space after a non-space
	for i, c := range referenceRow {
		if c != ' ' {
			inWord = true
			firstSpace = false
		} else {
			if inWord {
				if allowMultiWord {
					firstSpace = true
				} else {
					indexes = append(indexes, i)
				}
			} else {
				if allowMultiWord && firstSpace {
					indexes = append(indexes, i-1)
					firstSpace = false
				}
			}
			inWord = false
		}
	}
	return indexes
}

func split(line string, indexes []int) []string {
	var result []string
	if len(indexes) == 0 {
		result = append(result, line)
		return result
	}
	si := 0
	l := len(line)
	for _, idx := range indexes {
		if idx > l {
			break
		}
		result = append(result, line[si:idx])
		si = idx
	}
	rest := line[si:]
	if rest != "" {
		result = append(result, rest)
	}
	return result
}
