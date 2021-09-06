// TODO: comment
// TODO: note complexity b/c serving many uses: main CLI, .mlrrc, some verbs; OLH/man/docs autogen

package cliutil

import (
	"fmt"
	"sort"
	"strings"

	"mlr/src/colorizer"
)

// ----------------------------------------------------------------
// Callsites:
// * src/cli/mlrcli_parse.go
//   ParseCommandLine
//     MainOptions (--cpuprofile, --version, etc)
//     ParseReaderOptions
//     ParseWriterOptions
//     ParseReaderWriterOptions
//     ParseMiscOptions
//     help.ParseTerminalUsage
// * handleMlrrcLine
// * nest/tee/join/put/filter:
//     ParseReaderOptions
//     ParseWriterOptions
//   !! must use only cliutil package, not cli package

// sections:
//   how to factor
//   reader/writer/readerwriter/misc
//   -> split necessary for verbs & what they do/don't accept
//   vs
//   data-format options, --x2y, separators, compressed, comments in data,
//   csv-specific, number-formatting, other
//   -> split useful for on-line help

// ----------------------------------------------------------------
// TODO: comment
type FlagParser func(
	args []string,
	argc int,
	pargi *int,
	options *TOptions,
)

type Flag struct {
	// More common case: the flag has just one spelling, like "--ifs".
	name string

	// Less common case: the flag has more than one spelling, like "-h" and "--help",
	// or "-c" and "--csv".
	altNames []string

	help string

	parser FlagParser
	// TODO: comment
	// reader, writer, reader/writer, misc = neither
	forReader bool
	forWriter bool
}

type FlagSection struct {
	name string // TODO: lowercase? capcase? upper? make methods?
	// xxx common-info func
	flags []Flag
}

type FlagTable struct {
	sections []*FlagSection
}

// ----------------------------------------------------------------
// NoOpParse1 is a helper function for flags which take no argument and are
// backward-compatibility no-ops.
func NoOpParse1(args []string, argc int, pargi *int, options *TOptions) {
	*pargi += 1
}

var NoOpHelp string = "No-op pass-through for backward compatibility with Miller 5."

// ----------------------------------------------------------------
// Owns determines whether this object handles a command-line flag such as "--foo".
func (flag *Flag) Owns(input string) bool {
	if input == flag.name {
		return true
	}

	if flag.altNames != nil {
		for _, name := range flag.altNames {
			if input == name {
				return true
			}
		}
	}
	return false
}

// ----------------------------------------------------------------
// Sort organizes the flags in the section alphabetically, to make on-line help
// easier to read.

func (fs *FlagSection) Sort() {
	// Go sort API: for ascending sort, return true if element i < element j.
	sort.Slice(fs.flags, func(i, j int) bool {
		return strings.ToLower(fs.flags[i].name) < strings.ToLower(fs.flags[j].name)
	})
}

// ----------------------------------------------------------------
// Sort organizes the sections in the table alphabetically, to make on-line help
// easier to read.

func (ft *FlagTable) Sort() {
	// Go sort API: for ascending sort, return true if element i < element j.
	sort.Slice(ft.sections, func(i, j int) bool {
		return ft.sections[i].name < ft.sections[j].name
	})
}

// ----------------------------------------------------------------
func (ft *FlagTable) Parse(
	args []string,
	argc int,
	pargi *int,
	options *TOptions,
	// TODO forReader, forWriter
) bool {
	for _, section := range ft.sections {
		for _, flag := range section.flags {
			if flag.Owns(args[*pargi]) {
				// Let the flag-parser advance *pargi, depending on how many
				// arguments follow the flag. E.g. `--ifs pipe` will advance
				// *pargi by 2; `-I` will advance it by 1.
				flag.parser(args, argc, pargi, options)
				return true
			}
		}
	}
	return false
}

// ----------------------------------------------------------------
// TODO: more options for OLH
func (ft *FlagTable) ListTemp() {
	for i, section := range ft.sections {
		// TODO: colorize
		if i > 0 {
			fmt.Println()
		}
		fmt.Println(colorizer.MaybeColorizeHelp(strings.ToUpper(section.name), true))
		fmt.Println()
		section.ListTemp()
	}
}

// TODO: more options for OLH
func (fs *FlagSection) ListTemp() {
	for _, flag := range fs.flags {
		// TODO: colorize
		//if i > 0 {
		//fmt.Println()
		//}
		flag.ListTemp()
	}
}

func (flag *Flag) ListTemp() {
	// TODO: make method?
	displayNames := make([]string, 1)
	displayNames[0] = flag.name
	if flag.altNames != nil {
		displayNames = append(displayNames, flag.altNames...)
	}
	displayText := strings.Join(displayNames, " or ")
	// TODO: abend if flag.help == ""
	//displayText = fmt.Sprintf("%32s", displayText)
	//fmt.Printf("  %s: %s\n", colorizer.MaybeColorizeHelp(displayText, true), flag.help)
	displayText = fmt.Sprintf("%-32s", displayText)
	fmt.Printf("  %s %s\n", colorizer.MaybeColorizeHelp(displayText, true), flag.help)
}
