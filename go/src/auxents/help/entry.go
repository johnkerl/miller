// ================================================================
// TODO: comment
// ================================================================

// TODO: validate (args []string) non-empty by handlers that don't take them --
// maybe a bool in the LUT

package help

import (
	"fmt"
	"os"
	"path"

	"miller/src/cliutil"
	"miller/src/dsl/cst"
	"miller/src/lib"
	"miller/src/transformers"
)

type tHandlerFunc func(args []string)

type handlerInfo struct {
	name        string
	handlerFunc tHandlerFunc
}

// We get a Golang "initialization loop" if this is defined statically. So, we
// use a "package init" function.
var handlerLookupTable = []handlerInfo{}

func init() {
	handlerLookupTable = []handlerInfo{
		{name: "topics", handlerFunc: listTopics},
		{name: "auxents", handlerFunc: helpAuxents},
		{name: "comments-in-data", handlerFunc: helpCommentsInData},
		{name: "compressed-data", handlerFunc: helpCompressedDataOptions},
		{name: "csv-options", handlerFunc: helpCSVOptions},
		{name: "data-formats", handlerFunc: helpDataFormats},
		{name: "data-format-options", handlerFunc: helpDataFormatOptions},
		{name: "format-conversion", handlerFunc: helpFormatConversionKeystrokeSaverOptions},
		{name: "list-functions", handlerFunc: ListFunctions},
		{name: "list-keywords", handlerFunc: ListKeywords},
		{name: "list-verbs", handlerFunc: ListVerbs},
		// TODO: help for function
		// TODO: help for keyword
		// {name: "function", handlerFunc: HelpFunction},
		// {name: "keyword", handlerFunc: HelpKeyword},
		{name: "misc", handlerFunc: helpMiscOptions},
		{name: "mlrrc", handlerFunc: helpMlrrc},
		{name: "output-colorizations", handlerFunc: helpOutputColorization},
		// type-arithmetic-info
		//		printTypeArithmeticInfo(os.Stdout, lib.MlrExeName());
		// TODO
		//{name: "usage-functions", handlerFunc: UsageFunctions},
		//{name: "usage-keywords", handlerFunc: UsageKeywords},
		//{name: "usage-verbs", handlerFunc: UsageVerbs},
		// TODO: search
	}
}
		//listAllVerbs(os.Stdout, "")
		//help.ListBuiltinFunctions(os.Stdout)

// ================================================================
func HelpUsage(verbName string, o *os.File, exitCode int) {
	exeName := path.Base(os.Args[0])
	fmt.Printf("Usage: %s %s {TODO}\n", exeName, verbName)

	os.Exit(exitCode)
}

// Here the args are the full Miller command line: "mlr help foo bar".
func HelpMain(args []string) int {
	args = args[2:]

	// "mlr help" and nothing else
	if len(args) == 0 {
		handleDefault(args)
		return 0
	}

	// "mlr help something" where we recognize the something
	subcommand := args[0]
	for _, info := range handlerLookupTable {
		if info.name == subcommand {
			info.handlerFunc(args)
			return 0
		}
	}

	// "mlr help something" where we do not recognize the something
	listTopics(args)

	return 0
}

// ----------------------------------------------------------------
func MainUsage(o *os.File) {
	fmt.Fprintf(o,
		`Usage: mlr [I/O options] {verb} [verb-dependent options ...] {zero or more file names}
Output of one verb may be chained as input to another using "then", e.g.
  mlr stats1 -a min,mean,max -f flag,u,v -g color then sort -f color
Please see 'mlr help topics' for more information.
`)
}

// ----------------------------------------------------------------
func handleDefault(args []string) {
	MainUsage(os.Stdout)
}

// ----------------------------------------------------------------
func listTopics(args []string) {
	fmt.Println("Type 'mlr help {topic} for any of the following topics:")
	for _, info := range handlerLookupTable {
		fmt.Printf("  %s\n", info.name)
	}
}

// ----------------------------------------------------------------
func helpMlrrc(args []string) {
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
func helpOutputColorization(args []string) {
	fmt.Print(`Things having colors:
* Keys in CSV header lines, JSON keys, etc
* Values in CSV data lines, JSON scalar values, etc
 in regression-test output
* Some online-help strings

Rules for coloring:
* By default, colorize output only if writing to stdout and stdout is a TTY.
  * Example: color: mlr --csv cat foo.csv
  * Example: no color: mlr --csv cat foo.csv > bar.csv
  * Example: no color: mlr --csv cat foo.csv | less
* The default colors were chosen since they look OK with white or black terminal background,
  and are differentiable with common varieties of human color vision.

Mechanisms for coloring:
* Miller uses ANSI escape sequences only. This does not work on Windows except on Cygwin.
* Requires TERM environment variable to be set to non-empty string.
* Doesn't try to check to see whether the terminal is capable of 256-color
  ANSI vs 16-color ANSI. Note that if colors are in the range 0..15
  then 16-color ANSI escapes are used, so this is in the user's control.

How you can control colorization:
* Suppression/unsuppression:
  * Environment variable export MLR_NO_COLOR=true means don't color even if stdout+TTY.
  * Environment variable export MLR_ALWAYS_COLOR=true means do color even if not stdout+TTY.
    For example, you might want to use this when piping mlr output to less -r.
  * Command-line flags --no-color or -M, --always-color or -C.

* Color choices can be specified by using environment variables, or command-line flags,
  with values 0..255:
  * export MLR_KEY_COLOR=208, MLR_VALUE_COLOR-33, etc.:
    MLR_KEY_COLOR MLR_VALUE_COLOR MLR_PASS_COLOR MLR_FAIL_COLOR
    MLR_REPL_PS1_COLOR MLR_REPL_PS2_COLOR MLR_HELP_COLOR
  * Command-line flags --key-color 208, --value-color 33, etc.:
    --key-color --value-color --pass-color --fail-color
    --repl-ps1-color --repl-ps2-color --help-color
  * This is particularly useful if your terminal's background color clashes with current settings.

If environment-variable settings and command-line flags are both provided,the latter take precedence.

Please do mlr --list-colors to see the available color codes.
`)
}

// ----------------------------------------------------------------
// mlr --version
// = mlr version ?

// mlr --usage-data-format-examples             -> mlr help format-examples
// mlr --usage-examples                         -> mlr help ___
// mlr --usage-help-options                     -> mlr help ___ <-- should be help-of-help meta thing -- LUT
// mlr --usage-list-all-verbs                   -> mlr help ___
// mlr --usage-functions                        -> mlr help ___
// mlr --usage-data-format-options              -> mlr help ___
// mlr --usage-comments-in-data                 -> mlr help ___
// mlr --usage-format-conversion-keystroke-sa.. -> mlr help keystroke-savers
// mlr --usage-compressed-data-options          -> mlr help compressed-data
// mlr --usage-separator-options                -> mlr help separators
// mlr --usage-csv-options                      -> mlr help csv-options
// mlr --usage-double-quoting                   -> mlr help double-quoting
// mlr --usage-numerical-formatting             -> mlr help number-formatting
// mlr --usage-other-options                    -> mlr help ???
// mlr --usage-then-chaining                    -> mlr help then-chaining
// mlr --usage-auxents                          -> mlr help auxents
// mlr --list-all-verbs-raw                     -> mlr help list-verbs
// mlr #{verb} -h                               -> mlr help ????
// mlr --list-all-functions-raw                 -> mlr help list-functions
// mlr --help-function #{function}              -> mlr help ___
// mlr --list-all-keywords-raw                  -> mlr help list-keywords
// mlr --help-keyword                           -> mlr help '#{keyword}'

// mlr help {foo} w/ polymorphic lookup -- ?
// cst.BuiltinFunctionManagerInstance.TryListBuiltinFunctionUsage(arg, os.Stdout)
//
// ? integrate repl help w/ main help ... ?

// ----------------------------------------------------------------
// To keep:
//   -h or --help                 Show this message.
//   --version                    Show the software version.
//   {verb name} --help           Show verb-specific help.
//   --help-all-verbs             Show help on all verbs.
//   -l                           List only verb names.
//   -L                           List only verb names, one per line.
//   -f                           Show help on all built-in functions.
//   -F                           Show a bare listing of built-in functions by name.
//   -k or --help-all-keywords    Show help on all keywords.
//   -K                           Show a bare listing of keywords by name.

// mlr help verbs
// mlr help data-format-options
// mlr help comments-in-data
// mlr help other-options
// mlr help functions
// mlr help mlrrc
// mlr aux-list

// ================================================================
// ================================================================
// ================================================================

//func helpHelpOptions(o *os.File, argv0 string) {
//	fmt.Printf("  -h or --help                 Show this message.\n")
//	fmt.Printf("  --version                    Show the software version.\n")
//	fmt.Printf("  {verb name} --help           Show verb-specific help.\n")
//	fmt.Printf("  --help-all-verbs             Show help on all verbs.\n")
//	fmt.Printf("  -l or --list-all-verbs       List only verb names.\n")
//	fmt.Printf("  -L                           List only verb names, one per line.\n")
//	fmt.Printf("  -f or --help-all-functions   Show help on all built-in functions.\n")
//	fmt.Printf("  -F                           Show a bare listing of built-in functions by name.\n")
//	fmt.Printf("  -k or --help-all-keywords    Show help on all keywords.\n")
//	fmt.Printf("  -K                           Show a bare listing of keywords by name.\n")
//}

//func helpFunctions(o *os.File) {
//	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsRaw(os.Stdout)
//	fmt.Printf("Please use "mlr --help-function {function name}" for function-specific help.\n")
//}

// ----------------------------------------------------------------
func helpAuxents(args []string) {
	fmt.Print(`Miller has a few otherwise-standalone executables packaged within it.
They do not participate in any other parts of Miller.
Please "mlr aux-list" for more information.
`)
	// imports miller/src/auxents: import cycle not allowed
	//auxents.ShowAuxEntries(o)
}

// ----------------------------------------------------------------
func helpCommentsInData(args []string) {
	fmt.Printf(
		`--skip-comments                 Ignore commented lines (prefixed by "%s")
                                within the input.
--skip-comments-with {string}   Ignore commented lines within input, with
                                specified prefix.
--pass-comments                 Immediately print commented lines (prefixed by "%s")
                                within the input.
--pass-comments-with {string}   Immediately print commented lines within input, with
                                specified prefix.

Notes:
* Comments are only honored at the start of a line.
* In the absence of any of the above four options, comments are data like
  any other text.
* When pass-comments is used, comment lines are written to standard output
  immediately upon being read; they are not part of the record stream.  Results
  may be counterintuitive. A suggestion is to place comments at the start of
  data files.
`,
		cliutil.DEFAULT_COMMENT_STRING,
		cliutil.DEFAULT_COMMENT_STRING)
}

// ----------------------------------------------------------------
func helpCompressedDataOptions(args []string) {
	fmt.Print(`Decompression done within the Miller process itself:
--gzin  Uncompress gzip within the Miller process. Done by default if file ends in ".gz".
--bz2in Uncompress bz2ip within the Miller process. Done by default if file ends in ".bz2".
--zin   Uncompress zlib within the Miller process. Done by default if file ends in ".z".

Decompression done outside the Miller process:
--prepipe {command} You can, of course, already do without this for single input files,
  e.g. "gunzip < myfile.csv.gz | mlr ..."
--prepipex {command} Like --prepipe with one exception: doesn't insert '<' between
  command and filename at runtime. Useful for some commands like 'unzip -qc'
  which don't read standard input.

Using --prepipe and --prepipex you can specify an action to be taken on each
input file. This prepipe command must be able to read from standard input; it
will be invoked with {command} < {filename}.

Examples:
  mlr --prepipe gunzip
  mlr --prepipe zcat -cf
  mlr --prepipe xz -cd
  mlr --prepipe cat

Note that this feature is quite general and is not limited to decompression
utilities. You can use it to apply per-file filters of your choice.  For output
compression (or other) utilities, simply pipe the output:
mlr ... | {your compression command} > outputfilenamegoeshere

Lastly, note that if --prepipe or --prepipex is specified, it replaces any
decisions that might have been made based on the file suffix. Also,
--gzin/--bz2in/--zin are ignored if --prepipe is also specified.
`)
}

// ----------------------------------------------------------------
func helpCSVOptions(args []string) {
	fmt.Print(
		`  --implicit-csv-header Use 1,2,3,... as field labels, rather than from line 1
                     of input files. Tip: combine with "label" to recreate
                     missing headers.
  --no-implicit-csv-header Do not use --implicit-csv-header. This is the default
                     anyway -- the main use is for the flags to 'mlr join' if you have
                     main file(s) which are headerless but you want to join in on
                     a file which does have a CSV header. Then you could use
                     'mlr --csv --implicit-csv-header join --no-implicit-csv-header
                     -l your-join-in-with-header.csv ... your-headerless.csv'
  --allow-ragged-csv-input|--ragged If a data line has fewer fields than the header line,
                     fill remaining keys with empty string. If a data line has more
                     fields than the header line, use integer field labels as in
                     the implicit-header case.
  --headerless-csv-output   Print only CSV data lines.
  -N                 Keystroke-saver for --implicit-csv-header --headerless-csv-output.
`)
}

// ----------------------------------------------------------------
func helpDataFormats(args []string) {
	fmt.Printf(
		`CSV/CSV-lite: comma-separated values with separate header line
TSV: same but with tabs in places of commas
+---------------------+
| apple,bat,cog       |
| 1,2,3               | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
| 4,5,6               | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
+---------------------+

JSON (sequence or array of objects):
+---------------------+
| {                   |
|  "apple": 1,        | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
|  "bat": 2,          |
|  "cog": 3           |
| }                   |
| {                   |
|   "dish": {         | Record 2: "dish:egg" => "7", "dish:flint" => "8", "garlic" => ""
|     "egg": 7,       |
|     "flint": 8      |
|   },                |
|   "garlic": ""      |
| }                   |
+---------------------+

PPRINT: pretty-printed tabular
+---------------------+
| apple bat cog       |
| 1     2   3         | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
| 4     5   6         | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
+---------------------+

Markdown tabular (supported for output only):
+-----------------------+
| | apple | bat | cog | |
| | ---   | --- | --- | |
| | 1     | 2   | 3   | | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
| | 4     | 5   | 6   | | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
+-----------------------+

XTAB: pretty-printed transposed tabular
+---------------------+
| apple 1             | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
| bat   2             |
| cog   3             |
|                     |
| dish 7              | Record 2: "dish" => "7", "egg" => "8"
| egg  8              |
+---------------------+

DKVP: delimited key-value pairs (Miller default format)
+---------------------+
| apple=1,bat=2,cog=3 | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
| dish=7,egg=8,flint  | Record 2: "dish" => "7", "egg" => "8", "3" => "flint"
+---------------------+

NIDX: implicitly numerically indexed (Unix-toolkit style)
+---------------------+
| the quick brown     | Record 1: "1" => "the", "2" => "quick", "3" => "brown"
| fox jumped          | Record 2: "1" => "fox", "2" => "jumped"
+---------------------+
`)
}

// ----------------------------------------------------------------
func helpDataFormatOptions(args []string) {
	fmt.Printf(
		`--idkvp   --odkvp   --dkvp      Delimited key-value pairs, e.g "a=1,b=2"
                                 (Miller's default format).

--inidx   --onidx   --nidx      Implicitly-integer-indexed fields (Unix-toolkit style).
-T                              Synonymous with "--nidx --fs tab".

--icsv    --ocsv    --csv       Comma-separated value (or tab-separated with --fs tab, etc.)

--itsv    --otsv    --tsv       Keystroke-savers for "--icsv --ifs tab",
                                "--ocsv --ofs tab", "--csv --fs tab".
--iasv    --oasv    --asv       Similar but using ASCII FS %s and RS %s\n",
--iusv    --ousv    --usv       Similar but using Unicode FS %s\n",
                                and RS %s\n",

--icsvlite --ocsvlite --csvlite Comma-separated value (or tab-separated with --fs tab, etc.).
							    The 'lite' CSV does not handle RFC-CSV double-quoting rules; is
							    slightly faster and handles heterogeneity in the input stream via
							    empty newline followed by new header line. See also
								%s/file-formats.html#csv-tsv-asv-usv-etc

--itsvlite --otsvlite --tsvlite Keystroke-savers for "--icsvlite --ifs tab",
                                "--ocsvlite --ofs tab", "--csvlite --fs tab".
-t                              Synonymous with --tsvlite.
--iasvlite --oasvlite --asvlite Similar to --itsvlite et al. but using ASCII FS %s and RS %s\n",
--iusvlite --ousvlite --usvlite Similar to --itsvlite et al. but using Unicode FS %s\n",
                                and RS %s\n",

--ipprint --opprint --pprint    Pretty-printed tabular (produces no
                                output until all input is in).
                    --right     Right-justifies all fields for PPRINT output.
                    --barred    Prints a border around PPRINT output
                                (only available for output).

          --omd                 Markdown-tabular (only available for output).

--ixtab   --oxtab   --xtab      Pretty-printed vertical-tabular.
                    --xvright   Right-justifies values for XTAB format.

--ijson   --ojson   --json      JSON tabular: sequence or list of one-level
                                maps: {...}{...} or [{...},{...}].
                    --jvstack   Put one key-value pair per line for JSON output.
                 --no-jvstack   Put objects/arrays all on one line for JSON output.
              --jsonx --ojsonx  Keystroke-savers for --json --jvstack
              --jsonx --ojsonx  and --ojson --jvstack, respectively.
                    --jlistwrap Wrap JSON output in outermost [ ].
            --oflatsep {string} Separator for flattening multi-level JSON keys,
                                e.g. '{"a":{"b":3}}' becomes a:b => 3 for
                                non-JSON formats. Defaults to %s.\n",

-p is a keystroke-saver for --nidx --fs space --repifs

Examples: --csv for CSV-formatted input and output; --icsv --opprint for
CSV-formatted input and pretty-printed output.

Please use --iformat1 --oformat2 rather than --format1 --oformat2.
The latter sets up input and output flags for format1, not all of which
are overridden in all cases by setting output format to format2.`,

		cliutil.ASV_FS_FOR_HELP,
		cliutil.ASV_RS_FOR_HELP,
		cliutil.USV_FS_FOR_HELP,
		cliutil.USV_RS_FOR_HELP,
		lib.DOC_URL,
		cliutil.ASV_FS_FOR_HELP,
		cliutil.ASV_RS_FOR_HELP,
		cliutil.USV_FS_FOR_HELP,
		cliutil.USV_RS_FOR_HELP,
		cliutil.DEFAULT_JSON_FLATTEN_SEPARATOR,
	)
	fmt.Println()
}

// ----------------------------------------------------------------
////func helpDoubleQuoting(o *os.File, argv0 string) {
////	fmt.Printf("  --quote-all        Wrap all fields in double quotes\n")
////	fmt.Printf("  --quote-none       Do not wrap any fields in double quotes, even if they have\n")
////	fmt.Printf("                     OFS or ORS in them\n")
////	fmt.Printf("  --quote-minimal    Wrap fields in double quotes only if they have OFS or ORS\n")
////	fmt.Printf("                     in them (default)\n")
////	fmt.Printf("  --quote-numeric    Wrap fields in double quotes only if they have numbers\n")
////	fmt.Printf("                     in them\n")
////	fmt.Printf("  --quote-original   Wrap fields in double quotes if and only if they were\n")
////	fmt.Printf("                     quoted on input. This isn't sticky for computed fields:\n")
////	fmt.Printf("                     e.g. if fields a and b were quoted on input and you do\n")
////	fmt.Printf("                     "put '$c = $a . $b'" then field c won't inherit a or b's\n")
////	fmt.Printf("                     was-quoted-on-input flag.\n")
////}

// ----------------------------------------------------------------
func helpFormatConversionKeystrokeSaverOptions(args []string) {
	fmt.Print(`As keystroke-savers for format-conversion you may use the following:
--c2t --c2d --c2n --c2j --c2x --c2p --c2m
--t2c       --t2d --t2n --t2j --t2x --t2p --t2m
--d2c --d2t       --d2n --d2j --d2x --d2p --d2m
--n2c --n2t --n2d       --n2j --n2x --n2p --n2m
--j2c --j2t --j2d --j2n       --j2x --j2p --j2m
--x2c --x2t --x2d --x2n --x2j       --x2p --x2m
--p2c --p2t --p2d --p2n --p2j --p2x       --p2m
The letters c t d n j x p m refer to formats CSV, TSV, DKVP, NIDX, JSON, XTAB,
PPRINT, and markdown, respectively. Note that markdown format is available for
output only.
`)
}

// ----------------------------------------------------------------
func ListFunctions(args []string) {
	fmt.Println("TODO: list functions")
	//TODO
	//cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsages(os.Stdout)
	//cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsRaw(os.Stdout)
	//cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionUsage(args[argi+1], os.Stdout)
	//TODO: as-table
		//		fmgr_t* pfmgr = fmgr_alloc();
		//		fmgr_list_all_functions_as_table(pfmgr, os.Stdout);
		//		fmgr_free(pfmgr, nil);
		//		return true;
}

// ----------------------------------------------------------------
func HelpFunction(args []string) {
	fmt.Println("TODO: help for function")
}

// TODO: help all functions

// ----------------------------------------------------------------
func ListKeywords(args []string) {
	fmt.Println("TODO: list keywords")
}

// ----------------------------------------------------------------
func HelpKeyword(args []string) {
	fmt.Println("TODO: help for keyword")
}

// TODO: help all keywords

// ----------------------------------------------------------------
func ListVerbs(args []string) {
	transformers.ListAllVerbNamesAsParagraph()
}

func ListAllVerbNames() {
	transformers.ListAllVerbNames()
}

func ListAllVerbNamesAsParagraph() {
	transformers.ListAllVerbNamesAsParagraph()
}

// ----------------------------------------------------------------
func helpMiscOptions(args []string) {
	fmt.Printf(`  --seed {n} with n of the form 12345678 or 0xcafefeed. For put/filter
                     urand()/urandint()/urand32().
  --nr-progress-mod {m}, with m a positive integer: print filename and record
                     count to os.Stderr every m input records.
  --from {filename}  Use this to specify an input file before the verb(s),
                     rather than after. May be used more than once. Example:
                     "mlr --from a.dat --from b.dat cat" is the same as
                     "mlr cat a.dat b.dat".
  --mfrom {filenames} --  Use this to specify one of more input files before the verb(s),
                     rather than after. May be used more than once.
                     The list of filename must end with "--". This is useful
                     for example since "--from *.csv" doesn't do what you might
                     hope but "--mfrom *.csv --" does.
  --load {filename}  Load DSL script file for all put/filter operations on the command line.
                     If the name following --load is a directory, load all "*.mlr" files
                     in that directory. This is just like "put -f" and "filter -f"
                     except it's up-front on the command line, so you can do something like
                     alias mlr='mlr --load ~/myscripts' if you like.
  --mload {names} -- Like --load but works with more than one filename,
                     e.g. '--mload *.mlr --'.
  -n                 Process no input files, nor standard input either. Useful
                     for mlr put with begin/end statements only. (Same as --from
                     /dev/null.) Also useful in "mlr -n put -v '...'" for
                     analyzing abstract syntax trees (if that's your thing).
  -I                 Process files in-place. For each file name on the command
                     line, output is written to a temp file in the same
                     directory, which is then renamed over the original. Each
                     file is processed in isolation: if the output format is
                     CSV, CSV headers will be present in each output file
                     statistics are only over each file's own records; and so on.
`)
}

// ----------------------------------------------------------------
//func helpNumericalFormatting(o *os.File, argv0 string) {
//	fmt.Printf("  --ofmt {format}    E.g. %%.18f, %%.0f, %%9.6e. Please use sprintf-style codes for\n")
//	fmt.Printf("                     floating-point nummbers. If not specified, default formatting is used.\n")
//	fmt.Printf("                     See also the fmtnum function within mlr put (mlr --help-all-functions);\n")
//	fmt.Printf("                     see also the format-values function.\n")
//}

// ----------------------------------------------------------------
////func helpSeparatorOptions(args []string) {
////	fmt.Print(`Separator options:
////  --rs     --irs     --ors              Record separators, e.g. 'lf' or '\\r\\n'
////  --fs     --ifs     --ofs  --repifs    Field separators, e.g. comma
////  --ps     --ips     --ops              Pair separators, e.g. equals sign
////
////  Notes about line endings:
////  * Default line endings (--irs and --ors) are "auto" which means autodetect from
////    the input file format, as long as the input file(s) have lines ending in either
////    LF (also known as linefeed, '\\n', 0x0a, Unix-style) or CRLF (also known as
////    carriage-return/linefeed pairs, '\\r\\n', 0x0d 0x0a, Windows style).
////  * If both irs and ors are auto (which is the default) then LF input will lead to LF
////    output and CRLF input will lead to CRLF output, regardless of the platform you're
////    running on.
////  * The line-ending autodetector triggers on the first line ending detected in the input
////    stream. E.g. if you specify a CRLF-terminated file on the command line followed by an
////    LF-terminated file then autodetected line endings will be CRLF.
////  * If you use --ors {something else} with (default or explicitly specified) --irs auto
////    then line endings are autodetected on input and set to what you specify on output.
////  * If you use --irs {something else} with (default or explicitly specified) --ors auto
////    then the output line endings used are LF on Unix/Linux/BSD/MacOSX, and CRLF on Windows.
////
////  Notes about all other separators:
////  * IPS/OPS are only used for DKVP and XTAB formats, since only in these formats
////    do key-value pairs appear juxtaposed.
////  * IRS/ORS are ignored for XTAB format. Nominally IFS and OFS are newlines;
////    XTAB records are separated by two or more consecutive IFS/OFS -- i.e.
////    a blank line. Everything above about --irs/--ors/--rs auto becomes --ifs/--ofs/--fs
////    auto for XTAB format. (XTAB's default IFS/OFS are "auto".)
////  * OFS must be single-character for PPRINT format. This is because it is used
////    with repetition for alignment; multi-character separators would make
////    alignment impossible.
////  * OPS may be multi-character for XTAB format, in which case alignment is
////    disabled.
////  * TSV is simply CSV using tab as field separator ("--fs tab").
////  * FS/PS are ignored for markdown format; RS is used.
////  * All FS and PS options are ignored for JSON format, since they are not relevant
////    to the JSON format.
////  * You can specify separators in any of the following ways, shown by example:
////    - Type them out, quoting as necessary for shell escapes, e.g.
////      "--fs '|' --ips :"
////    - C-style escape sequences, e.g. "--rs '\\r\\n' --fs '\\t'".
////    - To avoid backslashing, you can use any of the following names:
////     ")
//////	lhmss_t* pmap = get_desc_to_chars_map()
//////	for (lhmsse_t* pe = pmap.phead; pe != nil; pe = pe.pnext) {
//// %s", pe.key)
//////	}
////
////  * Default separators by format:
////      %-12s %-8s %-8s %s\n", "File format", "RS", "FS", "PS")
//////	lhmss_t* default_rses = get_default_rses()
//////	lhmss_t* default_fses = get_default_fses()
//////	lhmss_t* default_pses = get_default_pses()
//////	for (lhmsse_t* pe = default_rses.phead; pe != nil; pe = pe.pnext) {
//////		char* filefmt = pe.key
//////		char* rs = pe.value
//////		char* fs = lhmss_get(default_fses, filefmt)
//////		char* ps = lhmss_get(default_pses, filefmt)
////      %-12s %-8s %-8s %s\n", filefmt, rebackslash(rs), rebackslash(fs), rebackslash(ps))
//////	}
////}

// ----------------------------------------------------------------
//func helpSeeAlso(o *os.File, argv0 string) {
//	fmt.Printf("For more information please see %s and/or\n", lib.DOC_URL)
//	fmt.Printf("http://github.com/johnkerl/miller.")
//	fmt.Printf(" This is Miller version %s.\n", version.STRING)
//}

// TODO port from src/cli
func ListBuiltinFunctions(o *os.File) {
	cst.BuiltinFunctionManagerInstance.ListBuiltinFunctionsRaw(os.Stdout)
	fmt.Fprintf(o, "Please use \"%s --help-function {function name}\" for function-specific help.\n", lib.MlrExeName())
}
