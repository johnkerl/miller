// xxx update comment

package lib

// TODO:
// * cst state for captures array
// * reset-hook for start of execution
//   o UTs for that
// * flesh out RegexCaptureBinaryFunctionCallsiteNode to do that

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TODO: comment
var captureDetector = regexp.MustCompile("\\\\[0-9]")
var captureSplitter = regexp.MustCompile("(\\\\[0-9])")

// CompileMillerRegex wraps Go regex-compile with some Miller-specific syntax
// which predate the port of Miller from C to Go.  Miller regexes use a final
// 'i' to indicate case-insensitivity; Go regexes use an initial "(?i)".  Also
// (TODO) I need to find all the right things to backslash-escape in Go.
//
// * If the regex_string is of the form a.*b, compiles it case-sensisitively.
// * If the regex_string is of the form "a.*b", compiles a.*b case-sensisitively.
// * If the regex_string is of the form "a.*b"i, compiles a.*b case-insensitively.
func CompileMillerRegex(regexString string) (*regexp.Regexp, error) {
	if !strings.HasPrefix(regexString, "\"") {
		return regexp.Compile(regexString)
	} else {
		n := len(regexString)
		if n < 2 {
			// This means the user entered "\"" which in the parser-to-AST was
			// presented as " which we need to handle as a single-character
			// regex.
			return regexp.Compile(regexString)
		}
		if strings.HasSuffix(regexString, "\"") {
			// TODO: rethink this. This will strip out things people have entered, e.g. "\"...\"".
			// The parser-to-AST will have stripped the outer and we'll strip the inner and the
			// user's intent will be lost.
			//
			// TODO: make separate functions for calling from parser-to-AST (string
			// literals) and from verbs (like cut -r or having-fields).
			return regexp.Compile(regexString[1 : n-1])
		} else if strings.HasSuffix(regexString, "\"i") {
			return regexp.Compile("(?i)" + regexString[1:n-2])
		} else {
			// The user can enter things like "\".." which comes in as ".. which is fine.
			return regexp.Compile(regexString)
		}
	}
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

// TODO: comment
func MakeEmptyRegexCaptures() []string {
	return nil
}

// xxx comment
// xxx MakeEmptyCaptures function for CST state
// xxx UT
// xxx more UT cases
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

// RegexMatches implements the =~ DSL operator. There is no
// with-captures/without-captures variant-pair since the captures are stored in
// DSL state and may be used by a DSL statement after the =~. For example, in
//
//   sub($a, "(..)_(...)", "\1:\2")
//
// the replacement string is an argument to sub and therefore the captures are
// confined to the implementation of the sub function.  Similarly for gsub. But
// for the match operator, people can do
//
//   if ($a =~ "(..)_(...)") {
//     $b = "\1:\2"
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
// benefit of subsequent references to "\1".."\9".
func RegexMatchesCompiled(
	input string,
	regex *regexp.Regexp,
) (bool, []string) {
	matrix := regex.FindAllSubmatchIndex([]byte(input), -1)
	if matrix == nil || len(matrix) == 0 {
		// xxx temp
		// return false, nil
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
		if start >= 0 && end >= 0 { // TODO: comment
			captures[di] = input[start:end]
		}
		di += 1
	}

	return true, captures
}

// TODO: comment
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

		// xxx comment
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

		// xxx already have split up replacement into its matrix, before entering this helper
		// xxx have another helper to iterate over, taking &buffer as arg ...

		nonMatchStartIndex = row[1]
		if breakOnFirst {
			break
		}
	}

	buffer.WriteString(input[nonMatchStartIndex:])
	return buffer.String()
}
