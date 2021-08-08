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

// ----------------------------------------------------------------
// Miller regexes use a final 'i' to indicate case-insensitivity; Go regexes
// use an initial "(?i)".  Also (TODO) I need to find all the right things to
// backslash-escape in Go.
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

func CompileMillerRegexOrDie(regexString string) *regexp.Regexp {
	regex, err := CompileMillerRegex(regexString)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	return regex
}

// ----------------------------------------------------------------
// API functions

// xxx ReplacementHasCaptures function

func RegexSubWithoutCaptures(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	return regexSubGsubWithCapturesAux(input, regex, replacement, true)
}

func RegexSubWithCaptures(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	return regexSubGsubWithCapturesAux(input, regex, replacement, true)
}

func RegexSubWithoutCapturesCompiled(
	input string,
	regex *regexp.Regexp,
	replacement string,
) string {
	return regexSubGsubWithCapturesAux(input, regex, replacement, true)
}

func RegexSubCompiledWithCaptures(
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

func RegexGsubWithoutCaptures(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	return regex.ReplaceAllString(input, replacement)
}

func RegexGsubWithCaptures(
	input string,
	sregex string,
	replacement string,
) string {
	regex := CompileMillerRegexOrDie(sregex)
	return regexSubGsubWithCapturesAux(input, regex, replacement, false)
}

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

	// xxx instantiate a RegexCaptures object

	// xxx update comment for FindAllSubmatchIndex

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

// xxx comment:
// xxx so we early-out
//
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

func regexMatchesAux(
	input string,
	regex *regexp.Regexp,
) (bool, []string) {
	matrix := regex.FindAllSubmatchIndex([]byte(input), -1)
	if matrix == nil || len(matrix) == 0 {
		return false, nil
	}

	// 0 is ""; 1..9 for "\1".."\9"
	captures := make([]string, 10)

	// xxx update comment for FindAllSubmatchIndex

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

	// xxx comment first outer-match only
	row := matrix[0]
	n := len(row)
	if n == 2 {
		// xxx comment no captures
		return true, nil
	}

	i := 1
	for j := 2; j < n; j += 2 {
		if i > 9 {
			break
		}
		start := row[j]
		end := row[j+1]
		captures[i] = input[start:end]
		i += 1
	}

	return true, captures
}
