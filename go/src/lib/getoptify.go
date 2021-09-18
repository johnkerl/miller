package lib

import (
	"regexp"
	"strings"
)

// Getoptify expands "-xyz" into "-x -y -z" while leaving "--xyz" intact. This
// is a keystroke-saver for the user.
//
// This is OK to do here globally since Miller is quite consistent (in main,
// verbs, and auxents) that multi-character options start with two dashes, e.g.
// "--csv". (The sole exception is the sort verb's -nf/-nr which are handled
// specially there.)
//
// Additionally, we split "--foo=bar" into "--foo" and "bar".
func Getoptify(inargs []string) []string {
	expandRegex := regexp.MustCompile("^-[a-zA-Z0-9]+$")
	splitRegex := regexp.MustCompile("^--[^=]+=.+$")
	outargs := make([]string, 0)
	for _, inarg := range inargs {
		if expandRegex.MatchString(inarg) {
			for _, c := range inarg[1:] {
				outargs = append(outargs, "-"+string(c))
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
