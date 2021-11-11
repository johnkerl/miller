// ================================================================
// Miller support for command-line flags.
//
// * Flags are used for several purposes:
//
//   o Command-line parsing the main mlr program.
//
//   o Record-reader and record-writer options for a few verbs such as join and
//     tee. E.g. `mlr --csv join -f foo.tsv --tsv ...`: the main input files are
//     CSV but the join-in file is TSV>
//
//   o Processing .mlrrc files.
//
//   o Autogenerating on-line help for `mlr help flags`.
//
//   o Autogenerating the manpage for `man mlr`.
//
//   o Autogenerating webdocs (mkdocs).
//
// * For these reasons, flags are organized into tables; for documentation
//   purposes, flags are organized into sections (see internal/pkg/cli/option_parse.go).
//
// * The Flag struct separates out flag name (e.g. `--csv`), any alternate
//   names (e.g. `-c`), any arguments the flag may take, a help string, and a
//   command-line parser function.
//
// * The tabular structure may seem overwrought; in fact it has been a blessing
//   to develop the tabular structure since these flags objects need to serve
//   so many roles as listed above.
//
// * I don't use Go flags for a few reasons. The most important one is that I
//   need to handle repeated flags, e.g. --from can be used more than once for
//   mlr, and -f/-n/-r etc can be used more than once for mlr sort, etc. I also
//   insist on total control of flag formatting including alphabetization of
//   flags for on-line help and documentation systems.
// ================================================================

package cli

import (
	"fmt"
	"sort"
	"strings"

	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------
// Data types used within the flags table.

// FlagParser is a function which takes a flag such as `--foo`.
// * It should assume that a flag.Owns method has already been invoked to be
//   sure that this function is indeed the right one to call for `--foo`.
// * The FlagParser function is responsible for advancing *pargi by 1 (if
//   `--foo`) or 2 (if `--foo bar`), checking to see if argc is long enough in
//   the latter case, and mutating the options struct.
// * Successful handling of the flag is indicated by this function making a
//   non-zero increment of *pargi.
type FlagParser func(
	args []string,
	argc int,
	pargi *int,
	options *TOptions,
)

// ----------------------------------------------------------------

// FlagTable holds all the flags for Miller, organized into sections.
type FlagTable struct {
	sections []*FlagSection
}

// FlagSection holds all the flags in a given cateogory, where these
// categories exist for documentation purposes.
//
// The name should be right-cased for webdocs. For on-line help and
// manpage use, it will get fully uppercased.
//
// The infoPrinter provides summary/overview for all flags in the
// section, for on-line help / webdocs.
type FlagSection struct {
	name        string
	infoPrinter func()
	flags       []Flag
}

// Flag is a container for all runtime as well as documentation information for
// a flag.
type Flag struct {
	// In most cases, the flag has just one spelling, like "--ifs".
	name string

	// In some cases, the flag has more than one spelling, like "-h" and
	// "--help", or "-c" and "--csv". The altNames field can be omitted from
	// struct initializers, which in Go means it will read as nil.
	altNames []string

	// If not "", a name for the flag's argument, for on-line help. E.g. the
	// "bar" in ""--foo {bar}". It should always be written in curly braces.
	arg string

	// Help string for `mlr help flags`, `man mlr`, and webdocs.
	// * It should be all one line within the source code. The text will be
	//   reformatted as a paragraph for on-line help / manpage, so there should
	//   be no attempt at line-breaking within the help string.
	// * Any code bits should be marked with backticks. These look OK for
	//   on-line help / manpage, and render marvelously for webdocs which
	//   take markdown.
	// * After changing flags you can run `make precommit` in the Miller
	//   repo base directory followed by `git diff` to see how the output
	//   looks. See also the README.md files in the docs and man directories
	//   for how to look at the autogenned docs pre-commit.
	help string

	// A function for parsing the command line, as described above.
	parser FlagParser

	// For format-conversion keystroke-savers, a matrix is plenty -- we don't
	// need to print a tedious 60-line list.
	suppressFlagEnumeration bool
}

// ================================================================
// FlagTable methods

// Sort organizes the sections in the table alphabetically, to make on-line
// help easier to read. This is done from func-init context so on-line help
// will always be easy to navigate.
func (ft *FlagTable) Sort() {
	// Go sort API: for ascending sort, return true if element i < element j.
	sort.Slice(ft.sections, func(i, j int) bool {
		return strings.ToLower(ft.sections[i].name) < strings.ToLower(ft.sections[j].name)
	})
}

// Parse is for parsing a flag on the command line. Given say `--foo`, if a
// Flag object is found which owns the flag, and if its parser accepts it (e.g.
// `bar` is present and spelt correctly if the flag-parser expects `--foo bar`)
// then the return value is true, else false.
func (ft *FlagTable) Parse(
	args []string,
	argc int,
	pargi *int,
	options *TOptions,
) bool {
	for _, section := range ft.sections {
		for _, flag := range section.flags {
			if flag.Owns(args[*pargi]) {
				// Let the flag-parser advance *pargi, depending on how many
				// arguments follow the flag. E.g. `--ifs pipe` will advance
				// *pargi by 2; `-I` will advance it by 1.
				oargi := *pargi
				flag.parser(args, argc, pargi, options)
				nargi := *pargi
				return nargi > oargi
			}
		}
	}
	return false
}

// ShowHelp prints all-in-one on-line help, nominally for `mlr help flags`.
func (ft *FlagTable) ShowHelp() {
	for i, section := range ft.sections {
		if i > 0 {
			fmt.Println()
		}
		fmt.Println(colorizer.MaybeColorizeHelp(strings.ToUpper(section.name), true))
		fmt.Println()
		section.PrintInfo()
		section.ShowHelpForFlags()
	}
}

// ListFlagSections exposes some of the flags-table structure, so Ruby autogen
// scripts for on-line help and webdocs can traverse the structure with looping
// inside their own code.
func (ft *FlagTable) ListFlagSections() {
	for _, section := range ft.sections {
		fmt.Println(section.name)
	}
}

// PrintInfoForSection exposes some of the flags-table structure, so Ruby
// autogen scripts for on-line help and webdocs can traverse the structure with
// looping inside their own code.
func (ft *FlagTable) ShowHelpForSection(sectionName string) bool {
	for _, section := range ft.sections {
		if sectionName == section.name {
			section.PrintInfo()
			section.ShowHelpForFlags()
			return true
		}
	}
	return false
}

// Sections are named like "CSV-only flags". `mlr help` uses `mlr help
// csv-only-flags`. The latter is downcased from the former, with spaces
// replaced by dashes -- hence "downdashed section name". Here we look up
// flag-section help given a downdashed section name.
func (ft *FlagTable) ShowHelpForSectionViaDowndash(downdashSectionName string) bool {
	for _, section := range ft.sections {
		if downdashSectionName == section.GetDowndashSectionName() {
			fmt.Println(colorizer.MaybeColorizeHelp(strings.ToUpper(section.name), true))
			section.PrintInfo()
			section.ShowHelpForFlags()
			return true
		}
	}
	return false
}

// PrintInfoForSection exposes some of the flags-table structure, so Ruby
// autogen scripts for on-line help and webdocs can traverse the structure with
// looping inside their own code.
func (ft *FlagTable) PrintInfoForSection(sectionName string) bool {
	for _, section := range ft.sections {
		if sectionName == section.name {
			section.PrintInfo()
			return true
		}
	}
	return false
}

// ListFlagsForSection exposes some of the flags-table structure, so Ruby
// autogen scripts for on-line help and webdocs can traverse the structure with
// looping inside their own code.
func (ft *FlagTable) ListFlagsForSection(sectionName string) bool {
	for _, section := range ft.sections {
		if sectionName == section.name {
			section.ListFlags()
			return true
		}
	}
	return false
}

// Given flag named `--foo`, altName `-f`, and argument spec `{bar}`, the
// headline is `--foo or -f {bar}`. This is the bit which is highlighted in
// on-line help; its length is also used for alignment decisions in the on-line
// help and the manapge.
func (ft *FlagTable) ShowHeadlineForFlag(flagName string) bool {
	for _, fs := range ft.sections {
		for _, flag := range fs.flags {
			if flag.Owns(flagName) {
				fmt.Println(flag.GetHeadline())
				return true
			}
		}
	}
	return false
}

// ShowHelpForFlag prints the flag's help-string all on one line.  This is for
// webdoc usage where the browser does dynamic line-wrapping, as the user
// resizes the browser window.
func (ft *FlagTable) ShowHelpForFlag(flagName string) bool {
	for _, fs := range ft.sections {
		for _, flag := range fs.flags {
			if flag.Owns(flagName) {
				fmt.Println(flag.GetHelpOneLine())
				return true
			}
		}
	}
	return false
}

// Map "CSV-only flags" to "csv-only-flags" etc. for the benefit of per-section
// help in `mlr help topics`.
func (ft *FlagTable) GetDowndashSectionNames() []string {
	downdashSectionNames := make([]string, len(ft.sections))
	for i, fs := range ft.sections {
		// Get names like "CSV-only flags" from the FLAG_TABLE.
		// Downcase and replace spaces with dashes to get names like
		// "csv-only-flags"
		downdashSectionNames[i] = fs.GetDowndashSectionName()
	}
	return downdashSectionNames
}

// NilCheck checks to see if any flag/section is missing help info. This arises
// since in Go you needn't specify all struct initializers, so for example a
// Flag struct-initializer which doesn't say `help: "..."` will have empty help
// string. This nil-checking doesn't need to be done on every Miller
// invocation, but rather, only at build time. The `mlr help` auxent has an
// entry point wherein a regression-test case can do `mlr help nil-check` and
// make this function exits cleanly.
func (ft *FlagTable) NilCheck() {
	lib.InternalCodingErrorWithMessageIf(ft.sections == nil, "Nil table sections")
	lib.InternalCodingErrorWithMessageIf(len(ft.sections) == 0, "Zero table sections")
	for _, fs := range ft.sections {
		fs.NilCheck()
	}
	fmt.Println("Flag-table nil check completed successfully.")
}

// ================================================================
// FlagSection methods

// Sort organizes the flags in the section alphabetically, to make on-line help
// easier to read.  This is done from func-init context so on-line help will
// always be easy to navigate.
func (fs *FlagSection) Sort() {
	// Go sort API: for ascending sort, return true if element i < element j.
	sort.Slice(fs.flags, func(i, j int) bool {
		return strings.ToLower(fs.flags[i].name) < strings.ToLower(fs.flags[j].name)
	})
}

// ShowHelpForFlags prints all-in-one on-line help, nominally for `mlr help
// flags`.
func (fs *FlagSection) ShowHelpForFlags() {
	for _, flag := range fs.flags {
		// For format-conversion keystroke-savers, a matrix is plenty -- we don't
		// need to print a tedious 60-line list.
		if flag.suppressFlagEnumeration {
			continue
		}
		flag.ShowHelp()
	}
}

// PrintInfo exposes some of the flags-table structure, so Ruby autogen scripts
// for on-line help and webdocs can traverse the structure with looping inside
// their own code.
func (fs *FlagSection) PrintInfo() {
	fs.infoPrinter()
	fmt.Println()
}

// ListFlags exposes some of the flags-table structure, so Ruby autogen scripts
// for on-line help and webdocs can traverse the structure with looping inside
// their own code.
func (fs *FlagSection) ListFlags() {
	for _, flag := range fs.flags {
		fmt.Println(flag.name)
	}
}

// Map "CSV-only flags" to "csv-only-flags" etc. for the benefit of per-section
// help in `mlr help topics`.
func (fs *FlagSection) GetDowndashSectionName() string {
	return strings.ReplaceAll(strings.ToLower(fs.name), " ", "-")
}

// See comments above FlagTable's NilCheck method.
func (fs *FlagSection) NilCheck() {
	lib.InternalCodingErrorWithMessageIf(fs.name == "", "Empty section name")
	lib.InternalCodingErrorWithMessageIf(fs.infoPrinter == nil, "Nil infoPrinter for section "+fs.name)
	lib.InternalCodingErrorWithMessageIf(fs.flags == nil, "Nil flags for section "+fs.name)
	lib.InternalCodingErrorWithMessageIf(len(fs.flags) == 0, "Zero flags for section "+fs.name)
	for _, flag := range fs.flags {
		flag.NilCheck()
	}
}

// ================================================================
// Flag methods

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

// ShowHelp produces formatting for `mlr help flags` and manpage use.
// Example:
// * Flag name is `--foo`
// * altName is `-f`
// * Argument spec is `{bar}`
// * Help string is "Lorem ipsum dolor sit amet, consectetur adipiscing elit,
//   sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
//   ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip
//   ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate
//   velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat
//   cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id
//   est laborum."
// * The headline (see the GetHeadline function) is `--foo or -f {bar}`.
// * We place the headline left in a 25-character column, colorized with the
//   help color.
// * We format the help text as 55-character lines and place them
//   to the right.
// * The result looks like
//
//   --foo or -f {bar}        Lorem ipsum dolor sit amet, consectetur adipiscing
//                            elit, sed do eiusmod tempor incididunt ut labore et
//                            dolore magna aliqua. Ut enim ad minim veniam, quis
//                            nostrud exercitation ullamco laboris nisi ut aliquip
//                            ex ea commodo consequat. Duis aute irure dolor in
//                            reprehenderit in voluptate velit esse cillum dolore
//                            eu fugiat nulla pariatur. Excepteur sint occaecat
//                            cupidatat non proident, sunt in culpa qui officia
//                            deserunt mollit anim id est laborum.
//
// * If the headline is too long we put the first help line a line below like this:
//
//   --foo-flag-is-very-very-long {bar}
//                            Lorem ipsum dolor sit amet, consectetur adipiscing
//                            elit, sed do eiusmod tempor incididunt ut labore et
//                            dolore magna aliqua. Ut enim ad minim veniam, quis
//                            nostrud exercitation ullamco laboris nisi ut aliquip
//                            ex ea commodo consequat. Duis aute irure dolor in
//                            reprehenderit in voluptate velit esse cillum dolore
//                            eu fugiat nulla pariatur. Excepteur sint occaecat
//                            cupidatat non proident, sunt in culpa qui officia
//                            deserunt mollit anim id est laborum.
//

func (flag *Flag) ShowHelp() {
	headline := flag.GetHeadline()
	displayHeadline := fmt.Sprintf("%-25s", headline)
	broken := len(headline) >= 25

	helpLines := lib.FormatAsParagraph(flag.help, 55)

	if broken {
		fmt.Printf("%s\n", colorizer.MaybeColorizeHelp(displayHeadline, true))
		for _, helpLine := range helpLines {
			fmt.Printf("%25s%s\n", " ", helpLine)
		}
	} else {
		fmt.Printf("%s", colorizer.MaybeColorizeHelp(displayHeadline, true))
		if len(helpLines) == 0 {
			fmt.Println()
		}
		for i, helpLine := range helpLines {
			if i == 0 {
				fmt.Printf("%s\n", helpLine)
			} else {
				fmt.Printf("%25s%s\n", " ", helpLine)
			}
		}
	}
}

// GetHeadline puts together the flag name, any altNames, and any argument spec
// into a single string for the left column of online help / manpage content.
// Given flag named `--foo`, altName `-f`, and argument spec `{bar}`, the
// headline is `--foo or -f {bar}`. This is the bit which is highlighted in
// on-line help; its length is also used for alignment decisions in the on-line
// help and the manapge.
func (flag *Flag) GetHeadline() string {
	displayNames := make([]string, 1)
	displayNames[0] = flag.name
	if flag.altNames != nil {
		displayNames = append(displayNames, flag.altNames...)
	}
	displayText := strings.Join(displayNames, " or ")
	if flag.arg != "" {
		displayText += " "
		displayText += flag.arg
	}
	return displayText
}

// Gets the help string all on one line (just in case anyone typed it in using
// multiline string-literal backtick notation in Go). This is suitable for
// webdoc use where we create all one line, and the browser dynamically
// line-wraps as the user resizes the window.
func (flag *Flag) GetHelpOneLine() string {
	return strings.Join(strings.Split(flag.help, "\n"), " ")
}

// See comments above FlagTable's NilCheck method.
func (flag *Flag) NilCheck() {
	lib.InternalCodingErrorWithMessageIf(flag.name == "", "Empty flag name")
	lib.InternalCodingErrorWithMessageIf(flag.help == "", "Empty flag help for flag "+flag.name)
	lib.InternalCodingErrorWithMessageIf(flag.parser == nil, "Nil parser help for flag "+flag.name)
}

// ================================================================
// Helper methods

// NoOpParse1 is a helper function for flags which take no argument and are
// backward-compatibility no-ops.
func NoOpParse1(args []string, argc int, pargi *int, options *TOptions) {
	*pargi += 1
}
