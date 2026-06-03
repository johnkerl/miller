package lib

import (
	"regexp"
	"strings"
)

// Getoptify expands "-xyz" into "-x -y -z" while leaving "--xyz" intact. This
// is a keystroke-saver for the user.
//
// This is OK to do here globally since Miller is quite consistent (in main,
// verbs, auxents, and terminals) that multi-character options start with two
// dashes, e.g.  "--csv". (The sole exception is the sort verb's -nf/-nr which
// are handled specially there.)
//
// Additionally, we split "--foo=bar" into "--foo" and "bar".
func Getoptify(inargs []string) []string {
	expandRegex := regexp.MustCompile("^-[a-zA-Z0-9]+$")
	splitRegex := regexp.MustCompile("^--[^=]+=.+$")
	numberRegex := regexp.MustCompile("^-[0-9]+$")
	outargs := []string{}
	for _, inarg := range inargs {
		if expandRegex.MatchString(inarg) {
			if numberRegex.MatchString(inarg) {
				// Don't expand things like '-12345' which are (likely!) numeric arguments to verbs.
				// Example: 'mlr unsparsify --fill-with -99999'.
				outargs = append(outargs, inarg)
			} else {
				rest := inarg[1:]
				for i := 0; i < len(rest); i++ {
					// Pass integers without a leading dash, so that negative integers can be represented.
					// Example: `head -n 4` and head -n4` can be differentiated from `head -n -4`.
					if rest[i] >= '0' && rest[i] <= '9' {
						outargs = append(outargs, rest[i:])
						break
					}
					outargs = append(outargs, "-"+string(rest[i]))
				}
			}
		} else if splitRegex.MatchString(inarg) {
			pair := strings.SplitN(inarg, "=", 2)
			InternalCodingErrorIf(len(pair) != 2)
			outargs = append(outargs, pair[0])
			outargs = append(outargs, pair[1])
		} else {
			outargs = append(outargs, inarg)
		}
	}
	return outargs
}
