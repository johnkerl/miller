package lib

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Miller regexes use a final 'i' to indicate case-insensitivity; Go regexes
// use an initial "(?i)".  Also (TODO) I need to find all the right things to
// backslash-escape in Go.
//
// * If the regex_string is of the form a.*b, compiles it using cflags without REG_ICASE.
// * If the regex_string is of the form "a.*b", compiles a.*b using cflags without REG_ICASE.
// * If the regex_string is of the form "a.*b"i, compiles a.*b using cflags with REG_ICASE.
func CompileMillerRegex(regexString string) (*regexp.Regexp, error) {
	if !strings.HasPrefix(regexString, "\"") {
		return regexp.Compile(regexString)
	} else {
		n := len(regexString)
		if n < 2 {
			return nil, errors.New(
				fmt.Sprintf(
					"%s: imbalanced double-quote in regex [%s].\n",
					MlrExeName(), regexString,
				),
			)
		}
		if strings.HasSuffix(regexString, "\"") {
			return regexp.Compile(regexString[1 : n-1])
		} else if strings.HasSuffix(regexString, "\"i") {
			return regexp.Compile("(?i)" + regexString[1:n-2])
		} else {
			return nil, errors.New(
				fmt.Sprintf(
					"%s: imbalanced double-quote in regex [%s].\n",
					MlrExeName(), regexString,
				),
			)
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
