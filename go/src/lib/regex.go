// ================================================================
// Support for regexes in Miller.
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
//   https://github.com/johnkerl/miller/issues/388 has a better suggestion
//   which would make the captures explicit as variables, rather than implicit
//   within CST state -- regardless, the current syntax will still be supprted
//   for backward compatability and so is here to stay.) Here we make use of Go
//   regexp-library functions to write to, and then later interpolate from, a
//   captures array which is stored within CST state. (See the `runtime.State`
//   object.)
//
// * "\0" is for a full match; "\1" .. "\9" are for submatch cqptures. E.g.
//   if $x is "foobarbaz" and the regex is "foo(.)(..)baz", then "\0" is
//   "foobarbaz", "\1" is "b", "\2" is "ar", and "\3".."\9" are "".
// ================================================================

package lib

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// captureDetector is used to see if a string literal interpolates previous
// captures (like "\2:\1") or not (like "2:1").
var captureDetector = regexp.MustCompile("\\\\[0-9]")

// captureSplitter is used to precompute an offsets matrix for strings like
// "\2:\1" so they don't need to be recomputed on every record.
var captureSplitter = regexp.MustCompile("(\\\\[0-9])")

// CompileMillerRegex wraps Go regex-compile with some Miller-specific syntax
// which predate the port of Miller from C to Go.  Miller regexes use a final
// 'i' to indicate case-insensitivity; Go regexes use an initial "(?i)".
//
// (See also mlr.bnf where we specify which things can be backslash-escaped
// without a syntax error at parse time.)
//
// * If the regex_string is of the form a.*b, compiles it case-sensisitively.
// * If the regex_string is of the form "a.*b", compiles a.*b case-sensisitively.
// * If the regex_string is of the form "a.*b"i, compiles a.*b case-insensitively.
func CompileMillerRegex(regexString string) (*regexp.Regexp, error) {
	n := len(regexString)
	if n < 2 {
		return regexp.Compile(regexString)
	}

	// TODO: rethink this. This will strip out things people have entered, e.g. "\"...\"".
	// The parser-to-AST will have stripped the outer and we'll strip the inner and the
	// user's intent will be lost.
	//
	// TODO: make separate functions for calling from parser-to-AST (string
	// literals) and from verbs (like cut -r or having-fields).

	if strings.HasPrefix(regexString, "\"") && strings.HasSuffix(regexString, "\"") {
		return regexp.Compile(regexString[1 : n-1])
	}
	if strings.HasPrefix(regexString, "/") && strings.HasSuffix(regexString, "/") {
		return regexp.Compile(regexString[1 : n-1])
	}

	if strings.HasPrefix(regexString, "\"") && strings.HasSuffix(regexString, "\"i") {
		return regexp.Compile("(?i)" + regexString[1:n-2])
	}
	if strings.HasPrefix(regexString, "/") && strings.HasSuffix(regexString, "/i") {
		return regexp.Compile("(?i)" + regexString[1:n-2])
	}

	return regexp.Compile(regexString)
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
func RegexSplitString(regex *regexp.Regexp, input string, n int) []string {
	if input == "" {
		return make([]string, 0)
	} else {
		return regex.Split(input, n)
	}
}

// MakeEmptyRegexCaptures is for initial CST state at the start of executing
// the DSL expression for the current record.  Even if '$x =~ "(..)_(...)" set
// "\1" and "\2" on the previous record, at start of processing for the current
// record we need to start with a clean slate.
func MakeEmptyRegexCaptures() []string {
	return nil
}

// RegexReplacementHasCaptures is used by the CST builder to see if
// string-literal is like "foo bar" or "foo \1 bar" -- in the latter case it
// needs to retain the compiled offsets-matrix information.
func RegexReplacementHasCaptures(
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

// RegexMatches implements the =~ DSL operator. The captures are stored in DSL
// state and may be used by a DSL statement after the =~. For example, in
//
//   sub($a, "(..)_(...)", "\1:\2")
//
// the replacement string is an argument to sub and therefore the captures are
// confined to the implementation of the sub function.  Similarly for gsub. But
// for the match operator, people can do
//
//   if ($x =~ "(..)_(...)") {
//     ... other lines of code ...
//     $y = "\2:\1"
//   }
//
// and the =~ callsite doesn't know if captures will be used or not. So,
// RegexMatches always returns the captures array. It is stored within the CST
// state.
func RegexMatches(
	input string,
	sregex string,
) (
	matches bool,
	capturesOneUp []string,
) {
	regex := CompileMillerRegexOrDie(sregex)
	return RegexMatchesCompiled(input, regex)
}

// RegexMatchesCompiled is the implementation for the =~ operator.  Without
// Miller-style regex captures this would a simple one-line
// regex.MatchString(input). However, we return the captures array for the
// benefit of subsequent references to "\0".."\9".
func RegexMatchesCompiled(
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

// InterpolateCaptures example:
// * Input $x is "ab_cde"
// * DSL expression
//     if ($x =~ "(..)_(...)") {
//       ... other lines of code ...
//       $y = "\2:\1";
//     }
// * InterpolateCaptures is used on the evaluation of "\2:\1"
// * replacementString is "\2:\1"
// * replacementMatrix contains precomputed/cached offsets for the "\2" and
//   "\1" substrings within "\2:\1"
// * captures has slot 0 being "ab_cde" (for "\0"), slot 1 being "ab" (for "\1"),
//   slot 2 being "cde" (for "\2"), and slots 3-9 being "".
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

// RegexSub implements the sub DSL function.
func RegexSub(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	_, replacementCaptureMatrix := RegexReplacementHasCaptures(replacement)
	return RegexSubCompiled(input, regex, replacement, replacementCaptureMatrix)
}

// RegexSubCompiled is the same as RegexSub but with compiled regex and
// replacement strings.
func RegexSubCompiled(
	input string,
	regex *regexp.Regexp,
	replacement string,
	replacementCaptureMatrix [][]int,
) string {
	return regexSubGsubCompiled(input, regex, replacement, replacementCaptureMatrix, true)
}

// RegexGsub implements the gsub DSL function.
func RegexGsub(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	_, replacementCaptureMatrix := RegexReplacementHasCaptures(replacement)
	return regexSubGsubCompiled(input, regex, replacement, replacementCaptureMatrix, false)
}

// regexSubGsubCompiled is the implementation for sub/gsub with compilex regex
// and replacement strings.
func regexSubGsubCompiled(
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
