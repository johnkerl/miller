// Entry point for the `mlr completion` terminal.

package completion

import (
	"fmt"
	"os"
	"strconv"
)

// CompletionMain is the handler for `mlr completion ...`, dispatched by the
// terminals framework.
//
//	mlr completion bash        Print a bash completion script.
//	mlr completion zsh         Print a zsh completion script.
//	mlr completion complete <cword> <word>...
//	                           Internal: emit completion candidates for the
//	                           given command-line words. Called by the shim
//	                           scripts above; not intended for direct use.
func CompletionMain(args []string) int {
	// args[0] is "completion".
	if len(args) < 2 {
		printUsage(os.Stdout)
		return 0
	}

	switch args[1] {
	case "bash":
		fmt.Print(bashScript)
		return 0
	case "zsh":
		fmt.Print(zshScript)
		return 0
	case "complete":
		return runComplete(args[2:])
	case "-h", "--help", "help":
		printUsage(os.Stdout)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "mlr completion: unrecognized subcommand %q.\n", args[1])
		fmt.Fprintf(os.Stderr, "Please run \"mlr completion --help\" for usage information.\n")
		return 1
	}
}

// runComplete handles `mlr completion complete <cword> <word>...`, printing a
// directive line followed by candidate words, one per line.
func runComplete(rest []string) int {
	if len(rest) < 1 {
		// Nothing to complete.
		fmt.Println(string(DirectiveFiles))
		return 0
	}

	cword, err := strconv.Atoi(rest[0])
	if err != nil {
		fmt.Println(string(DirectiveFiles))
		return 0
	}
	words := rest[1:]

	result := Complete(words, cword)

	fmt.Println(string(result.Directive))
	for _, candidate := range result.Candidates {
		fmt.Println(candidate)
	}
	return 0
}

func printUsage(o *os.File) {
	fmt.Fprintf(o, `Usage: mlr completion {bash|zsh}
Generates a shell tab-completion script for Miller.

Bash:
  Add to your ~/.bashrc:
    eval "$(mlr completion bash)"
  Or install system-wide:
    mlr completion bash > /etc/bash_completion.d/mlr
  Note: prefer 'eval' over 'source <(mlr completion bash)'. The latter
  silently fails on the bash 3.2 that ships with macOS, where sourcing from a
  process-substitution FIFO can read nothing.

Zsh:
  Add to your ~/.zshrc:
    eval "$(mlr completion zsh)"
  Or place the output on your $fpath, e.g.:
    mlr completion zsh > "${fpath[1]}/_mlr"
  The script initializes zsh's completion system (compinit) if your startup
  files have not done so already.

Completion is context-aware across Miller's then-chains: it offers main flags
and verb names before the first verb, the current verb's flags inside a verb,
verb names after 'then', and filenames where appropriate.
`)
}
