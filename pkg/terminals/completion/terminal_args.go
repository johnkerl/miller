// Completion for the arguments of terminal subcommands -- chiefly
// `mlr help <TAB>`, which completes help topics, and a topic's own argument
// (e.g. `mlr help verb <TAB>` -> verb names). `mlr completion <TAB>` completes
// bash/zsh. Other terminals have no argument completion.

package completion

import (
	"github.com/johnkerl/miller/v6/pkg/terminals/help"
	"github.com/johnkerl/miller/v6/pkg/terminals/registry"
)

// terminalNameSet is the set of terminal subcommand names, for quick lookup
// during the command-line walk.
var terminalNameSet = func() map[string]bool {
	m := make(map[string]bool, len(registry.Names))
	for _, name := range registry.Names {
		m[name] = true
	}
	return m
}()

func isTerminalName(s string) bool {
	return terminalNameSet[s]
}

// completeTerminalArgs produces candidates for the words following a terminal
// subcommand. args are the words already typed after the terminal name, before
// the cursor; cur is the word being completed.
func completeTerminalArgs(terminal string, args []string, cur string) Result {
	switch terminal {

	case registry.Help:
		switch len(args) {
		case 0:
			// `mlr help <TAB>`: the help topics.
			return Result{DirectiveCandidates, filterByPrefix(sortedUnion(help.GetTopicNames()), cur)}
		case 1:
			// `mlr help <topic> <TAB>`: the topic's own argument, if it takes
			// one (e.g. a verb, function, keyword, or flag name).
			if values := helpTopicArgValues(args[0]); values != nil {
				return Result{DirectiveCandidates, filterByPrefix(values, cur)}
			}
		}
		return Result{DirectiveCandidates, nil}

	case registry.Completion:
		if len(args) == 0 {
			return Result{DirectiveCandidates, filterByPrefix([]string{"bash", "zsh"}, cur)}
		}
		return Result{DirectiveCandidates, nil}
	}

	// Other terminals (version, repl, regtest, script, terminal-list) have no
	// argument completion.
	return Result{DirectiveCandidates, nil}
}

// helpTopicArgValues returns the candidate values for the argument of a
// `mlr help {topic}` whose topic takes a name argument, or nil otherwise.
func helpTopicArgValues(topic string) []string {
	switch topic {
	case "verb":
		return verbNames()
	case "flag":
		return mainFlagNames()
	case "function":
		return sortedUnion(help.GetFunctionNames())
	case "keyword":
		return sortedUnion(help.GetKeywordNames())
	}
	return nil
}
