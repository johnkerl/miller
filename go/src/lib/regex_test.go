// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.
// ================================================================

package lib

import (
	"testing"
)

// ----------------------------------------------------------------
type tDataForSubGsub struct {
	input          string
	sregex         string
	replacement    string
	expectedOutput string
}

type tDataForMatches struct {
	input          string
	sregex         string
	expectedOutput bool
}

// ----------------------------------------------------------------
var dataForSub = []tDataForSubGsub{
	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "Xbcde"},
	{"abcde", "[A-Z]", "X", "abcde"},
}

var dataForGsub = []tDataForSubGsub{
	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "XXXXX"},
	{"abcde", "[A-Z]", "X", "abcde"},
	{"abcde", "[c-d]", "X", "abXXe"},
}

var dataForMatches = []tDataForMatches{
	{"abcde", "[a-z]", true},
	{"abcde", "[A-Z]", false},
}

// ----------------------------------------------------------------
func TestRegexSub(t *testing.T) {
	for _, entry := range dataForSub {
		actualOutput := RegexSub(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexGsub(t *testing.T) {
	for _, entry := range dataForGsub {
		actualOutput := RegexGsub(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexMatches(t *testing.T) {
	for _, entry := range dataForMatches {
		actualOutput := RegexMatches(entry.input, entry.sregex)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("input \"%s\" sregex \"%s\" expected %v got %v\n",
				entry.input, entry.sregex, entry.expectedOutput, actualOutput,
			)
		}
	}
}
