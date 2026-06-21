// Package registry is the single source of truth for the names of Miller's
// "terminals" -- the top-level subcommands like `mlr help` and `mlr version`
// -- and for the top-level version flags.
//
// It exists as its own leaf package (importing nothing within Miller) so that
// both pkg/terminals (which builds the dispatch table) and
// pkg/terminals/completion (which offers these as tab-completion candidates)
// can import it without an import cycle. pkg/terminals imports
// pkg/terminals/completion, so completion cannot import pkg/terminals directly.
package registry

// Terminal subcommand names, in display order. pkg/terminals builds its
// dispatch table from these constants.
const (
	TerminalList = "terminal-list"
	Completion   = "completion"
	Help         = "help"
	Regtest      = "regtest"
	Repl         = "repl"
	Script       = "script"
	Version      = "version"
)

// Names is the ordered list of all terminal subcommand names.
var Names = []string{
	TerminalList,
	Completion,
	Help,
	Regtest,
	Repl,
	Script,
	Version,
}

// Top-level version flags, handled in pkg/climain before normal command-line
// parsing.
const (
	VersionFlag     = "--version"
	BareVersionFlag = "--bare-version"
)

// VersionFlagNames is the list of top-level version flags.
var VersionFlagNames = []string{VersionFlag, BareVersionFlag}
