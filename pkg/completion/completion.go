// Package completion handles Shell completion
package completion

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/transformers"
)

func DoCompletion() {
	if os.Args[1] != "_complete_bash" {
		return
	}
	if len(os.Args) < 5 {
		Debug()
		return
	}
	// See: https://www.gnu.org/software/bash/manual/html_node/Programmable-Completion-Builtins.html#index-complete
	// Bash completion calls with three arguments: $1 is the name of the
	// command whose arguments are being completed, $2 is the word being
	// completed, and $3 is the word preceding the word being completed. Since
	// we already set one argument, the rest of them are shifted by one
	// position i.e.  `mlr _complete_bash <mlr> <last> <prev>`,
	last := os.Args[3]
	prev := os.Args[4]
	if prev == "then" {
		matches := GetMatchingVerbs(last)
		// See: https://www.gnu.org/software/bash/manual/html_node/Bash-Variables.html#index-COMP_005fTYPE
		// When tab is hit two times, bash sets COMP_TYPE to ascii value of `?` i.e. 63
		if len(matches) == 1 && matches[0].Verb == last && os.Getenv("COMP_TYPE") == "63" {
			v := matches[0]
			v.UsageFunc(os.Stdout)
		} else {
			sort.Slice(matches, func(i, j int) bool { return matches[i].Verb < matches[j].Verb })
			for _, verb := range matches {
				fmt.Println(verb.Verb)
			}
		}
	}
}

func GetMatchingVerbs(partVerb string) []*transformers.TransformerSetup {
	var matches []*transformers.TransformerSetup
	for _, verb := range transformers.TRANSFORMER_LOOKUP_TABLE {
		localv := verb
		if strings.HasPrefix(verb.Verb, partVerb) {
			matches = append(matches, &localv)
		}
	}
	return matches
}

func Debug() {
	for i, arg := range os.Args {
		fmt.Fprintln(os.Stderr, i, arg)
	}
	for _, val := range os.Environ() {
		if strings.HasPrefix(val, "COMP") {
			fmt.Fprintln(os.Stderr, val)
		}
	}
}
