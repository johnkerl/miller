package lib

// TODO:
// * cst state for captures array
// * reset-hook for start of execution
//   o UTs for that
// * flesh out RegexCaptureBinaryFunctionCallsiteNode to do that

// xxx update comment

// The key is the Go library's regex.FindAllStringIndex.  It gives us start
// (inclusive) and end (exclusive) indices for matches.
//
// Example: for pattern "foo" and input "abc foo def foo ghi" we'll have
// matrix [[4 7] [12 15]] which indicates matches from positions 4-6 and
// 12-14.  We simply need to concatenate
// *  0-3  "abc "  not matching
// *  4-6  "foo"   matching
// *  7-11 " def " not matching
// * 12-14 "foo"   matching
// * 15-18 " ghi"  not matching
//
// Example: with pattern "f.*o" and input "abc foo def foo ghi" we'll have
// matrix [[4 15]] so "foo def foo" will be a matched substring.

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ================================================================
// API functions

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

// xxx MakeEmptyCaptures function for CST state

// xxx ReplacementHasCaptures function

// RegexSubWithoutCaptures implements the sub DSL function when the replacement
// string has none of "\1".."\9".
func RegexSubWithoutCaptures(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	// to do
	return RegexSubWithoutCapturesCompiled(input, regex, replacement)
}

// RegexSubWithoutCapturesCompiled is the same as RegexSubWithoutCaptures but
// with compiled regex instead of regex-as-string.
func RegexSubWithoutCapturesCompiled(
	input string,
	regex *regexp.Regexp,
	replacement string,
) string {
	onFirst := true
	output := regex.ReplaceAllStringFunc(input, func(s string) string {
		if !onFirst {
			return s
		}
		onFirst = false
		return regex.ReplaceAllString(s, replacement)
	})
	return output
}

// RegexSubWithCaptures implements the sub DSL function when the replacement
// string has one or more of "\1".."\9".
func RegexSubWithCaptures(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	return RegexSubCompiledWithCaptures(input, regex, replacement)
}

// RegexSubCompiledWithCaptures is the same as RegexSubWithCaptures but
// with compiled regex instead of regex-as-string.
func RegexSubCompiledWithCaptures(
	input string,
	regex *regexp.Regexp,
	replacement string,
) string {
	return regexSubGsubWithCapturesAux(input, regex, replacement, true)
}

// RegexGsubWithoutCaptures implements the gsub DSL function when the replacement
// string has none of "\1".."\9".
func RegexGsubWithoutCaptures(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	return regex.ReplaceAllString(input, replacement)
}

// RegexGsubWithoutCaptures implements the gsub DSL function when the replacement
// string has one or more of "\1".."\9".
func RegexGsubWithCaptures(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	return regexSubGsubWithCapturesAux(input, regex, replacement, false)
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
func RegexMatches(input string, sregex string) (matches bool, capturesOneUp []string) {
	regex := CompileMillerRegexOrDie(sregex)
	return regexMatchesAux(input, regex)
}

// ----------------------------------------------------------------
// Package-internal/implementation functions

// xxx:
// $ go run foo.go "ab_cde  ab_cde" "(ab)_(cde)"
// MATRIX [][]int{[]int{0, 6, 0, 2, 3, 6}, []int{8, 14, 8, 10, 11, 14}}
// 0 6 "ab_cde"
// n:6
//   2 0 2 "ab"
//   4 3 6 "cde"
// 8 14 "ab_cde"
// n:6
//   2 8 10 "ab"
//   4 11 14 "cde"

// regexSubGsubWithCapturesAux is the implementation for sub/gsub when the
// replacement string uses captures in the form "\1".."\9".
func regexSubGsubWithCapturesAux(
	input string,
	regex *regexp.Regexp,
	replacement string,
	breakOnFirst bool,
) string {
	matrix := regex.FindAllSubmatchIndex([]byte(input), -1)
	if matrix == nil || len(matrix) == 0 {
		return input
	}

	var buffer bytes.Buffer // Faster since os.Stdout is unbuffered
	nonMatchStartIndex := 0

	for _, startEnd := range matrix {
		buffer.WriteString(input[nonMatchStartIndex:startEnd[0]])
		buffer.WriteString(replacement)
		nonMatchStartIndex = startEnd[1]
		if breakOnFirst {
			break
		}
	}

	buffer.WriteString(input[nonMatchStartIndex:])
	return buffer.String()
}

// regexMatchesAux is the implementation for the =~ operator.
func regexMatchesAux(
	input string,
	regex *regexp.Regexp,
) (bool, []string) {
	matrix := regex.FindAllSubmatchIndex([]byte(input), -1)
	if matrix == nil || len(matrix) == 0 {
		return false, nil
	}

	// Slot 0 is ""; then slots 1..9 for "\1".."\9".
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
	if n == 2 {
		// There were no regex captures like "(..)" within the regex.
		return true, nil
	}

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

	di := 1
	for si := 2; si < n; si += 2 {
		if di > 9 {
			break
		}
		start := row[si]
		end := row[si+1]
		captures[di] = input[start:end]
		di += 1
	}

	return true, captures
}
