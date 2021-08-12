// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.
// ================================================================

package lib

import (
	"testing"
)

type tDataForUnbackslash struct {
	input          string
	expectedOutput string
}

// Note we are here dealing with Go's backslashing conventions.
// At the Miller user-space level this is simply "\t" -> TAB, etc.
var dataForUnbackslash = []tDataForUnbackslash{
	{"", ""},
	{"abcde", "abcde"},
	{"\\1", "\\1"},
	{"a\\tb\\tc", "a\tb\tc"},
	{"a\\fb\\rc", "a\fb\rc"},
	{"a\"b\"c", "a\"b\"c"},
	{"a\\\"b\\\"c", "a\"b\"c"},
	{"a\\'b\\'c", "a'b'c"},
	{"a\102c", "aBc"},
	{"a\x42c", "aBc"},
	{"[\101\102\103]", "[ABC]"},
	{"[\x44\x45\x46]", "[DEF]"},
}

func TestUnbackslash(t *testing.T) {
	for i, entry := range dataForUnbackslash {
		actualOutput := UnbackslashStringLiteral(entry.input)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.expectedOutput, actualOutput,
			)
		}
	}
}
