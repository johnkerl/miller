// Most Miller tests (thousands of them) are command-line-driven via
// mlr regtest. Here are some cases needing special focus.

package lib

import (
	"fmt"
	"slices"
	"testing"
)

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
		actualHasCaptures, actualMatrix := ReplacementHasCaptures(entry.replacement)
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
		actualOutput := RegexStringSub(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexGsub(t *testing.T) {
	for i, entry := range dataForGsub {
		actualOutput := RegexStringGsub(entry.input, entry.sregex, entry.replacement)
		if actualOutput != entry.expectedOutput {
			t.Fatalf("case %d input \"%s\" sregex \"%s\" replacement \"%s\" expected \"%s\" got \"%s\"\n",
				i, entry.input, entry.sregex, entry.replacement, entry.expectedOutput, actualOutput,
			)
		}
	}
}

func TestRegexMatches(t *testing.T) {
	for i, entry := range dataForMatches {
		actualOutput, actualCaptures := RegexStringMatchWithCaptures(entry.input, entry.sregex)
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

func TestRegexStringMatchWithMapResultsUTF8(t *testing.T) {
	// Starts and ends are 1-up character positions, as with substr and index,
	// not byte offsets.
	matched, captures, starts, ends := RegexStringMatchWithMapResults("héllo:123", "([a-z]+):([0-9]+)")
	if !matched {
		t.Fatalf("expected a match")
	}
	if !compareCaptures(captures, []string{"llo:123", "llo", "123"}) {
		t.Fatalf("captures: got %#v", captures)
	}
	if !slices.Equal(starts, []int{3, 3, 7}) {
		t.Fatalf("starts: expected %v got %v", []int{3, 3, 7}, starts)
	}
	if !slices.Equal(ends, []int{9, 5, 9}) {
		t.Fatalf("ends: expected %v got %v", []int{9, 5, 9}, ends)
	}
}

// withEmptyRegexpCache empties the regex cache for the duration of a test,
// restoring the package-level cache afterward so other tests are unaffected.
func withEmptyRegexpCache(t *testing.T) {
	cacheMutex.Lock()
	saved := regexpCache
	regexpCache = nil
	cacheMutex.Unlock()
	t.Cleanup(func() {
		cacheMutex.Lock()
		regexpCache = saved
		cacheMutex.Unlock()
	})
}

func regexpCacheLookup(s string) bool {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	_, ok := regexpCache[s]
	return ok
}

func regexpCacheSize() int {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()
	return len(regexpCache)
}

func TestRegexpCompileCachedReusesCompiledRegex(t *testing.T) {
	withEmptyRegexpCache(t)

	for _, sregex := range []string{"^([0-9]{4})-([0-9]{2})", "\"(..)_(...)\"", "/a.*b/"} {
		first, err := CompileMillerRegex(sregex)
		if err != nil {
			t.Fatalf("sregex %s: unexpected error %v", sregex, err)
		}
		second, err := CompileMillerRegex(sregex)
		if err != nil {
			t.Fatalf("sregex %s: unexpected error %v", sregex, err)
		}
		if first != second {
			t.Fatalf("sregex %s: expected the cached *regexp.Regexp on the second call", sregex)
		}
	}
}

func TestRegexpCompileCachedCaseInsensitive(t *testing.T) {
	withEmptyRegexpCache(t)

	// "..."i and /.../i are compiled as "(?i)...", so they share a cache entry
	// with each other and with an explicit "(?i)..." but not with the
	// case-sensitive form.
	quoted, err := CompileMillerRegex("\"ab+c\"i")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	slashed, err := CompileMillerRegex("/ab+c/i")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	explicit, err := CompileMillerRegex("(?i)ab+c")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	sensitive, err := CompileMillerRegex("\"ab+c\"")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	if quoted != slashed || quoted != explicit {
		t.Fatalf("expected one cached regex for \"ab+c\"i, /ab+c/i, and (?i)ab+c")
	}
	if quoted == sensitive {
		t.Fatalf("case-insensitive and case-sensitive regexes must not share a cache entry")
	}
	if quoted.String() != "(?i)ab+c" || sensitive.String() != "ab+c" {
		t.Fatalf("got %q and %q", quoted.String(), sensitive.String())
	}
	if !quoted.MatchString("xABBCx") || sensitive.MatchString("xABBCx") {
		t.Fatalf("case-insensitive regex should match \"xABBCx\" and case-sensitive should not")
	}

	// Same outputs on a cold and a warm cache.
	for pass := 0; pass < 2; pass++ {
		if output := RegexStringSub("ABBC-abc", "\"ab+c\"i", "X"); output != "X-abc" {
			t.Fatalf("pass %d: sub expected \"X-abc\" got \"%s\"", pass, output)
		}
		if output := RegexStringGsub("ABBC-abc", "\"ab+c\"i", "X"); output != "X-X" {
			t.Fatalf("pass %d: gsub expected \"X-X\" got \"%s\"", pass, output)
		}
		if output := RegexStringGsub("ABBC-abc", "\"ab+c\"", "X"); output != "ABBC-X" {
			t.Fatalf("pass %d: gsub expected \"ABBC-X\" got \"%s\"", pass, output)
		}
		matched, captures := RegexStringMatchWithCaptures("xxAB_CDE", "\"(ab)_(c..)\"i")
		if !matched || !compareCaptures(captures, []string{"AB_CDE", "AB", "CDE", "", "", "", "", "", "", ""}) {
			t.Fatalf("pass %d: =~ got %v %#v", pass, matched, captures)
		}
	}
}

func TestRegexpCompileCachedDoesNotCacheErrors(t *testing.T) {
	withEmptyRegexpCache(t)

	for i := 0; i < 2; i++ {
		regex, err := CompileMillerRegex("\"a(b\"")
		if err == nil {
			t.Fatalf("call %d: expected an error for an unbalanced parenthesis", i)
		}
		if regex != nil {
			t.Fatalf("call %d: expected a nil regex on error", i)
		}
	}
	if regexpCacheLookup("a(b") {
		t.Fatalf("a regex which failed to compile should not be cached")
	}
	if n := regexpCacheSize(); n != 0 {
		t.Fatalf("expected an empty cache, got %d entries", n)
	}
}

func TestRegexpCompileCachedSizeLimit(t *testing.T) {
	withEmptyRegexpCache(t)

	for i := range cacheMaxSize + 10 {
		sregex := fmt.Sprintf("^k%d$", i)
		regex, err := CompileMillerRegex(sregex)
		if err != nil {
			t.Fatalf("sregex %s: unexpected error %v", sregex, err)
		}
		if !regex.MatchString(fmt.Sprintf("k%d", i)) {
			t.Fatalf("sregex %s: expected a match", sregex)
		}
	}
	if n := regexpCacheSize(); n != cacheMaxSize {
		t.Fatalf("expected %d cached regexes, got %d", cacheMaxSize, n)
	}

	// Regexes cached before the limit was reached are still reused.
	first, _ := CompileMillerRegex("^k0$")
	second, _ := CompileMillerRegex("^k0$")
	if first != second {
		t.Fatalf("expected the cached regex for ^k0$")
	}

	// Regexes past the limit are compiled each time and not cached.
	sregex := fmt.Sprintf("^k%d$", cacheMaxSize+5)
	first, _ = CompileMillerRegex(sregex)
	second, _ = CompileMillerRegex(sregex)
	if first == second || regexpCacheLookup(sregex) {
		t.Fatalf("sregex %s: expected no caching past the size limit", sregex)
	}
	if !first.MatchString(fmt.Sprintf("k%d", cacheMaxSize+5)) {
		t.Fatalf("sregex %s: expected a match", sregex)
	}
	if n := regexpCacheSize(); n != cacheMaxSize {
		t.Fatalf("expected %d cached regexes, got %d", cacheMaxSize, n)
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
