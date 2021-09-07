// ================================================================
// Online help
// ================================================================

package help

import (
	"fmt"
	"os"

	"github.com/mattn/go-isatty"

	"mlr/src/cli"
	"mlr/src/dsl/cst"
	"mlr/src/lib"
	"mlr/src/transformers"
	"mlr/src/types"
)

// ================================================================
type tZaryHandlerFunc func()
type tUnaryHandlerFunc func(arg string)

type shorthandInfo struct {
	shorthand string
	longhand  string
}

type handlerInfo struct {
	name             string
	zaryHandlerFunc  tZaryHandlerFunc
	unaryHandlerFunc tUnaryHandlerFunc

	// Some handlers are used only for webdoc/manpage autogen and needn't
	// clutter up the on-line help experience for the interactive user
	internal bool
}

// We get a Golang "initialization loop" if this is defined statically. So, we
// use a "package init" function.
var shorthandLookupTable = []shorthandInfo{}
var handlerLookupTable = []handlerInfo{}

func init() {
	// For things like 'mlr -f', invoked through the CLI parser which does not
	// go through our HelpMain().
	shorthandLookupTable = []shorthandInfo{
		{shorthand: "-l", longhand: "list-verbs"},
		{shorthand: "-L", longhand: "usage-verbs"},
		{shorthand: "-f", longhand: "list-functions"},
		{shorthand: "-F", longhand: "usage-functions"},
		{shorthand: "-k", longhand: "list-keywords"},
		{shorthand: "-K", longhand: "usage-keywords"},
	}

	// For things like 'mlr help foo', invoked through the auxent framework
	// which goes through our HelpMain().
	handlerLookupTable = []handlerInfo{
		{name: "topics", zaryHandlerFunc: listTopics},
		{name: "auxents", zaryHandlerFunc: helpAuxents},
		{name: "basic-examples", zaryHandlerFunc: helpBasicExamples},
		{name: "data-formats", zaryHandlerFunc: helpDataFormats},
		{name: "function", unaryHandlerFunc: helpForFunction},
		{name: "keyword", unaryHandlerFunc: helpForKeyword},
		{name: "list-functions", zaryHandlerFunc: listFunctions},
		{name: "list-function-classes", zaryHandlerFunc: listFunctionClasses},
		{name: "list-functions-in-class", unaryHandlerFunc: listFunctionsInClass},
		{name: "list-functions-as-paragraph", zaryHandlerFunc: listFunctionsAsParagraph},
		{name: "list-keywords", zaryHandlerFunc: listKeywords},
		{name: "list-keywords-as-paragraph", zaryHandlerFunc: listKeywordsAsParagraph},
		{name: "list-verbs", zaryHandlerFunc: listVerbs},
		{name: "list-verbs-as-paragraph", zaryHandlerFunc: listVerbsAsParagraph},
		{name: "mlrrc", zaryHandlerFunc: helpMlrrc},
		{name: "number-formatting", zaryHandlerFunc: helpNumberFormatting},
		{name: "type-arithmetic-info", zaryHandlerFunc: helpTypeArithmeticInfo},
		{name: "usage-functions", zaryHandlerFunc: usageFunctions},
		{name: "usage-functions-by-class", zaryHandlerFunc: usageFunctionsByClass},
		{name: "usage-keywords", zaryHandlerFunc: usageKeywords},
		{name: "usage-verbs", zaryHandlerFunc: usageVerbs},
		{name: "verb", unaryHandlerFunc: helpForVerb},

		// TODO: to flags-sections
		{name: "comments-in-data", zaryHandlerFunc: helpCommentsInData},
		{name: "compressed-data", zaryHandlerFunc: helpCompressedDataOptions},
		{name: "data-format-options", zaryHandlerFunc: helpDataFormatOptions},
		{name: "double-quoting", zaryHandlerFunc: helpDoubleQuoting},
		{name: "format-conversion", zaryHandlerFunc: helpFormatConversionKeystrokeSaverOptions},
		{name: "separator-options", zaryHandlerFunc: helpSeparatorOptions},

		// Internal-only
		{name: "list-functions-as-table", zaryHandlerFunc: listFunctionsAsTable, internal: true},
		{name: "list-flag-sections", zaryHandlerFunc: listFlagSections, internal: true},
		{name: "print-info-for-section", unaryHandlerFunc: printInfoForSection, internal: true},
		{name: "list-flags-for-section", unaryHandlerFunc: listFlagsForSection, internal: true},
		{name: "show-headline-for-flag", unaryHandlerFunc: showHeadlineForFlag, internal: true},
		{name: "show-help-for-flag", unaryHandlerFunc: showHelpForFlag, internal: true},

		// TBD: have an info-only handler in addition to flags-section
		{name: "output-colorization", zaryHandlerFunc: helpOutputColorization},
	}
}

// ================================================================
// For things like 'mlr help foo', invoked through the auxent framework which
// goes through our HelpMain().  Here, the args are the full Miller command
// line: "mlr help foo bar".
func HelpMain(args []string) int {
	args = args[2:]

	// "mlr help" and nothing else
	if len(args) == 0 {
		handleDefault()
		return 0
	}

	// "mlr help something" where we recognize the something
	name := args[0]
	for _, info := range handlerLookupTable {
		if info.name == name {
			if info.zaryHandlerFunc != nil {
				if len(args) != 1 {
					fmt.Printf("mlr help %s takes no additional argument.\n", name)
					return 0
				}
				info.zaryHandlerFunc()
				return 0
			}
			if info.unaryHandlerFunc != nil {
				if len(args) < 2 {
					fmt.Printf("mlr help %s takes at least one required argument.\n", name)
					return 0
				}
				for _, arg := range args[1:] {
					info.unaryHandlerFunc(arg)
				}
				return 0
			}
		}
	}

	// "mlr help something" where we do not recognize the something
	listTopics()

	return 0
}

// ----------------------------------------------------------------
func MainUsage(o *os.File) {
	fmt.Fprintf(o,
		`Usage: mlr [flags] {verb} [verb-dependent options ...] {zero or more file names}
Output of one verb may be chained as input to another using "then", e.g.
  mlr stats1 -a min,mean,max -f flag,u,v -g color then sort -f color
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
	for _, sinfo := range shorthandLookupTable {
		if sinfo.shorthand == arg {
			for _, info := range handlerLookupTable {
				if info.name == sinfo.longhand {
					info.zaryHandlerFunc()
					return true
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
	for _, info := range handlerLookupTable {
		if !info.internal {
			fmt.Printf("  mlr help %s\n", info.name)
		}
	}
	fmt.Println("Shorthands:")
	for _, info := range shorthandLookupTable {
		fmt.Printf("  mlr %s = mlr help %s\n", info.shorthand, info.longhand)
	}
}

// ----------------------------------------------------------------
func helpAuxents() {
	fmt.Print(`Miller has a few otherwise-standalone executables packaged within it.
They do not participate in any other parts of Miller.
Please "mlr aux-list" for more information.
`)
	// imports mlr/src/auxents: import cycle not allowed
	// auxents.ShowAuxEntries(o)
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
func helpCommentsInData() {
	cli.CommentsInDataPrintInfo()
}

// ----------------------------------------------------------------
func helpCompressedDataOptions() {
	cli.CompressedDataPrintInfo()
}

// ----------------------------------------------------------------
func helpDataFormats() {
	fmt.Printf(
		`CSV/CSV-lite: comma-separated values with separate header line
TSV: same but with tabs in places of commas
+---------------------+
| apple,bat,cog       |
| 1,2,3               | Record 1: "apple":"1", "bat":"2", "cog":"3"
| 4,5,6               | Record 2: "apple":"4", "bat":"5", "cog":"6"
+---------------------+

JSON (sequence or array of objects):
+---------------------+
| {                   |
|  "apple": 1,        | Record 1: "apple":"1", "bat":"2", "cog":"3"
|  "bat": 2,          |
|  "cog": 3           |
| }                   |
| {                   |
|   "dish": {         | Record 2: "dish:egg":"7",
|     "egg": 7,       | "dish:flint":"8", "garlic":""
|     "flint": 8      |
|   },                |
|   "garlic": ""      |
| }                   |
+---------------------+

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
func helpDataFormatOptions() {
	cli.FileFormatPrintInfo()
}

// ----------------------------------------------------------------
// TBD FOR MILLER 6:

func helpDoubleQuoting() {
	fmt.Printf("THIS IS STILL WIP FOR MILLER 6\n")
	fmt.Println(
		`--quote-all        Wrap all fields in double quotes
--quote-none       Do not wrap any fields in double quotes, even if they have
                   OFS or ORS in them
--quote-minimal    Wrap fields in double quotes only if they have OFS or ORS
                   in them (default)
--quote-numeric    Wrap fields in double quotes only if they have numbers
                   in them
--quote-original   Wrap fields in double quotes if and only if they were
                   quoted on input. This isn't sticky for computed fields:
                   e.g. if fields a and b were quoted on input and you do
                   "put '$c = $a . $b'" then field c won't inherit a or b's
                   was-quoted-on-input flag.`)
}

// ----------------------------------------------------------------
func helpFormatConversionKeystrokeSaverOptions() {
	cli.FileFormatPrintInfo()
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

See also:
https://miller.readthedocs.io/en/latest/customization.html
`)
}

// ----------------------------------------------------------------
func helpOutputColorization() {
	cli.OutputColorizationPrintInfo()
}

// ----------------------------------------------------------------
// TBD FOR MILLER 6:

func helpNumberFormatting() {
	fmt.Printf("THIS IS STILL WIP FOR MILLER 6\n")
	fmt.Printf("  --ofmt {format}    E.g. %%.18f, %%.0f, %%9.6e. Please use sprintf-style codes for\n")
	fmt.Printf("                     floating-point nummbers. If not specified, default formatting is used.\n")
	fmt.Printf("                     See also the fmtnum function within mlr put (mlr --help-all-functions);\n")
	fmt.Printf("                     see also the format-values function.\n")
}

// ----------------------------------------------------------------
// TBD FOR MILLER 6:

func helpSeparatorOptions() {
	cli.SeparatorPrintInfo()
}

// ----------------------------------------------------------------
func helpTypeArithmeticInfo() {
	mlrvals := []*types.Mlrval{
		types.MlrvalPointerFromInt(1),
		types.MlrvalPointerFromFloat64(2.5),
		types.MLRVAL_ABSENT,
		types.MLRVAL_ERROR,
	}

	n := len(mlrvals)

	for i := -2; i < n; i++ {
		if i == -2 {
			fmt.Printf("%-10s |", "(+)")
		} else if i == -1 {
			fmt.Printf("%-10s +", "------")
		} else {
			fmt.Printf("%-10s |", mlrvals[i].String())
		}
		for j := 0; j < n; j++ {
			if i == -2 {
				fmt.Printf(" %-10s", mlrvals[j].String())
			} else if i == -1 {
				fmt.Printf(" %-10s", "------")
			} else {
				sum := types.MlrvalBinaryPlus(mlrvals[i], mlrvals[j])
				fmt.Printf(" %-10s", sum.String())
			}
		}
		fmt.Println()
	}

}

// ----------------------------------------------------------------
// listFlagSections is for webdoc/manpage autogen.

func listFlagSections() {
	// xxx temp factorization
	cli.FLAG_TABLE.ListFlagSections()
}

func printInfoForSection(sectionName string) {
	if !cli.FLAG_TABLE.PrintInfoForSection(sectionName) {
		fmt.Printf(
			"mlr: flag-section \"%s\" not found. Please use \"mlr help list-flag-sections\" for a list.\n",
			sectionName)
	}
}

func listFlagsForSection(sectionName string) {
	if !cli.FLAG_TABLE.ListFlagsForSection(sectionName) {
		fmt.Printf(
			"mlr: flag-section \"%s\" not found. Please use \"mlr help list-flag-sections\" for a list.\n",
			sectionName)
	}
}

func showHeadlineForFlag(flagName string) {
	if !cli.FLAG_TABLE.ShowHeadlineForFlag(flagName) {
		fmt.Printf("mlr: flag \"%s\" not found..\n", flagName)
	}
}

func showHelpForFlag(flagName string) {
	if !cli.FLAG_TABLE.ShowHelpForFlag(flagName) {
		fmt.Printf("mlr: flag \"%s\" not found..\n", flagName)
	}
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

func listFunctionsInClass(class string) {
	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsInClass(class)
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

func helpForFunction(arg string) {
	cst.BuiltinFunctionManagerInstance.TryListBuiltinFunctionUsage(arg)
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

func helpForKeyword(arg string) {
	cst.UsageForKeyword(arg)
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

func helpForVerb(arg string) {
	transformerSetup := transformers.LookUp(arg)
	if transformerSetup != nil {
		transformerSetup.UsageFunc(os.Stdout, true, 0)
	} else {
		fmt.Printf(
			"mlr: verb \"%s\" not found. Please use \"mlr help list-verbs\" for a list.\n",
			arg)
	}
}

func usageVerbs() {
	transformers.UsageVerbs()
}
