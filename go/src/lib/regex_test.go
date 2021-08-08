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

	{"ab_cde", "(..)_(...)", "\\2\\1", "cdeab"},
	{"ab_cde", "(..)_(...)", "\\2-\\1", "cde-ab"},
	{"ab_cde", "(..)_(...)", "X\\2Y\\1Z", "XcdeYabZ"},

	{"foofoofoo", "(f.o)", "b\\1r", "bfoorfoofoo"},
	{"foofoofoo", "(f.*o)", "b\\1r", "bfoofoofoor"},
	{"foofoofoo", "(f.o)", "b\\2r", "brfoofoo"},
	{"foofoofoo", "(f.*o)", "b\\2r", "br"},
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

	{"abacad", "a(.)", "<\\1>", "<b><c><d>"},
	{"abacad", "a(.)", "<\\2>", "<><><>"},
}

var dataForMatches = []tDataForMatches{
	{"abcde", "[A-Z]", false, nil},
	{"abcde", "[a-z]", true, nil},
	{"...ab_cde...", "(..)_(...)", true, []string{"", "ab", "cde", "", "", "", "", "", "", ""}},
	{"...ab_cde...fg_hij...", "(..)_(...)", true, []string{"", "ab", "cde", "", "", "", "", "", "", ""}},
	{"foofoofoo", "(f.o)", true, []string{"", "foo", "", "", "", "", "", "", "", ""}},
	{"foofoofoo", "(f.*o)", true, []string{"", "foofoofoo", "", "", "", "", "", "", "", ""}},
}

// ----------------------------------------------------------------
func TestRegexSubWithoutCaptures(t *testing.T) {
	for i, entry := range dataForSubWithoutCaptures {
		actualOutput := RegexSubWithoutCaptures(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexSubWithCaptures(t *testing.T) {
	for i, entry := range dataForSubWithCaptures {
		actualOutput := RegexSubWithCaptures(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexGsubWithoutCaptures(t *testing.T) {
	for i, entry := range dataForGsubWithoutCaptures {
		actualOutput := RegexGsubWithoutCaptures(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexGsubWithCaptures(t *testing.T) {
	for i, entry := range dataForGsubWithCaptures {
		actualOutput := RegexGsubWithCaptures(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexMatches(t *testing.T) {
	for i, entry := range dataForMatches {
		actualOutput, actualCaptures := RegexMatches(entry.input, entry.sregex)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" expected %v got %v\n",
				i, entry.input, entry.sregex, entry.expectedOutput, actualOutput,
			)
		}
		if !compareCaptures(actualCaptures, entry.expectedCaptures) {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" expected captures %#v got %#v\n",
				i, entry.input, entry.sregex, entry.expectedCaptures, actualCaptures,
			)
		}
	}
}

func compareCaptures(
	actualCaptures []string,
	expectedCaptures []string,
) bool {
	if actualCaptures == nil && expectedCaptures == nil {
		return true
	}
	if actualCaptures == nil || expectedCaptures == nil {
		return false
	}
	if len(actualCaptures) != len(expectedCaptures) {
		return false
	}
	for i, _ := range expectedCaptures {
		if actualCaptures[i] != expectedCaptures[i] {
			return false
		}
	}
	return true
}
