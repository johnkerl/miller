package lib

import (
	"regexp"
)

// Getoptify expands "-xyz" into "-x -y -z" while leaving "--xyz" intact. This
// is a keystroke-saver for the user.
//
// This is OK to do here globally since Miller is quite consistent (in main,
// verbs, and auxents) that multi-character options start with two dashes, e.g.
// "--csv". (The sole exception is the sort verb's -nf/-nr which are handled
// specially there.)
func Getoptify(inargs []string) []string {
	regex := regexp.MustCompile("^-[a-zA-Z0-9]+$")
	outargs := make([]string, 0)
	for _, inarg := range inargs {
		if regex.MatchString(inarg) {
			for _, c := range inarg[1:] {
				outargs = append(outargs, "-"+string(c))
			}
		} else {
			outargs = append(outargs, inarg)
		}
	}
	return outargs
}
