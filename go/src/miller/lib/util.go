package lib

import (
	"strings"
)

func Plural(n int) string {
	if n == 1 {
		return ""
	} else {
		return "s"
	}
}

// In Go as in all languages I'm aware of with a string-split, "a,b,c" splits
// on "," to ["a", "b", "c" and "a" splits to ["a"], both of which are fine --
// but "" splits to [""] when I wish it were []. This function does the latter.
func SplitString(input string, separator string) []string {
	if input == "" {
		return make([]string, 0)
	} else {
		return strings.Split(input, separator)
	}
}

func StringListToSet(stringList []string) map[string]bool {
	stringSet := make(map[string]bool)
	for _, s := range stringList {
		stringSet[s] = true
	}
	return stringSet
}

func IntMin2(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
