// Value completion for arg-taking main flags whose argument is a known,
// enumerable set -- e.g. file formats for -i/-o/--io and named separator
// aliases for --ifs and friends. Flags not listed here fall back to filename
// completion.
//
// The candidate sets come from cli (GetFileFormatNames, GetSeparatorAliasNames,
// GetSeparatorRegexAliasNames), which read the same maps Miller uses at
// runtime, so there is no separate list to keep in sync.

package completion

import (
	"github.com/johnkerl/miller/v6/pkg/cli"
)

// formatFlagNames are the main flags that take a file-format name.
var formatFlagNames = map[string]bool{
	"-i":   true,
	"-o":   true,
	"--io": true,
}

// separatorFlagNames are the main flags that take a (possibly aliased)
// separator string.
var separatorFlagNames = map[string]bool{
	"--fs":  true,
	"--ifs": true,
	"--ofs": true,
	"--ps":  true,
	"--ips": true,
	"--ops": true,
	"--rs":  true,
	"--irs": true,
	"--ors": true,
}

// separatorRegexFlagNames are the main flags that take a regex-separator,
// which has its own alias set.
var separatorRegexFlagNames = map[string]bool{
	"--ifs-regex": true,
	"--ips-regex": true,
}

// flagValueCandidates returns the enumerated values a flag's argument may take,
// or nil if the flag's value is not an enumerable set (in which case the caller
// falls back to filename completion).
func flagValueCandidates(flag string) []string {
	switch {
	case formatFlagNames[flag]:
		return cli.GetFileFormatNames()
	case separatorFlagNames[flag]:
		return cli.GetSeparatorAliasNames()
	case separatorRegexFlagNames[flag]:
		return cli.GetSeparatorRegexAliasNames()
	}
	return nil
}
