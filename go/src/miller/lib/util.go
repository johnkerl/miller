package lib

import (
	"sort"
	"strconv"
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

func SortStrings(strings []string) {
	// Go sort API: for ascending sort, return true if element i < element j.
	sort.Slice(strings, func(i, j int) bool {
		return strings[i] < strings[j]
	})
}

func SortedStrings(strings []string) []string {
	copy := make([]string, len(strings))
	for i, s := range(strings) {
		copy[i] = s
	}
	// Go sort API: for ascending sort, return true if element i < element j.
	sort.Slice(copy, func(i, j int) bool {
		return copy[i] < copy[j]
	})
	return copy
}

func IntMin2(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

// Tries decimal, hex, octal, and binary.
func TryInt64FromString(input string) (int64, bool) {
	ival, err := strconv.ParseInt(input, 0 /* check all*/, 64)
	if err == nil {
		return ival, true
	} else {
		return 0, false
	}
}
