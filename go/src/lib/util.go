package lib

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

func BooleanXOR(a, b bool) bool {
	return a != b
}

func BoolToInt(b bool) int {
	if b == false {
		return 0
	} else {
		return 1
	}
}

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
	if stringList == nil {
		return nil
	}

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

func ReverseStringList(strings []string) {
	n := len(strings)
	i := 0
	j := n - 1
	for i < j {
		temp := strings[i]
		strings[i] = strings[j]
		strings[j] = temp
		i++
		j--
	}
}

func SortedStrings(strings []string) []string {
	copy := make([]string, len(strings))
	for i, s := range strings {
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
func TryIntFromString(input string) (int, bool) {
	ival, err := strconv.ParseInt(input, 0 /* check all*/, 64)
	if err == nil {
		return int(ival), true
	} else {
		return 0, false
	}
}

func TryFloat64FromString(input string) (float64, bool) {
	fval, err := strconv.ParseFloat(input, 64)
	if err == nil {
		return fval, true
	} else {
		return 0, false
	}
}

func TryBoolFromBoolString(input string) (bool, bool) {
	if input == "true" {
		return true, true
	} else if input == "false" {
		return false, true
	} else {
		return false, false
	}
}

// Go doesn't preserve insertion order in its arrays, so here we make an
// accessor for getting the keys in sorted order for the benefit of
// map-printers.
func GetArrayKeysSorted(input map[string]string) []string {
	keys := make([]string, len(input))
	i := 0
	for key := range input {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// WriteTempFile places the contents string into a temp file, which the caller
// must remove.
func WriteTempFileOrDie(contents string) string {
	// Use "" as first argument to ioutil.TempFile to use default directory.
	// Nominally "/tmp" or somesuch on all unix-like systems, but not for Windows.
	handle, err := ioutil.TempFile("", "mlr-temp")
	if err != nil {
		fmt.Printf("mlr: could not create temp file.\n")
		os.Exit(1)
	}

	_, err = handle.WriteString(contents)
	if err != nil {
		fmt.Printf("mlr: could not populate temp file.\n")
		os.Exit(1)
	}

	err = handle.Close()
	if err != nil {
		fmt.Printf("mlr: could not finish write of  temp file.\n")
		os.Exit(1)
	}
	return handle.Name()
}
