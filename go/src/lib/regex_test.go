// ================================================================
// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.
// ================================================================

package lib

import (
	"testing"
)

// ----------------------------------------------------------------
type tDataForHasCaptures struct {
	replacement         string
	expectedHasCaptures bool
	expectedMatrix      [][]int
}

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
var dataForHasCaptures = []tDataForHasCaptures{
	{"foo", false, nil},
	{"\\0", true, [][]int{{0, 2, 0, 2}}},
	{"\\3", true, [][]int{{0, 2, 0, 2}}},
	{"\\34", true, [][]int{{0, 2, 0, 2}}},
	{"abc\\1def\\2ghi", true, [][]int{{3, 5, 3, 5}, {8, 10, 8, 10}}},
}

var dataForSub = []tDataForSubGsub{
	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "Xbcde"},
	{"abcde", "[A-Z]", "X", "abcde"},

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

var dataForGsub = []tDataForSubGsub{
	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "XXXXX"},
	{"abcde", "[A-Z]", "X", "abcde"},
	{"abcde", "[c-d]", "X", "abXXe"},

	{"abcde", "c", "X", "abXde"},
	{"abcde", "z", "X", "abcde"},
	{"abcde", "[a-z]", "X", "XXXXX"},
	{"abcde", "[A-Z]", "X", "abcde"},
	{"abcde", "[c-d]", "X", "abXXe"},

	{"abacad", "a(.)", "<\\1>", "<b><c><d>"},
	{"abacad", "a(.)", "<\\2>", "<><><>"},
}

var dataForMatches = []tDataForMatches{
	{"abcde", "[A-Z]", false, []string{"", "", "", "", "", "", "", "", "", ""}},
	{"abcde", "[a-z]", true, []string{"a", "", "", "", "", "", "", "", "", ""}},
	{"...ab_cde...", "(..)_(...)", true, []string{"ab_cde", "ab", "cde", "", "", "", "", "", "", ""}},
	{"...ab_cde...fg_hij...", "(..)_(...)", true, []string{"ab_cde", "ab", "cde", "", "", "", "", "", "", ""}},
	{"foofoofoo", "(f.o)", true, []string{"foo", "foo", "", "", "", "", "", "", "", ""}},
	{"foofoofoo", "(f.*o)", true, []string{"foofoofoo", "foofoofoo", "", "", "", "", "", "", "", ""}},
}

func TestRegexReplacementHasCaptures(t *testing.T) {
	for i, entry := range dataForHasCaptures {
		actualHasCaptures, actualMatrix := RegexReplacementHasCaptures(entry.replacement)
		if actualHasCaptures != entry.expectedHasCaptures {
			t.Fatalf("case %d replacement \"%s\" expected %v got %v\n",
				i, entry.replacement, entry.expectedHasCaptures, actualHasCaptures,
			)
		}
		if !compareMatrices(actualMatrix, entry.expectedMatrix) {
			t.Fatalf("case %d replacement \"%s\" expected matrix %#v got %#v\n",
				i, entry.replacement, entry.expectedMatrix, actualMatrix,
			)
		}
	}
}

func TestRegexSub(t *testing.T) {
	for i, entry := range dataForSub {
		actualOutput := RegexSub(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexGsub(t *testing.T) {
	for i, entry := range dataForGsub {
		actualOutput := RegexGsub(entry.input, entry.sregex, entry.replacement)
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

func compareMatrices(
	actualMatrix [][]int,
	expectedMatrix [][]int,
) bool {
	if actualMatrix == nil && expectedMatrix == nil {
		return true
	}
	if actualMatrix == nil || expectedMatrix == nil {
		return false
	}
	if len(actualMatrix) != len(expectedMatrix) {
		return false
	}
	for i := range expectedMatrix {
		actualRow := actualMatrix[i]
		expectedRow := expectedMatrix[i]
		if len(actualRow) != len(expectedRow) {
			return false
		}
		for j := range expectedRow {
			if actualRow[j] != expectedRow[j] {
				return false
			}
		}
	}
	return true
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
	for i := range expectedCaptures {
		if actualCaptures[i] != expectedCaptures[i] {
			return false
		}
	}
	return true
}
