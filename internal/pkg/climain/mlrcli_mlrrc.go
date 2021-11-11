package climain

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"mlr/internal/pkg/cli"
)

// loadMlrrcOrDie rule: If $MLRRC is set, use it and only it.  Otherwise try
// first $HOME/.mlrrc and then ./.mlrrc but let them stack: e.g. $HOME/.mlrrc
// is lots of settings and maybe in one subdir you want to override just a
// setting or two.
func loadMlrrcOrDie(
	options *cli.TOptions,
) {
	env_mlrrc := os.Getenv("MLRRC")

	if env_mlrrc != "" {
		if env_mlrrc == "__none__" {
			return
		}
		if tryLoadMlrrc(options, env_mlrrc) {
			return
		}
	}

	env_home := os.Getenv("HOME")
	if env_home != "" {
		path := env_home + "/.mlrrc"
		tryLoadMlrrc(options, path)
	}

	tryLoadMlrrc(options, "./.mlrrc")
}

// tryLoadMlrrc is a helper function for loadMlrrcOrDie.
func tryLoadMlrrc(
	options *cli.TOptions,
	path string,
) bool {
	handle, err := os.Open(path)
	if err != nil {
		return false
	}
	defer handle.Close()

	lineReader := bufio.NewReader(handle)

	eof := false
	lineno := 0
	for !eof {
		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
			break
		}
		lineno++

		if err != nil {
			fmt.Fprintln(os.Stderr, "mlr", err)
			os.Exit(1)
			return false
		}

		// This is how to do a chomp:
		// TODO: handle \r\n with libified solution.
		line = strings.TrimRight(line, "\n")

		if !handleMlrrcLine(options, line) {
			fmt.Fprintf(os.Stderr, "%s: parse error at file \"%s\" line %d: %s\n",
				"mlr", path, lineno, line,
			)
			os.Exit(1)
		}
	}

	return true
}

// handleMlrrcLine is a helper function for loadMlrrcOrDie.
func handleMlrrcLine(
	options *cli.TOptions,
	line string,
) bool {

	// Comment-strip
	re := regexp.MustCompile("#.*")
	line = re.ReplaceAllString(line, "")

	// Left-trim / right-trim
	line = strings.TrimSpace(line)

	if line == "" { // line was whitespace-only
		return true
	}

	// Prepend initial "--" if it's not already there
	if !strings.HasPrefix(line, "-") {
		line = "--" + line
	}

	// Split line into args array
	args := strings.Fields(line)
	argi := 0
	argc := len(args)

	if args[0] == "--prepipe" || args[0] == "--prepipex" {
		// Don't allow code execution via .mlrrc
		return false
	} else if args[0] == "--load" || args[0] == "--mload" {
		// Don't allow code execution via .mlrrc
		return false
	} else if cli.FLAG_TABLE.Parse(args, argc, &argi, options) {
		// handled
	} else {
		return false
	}

	return true
}
