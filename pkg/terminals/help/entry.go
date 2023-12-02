// ================================================================
// Online help
// ================================================================

package help

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"

	"github.com/johnkerl/miller/pkg/auxents"
	"github.com/johnkerl/miller/pkg/bifs"
	"github.com/johnkerl/miller/pkg/cli"
	"github.com/johnkerl/miller/pkg/dsl/cst"
	"github.com/johnkerl/miller/pkg/lib"
	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/runtime"
	"github.com/johnkerl/miller/pkg/transformers"
)

// ================================================================
type tZaryHandlerFunc func()
type tVarArgHandlerFunc func(args []string)

type tHandlerLookupTable struct {
	sections []tHandlerInfoSection
}

type tHandlerInfoSection struct {
	name         string
	handlerInfos []tHandlerInfo

	// Some handlers are used only for webdoc/manpage autogen and needn't
	// clutter up the on-line help experience for the interactive user
	internal bool
}

type tHandlerInfo struct {
	name              string
	zaryHandlerFunc   tZaryHandlerFunc
	varArgHandlerFunc tVarArgHandlerFunc
}

type tShorthandTable struct {
	shorthandInfos []tShorthandInfo
}

type tShorthandInfo struct {
	shorthand string
	longhand  string
}

// We get a Golang "initialization loop" if this is defined statically. So, we
// use a "package init" function.
var handlerLookupTable = tHandlerLookupTable{}
var shorthandLookupTable = tShorthandTable{}

func init() {
	// For things like 'mlr help foo', invoked through the terminals framework
	// which goes through our HelpMain().
	handlerLookupTable = tHandlerLookupTable{
		sections: []tHandlerInfoSection{
			{
				name: "Essentials",
				handlerInfos: []tHandlerInfo{
					{name: "topics", zaryHandlerFunc: listTopics},
					{name: "basic-examples", zaryHandlerFunc: helpBasicExamples},
					{name: "file-formats", zaryHandlerFunc: helpFileFormats},
				},
			},
			{
				name: "Flags",
				handlerInfos: []tHandlerInfo{
					{name: "flags", zaryHandlerFunc: showFlagsHelp},
					{name: "flag", varArgHandlerFunc: helpForFlag},
					{name: "list-separator-aliases", zaryHandlerFunc: listSeparatorAliases},
					{name: "list-separator-regex-aliases", zaryHandlerFunc: listSeparatorRegexAliases},
					// Per-section entries will be computed and installed below
				},
			},
			{
				name: "Verbs",
				handlerInfos: []tHandlerInfo{
					{name: "list-verbs", zaryHandlerFunc: listVerbs},
					{name: "usage-verbs", zaryHandlerFunc: usageVerbs},
					{name: "verb", varArgHandlerFunc: helpForVerb},
				},
			},
			{
				name: "Functions",
				handlerInfos: []tHandlerInfo{
					{name: "list-functions", zaryHandlerFunc: listFunctions},
					{name: "list-function-classes", zaryHandlerFunc: listFunctionClasses},
					{name: "list-functions-in-class", varArgHandlerFunc: listFunctionsInClass},
					{name: "usage-functions", zaryHandlerFunc: usageFunctions},
					{name: "usage-functions-by-class", zaryHandlerFunc: usageFunctionsByClass},
					{name: "function", varArgHandlerFunc: helpForFunction},
				},
			},
			{
				name: "Keywords",
				handlerInfos: []tHandlerInfo{
					{name: "list-keywords", zaryHandlerFunc: listKeywords},
					{name: "usage-keywords", zaryHandlerFunc: usageKeywords},
					{name: "keyword", varArgHandlerFunc: helpForKeyword},
				},
			},
			{
				name: "Other",
				handlerInfos: []tHandlerInfo{
					{name: "auxents", zaryHandlerFunc: helpAuxents},
					{name: "terminals", zaryHandlerFunc: helpTerminals},
					{name: "mlrrc", zaryHandlerFunc: helpMlrrc},
					{name: "output-colorization", zaryHandlerFunc: helpOutputColorization},
					{name: "type-arithmetic-info", zaryHandlerFunc: helpTypeArithmeticInfo},
					{name: "type-arithmetic-info-extended", zaryHandlerFunc: helpTypeArithmeticInfoExtended},
				},
			},
			{
				name:     "Internal/docgen",
				internal: true,
				handlerInfos: []tHandlerInfo{
					{name: "flag-table-nil-check", zaryHandlerFunc: flagTableNilCheck},
					{name: "list-flag-sections", zaryHandlerFunc: listFlagSections},
					{name: "list-flags-for-section", varArgHandlerFunc: listFlagsForSection},
					{name: "list-functions-as-paragraph", zaryHandlerFunc: listFunctionsAsParagraph},
					{name: "list-functions-as-table", zaryHandlerFunc: listFunctionsAsTable},
					{name: "list-keywords-as-paragraph", zaryHandlerFunc: listKeywordsAsParagraph},
					{name: "list-verbs-as-paragraph", zaryHandlerFunc: listVerbsAsParagraph},
					{name: "print-info-for-section", varArgHandlerFunc: printInfoForSection},
					{name: "show-headline-for-flag", varArgHandlerFunc: showHeadlineForFlag},
					{name: "show-help-for-flag", varArgHandlerFunc: showHelpForFlag},
					{name: "show-help-for-section", varArgHandlerFunc: showHelpForSection},
					{name: "show-help-for-section-via-downdash", varArgHandlerFunc: showHelpForSectionViaDowndash},
				},
			},
		},
	}

	// This is a wee bit clever. The rest of the topics in the table have names
	// manually keyed in. But we want to produce `mlr help csv-only-flags` for
	// flag-section named "CSV-only flags", etc. Here we can't key in the names
	// since we want to compute them dynamically from cli.FLAG_TABLE which is
	// Miller's wqy of tracking command-line flags.

	// For this file's topic-lookup table, find and extend the section called "Flags".
	for i, section := range handlerLookupTable.sections {
		if section.name != "Flags" {
			continue
		}

		// Ask the flags table for a list of flag-section names, downcased and
		// with spaces replaced with dashes -- "downdashed" -- making the
		// punctuation/casing style for online help.
		downdashSectionNames := cli.FLAG_TABLE.GetDowndashSectionNames()
		// Note: `j, _` rather than `_, downdashSectionName` since the latter
		// is a data copy while the former allows us to do a reference. The
		// former won't produce correct lookup-table data.
		for j := range downdashSectionNames {
			downdashSectionName := downdashSectionNames[j]
			// Patch a new entry into the "Flags" section of our lookup table.
			entry := tHandlerInfo{
				name: downdashSectionName,
				// Make a function which passes in "csv-only-flags" etc. to the FLAG_TABLE.
				zaryHandlerFunc: func() {
					showHelpForSectionViaDowndash([]string{downdashSectionName})
				},
			}
			handlerLookupTable.sections[i].handlerInfos = append(handlerLookupTable.sections[i].handlerInfos, entry)
		}
	}

	// For things like 'mlr -f', invoked through the CLI parser which does not
	// go through our HelpMain().
	shorthandLookupTable = tShorthandTable{
		shorthandInfos: []tShorthandInfo{
			{shorthand: "-g", longhand: "flags"},
			{shorthand: "-l", longhand: "list-verbs"},
			{shorthand: "-L", longhand: "usage-verbs"},
			{shorthand: "-f", longhand: "list-functions"},
			{shorthand: "-F", longhand: "usage-functions"},
			{shorthand: "-k", longhand: "list-keywords"},
			{shorthand: "-K", longhand: "usage-keywords"},
		},
	}
}

// ================================================================
// For things like 'mlr help foo', invoked through the terminals framework which
// goes through our HelpMain().  Here, the args are the terminal part of the full
// Miller command line: if the latter was "mlr --some-flag help foo bar" then
// the former is "help foo bar".
func HelpMain(args []string) int {
	args = args[1:]

	// "mlr help" and nothing else
	if len(args) == 0 {
		handleDefault()
		return 0
	}

	// 'mlr help find x' searches for all things (flags, transformers,
	// functions, keywords) with an "x" in their name, as a substring.
	if args[0] == "find" {
		args = args[1:]
		if len(args) > 0 {
			helpByApproximateSearch(args)
			return 0
		} else {
			fmt.Printf("mlr help find: need one or more things to search for.\n")
			return 1
		}
	}

	// "mlr help something" where we recognize the something
	name := args[0]
	for _, section := range handlerLookupTable.sections {
		for _, info := range section.handlerInfos {
			if info.name == name {
				if info.zaryHandlerFunc != nil {
					if len(args) != 1 {
						fmt.Printf("mlr help %s takes no additional argument.\n", name)
						return 0
					}
					info.zaryHandlerFunc()
					return 0
				}
				if info.varArgHandlerFunc != nil {
					if len(args) < 2 {
						fmt.Printf("mlr help %s takes at least one required argument.\n", name)
						return 0
					}
					info.varArgHandlerFunc(args[1:])
					return 0
				}
			}
		}
	}

	// 'mlr help x' searches for all things (flags, transformers, functions, keywords) named "x".
	if helpByExactSearch(args) {
		return 0
	}

	return 0
}

// ----------------------------------------------------------------
func MainUsage(o *os.File) {
	fmt.Fprintf(o,
		`Usage: mlr [flags] {verb} [verb-dependent options ...] {zero or more file names}

If zero file names are provided, standard input is read, e.g.
  mlr --csv sort -f shape example.csv

Output of one verb may be chained as input to another using "then", e.g.
  mlr --csv stats1 -a min,mean,max -f quantity then sort -f color example.csv

Please see 'mlr help topics' for more information.
`)
	fmt.Fprintf(o, "Please also see %s\n", lib.DOC_URL)
}

// ----------------------------------------------------------------
// For things like 'mlr -F', invoked through the CLI parser which does not
// go through our HelpMain().
func ParseTerminalUsage(arg string) bool {
	if arg == "-h" || arg == "--help" {
		handleDefault()
		return true
	}
	// "mlr -l" is shorthand for "mlr help list-verbs", etc.
	for _, sinfo := range shorthandLookupTable.shorthandInfos {
		if sinfo.shorthand == arg {
			for _, section := range handlerLookupTable.sections {
				for _, info := range section.handlerInfos {
					if info.name == sinfo.longhand {
						info.zaryHandlerFunc()
						return true
					}
				}
			}
		}
	}
	return false
}

// ================================================================
func handleDefault() {
	MainUsage(os.Stdout)
}

// ----------------------------------------------------------------
func listTopics() {
	fmt.Println("Type 'mlr help {topic}' for any of the following:")
	for _, section := range handlerLookupTable.sections {
		if !section.internal {
			fmt.Printf("%s:\n", section.name)
			for _, info := range section.handlerInfos {
				fmt.Printf("  mlr help %s\n", info.name)
			}
		}
	}
	fmt.Println("Shorthands:")
	for _, info := range shorthandLookupTable.shorthandInfos {
		fmt.Printf("  mlr %s = mlr help %s\n", info.shorthand, info.longhand)
	}
	fmt.Printf("Lastly, 'mlr help ...' will search for your exact text '...' using the sources of\n")
	fmt.Printf("'mlr help flag', 'mlr help verb', 'mlr help function', and 'mlr help keyword'.\n")
	fmt.Printf("Use 'mlr help find ...' for approximate (substring) matches, e.g. 'mlr help find map'\n")
	fmt.Printf("for all things with \"map\" in their names.\n")
}

// ----------------------------------------------------------------
func showFlagsHelp() {
	cli.FLAG_TABLE.ShowHelp()
}

func helpForFlag(args []string) {
	for i, arg := range args {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s:\n", arg)
		if !cli.FLAG_TABLE.ShowHelpForFlag(arg) {
			fmt.Println("Not found.")
		}
	}
}

func listSeparatorAliases() {
	cli.ListSeparatorAliasesForOnlineHelp()
}

func listSeparatorRegexAliases() {
	cli.ListSeparatorRegexAliasesForOnlineHelp()
}

func helpAuxents() {
	fmt.Print(`Miller has a few otherwise-standalone executables packaged within it.
They do not participate in any other parts of Miller.
`)
	fmt.Println()
	auxents.ShowAuxEntries(os.Stdout)
}

func helpTerminals() {
	fmt.Println("Terminals include mlr help, the regression-test entry point mlr regtest, and the REPL mlr repl.")
	// We can't invoke the terminal-lister since that would create a cyclic package reference.
}

// ----------------------------------------------------------------
func helpBasicExamples() {
	fmt.Print(
		`mlr --icsv --opprint cat example.csv
mlr --icsv --opprint sort -f shape example.csv
mlr --icsv --opprint sort -f shape -nr index example.csv
mlr --icsv --opprint cut -f flag,shape example.csv
mlr --csv filter '$color == "red"' example.csv
mlr --icsv --ojson put '$ratio = $quantity / $rate' example.csv
mlr --icsv --opprint --from example.csv sort -nr index then cut -f shape,quantity
`)
}

// ----------------------------------------------------------------
func helpFileFormats() {
	fmt.Printf(
		`CSV/CSV-lite: comma-separated values with separate header line
TSV: same but with tabs in places of commas
+---------------------+
| apple,bat,cog       |
| 1,2,3               | Record 1: "apple":"1", "bat":"2", "cog":"3"
| 4,5,6               | Record 2: "apple":"4", "bat":"5", "cog":"6"
+---------------------+

JSON (array of objects):
+---------------------+
| [                   |
| {                   |
|  "apple": 1,        | Record 1: "apple":"1", "bat":"2", "cog":"3"
|  "bat": 2,          |
|  "cog": 3           |
| },                  |
| {                   |
|   "dish": {         | Record 2: "dish.egg":"7",
|     "egg": 7,       | "dish.flint":"8", "garlic":""
|     "flint": 8      |
|   },                |
|   "garlic": ""      |
| }                   |
| ]                   |
+---------------------+

JSON Lines (sequence of one-line objects):
+------------------------------------------------+
| {"apple": 1, "bat": 2, "cog": 3}               |
| {"dish": {"egg": 7, "flint": 8}, "garlic": ""} |
+------------------------------------------------+
  Record 1: "apple":"1", "bat":"2", "cog":"3"
  Record 2: "dish:egg":"7", "dish:flint":"8", "garlic":""

PPRINT: pretty-printed tabular
+---------------------+
| apple bat cog       |
| 1     2   3         | Record 1: "apple:"1", "bat":"2", "cog":"3"
| 4     5   6         | Record 2: "apple":"4", "bat":"5", "cog":"6"
+---------------------+

Markdown tabular (supported for output only):
+-----------------------+
| | apple | bat | cog | |
| | ---   | --- | --- | |
| | 1     | 2   | 3   | | Record 1: "apple:"1", "bat":"2", "cog":"3"
| | 4     | 5   | 6   | | Record 2: "apple":"4", "bat":"5", "cog":"6"
+-----------------------+

XTAB: pretty-printed transposed tabular
+---------------------+
| apple 1             | Record 1: "apple":"1", "bat":"2", "cog":"3"
| bat   2             |
| cog   3             |
|                     |
| dish 7              | Record 2: "dish":"7", "egg":"8"
| egg  8              |
+---------------------+

DKVP: delimited key-value pairs (Miller default format)
+---------------------+
| apple=1,bat=2,cog=3 | Record 1: "apple":"1", "bat":"2", "cog":"3"
| dish=7,egg=8,flint  | Record 2: "dish":"7", "egg":"8", "3":"flint"
+---------------------+

NIDX: implicitly numerically indexed (Unix-toolkit style)
+---------------------+
| the quick brown     | Record 1: "1":"the", "2":"quick", "3":"brown"
| fox jumped          | Record 2: "1":"fox", "2":"jumped"
+---------------------+
`)
}

// ----------------------------------------------------------------
func helpMlrrc() {
	fmt.Print(
		`You can set up personal defaults via a $HOME/.mlrrc and/or ./.mlrrc.
For example, if you usually process CSV, then you can put "--csv" in your .mlrrc file
and that will be the default input/output format unless otherwise specified on the command line.

The .mlrrc file format is one "--flag" or "--option value" per line, with the leading "--" optional.
Hash-style comments and blank lines are ignored.

Sample .mlrrc:
# Input and output formats are CSV by default (unless otherwise specified
# on the mlr command line):
csv
# These are no-ops for CSV, but when I do use JSON output, I want these
# pretty-printing options to be used:
jvstack
jlistwrap

How to specify location of .mlrrc:
* If $MLRRC is set:
  o If its value is "__none__" then no .mlrrc files are processed.
  o Otherwise, its value (as a filename) is loaded and processed. If there are syntax
    errors, they abort mlr with a usage message (as if you had mistyped something on the
    command line). If the file can't be loaded at all, though, it is silently skipped.
  o Any .mlrrc in your home directory or current directory is ignored whenever $MLRRC is
    set in the environment.
* Otherwise:
  o If $HOME/.mlrrc exists, it's then processed as above.
  o If ./.mlrrc exists, it's then also processed as above.
  (I.e. current-directory .mlrrc defaults are stacked over home-directory .mlrrc defaults.)
* The command-line flag "--norc" can be used to suppress loading the .mlrrc file even when other
  conditions are met.
See also:
https://miller.readthedocs.io/en/latest/customization.html
`)
}

// ----------------------------------------------------------------
func helpOutputColorization() {
	cli.OutputColorizationPrintInfo()
}

// ----------------------------------------------------------------
func helpTypeArithmeticInfo() {
	helpTypeArithmeticInfoAux(false)
}

func helpTypeArithmeticInfoExtended() {
	helpTypeArithmeticInfoAux(true)
}

func helpTypeArithmeticInfoAux(extended bool) {
	mlrvals := []*mlrval.Mlrval{
		mlrval.FromInt(1),
		mlrval.FromFloat(2.5),
		mlrval.FromBool(true),
		mlrval.VOID,
		mlrval.ABSENT,
		mlrval.FromAnonymousError(),
	}

	n := len(mlrvals)

	for i := -2; i < n; i++ {
		if i == -2 {
			fmt.Printf("%-10s |", "(+)")
		} else if i == -1 {
			fmt.Printf("%-10s +", "------")
		} else if mlrvals[i].IsVoid() {
			fmt.Printf("%-10s |", "(empty)")
		} else {
			fmt.Printf("%-10s |", mlrvals[i].String())
		}
		for j := 0; j < n; j++ {
			if i == -2 {
				if mlrvals[j].IsVoid() {
					fmt.Printf("%-10s", "(empty)")
				} else {
					fmt.Printf(" %-10s", mlrvals[j].String())
				}
			} else if i == -1 {
				fmt.Printf(" %-10s", "------")
			} else {
				sum := bifs.BIF_plus_binary(mlrvals[i], mlrvals[j])
				if sum.IsVoid() {
					fmt.Printf(" %-10s", "(empty)")
				} else {
					fmt.Printf(" %-10s", sum.String())
				}
			}
		}
		fmt.Println()
	}

	if !extended {
		return
	}

	mlrvals = []*mlrval.Mlrval{
		mlrval.FromBool(true),
		mlrval.FromBool(false),
		mlrval.FromInt(3),
		mlrval.VOID,
		mlrval.ABSENT,
		mlrval.FromAnonymousError(),
	}

	n = len(mlrvals)

	state := runtime.NewEmptyState(cli.DefaultOptions(), false)

	descs := []string{"(&&)", "(||)"}
	for k, desc := range descs {

		fmt.Println()
		for i := -2; i < n; i++ {
			if i == -2 {
				fmt.Printf("%-10s |", desc)
			} else if i == -1 {
				fmt.Printf("%-10s +", "------")
			} else if mlrvals[i].IsVoid() {
				fmt.Printf("%-10s |", "(empty)")
			} else {
				fmt.Printf("%-10s |", mlrvals[i].String())
			}
			for j := 0; j < n; j++ {
				if i == -2 {
					if mlrvals[j].IsVoid() {
						fmt.Printf("%-10s", "(empty)")
					} else {
						fmt.Printf(" %-10s", mlrvals[j].String())
					}
				} else if i == -1 {
					fmt.Printf(" %-10s", "------")
				} else {

					inode := cst.BuildMlrvalLiteralNode(mlrvals[i])
					jnode := cst.BuildMlrvalLiteralNode(mlrvals[j])

					var binary_node cst.IEvaluable
					if k == 0 {
						binary_node = cst.BuildLogicalANDOperatorNode(inode, jnode)
					} else {
						binary_node = cst.BuildLogicalOROperatorNode(inode, jnode)
					}

					output := binary_node.Evaluate(state)

					if output.IsVoid() {
						fmt.Printf(" %-10s", "(empty)")
					} else {
						fmt.Printf(" %-10s", output.String())
					}
				}
			}
			fmt.Println()
		}
	}
}

// ----------------------------------------------------------------
// listFlagSections et al. are for webdoc/manpage autogen in the miller/docs
// and miller/man subdirectories. Unlike showFlagsHelp where all looping over
// the flags table, its sections, and flags within each section is done within
// this Go program, by contrast the following few methods expose the hierarchy
// to standard output, letting the calling programs (nominally Ruby autogen
// scripts) control their own looping and formatting.

func listFlagSections() {
	cli.FLAG_TABLE.ListFlagSections()
}

func printInfoForSection(sectionNames []string) {
	for _, sectionName := range sectionNames {
		if !cli.FLAG_TABLE.PrintInfoForSection(sectionName) {
			fmt.Printf(
				"mlr: flag-section \"%s\" not found. Please use \"mlr help list-flag-sections\" for a list.\n",
				sectionName)
		}
	}
}

func listFlagsForSection(sectionNames []string) {
	for _, sectionName := range sectionNames {
		if !cli.FLAG_TABLE.ListFlagsForSection(sectionName) {
			fmt.Printf(
				"mlr: flag-section \"%s\" not found. Please use \"mlr help list-flag-sections\" for a list.\n",
				sectionName)
		}
	}
}

// For manpage autogen: just produce text
func showHelpForSection(sectionNames []string) {
	for _, sectionName := range sectionNames {
		if !cli.FLAG_TABLE.ShowHelpForSection(sectionName) {
			fmt.Printf(
				"mlr: flag-section \"%s\" not found. Please use \"mlr help list-flag-sections\" for a list.\n",
				sectionName)
		}
	}
}

// For on-the-fly `mlr help foo-bar-flags` where `Foo-bar flags` is the name of
// a section in the FLAG_TABLE. See the func-init block at the top of this
// file.
func showHelpForSectionViaDowndash(downdashSectionNames []string) {
	for _, downdashSectionName := range downdashSectionNames {
		if !cli.FLAG_TABLE.ShowHelpForSectionViaDowndash(downdashSectionName) {
			fmt.Printf("mlr: flag-section \"%s\" not found.\n", downdashSectionName)
		}
	}
}

// For webdocs autogen: we want the headline separately so we can backtick it.
func showHeadlineForFlag(flagNames []string) {
	for _, flagName := range flagNames {
		if !cli.FLAG_TABLE.ShowHeadlineForFlag(flagName) {
			fmt.Printf("mlr: flag \"%s\" not found..\n", flagName)
		}
	}
}

// For webdocs autogen
func showHelpForFlag(flagNames []string) {
	for _, flagName := range flagNames {
		if !cli.FLAG_TABLE.ShowHelpForFlag(flagName) {
			fmt.Printf("mlr: flag \"%s\" not found..\n", flagName)
		}
	}
}

// ----------------------------------------------------------------
func listVerbs() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		transformers.ListVerbNamesAsParagraph()
	} else {
		transformers.ListVerbNamesVertically()
	}
}

func listVerbsAsParagraph() {
	transformers.ListVerbNamesAsParagraph()
}

func helpForVerb(args []string) {
	for _, arg := range args {
		if !transformers.ShowHelpForTransformer(arg) {
			fmt.Printf(
				"mlr: verb \"%s\" not found. Please use \"mlr help list-verbs\" for a list.\n",
				arg)
		}
	}
}

func usageVerbs() {
	transformers.UsageVerbs()
}

// ----------------------------------------------------------------
func listFunctions() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionNamesAsParagraph()
	} else {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionNamesVertically()
	}
}

func listFunctionClasses() {
	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionClasses()
}

func listFunctionsInClass(classes []string) {
	for _, class := range classes {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsInClass(class)
	}
}

func listFunctionsAsParagraph() {
	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionNamesAsParagraph()
}

func listFunctionsAsTable() {
	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsAsTable()
}

func usageFunctions() {
	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsages()
}

func usageFunctionsByClass() {
	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsagesByClass()
}

func helpForFunction(args []string) {
	for _, arg := range args {
		cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsage(arg)
	}
}

func helpByExactSearch(things []string) bool {
	foundAny := false
	for _, thing := range things {
		foundThisOne := helpByExactSearchOne(thing)
		foundAny = foundAny || foundThisOne
		if !foundThisOne {
			fmt.Printf("No help found for \"%s\". Please try 'mlr help find %s' for approximate match.\n", thing, thing)
			fmt.Printf("See also 'mlr help topics'.\n")
		}
	}

	return foundAny
}

// We need to look various places, e.g. "sec2gmt" is the name of a verb as well
// as a DSL function.
func helpByExactSearchOne(thing string) bool {
	found := false

	// flag
	if cli.FLAG_TABLE.ShowHelpForFlagWithName(thing) {
		found = true
	}

	// verb
	if transformers.ShowHelpForTransformer(thing) {
		found = true
	}

	// function
	if cst.BuiltinFunctionManagerInstance.TryListBuiltinFunctionUsage(thing) {
		found = true
	}

	// keyword
	if cst.TryUsageForKeyword(thing) {
		found = true
	}

	return found
}

func helpByApproximateSearch(things []string) bool {
	foundAny := false
	for _, thing := range things {
		foundThisOne := helpByApproximateSearchOne(thing)
		foundAny = foundAny || foundThisOne
	}

	return foundAny
}

func helpByApproximateSearchOne(thing string) bool {
	found := false

	// flag
	if cli.FLAG_TABLE.ShowHelpForFlagApproximateWithName(thing) {
		found = true
	}

	// verb
	if transformers.ShowHelpForTransformerApproximate(thing) {
		found = true
	}

	// function
	if cst.BuiltinFunctionManagerInstance.TryListBuiltinFunctionUsageApproximate(thing) {
		found = true
	}

	// keyword
	if cst.TryUsageForKeywordApproximate(thing) {
		found = true
	}

	if !found {
		fmt.Printf("No help found for \"%s\". See also 'mlr help topics'.\n", thing)
	}
	return found
}

// ----------------------------------------------------------------
func listKeywords() {
	if isatty.IsTerminal(os.Stdout.Fd()) {
		cst.ListKeywordsAsParagraph()
	} else {
		cst.ListKeywordsVertically()
	}
}

func listKeywordsAsParagraph() {
	cst.ListKeywordsAsParagraph()
}

func usageKeywords() {
	cst.UsageKeywords()
}

func helpForKeyword(args []string) {
	for _, arg := range args {
		cst.UsageForKeyword(arg)
	}
}

// ----------------------------------------------------------------
// flagTableNilCheckflagTableNilCheck is invoked by an internal-only
// command-handler. It's intended to be invoked from a regression-test context.
// It makes sure (at build time) that the flags-table isn't missing help strings
// for any flags, etc.
func flagTableNilCheck() {
	cli.FLAG_TABLE.NilCheck()
}
