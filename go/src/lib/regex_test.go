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
	input            string
	sregex           string
	expectedOutput   bool
	expectedCaptures []string
}

// ----------------------------------------------------------------
var dataForSubWithoutCaptures = []tDataForSubGsub{
	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "Xbcde"},
	{"abcde", "[A-Z]", "X", "abcde"},
	{"abcde", "ab(.)de", "X\\1Y", "X\\1Y"},
}

var dataForSubWithCaptures = []tDataForSubGsub{
	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "Xbcde"},
	{"abcde", "[A-Z]", "X", "abcde"},

	//{"ab_cde", "(..)_(...)e", "\\2\\1", "cdeab"},
	//{"ab_cde", "(..)_(...)e", "\\2-\\1", "cde-ab"},
	//{"ab_cde", "(..)_(...)e", "X\\2Y\\1Z", "XcdeYabZ"},
}

var dataForGsubWithoutCaptures = []tDataForSubGsub{
	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "XXXXX"},
	{"abcde", "[A-Z]", "X", "abcde"},
	{"abcde", "[c-d]", "X", "abXXe"},
}

var dataForGsubWithCaptures = []tDataForSubGsub{
	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "XXXXX"},
	{"abcde", "[A-Z]", "X", "abcde"},
	{"abcde", "[c-d]", "X", "abXXe"},

	//{"abacad", "a(.)", "<\\2>", "<b><c><d>"},
}

// xxx needs expected-capture data
var dataForMatches = []tDataForMatches{
	{"abcde", "[a-z]", true, nil},
	{"abcde", "[A-Z]", false, nil},
}

// ----------------------------------------------------------------
func TestRegexSubWithoutCaptures(t *testing.T) {
	for _, entry := range dataForSubWithoutCaptures {
		actualOutput := RegexSubWithoutCaptures(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexSubWithCaptures(t *testing.T) {
	for _, entry := range dataForSubWithCaptures {
		actualOutput := RegexSubWithCaptures(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexGsubWithoutCaptures(t *testing.T) {
	for _, entry := range dataForGsubWithoutCaptures {
		actualOutput := RegexGsubWithoutCaptures(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexGsubWithCaptures(t *testing.T) {
	for _, entry := range dataForGsubWithCaptures {
		actualOutput := RegexGsubWithCaptures(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexMatches(t *testing.T) {
	for _, entry := range dataForMatches {
		actualOutput, actualCaptures := RegexMatches(entry.input, entry.sregex)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("input \"%s\" sregex \"%s\" expected %v got %v\n",
				entry.input, entry.sregex, entry.expectedOutput, actualOutput,
			)
		}
		// xxx compare actual/expected captures
		// xxx make a comparator function
		if actualCaptures == nil {
		}
	}
}
