package climain

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/johnkerl/miller/v6/pkg/cli"
)

// loadMlrrcOrDie is a fatal-error wrapper around loadMlrrc.
func loadMlrrcOrDie(
	options *cli.TOptions,
	profileName string,
) {
	err := loadMlrrc(options, profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "mlr: %v\n", err)
		os.Exit(1)
	}
}

// loadMlrrc rule: If $MLRRC is set, use it and only it.  Otherwise try first
// $HOME/.mlrrc and then ./.mlrrc but let them stack: e.g. $HOME/.mlrrc is
// lots of settings and maybe in one subdir you want to override just a
// setting or two.
//
// The profileName argument comes from the --profile {name} / -P {name} main
// flag. Empty string means no profile was requested: only global (pre-section)
// lines of the .mlrrc file(s) are applied, and any [section] blocks are
// skipped. Non-empty means global lines are applied first, then the lines in
// any [profileName] section. It's a fatal error if a profile was requested
// but no matching section exists in any .mlrrc file processed.
func loadMlrrc(
	options *cli.TOptions,
	profileName string,
) error {
	foundProfile := false
	loadedPaths := []string{}

	env_mlrrc := os.Getenv("MLRRC")

	if env_mlrrc != "" {
		if env_mlrrc == "__none__" {
			if profileName != "" {
				return fmt.Errorf(
					"--profile \"%s\" was specified, but .mlrrc processing is disabled since the MLRRC environment variable is set to \"__none__\"",
					profileName,
				)
			}
			return nil
		}
		loaded, err := tryLoadMlrrc(options, env_mlrrc, profileName, &foundProfile)
		if err != nil {
			return err
		}
		if loaded {
			return checkMlrrcProfileWasFound(profileName, foundProfile, []string{env_mlrrc})
		}
	}

	env_home := os.Getenv("HOME")
	if env_home != "" {
		path := env_home + "/.mlrrc"
		loaded, err := tryLoadMlrrc(options, path, profileName, &foundProfile)
		if err != nil {
			return err
		}
		if loaded {
			loadedPaths = append(loadedPaths, path)
		}
	}

	loaded, err := tryLoadMlrrc(options, "./.mlrrc", profileName, &foundProfile)
	if err != nil {
		return err
	}
	if loaded {
		loadedPaths = append(loadedPaths, "./.mlrrc")
	}

	return checkMlrrcProfileWasFound(profileName, foundProfile, loadedPaths)
}

// checkMlrrcProfileWasFound is a helper function for loadMlrrc: if a profile
// was requested via --profile {name} / -P {name}, there must be a matching
// [name] section in at least one processed .mlrrc file.
func checkMlrrcProfileWasFound(
	profileName string,
	foundProfile bool,
	loadedPaths []string,
) error {
	if profileName == "" || foundProfile {
		return nil
	}
	if len(loadedPaths) == 0 {
		return fmt.Errorf(
			"--profile \"%s\" was specified, but no .mlrrc file was found",
			profileName,
		)
	}
	return fmt.Errorf(
		"--profile \"%s\" was specified, but no [%s] section was found in %s",
		profileName, profileName, strings.Join(loadedPaths, ", "),
	)
}

// tryLoadMlrrc is a helper function for loadMlrrc. The first return value is
// whether the file could be opened at all: an unopenable file is not an error
// (that's the normal case when no .mlrrc file exists). The second is any
// parse error within an opened file.
func tryLoadMlrrc(
	options *cli.TOptions,
	path string,
	profileName string,
	pFoundProfile *bool,
) (bool, error) {
	handle, err := os.Open(path)
	if err != nil {
		return false, nil
	}
	defer func() { _ = handle.Close() }()

	lineReader := bufio.NewReader(handle)

	// Empty string means we're before any [section] header: the global part
	// of the file which is applied unconditionally.
	currentSectionName := ""

	lineno := 0
	for {
		line, err := lineReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		lineno++

		if err != nil {
			return true, err
		}

		// This is how to do a chomp:
		// TODO: handle \r\n with libified solution.
		line = strings.TrimRight(line, "\n")

		// Comment-strip, then left-trim / right-trim.
		stripped := stripMlrrcLine(line)

		if stripped == "" { // line was whitespace-only, or comment-only
			continue
		}

		if strings.HasPrefix(stripped, "[") {
			sectionName, ok := parseMlrrcSectionHeader(stripped)
			if !ok {
				return true, fmt.Errorf(
					"parse error at file \"%s\" line %d: %s", path, lineno, line,
				)
			}
			currentSectionName = sectionName
			if profileName != "" && sectionName == profileName {
				*pFoundProfile = true
			}
			continue
		}

		// Global (pre-section) lines are always applied. Section lines are
		// applied only if their section is the requested profile; lines in
		// other sections are skipped entirely (not even parsed), so that a
		// typo within an unused profile doesn't break every mlr invocation.
		if currentSectionName != "" && currentSectionName != profileName {
			continue
		}

		if !handleMlrrcLine(options, stripped) {
			return true, fmt.Errorf(
				"parse error at file \"%s\" line %d: %s", path, lineno, line,
			)
		}
	}

	return true, nil
}

// stripMlrrcLine removes any comment (from '#' to end of line) and
// surrounding whitespace from a .mlrrc line.
func stripMlrrcLine(line string) string {
	re := regexp.MustCompile("#.*")
	line = re.ReplaceAllString(line, "")
	return strings.TrimSpace(line)
}

// parseMlrrcSectionHeader parses an INI-style section header like "[name]".
// The input must already be comment-stripped and whitespace-trimmed, and
// start with '['. Whitespace within the brackets is allowed: "[ name ]" is
// the same as "[name]". Returns the section name, and false if the line is
// not a well-formed section header.
func parseMlrrcSectionHeader(line string) (string, bool) {
	if !strings.HasSuffix(line, "]") {
		return "", false
	}
	sectionName := strings.TrimSpace(line[1 : len(line)-1])
	if sectionName == "" || strings.ContainsAny(sectionName, "[]") {
		return "", false
	}
	return sectionName, true
}

// handleMlrrcLine is a helper function for tryLoadMlrrc, handling a single
// (comment-stripped, whitespace-trimmed, non-empty) settings line.
func handleMlrrcLine(
	options *cli.TOptions,
	line string,
) bool {
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
	} else if args[0] == "--profile" || args[0] == "-P" {
		// Profiles are selected on the mlr command line, not from within a
		// .mlrrc file
		return false
	} else if cli.FLAG_TABLE.Parse(args, argc, &argi, options) {
		// handled
	} else {
		return false
	}

	return true
}
