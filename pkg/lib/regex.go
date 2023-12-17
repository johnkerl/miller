// ================================================================
// Support for regular expressions in Miller.
//
// * By and large we use the Go library.
//
// * There is (for historical reasons) a DSL syntax "[a-z]"i (note the trailing i)
//   for case-insensitive regular expressions which we map into Go syntax for
//   regex-compilation.
//
// * Also for historical reasons, we allow things like
//     if ($x =~ "(..)_(...)") {
//       ... other lines of code ...
//       $y = "\2:\1";
//     }
//   where the '=~' sets the captures and the "\2:\1" uses them.  (Note that
//   https://github.com/johnkerl/miller/issues/388 has a better suggestion which would make the
//   captures explicit as variables, rather than implicit within CST state: this is implemented by
//   the `match` and `matchx` DSL functions.  Regardless, the `=~` syntax will still be supported
//   for backward compatibility and so is here to stay.) Here we make use of Go regexp-library
//   functions to write to, and then later interpolate from, a captures array which is stored within
//   CST state. (See the `runtime.State` object.)
//
// * "\0" is for a full match; "\1" .. "\9" are for submatch cqptures. E.g.
//   if $x is "foobarbaz" and the regex is "foo(.)(..)baz", then "\0" is
//   "foobarbaz", "\1" is "b", "\2" is "ar", and "\3".."\9" are "".
//
// * Naming:
//
//   o "regexp" and "Regexp" are used for the Go library and its data structure, respectively;
//
//   o "regex" is used for regular-expression strings following Miller's idiosyncratic syntax and
//     semantics as described above.
//
// ================================================================

package lib

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
)

// captureDetector is used to see if a string literal interpolates previous
// captures (like "\2:\1") or not (like "2:1").
var captureDetector = regexp.MustCompile(`\\[0-9]`)

// captureSplitter is used to precompute an offsets matrix for strings like
// "\2:\1" so they don't need to be recomputed on every record.
var captureSplitter = regexp.MustCompile(`(\\[0-9])`)

// See regexpCompileCached
var regexpCache map[string]*regexp.Regexp

const cacheMaxSize = 1000

var cacheMutex sync.Mutex

// regexpCompileCached keeps a cache of compiled regexes, so that the caller has the flexibility to
// only pass in strings while getting the benefits of compilation avoidance.
//
// Regarding cache size: in nominal use, regexp strings are within Miller DSL code statements, and
// there will be a handful. These will all get re-used after their first application, and the cache
// will remain bounded by the size of the user's DSL code. However, it is possible to have regex
// strings contained within Miller record-field data.
//
// We could solve this by using an LRU cache. However, for simplicity, we limit the number of
// cached compiles, and for any extras that appear during record processing, we simply recompile
// each time.
func regexpCompileCached(s string) (*regexp.Regexp, error) {
	if len(regexpCache) > cacheMaxSize {
		return regexp.Compile(s)
	}
	r, err := regexp.Compile(s)
	if err == nil {
		cacheMutex.Lock()
		if regexpCache == nil {
			regexpCache = make(map[string]*regexp.Regexp)
		}
		regexpCache[s] = r
		cacheMutex.Unlock()
	}
	return r, err
}

// CompileMillerRegex wraps Go regex-compile with some Miller-specific syntax which predates the
// port of Miller from C to Go.  Miller regexes use a final 'i' to indicate case-insensitivity; Go
// regexes use an initial "(?i)".
//
// (See also mlr.bnf where we specify which things can be backslash-escaped without a syntax error
// at parse time.)
//
// * If the regex_string is of the form a.*b, compiles it case-sensitively.
// * If the regex_string is of the form "a.*b", compiles a.*b case-sensitively.
// * If the regex_string is of the form "a.*b"i, compiles a.*b case-insensitively.
func CompileMillerRegex(regexString string) (*regexp.Regexp, error) {
	n := len(regexString)
	if n < 2 {
		return regexpCompileCached(regexString)
	}

	// TODO: rethink this. This will strip out things people have entered, e.g. "\"...\"".
	// The parser-to-AST will have stripped the outer and we'll strip the inner and the
	// user's intent will be lost.
	//
	// TODO: make separate functions for calling from parser-to-AST (string
	// literals) and from verbs (like cut -r or having-fields).

	if strings.HasPrefix(regexString, "\"") && strings.HasSuffix(regexString, "\"") {
		return regexpCompileCached(regexString[1 : n-1])
	}
	if strings.HasPrefix(regexString, "/") && strings.HasSuffix(regexString, "/") {
		return regexpCompileCached(regexString[1 : n-1])
	}

	if strings.HasPrefix(regexString, "\"") && strings.HasSuffix(regexString, "\"i") {
		return regexpCompileCached("(?i)" + regexString[1:n-2])
	}
	if strings.HasPrefix(regexString, "/") && strings.HasSuffix(regexString, "/i") {
		return regexpCompileCached("(?i)" + regexString[1:n-2])
	}

	return regexpCompileCached(regexString)
}

// CompileMillerRegexOrDie wraps CompileMillerRegex. Usually in Go we want to
// return a second error argument rather than fataling. However, if there's a
// malformed regex we really cannot continue so it's simpler to just fatal.
func CompileMillerRegexOrDie(regexString string) *regexp.Regexp {
	regex, err := CompileMillerRegex(regexString)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return regex
}

// CompileMillerRegexesOrDie is a convenenience looper over CompileMillerRegexOrDie.
func CompileMillerRegexesOrDie(regexStrings []string) []*regexp.Regexp {
	regexes := make([]*regexp.Regexp, len(regexStrings))

	for i, regexString := range regexStrings {
		regexes[i] = CompileMillerRegexOrDie(regexString)
	}

	return regexes
}

// In Go as in all languages I'm aware of with a string-split, "a,b,c" splits
// on "," to ["a", "b", "c" and "a" splits to ["a"], both of which are fine --
// but "" splits to [""] when I wish it were []. This function does the latter.
func RegexCompiledSplitString(regex *regexp.Regexp, input string, n int) []string {
	if input == "" {
		return make([]string, 0)
	} else {
		return regex.Split(input, n)
	}
}

// RegexStringSub implements the sub DSL function.
func RegexStringSub(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	_, replacementCaptureMatrix := ReplacementHasCaptures(replacement)
	return RegexCompiledSub(input, regex, replacement, replacementCaptureMatrix)
}

// RegexCompiledSub is the same as RegexStringSub but with compiled regex and
// replacement strings.
func RegexCompiledSub(
	input string,
	regex *regexp.Regexp,
	replacement string,
	replacementCaptureMatrix [][]int,
) string {
	return regexCompiledSubOrGsub(input, regex, replacement, replacementCaptureMatrix, true)
}

// RegexStringGsub implements the `gsub` DSL function.
func RegexStringGsub(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	_, replacementCaptureMatrix := ReplacementHasCaptures(replacement)
	return regexCompiledSubOrGsub(input, regex, replacement, replacementCaptureMatrix, false)
}

// regexCompiledSubOrGsub is the implementation for `sub`/`gsub` with compilex regex
// and replacement strings.
func regexCompiledSubOrGsub(
	input string,
	regex *regexp.Regexp,
	replacement string,
	replacementCaptureMatrix [][]int,
	breakOnFirst bool,
) string {
	matrix := regex.FindAllSubmatchIndex([]byte(input), -1)
	if matrix == nil || len(matrix) == 0 {
		return input
	}

	// Example return value from FindAllSubmatchIndex with input
	// "...ab_cde...fg_hij..." and regex "(..)_(...)":
	//
	// Matrix is [][]int{
	//   []int{3, 9, 3, 5, 6, 9},
	//   []int{12, 18, 12, 14, 15, 18},
	// }
	//
	// * 3-9 is for the entire match "ab_cde"
	// * 3-5 is for the first capture "ab"
	// * 6-9 is for the second capture "cde"
	//
	// * 12-18 is for the entire match "fg_hij"
	// * 12-14 is for the first capture "fg"
	// * 15-18 is for the second capture "hij"

	var buffer bytes.Buffer
	nonMatchStartIndex := 0

	for _, row := range matrix {
		buffer.WriteString(input[nonMatchStartIndex:row[0]])

		// "\0" .. "\9"
		captures := make([]string, 10)
		di := 0
		n := len(row)
		for si := 0; si < n && di <= 9; si += 2 {
			start := row[si]
			end := row[si+1]
			if start >= 0 && end >= 0 {
				captures[di] = input[start:end]
			}
			di += 1
		}

		// If the replacement had no captures, e.g. "xyz", we would insert it
		//
		//   "..."     -> "..."
		//   "ab_cde"  -> "xyz"   --- here
		//   "..."     -> "..."
		//   "fg_hij"  -> "xyz"   --- and here
		//   "..."     -> "..."
		//
		// using buffer.WriteString(replacement). However, this function exists
		// to handle the case when the replacement string has captures like
		// "\2:\1", so we need to produce
		//
		//   "..."     -> "..."
		//   "ab_cde"  -> "cde:ab"   --- here
		//   "..."     -> "..."
		//   "fg_hij"  -> "hij:fg"   --- and here
		//   "..."     -> "..."
		updatedReplacement := InterpolateCaptures(
			replacement,
			replacementCaptureMatrix,
			captures,
		)
		buffer.WriteString(updatedReplacement)

		nonMatchStartIndex = row[1]
		if breakOnFirst {
			break
		}
	}

	buffer.WriteString(input[nonMatchStartIndex:])
	return buffer.String()
}

// RegexStringMatchSimple is for simple boolean return without any substring captures.
func RegexStringMatchSimple(
	input string,
	sregex string,
) bool {
	regex := CompileMillerRegexOrDie(sregex)
	return RegexCompiledMatchSimple(input, regex)
}

// RegexCompiledMatchSimple is for simple boolean return without any substring captures.
func RegexCompiledMatchSimple(
	input string,
	regex *regexp.Regexp,
) bool {
	return regex.Match([]byte(input))
}

// RegexStringMatchWithCaptures implements the =~ DSL operator. The captures are stored in DSL
// state and may be used by a DSL statement after the =~. For example, in
//
//	sub($a, "(..)_(...)", "\1:\2")
//
// the replacement string is an argument to sub and therefore the captures are
// confined to the implementation of the sub function.  Similarly for gsub. But
// for the match operator, people can do
//
//	if ($x =~ "(..)_(...)") {
//	  ... other lines of code ...
//	  $y = "\2:\1"
//	}
//
// and the =~ callsite doesn't know if captures will be used or not. So,
// RegexStringMatchWithCaptures always returns the captures array. It is stored within the CST
// state.
func RegexStringMatchWithCaptures(
	input string,
	sregex string,
) (
	matches bool,
	capturesOneUp []string,
) {
	regex := CompileMillerRegexOrDie(sregex)
	return RegexCompiledMatchWithCaptures(input, regex)
}

// RegexCompiledMatchWithCaptures is the implementation for the =~ operator.  Without
// Miller-style regex captures this would a simple one-line
// regex.MatchString(input). However, we return the captures array for the
// benefit of subsequent references to "\0".."\9".
func RegexCompiledMatchWithCaptures(
	input string,
	regex *regexp.Regexp,
) (bool, []string) {
	matrix := regex.FindAllSubmatchIndex([]byte(input), -1)
	if matrix == nil || len(matrix) == 0 {
		// Set all captures to ""
		return false, make([]string, 10)
	}

	// "\0" .. "\9"
	captures := make([]string, 10)

	// If there are multiple matches -- e.g. input is
	//
	//   "...ab_cde...fg_hij..."
	//
	// with regex
	//
	//   "(..)_(...)"
	//
	// -- then we only consider the first match: boolean return value is true
	// (the input string matched the regex), and the captures array will map
	// "\1" to "ab" and "\2" to "cde".
	row := matrix[0]
	n := len(row)

	// Example return value from FindAllSubmatchIndex with input
	// "...ab_cde...fg_hij..." and regex "(..)_(...)":
	//
	// Matrix is [][]int{
	//   []int{3, 9, 3, 5, 6, 9},
	//   []int{12, 18, 12, 14, 15, 18},
	// }
	//
	// As noted above we look at only the first row.
	//
	// * 3-9 is for the entire match "ab_cde"
	// * 3-5 is for the first capture "ab"
	// * 6-9 is for the second capture "cde"

	di := 0
	for si := 0; si < n && di <= 9; si += 2 {
		start := row[si]
		end := row[si+1]
		if start >= 0 && end >= 0 {
			captures[di] = input[start:end]
		}
		di += 1
	}

	return true, captures
}

// MakeEmptyCaptures is for initial CST state at the start of executing the DSL expression for the
// current record.  Even if '$x =~ "(..)_(...)" set "\1" and "\2" on the previous record, at start
// of processing for the current record we need to start with a clean slate. This is in support of
// CST state, which `=~` semantics requires.
func MakeEmptyCaptures() []string {
	return nil
}

// ReplacementHasCaptures is used by the CST builder to see if string-literal is like "foo bar" or
// "foo \1 bar" -- in the latter case it needs to retain the compiled offsets-matrix information.
// This is in support of CST state, which `=~` semantics requires.
func ReplacementHasCaptures(
	replacement string,
) (
	hasCaptures bool,
	matrix [][]int,
) {
	if captureDetector.MatchString(replacement) {
		return true, captureSplitter.FindAllSubmatchIndex([]byte(replacement), -1)
	} else {
		return false, nil
	}
}

// InterpolateCaptures example:
//
// * Input $x is "ab_cde"
//
//   - DSL expression
//     if ($x =~ "(..)_(...)") {
//     ... other lines of code ...
//     $y = "\2:\1";
//     }
//
// * InterpolateCaptures is used on the evaluation of "\2:\1"
//
// * replacementString is "\2:\1"
//
//   - replacementMatrix contains precomputed/cached offsets for the "\2" and
//     "\1" substrings within "\2:\1"
//
//   - captures has slot 0 being "ab_cde" (for "\0"), slot 1 being "ab" (for "\1"),
//     slot 2 being "cde" (for "\2"), and slots 3-9 being "".
func InterpolateCaptures(
	replacementString string,
	replacementMatrix [][]int,
	captures []string,
) string {
	if replacementMatrix == nil || captures == nil {
		return replacementString
	}
	var buffer bytes.Buffer

	nonMatchStartIndex := 0

	for _, row := range replacementMatrix {
		start := row[0]
		buffer.WriteString(replacementString[nonMatchStartIndex:row[0]])

		// Map "\0".."\9" to integer index 0..9
		index := replacementString[start+1] - '0'
		buffer.WriteString(captures[index])

		nonMatchStartIndex = row[1]
	}

	buffer.WriteString(replacementString[nonMatchStartIndex:])

	return buffer.String()
}
