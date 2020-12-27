package cli

import (
	"fmt"
	"os"

	"miller/version"
)

// ----------------------------------------------------------------
func mainUsageShort() {
	fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for detailed usage information.\n", os.Args[0])
	os.Exit(1)
}

// ----------------------------------------------------------------
// The mainUsageLong() function is split out into subroutines in support of the
// manpage autogenerator.

func mainUsageLong(o *os.File, argv0 string) {
	mainUsageSynopsis(o, argv0)
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "COMMAND-LINE-SYNTAX EXAMPLES:\n")
	mainUsageExamples(o, argv0, "  ")
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "DATA-FORMAT EXAMPLES:\n")
	mainUsageDataFormatExamples(o, argv0)
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "HELP OPTIONS:\n")
	mainUsageHelpOptions(o, argv0)
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "CUSTOMIZATION VIA .MLRRC:\n")
	mainUsageMlrrc(o, argv0)
	fmt.Fprintf(o, "\n")

	fmt.Fprintf(o, "VERBS:\n")
	listAllVerbs(o, "  ")
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "FUNCTIONS FOR THE FILTER AND PUT VERBS:\n");
	//	mainUsageFunctions(o, argv0, "  ");
	//	fmt.Fprintf(o, "\n");

	fmt.Fprintf(o, "DATA-FORMAT OPTIONS, FOR INPUT, OUTPUT, OR BOTH:\n")
	mainUsageDataFormatOptions(o, argv0)
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "COMMENTS IN DATA:\N");
	//	mainUsageCommentsInData(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	fmt.Fprintf(o, "FORMAT-CONVERSION KEYSTROKE-SAVER OPTIONS:\n")
	mainUsageFormatConversionKeystrokeSaverOptions(o, argv0)
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "COMPRESSED-DATA OPTIONS:\N");
	//	mainUsageCompressedDataOptions(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "SEPARATOR OPTIONS:\n");
	//	mainUsageSeparatorOptions(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "RELEVANT TO CSV/CSV-LITE INPUT ONLY:\n");
	//	mainUsageCsvOptions(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "DOUBLE-QUOTING FOR CSV OUTPUT:\n");
	//	mainUsageDoubleQuoting(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "NUMERICAL FORMATTING:\n");
	//	mainUsageNumericalFormatting(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	//	fmt.Fprintf(o, "OTHER OPTIONS:\n");
	//	mainUsageOtherOptions(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	fmt.Fprintf(o, "THEN-CHAINING:\n")
	mainUsageThenChaining(o, argv0)
	fmt.Fprintf(o, "\n")

	//	fmt.Fprintf(o, "AUXILIARY COMMANDS:\N");
	//	mainUsageAuxents(o, argv0);
	//	fmt.Fprintf(o, "\n");
	//
	fmt.Fprintf(o, "SEE ALSO:\n")
	mainUsageSeeAlso(o, argv0)
}

func mainUsageSynopsis(o *os.File, argv0 string) {
	fmt.Fprintf(o, "Usage: %s [I/O options] {verb} [verb-dependent options ...] {zero or more file names}\n", argv0)
}

func mainUsageExamples(o *os.File, argv0 string, leader string) {
	fmt.Fprintf(o, "%s%s --csv cut -f hostname,uptime mydata.csv\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --tsv --rs lf filter '$status != \"down\" && $upsec >= 10000' *.tsv\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.dat\n", leader, argv0)
	fmt.Fprintf(o, "%sgrep -v '^#' /etc/group | %s --ifs : --nidx --opprint label group,pass,gid,member then sort -f group\n", leader, argv0)
	fmt.Fprintf(o, "%s%s join -j account_id -f accounts.dat then group-by account_name balances.dat\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --json put '$attr = sub($attr, \"([0-9]+)_([0-9]+)_.*\", \"\\1:\\2\")' data/*.json\n", leader, argv0)
	fmt.Fprintf(o, "%s%s stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*\n", leader, argv0)
	fmt.Fprintf(o, "%s%s stats2 -a linreg-pca -f u,v -g shape data/*\n", leader, argv0)
	fmt.Fprintf(o, "%s%s put -q '@sum[$a][$b] += $x; end {emit @sum, \"a\", \"b\"}' data/*\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from estimates.tbl put '\n", leader, argv0)
	fmt.Fprintf(o, "  for (k,v in $*) {\n")
	fmt.Fprintf(o, "    if (is_numeric(v) && k =~ \"^[t-z].*$\") {\n")
	fmt.Fprintf(o, "      $sum += v; $count += 1\n")
	fmt.Fprintf(o, "    }\n")
	fmt.Fprintf(o, "  }\n")
	fmt.Fprintf(o, "  $mean = $sum / $count # no assignment if count unset'\n")
	fmt.Fprintf(o, "%s%s --from infile.dat put -f analyze.mlr\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from infile.dat put 'tee > \"./taps/data-\".$a.\"-\".$b, $*'\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from infile.dat put 'tee | \"gzip > ./taps/data-\".$a.\"-\".$b.\".gz\", $*'\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from infile.dat put -q '@v=$*; dump | \"jq .[]\"'\n", leader, argv0)
	fmt.Fprintf(o, "%s%s --from infile.dat put  '(NR %% 1000 == 0) { print > os.Stderr, \"Checkpoint \".NR}'\n",
		leader, argv0)
}

func mainUsageHelpOptions(o *os.File, argv0 string) {
	fmt.Fprintf(o, "  -h or --help                 Show this message.\n")
	fmt.Fprintf(o, "  --version                    Show the software version.\n")
	fmt.Fprintf(o, "  {verb name} --help           Show verb-specific help.\n")
	fmt.Fprintf(o, "  --help-all-verbs             Show help on all verbs.\n")
	fmt.Fprintf(o, "  -l or --list-all-verbs       List only verb names.\n")
	fmt.Fprintf(o, "  -L                           List only verb names, one per line.\n")
	fmt.Fprintf(o, "  -f or --help-all-functions   Show help on all built-in functions.\n")
	fmt.Fprintf(o, "  -F                           Show a bare listing of built-in functions by name.\n")
	fmt.Fprintf(o, "  -k or --help-all-keywords    Show help on all keywords.\n")
	fmt.Fprintf(o, "  -K                           Show a bare listing of keywords by name.\n")
}

func mainUsageMlrrc(o *os.File, argv0 string) {
	fmt.Fprintf(o, "You can set up personal defaults via a $HOME/.mlrrc and/or ./.mlrrc.\n")
	fmt.Fprintf(o, "For example, if you usually process CSV, then you can put \"--csv\" in your .mlrrc file\n")
	fmt.Fprintf(o, "and that will be the default input/output format unless otherwise specified on the command line.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "The .mlrrc file format is one \"--flag\" or \"--option value\" per line, with the leading \"--\" optional.\n")
	fmt.Fprintf(o, "Hash-style comments and blank lines are ignored.\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "Sample .mlrrc:\n")
	fmt.Fprintf(o, "# Input and output formats are CSV by default (unless otherwise specified\n")
	fmt.Fprintf(o, "# on the mlr command line):\n")
	fmt.Fprintf(o, "csv\n")
	fmt.Fprintf(o, "# These are no-ops for CSV, but when I do use JSON output, I want these\n")
	fmt.Fprintf(o, "# pretty-printing options to be used:\n")
	fmt.Fprintf(o, "jvstack\n")
	fmt.Fprintf(o, "jlistwrap\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "How to specify location of .mlrrc:\n")
	fmt.Fprintf(o, "* If $MLRRC is set:\n")
	fmt.Fprintf(o, "  o If its value is \"__none__\" then no .mlrrc files are processed.\n")
	fmt.Fprintf(o, "  o Otherwise, its value (as a filename) is loaded and processed. If there are syntax\n")
	fmt.Fprintf(o, "    errors, they abort mlr with a usage message (as if you had mistyped something on the\n")
	fmt.Fprintf(o, "    command line). If the file can't be loaded at all, though, it is silently skipped.\n")
	fmt.Fprintf(o, "  o Any .mlrrc in your home directory or current directory is ignored whenever $MLRRC is\n")
	fmt.Fprintf(o, "    set in the environment.\n")
	fmt.Fprintf(o, "* Otherwise:\n")
	fmt.Fprintf(o, "  o If $HOME/.mlrrc exists, it's then processed as above.\n")
	fmt.Fprintf(o, "  o If ./.mlrrc exists, it's then also processed as above.\n")
	fmt.Fprintf(o, "  (I.e. current-directory .mlrrc defaults are stacked over home-directory .mlrrc defaults.)\n")
	fmt.Fprintf(o, "\n")
	fmt.Fprintf(o, "See also:\n")
	fmt.Fprintf(o, "https://miller.readthedocs.io/en/latest/customization.html\n")
}

//func mainUsageFunctions(o *os.File, argv0 string, char* leader) {
//	fmgr_t* pfmgr = fmgr_alloc();
//	fmgr_list_functions(pfmgr, o, leader);
//	fmgr_free(pfmgr, nil);
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "Please use \"%s --help-function {function name}\" for function-specific help.\n", argv0);
//}

func mainUsageDataFormatExamples(o *os.File, argv0 string) {
	fmt.Fprintf(o,
		`  DKVP: delimited key-value pairs (Miller default format)
  +---------------------+
  | apple=1,bat=2,cog=3 | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
  | dish=7,egg=8,flint  | Record 2: "dish" => "7", "egg" => "8", "3" => "flint"
  +---------------------+

  NIDX: implicitly numerically indexed (Unix-toolkit style)
  +---------------------+
  | the quick brown     | Record 1: "1" => "the", "2" => "quick", "3" => "brown"
  | fox jumped          | Record 2: "1" => "fox", "2" => "jumped"
  +---------------------+

  CSV/CSV-lite: comma-separated values with separate header line
  +---------------------+
  | apple,bat,cog       |
  | 1,2,3               | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
  | 4,5,6               | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
  +---------------------+

  Tabular JSON: nested objects are supported, although arrays within them are not:
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

  XTAB: pretty-printed transposed tabular
  +---------------------+
  | apple 1             | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
  | bat   2             |
  | cog   3             |
  |                     |
  | dish 7              | Record 2: "dish" => "7", "egg" => "8"
  | egg  8              |
  +---------------------+

  Markdown tabular (supported for output only):
  +-----------------------+
  | | apple | bat | cog | |
  | | ---   | --- | --- | |
  | | 1     | 2   | 3   | | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
  | | 4     | 5   | 6   | | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
  +-----------------------+
`)
}

func mainUsageDataFormatOptions(o *os.File, argv0 string) {
	fmt.Fprintln(o,
		`
	  --idkvp   --odkvp   --dkvp      Delimited key-value pairs, e.g "a=1,b=2"
	                                  (this is Miller's default format).

	  --inidx   --onidx   --nidx      Implicitly-integer-indexed fields
	                                  (Unix-toolkit style).
	  -T                              Synonymous with "--nidx --fs tab".

	  --icsv    --ocsv    --csv       Comma-separated value (or tab-separated
	                                  with --fs tab, etc.)

	  --itsv    --otsv    --tsv       Keystroke-savers for "--icsv --ifs tab",
	                                  "--ocsv --ofs tab", "--csv --fs tab".
	  --iasv    --oasv    --asv       Similar but using ASCII FS %s and RS %s\n",
		ASV_FS_FOR_HELP, ASV_RS_FOR_HELP);
	  --iusv    --ousv    --usv       Similar but using Unicode FS %s\n",
		USV_FS_FOR_HELP);
	                                  and RS %s\n",
		USV_RS_FOR_HELP);

	  --icsvlite --ocsvlite --csvlite Comma-separated value (or tab-separated
	                                  with --fs tab, etc.). The 'lite' CSV does not handle
	                                  RFC-CSV double-quoting rules; is slightly faster;
	                                  and handles heterogeneity in the input stream via
	                                  empty newline followed by new header line. See also
	                                  http://johnkerl.org/miller/doc/file-formats.html#CSV/TSV/etc.

	  --itsvlite --otsvlite --tsvlite Keystroke-savers for "--icsvlite --ifs tab",
	                                  "--ocsvlite --ofs tab", "--csvlite --fs tab".
	  -t                              Synonymous with --tsvlite.
	  --iasvlite --oasvlite --asvlite Similar to --itsvlite et al. but using ASCII FS %s and RS %s\n",
		ASV_FS_FOR_HELP, ASV_RS_FOR_HELP);
	  --iusvlite --ousvlite --usvlite Similar to --itsvlite et al. but using Unicode FS %s\n",
		USV_FS_FOR_HELP);
	                                  and RS %s\n",
		USV_RS_FOR_HELP);

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
	    --json-map-arrays-on-input    JSON arrays are unmillerable. --json-map-arrays-on-input
	    --json-skip-arrays-on-input   is the default: arrays are converted to integer-indexed
	    --json-fatal-arrays-on-input  maps. The other two options cause them to be skipped, or
	                                  to be treated as errors.  Please use the jq tool for full
	                                  JSON (pre)processing.
	                      --jvstack   Put one key-value pair per line for JSON
	                                  output.
	                --jsonx --ojsonx  Keystroke-savers for --json --jvstack
	                --jsonx --ojsonx  and --ojson --jvstack, respectively.
	                      --jlistwrap Wrap JSON output in outermost [ ].
	                    --jknquoteint Do not quote non-string map keys in JSON output.
	                     --jvquoteall Quote map values in JSON output, even if they're
	                                  numeric.
	              --jflatsep {string} Separator for flattening multi-level JSON keys,
	                                  e.g. '{"a":{"b":3}}' becomes a:b => 3 for
	                                  non-JSON formats. Defaults to %s.\n",
		DEFAULT_JSON_FLATTEN_SEPARATOR);

	  -p is a keystroke-saver for --nidx --fs space --repifs

	  Examples: --csv for CSV-formatted input and output; --idkvp --opprint for
	  DKVP-formatted input and pretty-printed output.

	  Please use --iformat1 --oformat2 rather than --format1 --oformat2.
	  The latter sets up input and output flags for format1, not all of which
	  are overridden in all cases by setting output format to format2.

`)
}

//func mainUsageCommentsInData(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --skip-comments                 Ignore commented lines (prefixed by \"%s\")\n",
//		DEFAULT_COMMENT_STRING);
//	fmt.Fprintf(o, "                                  within the input.\n");
//	fmt.Fprintf(o, "  --skip-comments-with {string}   Ignore commented lines within input, with\n");
//	fmt.Fprintf(o, "                                  specified prefix.\n");
//	fmt.Fprintf(o, "  --pass-comments                 Immediately print commented lines (prefixed by \"%s\")\n",
//		DEFAULT_COMMENT_STRING);
//	fmt.Fprintf(o, "                                  within the input.\n");
//	fmt.Fprintf(o, "  --pass-comments-with {string}   Immediately print commented lines within input, with\n");
//	fmt.Fprintf(o, "                                  specified prefix.\n");
//	fmt.Fprintf(o, "Notes:\n");
//	fmt.Fprintf(o, "* Comments are only honored at the start of a line.\n");
//	fmt.Fprintf(o, "* In the absence of any of the above four options, comments are data like\n");
//	fmt.Fprintf(o, "  any other text.\n");
//	fmt.Fprintf(o, "* When pass-comments is used, comment lines are written to standard output\n");
//	fmt.Fprintf(o, "  immediately upon being read; they are not part of the record stream.\n");
//	fmt.Fprintf(o, "  Results may be counterintuitive. A suggestion is to place comments at the\n");
//	fmt.Fprintf(o, "  start of data files.\n");
//}

func mainUsageFormatConversionKeystrokeSaverOptions(o *os.File, argv0 string) {
	fmt.Fprintf(o, "As keystroke-savers for format-conversion you may use the following:\n")
	fmt.Fprintf(o, "        --c2t --c2d --c2n --c2j --c2x --c2p --c2m\n")
	fmt.Fprintf(o, "  --t2c       --t2d --t2n --t2j --t2x --t2p --t2m\n")
	fmt.Fprintf(o, "  --d2c --d2t       --d2n --d2j --d2x --d2p --d2m\n")
	fmt.Fprintf(o, "  --n2c --n2t --n2d       --n2j --n2x --n2p --n2m\n")
	fmt.Fprintf(o, "  --j2c --j2t --j2d --j2n       --j2x --j2p --j2m\n")
	fmt.Fprintf(o, "  --x2c --x2t --x2d --x2n --x2j       --x2p --x2m\n")
	fmt.Fprintf(o, "  --p2c --p2t --p2d --p2n --p2j --p2x       --p2m\n")
	fmt.Fprintf(o, "The letters c t d n j x p m refer to formats CSV, TSV, DKVP, NIDX, JSON, XTAB,\n")
	fmt.Fprintf(o, "PPRINT, and markdown, respectively. Note that markdown format is available for\n")
	fmt.Fprintf(o, "output only.\n")
}

//func mainUsageCompressedDataOptions(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --prepipe {command} This allows Miller to handle compressed inputs. You can do\n");
//	fmt.Fprintf(o, "  without this for single input files, e.g. \"gunzip < myfile.csv.gz | %s ...\".\n",
//		argv0);
//	fmt.Fprintf(o, "  However, when multiple input files are present, between-file separations are\n");
//	fmt.Fprintf(o, "  lost; also, the FILENAME variable doesn't iterate. Using --prepipe you can\n");
//	fmt.Fprintf(o, "  specify an action to be taken on each input file. This pre-pipe command must\n");
//	fmt.Fprintf(o, "  be able to read from standard input; it will be invoked with\n");
//	fmt.Fprintf(o, "    {command} < {filename}.\n");
//	fmt.Fprintf(o, "  Examples:\n");
//	fmt.Fprintf(o, "    %s --prepipe 'gunzip'\n", argv0);
//	fmt.Fprintf(o, "    %s --prepipe 'zcat -cf'\n", argv0);
//	fmt.Fprintf(o, "    %s --prepipe 'xz -cd'\n", argv0);
//	fmt.Fprintf(o, "    %s --prepipe cat\n", argv0);
//	fmt.Fprintf(o, "  Note that this feature is quite general and is not limited to decompression\n");
//	fmt.Fprintf(o, "  utilities. You can use it to apply per-file filters of your choice.\n");
//	fmt.Fprintf(o, "  For output compression (or other) utilities, simply pipe the output:\n");
//	fmt.Fprintf(o, "    %s ... | {your compression command}\n", argv0);
//}

//func mainUsageSeparatorOptions(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --rs     --irs     --ors              Record separators, e.g. 'lf' or '\\r\\n'\n");
//	fmt.Fprintf(o, "  --fs     --ifs     --ofs  --repifs    Field separators, e.g. comma\n");
//	fmt.Fprintf(o, "  --ps     --ips     --ops              Pair separators, e.g. equals sign\n");
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "  Notes about line endings:\n");
//	fmt.Fprintf(o, "  * Default line endings (--irs and --ors) are \"auto\" which means autodetect from\n");
//	fmt.Fprintf(o, "    the input file format, as long as the input file(s) have lines ending in either\n");
//	fmt.Fprintf(o, "    LF (also known as linefeed, '\\n', 0x0a, Unix-style) or CRLF (also known as\n");
//	fmt.Fprintf(o, "    carriage-return/linefeed pairs, '\\r\\n', 0x0d 0x0a, Windows style).\n");
//	fmt.Fprintf(o, "  * If both irs and ors are auto (which is the default) then LF input will lead to LF\n");
//	fmt.Fprintf(o, "    output and CRLF input will lead to CRLF output, regardless of the platform you're\n");
//	fmt.Fprintf(o, "    running on.\n");
//	fmt.Fprintf(o, "  * The line-ending autodetector triggers on the first line ending detected in the input\n");
//	fmt.Fprintf(o, "    stream. E.g. if you specify a CRLF-terminated file on the command line followed by an\n");
//	fmt.Fprintf(o, "    LF-terminated file then autodetected line endings will be CRLF.\n");
//	fmt.Fprintf(o, "  * If you use --ors {something else} with (default or explicitly specified) --irs auto\n");
//	fmt.Fprintf(o, "    then line endings are autodetected on input and set to what you specify on output.\n");
//	fmt.Fprintf(o, "  * If you use --irs {something else} with (default or explicitly specified) --ors auto\n");
//	fmt.Fprintf(o, "    then the output line endings used are LF on Unix/Linux/BSD/MacOSX, and CRLF on Windows.\n");
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "  Notes about all other separators:\n");
//	fmt.Fprintf(o, "  * IPS/OPS are only used for DKVP and XTAB formats, since only in these formats\n");
//	fmt.Fprintf(o, "    do key-value pairs appear juxtaposed.\n");
//	fmt.Fprintf(o, "  * IRS/ORS are ignored for XTAB format. Nominally IFS and OFS are newlines;\n");
//	fmt.Fprintf(o, "    XTAB records are separated by two or more consecutive IFS/OFS -- i.e.\n");
//	fmt.Fprintf(o, "    a blank line. Everything above about --irs/--ors/--rs auto becomes --ifs/--ofs/--fs\n");
//	fmt.Fprintf(o, "    auto for XTAB format. (XTAB's default IFS/OFS are \"auto\".)\n");
//	fmt.Fprintf(o, "  * OFS must be single-character for PPRINT format. This is because it is used\n");
//	fmt.Fprintf(o, "    with repetition for alignment; multi-character separators would make\n");
//	fmt.Fprintf(o, "    alignment impossible.\n");
//	fmt.Fprintf(o, "  * OPS may be multi-character for XTAB format, in which case alignment is\n");
//	fmt.Fprintf(o, "    disabled.\n");
//	fmt.Fprintf(o, "  * TSV is simply CSV using tab as field separator (\"--fs tab\").\n");
//	fmt.Fprintf(o, "  * FS/PS are ignored for markdown format; RS is used.\n");
//	fmt.Fprintf(o, "  * All FS and PS options are ignored for JSON format, since they are not relevant\n");
//	fmt.Fprintf(o, "    to the JSON format.\n");
//	fmt.Fprintf(o, "  * You can specify separators in any of the following ways, shown by example:\n");
//	fmt.Fprintf(o, "    - Type them out, quoting as necessary for shell escapes, e.g.\n");
//	fmt.Fprintf(o, "      \"--fs '|' --ips :\"\n");
//	fmt.Fprintf(o, "    - C-style escape sequences, e.g. \"--rs '\\r\\n' --fs '\\t'\".\n");
//	fmt.Fprintf(o, "    - To avoid backslashing, you can use any of the following names:\n");
//	fmt.Fprintf(o, "     ");
//	lhmss_t* pmap = get_desc_to_chars_map();
//	for (lhmsse_t* pe = pmap.phead; pe != nil; pe = pe.pnext) {
//		fmt.Fprintf(o, " %s", pe.key);
//	}
//	fmt.Fprintf(o, "\n");
//	fmt.Fprintf(o, "  * Default separators by format:\n");
//	fmt.Fprintf(o, "      %-12s %-8s %-8s %s\n", "File format", "RS", "FS", "PS");
//	lhmss_t* default_rses = get_default_rses();
//	lhmss_t* default_fses = get_default_fses();
//	lhmss_t* default_pses = get_default_pses();
//	for (lhmsse_t* pe = default_rses.phead; pe != nil; pe = pe.pnext) {
//		char* filefmt = pe.key;
//		char* rs = pe.value;
//		char* fs = lhmss_get(default_fses, filefmt);
//		char* ps = lhmss_get(default_pses, filefmt);
//		fmt.Fprintf(o, "      %-12s %-8s %-8s %s\n", filefmt, rebackslash(rs), rebackslash(fs), rebackslash(ps));
//	}
//}

//func mainUsageCsvOptions(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --implicit-csv-header Use 1,2,3,... as field labels, rather than from line 1\n");
//	fmt.Fprintf(o, "                     of input files. Tip: combine with \"label\" to recreate\n");
//	fmt.Fprintf(o, "                     missing headers.\n");
//	fmt.Fprintf(o, "  --allow-ragged-csv-input|--ragged If a data line has fewer fields than the header line,\n");
//	fmt.Fprintf(o, "                     fill remaining keys with empty string. If a data line has more\n");
//	fmt.Fprintf(o, "                     fields than the header line, use integer field labels as in\n");
//	fmt.Fprintf(o, "                     the implicit-header case.\n");
//	fmt.Fprintf(o, "  --headerless-csv-output   Print only CSV data lines.\n");
//	fmt.Fprintf(o, "  -N                 Keystroke-saver for --implicit-csv-header --headerless-csv-output.\n");
//}

//func mainUsageDoubleQuoting(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --quote-all        Wrap all fields in double quotes\n");
//	fmt.Fprintf(o, "  --quote-none       Do not wrap any fields in double quotes, even if they have\n");
//	fmt.Fprintf(o, "                     OFS or ORS in them\n");
//	fmt.Fprintf(o, "  --quote-minimal    Wrap fields in double quotes only if they have OFS or ORS\n");
//	fmt.Fprintf(o, "                     in them (default)\n");
//	fmt.Fprintf(o, "  --quote-numeric    Wrap fields in double quotes only if they have numbers\n");
//	fmt.Fprintf(o, "                     in them\n");
//	fmt.Fprintf(o, "  --quote-original   Wrap fields in double quotes if and only if they were\n");
//	fmt.Fprintf(o, "                     quoted on input. This isn't sticky for computed fields:\n");
//	fmt.Fprintf(o, "                     e.g. if fields a and b were quoted on input and you do\n");
//	fmt.Fprintf(o, "                     \"put '$c = $a . $b'\" then field c won't inherit a or b's\n");
//	fmt.Fprintf(o, "                     was-quoted-on-input flag.\n");
//}

//func mainUsageNumericalFormatting(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --ofmt {format}    E.g. %%.18lf, %%.0lf. Please use sprintf-style codes for\n");
//	fmt.Fprintf(o, "                     double-precision. Applies to verbs which compute new\n");
//	fmt.Fprintf(o, "                     values, e.g. put, stats1, stats2. See also the fmtnum\n");
//	fmt.Fprintf(o, "                     function within mlr put (mlr --help-all-functions).\n");
//	fmt.Fprintf(o, "                     Defaults to %s.\n", DEFAULT_OFMT);
//}

//func mainUsageOtherOptions(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "  --seed {n} with n of the form 12345678 or 0xcafefeed. For put/filter\n");
//	fmt.Fprintf(o, "                     urand()/urandint()/urand32().\n");
//	fmt.Fprintf(o, "  --nr-progress-mod {m}, with m a positive integer: print filename and record\n");
//	fmt.Fprintf(o, "                     count to os.Stderr every m input records.\n");
//	fmt.Fprintf(o, "  --from {filename}  Use this to specify an input file before the verb(s),\n");
//	fmt.Fprintf(o, "                     rather than after. May be used more than once. Example:\n");
//	fmt.Fprintf(o, "                     \"%s --from a.dat --from b.dat cat\" is the same as\n", argv0);
//	fmt.Fprintf(o, "                     \"%s cat a.dat b.dat\".\n", argv0);
//	fmt.Fprintf(o, "  -n                 Process no input files, nor standard input either. Useful\n");
//	fmt.Fprintf(o, "                     for %s put with begin/end statements only. (Same as --from\n", argv0);
//	fmt.Fprintf(o, "                     /dev/null.) Also useful in \"%s -n put -v '...'\" for\n", argv0);
//	fmt.Fprintf(o, "                     analyzing abstract syntax trees (if that's your thing).\n");
//	fmt.Fprintf(o, "  -I                 Process files in-place. For each file name on the command\n");
//	fmt.Fprintf(o, "                     line, output is written to a temp file in the same\n");
//	fmt.Fprintf(o, "                     directory, which is then renamed over the original. Each\n");
//	fmt.Fprintf(o, "                     file is processed in isolation: if the output format is\n");
//	fmt.Fprintf(o, "                     CSV, CSV headers will be present in each output file;\n");
//	fmt.Fprintf(o, "                     statistics are only over each file's own records; and so on.\n");
//}

func mainUsageThenChaining(o *os.File, argv0 string) {
	fmt.Fprintf(o, "Output of one verb may be chained as input to another using \"then\", e.g.\n")
	fmt.Fprintf(o, "  %s stats1 -a min,mean,max -f flag,u,v -g color then sort -f color\n", argv0)
}

//func mainUsageAuxents(o *os.File, argv0 string) {
//	fmt.Fprintf(o, "Miller has a few otherwise-standalone executables packaged within it.\n");
//	fmt.Fprintf(o, "They do not participate in any other parts of Miller.\n");
//	show_aux_entries(o);
//}

func mainUsageSeeAlso(o *os.File, argv0 string) {
	fmt.Fprintf(o, "For more information please see http://johnkerl.org/miller/doc and/or\n")
	fmt.Fprintf(o, "http://github.com/johnkerl/miller.")
	fmt.Fprintf(o, " This is Miller version %s.\n", version.STRING)
}

func usageUnrecognizedVerb(argv0 string, arg string) {
	fmt.Fprintf(os.Stderr, "%s: option \"%s\" not recognized.\n", argv0, arg)
	fmt.Fprintf(os.Stderr, "Please run \"%s --help\" for usage information.\n", argv0)
	os.Exit(1)
}
