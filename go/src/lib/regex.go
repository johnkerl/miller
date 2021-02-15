package lib

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

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
