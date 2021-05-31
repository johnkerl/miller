..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Manpage
================================================================

This is simply a copy of what you should see on running **man mlr** at a command prompt, once Miller is installed on your system.

.. code-block::

    MILLER(1)							     MILLER(1)
    
    
    
    NAME
           miller - like awk, sed, cut, join, and sort for name-indexed data such
           as CSV and tabular JSON.
    
    SYNOPSIS
           Usage: mlr [I/O options] {verb} [verb-dependent options ...] {zero or
           more file names}
    
    
    DESCRIPTION
           Miller operates on key-value-pair data while the familiar Unix tools
           operate on integer-indexed fields: if the natural data structure for
           the latter is the array, then Miller's natural data structure is the
           insertion-ordered hash map.  This encompasses a variety of data
           formats, including but not limited to the familiar CSV, TSV, and JSON.
           (Miller can handle positionally-indexed data as a special case.) This
           manpage documents Miller v5.10.1.
    
    EXAMPLES
       COMMAND-LINE SYNTAX
           mlr --csv cut -f hostname,uptime mydata.csv
           mlr --tsv --rs lf filter '$status != "down" && $upsec >= 10000' *.tsv
           mlr --nidx put '$sum = $7 < 0.0 ? 3.5 : $7 + 2.1*$8' *.dat
           grep -v '^#' /etc/group | mlr --ifs : --nidx --opprint label group,pass,gid,member then sort -f group
           mlr join -j account_id -f accounts.dat then group-by account_name balances.dat
           mlr --json put '$attr = sub($attr, "([0-9]+)_([0-9]+)_.*", "\1:\2")' data/*.json
           mlr stats1 -a min,mean,max,p10,p50,p90 -f flag,u,v data/*
           mlr stats2 -a linreg-pca -f u,v -g shape data/*
           mlr put -q '@sum[$a][$b] += $x; end {emit @sum, "a", "b"}' data/*
           mlr --from estimates.tbl put '
    	 for (k,v in $*) {
    	   if (is_numeric(v) && k =~ "^[t-z].*$") {
    	     $sum += v; $count += 1
    	   }
    	 }
    	 $mean = $sum / $count # no assignment if count unset'
           mlr --from infile.dat put -f analyze.mlr
           mlr --from infile.dat put 'tee > "./taps/data-".$a."-".$b, $*'
           mlr --from infile.dat put 'tee | "gzip > ./taps/data-".$a."-".$b.".gz", $*'
           mlr --from infile.dat put -q '@v=$*; dump | "jq .[]"'
           mlr --from infile.dat put  '(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'
    
       DATA FORMATS
    	 DKVP: delimited key-value pairs (Miller default format)
    	 +---------------------+
    	 | apple=1,bat=2,cog=3 | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
    	 | dish=7,egg=8,flint  | Record 2: "dish" => "7", "egg" => "8", "3" => "flint"
    	 +---------------------+
    
    	 NIDX: implicitly numerically indexed (Unix-toolkit style)
    	 +---------------------+
    	 | the quick brown     | Record 1: "1" => "the", "2" => "quick", "3" => "brown"
    	 | fox jumped	       | Record 2: "1" => "fox", "2" => "jumped"
    	 +---------------------+
    
    	 CSV/CSV-lite: comma-separated values with separate header line
    	 +---------------------+
    	 | apple,bat,cog       |
    	 | 1,2,3	       | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
    	 | 4,5,6	       | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
    	 +---------------------+
    
    	 Tabular JSON: nested objects are supported, although arrays within them are not:
    	 +---------------------+
    	 | {		       |
    	 |  "apple": 1,        | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
    	 |  "bat": 2,	       |
    	 |  "cog": 3	       |
    	 | }		       |
    	 | {		       |
    	 |   "dish": {	       | Record 2: "dish:egg" => "7", "dish:flint" => "8", "garlic" => ""
    	 |     "egg": 7,       |
    	 |     "flint": 8      |
    	 |   }, 	       |
    	 |   "garlic": ""      |
    	 | }		       |
    	 +---------------------+
    
    	 PPRINT: pretty-printed tabular
    	 +---------------------+
    	 | apple bat cog       |
    	 | 1	 2   3	       | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
    	 | 4	 5   6	       | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
    	 +---------------------+
    
    	 XTAB: pretty-printed transposed tabular
    	 +---------------------+
    	 | apple 1	       | Record 1: "apple" => "1", "bat" => "2", "cog" => "3"
    	 | bat	 2	       |
    	 | cog	 3	       |
    	 |		       |
    	 | dish 7	       | Record 2: "dish" => "7", "egg" => "8"
    	 | egg	8	       |
    	 +---------------------+
    
    	 Markdown tabular (supported for output only):
    	 +-----------------------+
    	 | | apple | bat | cog | |
    	 | | ---   | --- | --- | |
    	 | | 1	   | 2	 | 3   | | Record 1: "apple => "1", "bat" => "2", "cog" => "3"
    	 | | 4	   | 5	 | 6   | | Record 2: "apple" => "4", "bat" => "5", "cog" => "6"
    	 +-----------------------+
    
    OPTIONS
           In the following option flags, the version with "i" designates the
           input stream, "o" the output stream, and the version without prefix
           sets the option for both input and output stream. For example: --irs
           sets the input record separator, --ors the output record separator, and
           --rs sets both the input and output separator to the given value.
    
       HELP OPTIONS
    	 -h or --help		      Show this message.
    	 --version		      Show the software version.
    	 {verb name} --help	      Show verb-specific help.
    	 --help-all-verbs	      Show help on all verbs.
    	 -l or --list-all-verbs       List only verb names.
    	 -L			      List only verb names, one per line.
    	 -f or --help-all-functions   Show help on all built-in functions.
    	 -F			      Show a bare listing of built-in functions by name.
    	 -k or --help-all-keywords    Show help on all keywords.
    	 -K			      Show a bare listing of keywords by name.
    
       VERB LIST
    	altkv bar bootstrap cat check clean-whitespace count count-distinct
    	count-similar cut decimate fill-down filter format-values fraction grep
    	group-by group-like having-fields head histogram join label least-frequent
    	merge-fields most-frequent nest nothing put regularize remove-empty-columns
    	rename reorder repeat reshape sample sec2gmt sec2gmtdate seqgen shuffle
    	skip-trivial-records sort sort-within-records stats1 stats2 step tac tail tee
    	top uniq unsparsify
    
       FUNCTION LIST
    	+ + - - * / // .+ .+ .- .- .* ./ .// % ** | ^ & ~ << >> bitcount == != =~ !=~
    	> >= < <= && || ^^ ! ? : . gsub regextract regextract_or_else strlen sub ssub
    	substr tolower toupper truncate capitalize lstrip rstrip strip
    	collapse_whitespace clean_whitespace system abs acos acosh asin asinh atan
    	atan2 atanh cbrt ceil cos cosh erf erfc exp expm1 floor invqnorm log log10
    	log1p logifit madd max mexp min mmul msub pow qnorm round roundm sgn sin sinh
    	sqrt tan tanh urand urandrange urand32 urandint dhms2fsec dhms2sec fsec2dhms
    	fsec2hms gmt2sec localtime2sec hms2fsec hms2sec sec2dhms sec2gmt sec2gmt
    	sec2gmtdate sec2localtime sec2localtime sec2localdate sec2hms strftime
    	strftime_local strptime strptime_local systime is_absent is_bool is_boolean
    	is_empty is_empty_map is_float is_int is_map is_nonempty_map is_not_empty
    	is_not_map is_not_null is_null is_numeric is_present is_string
    	asserting_absent asserting_bool asserting_boolean asserting_empty
    	asserting_empty_map asserting_float asserting_int asserting_map
    	asserting_nonempty_map asserting_not_empty asserting_not_map
    	asserting_not_null asserting_null asserting_numeric asserting_present
    	asserting_string boolean float fmtnum hexfmt int string typeof depth haskey
    	joink joinkv joinv leafcount length mapdiff mapexcept mapselect mapsum splitkv
    	splitkvx splitnv splitnvx
    
           Please use "mlr --help-function {function name}" for function-specific help.
    
       I/O FORMATTING
    	 --idkvp   --odkvp   --dkvp	 Delimited key-value pairs, e.g "a=1,b=2"
    					 (this is Miller's default format).
    
    	 --inidx   --onidx   --nidx	 Implicitly-integer-indexed fields
    					 (Unix-toolkit style).
    	 -T				 Synonymous with "--nidx --fs tab".
    
    	 --icsv    --ocsv    --csv	 Comma-separated value (or tab-separated
    					 with --fs tab, etc.)
    
    	 --itsv    --otsv    --tsv	 Keystroke-savers for "--icsv --ifs tab",
    					 "--ocsv --ofs tab", "--csv --fs tab".
    	 --iasv    --oasv    --asv	 Similar but using ASCII FS 0x1f and RS 0x1e
    	 --iusv    --ousv    --usv	 Similar but using Unicode FS U+241F (UTF-8 0xe2909f)
    					 and RS U+241E (UTF-8 0xe2909e)
    
    	 --icsvlite --ocsvlite --csvlite Comma-separated value (or tab-separated
    					 with --fs tab, etc.). The 'lite' CSV does not handle
    					 RFC-CSV double-quoting rules; is slightly faster;
    					 and handles heterogeneity in the input stream via
    					 empty newline followed by new header line. See also
    					 http://johnkerl.org/miller/doc/file-formats.html#CSV/TSV/etc.
    
    	 --itsvlite --otsvlite --tsvlite Keystroke-savers for "--icsvlite --ifs tab",
    					 "--ocsvlite --ofs tab", "--csvlite --fs tab".
    	 -t				 Synonymous with --tsvlite.
    	 --iasvlite --oasvlite --asvlite Similar to --itsvlite et al. but using ASCII FS 0x1f and RS 0x1e
    	 --iusvlite --ousvlite --usvlite Similar to --itsvlite et al. but using Unicode FS U+241F (UTF-8 0xe2909f)
    					 and RS U+241E (UTF-8 0xe2909e)
    
    	 --ipprint --opprint --pprint	 Pretty-printed tabular (produces no
    					 output until all input is in).
    			     --right	 Right-justifies all fields for PPRINT output.
    			     --barred	 Prints a border around PPRINT output
    					 (only available for output).
    
    		   --omd		 Markdown-tabular (only available for output).
    
    	 --ixtab   --oxtab   --xtab	 Pretty-printed vertical-tabular.
    			     --xvright	 Right-justifies values for XTAB format.
    
    	 --ijson   --ojson   --json	 JSON tabular: sequence or list of one-level
    					 maps: {...}{...} or [{...},{...}].
    	   --json-map-arrays-on-input	 JSON arrays are unmillerable. --json-map-arrays-on-input
    	   --json-skip-arrays-on-input	 is the default: arrays are converted to integer-indexed
    	   --json-fatal-arrays-on-input  maps. The other two options cause them to be skipped, or
    					 to be treated as errors.  Please use the jq tool for full
    					 JSON (pre)processing.
    			     --jvstack	 Put one key-value pair per line for JSON
    					 output.
    		       --jsonx --ojsonx  Keystroke-savers for --json --jvstack
    		       --jsonx --ojsonx  and --ojson --jvstack, respectively.
    			     --jlistwrap Wrap JSON output in outermost [ ].
    			   --jknquoteint Do not quote non-string map keys in JSON output.
    			    --jvquoteall Quote map values in JSON output, even if they're
    					 numeric.
    		     --jflatsep {string} Separator for flattening multi-level JSON keys,
    					 e.g. '{"a":{"b":3}}' becomes a:b => 3 for
    					 non-JSON formats. Defaults to :.
    
    	 -p is a keystroke-saver for --nidx --fs space --repifs
    
    	 Examples: --csv for CSV-formatted input and output; --idkvp --opprint for
    	 DKVP-formatted input and pretty-printed output.
    
    	 Please use --iformat1 --oformat2 rather than --format1 --oformat2.
    	 The latter sets up input and output flags for format1, not all of which
    	 are overridden in all cases by setting output format to format2.
    
       COMMENTS IN DATA
    	 --skip-comments		 Ignore commented lines (prefixed by "#")
    					 within the input.
    	 --skip-comments-with {string}	 Ignore commented lines within input, with
    					 specified prefix.
    	 --pass-comments		 Immediately print commented lines (prefixed by "#")
    					 within the input.
    	 --pass-comments-with {string}	 Immediately print commented lines within input, with
    					 specified prefix.
           Notes:
           * Comments are only honored at the start of a line.
           * In the absence of any of the above four options, comments are data like
    	 any other text.
           * When pass-comments is used, comment lines are written to standard output
    	 immediately upon being read; they are not part of the record stream.
    	 Results may be counterintuitive. A suggestion is to place comments at the
    	 start of data files.
    
       FORMAT-CONVERSION KEYSTROKE-SAVERS
           As keystroke-savers for format-conversion you may use the following:
    	       --c2t --c2d --c2n --c2j --c2x --c2p --c2m
    	 --t2c	     --t2d --t2n --t2j --t2x --t2p --t2m
    	 --d2c --d2t	   --d2n --d2j --d2x --d2p --d2m
    	 --n2c --n2t --n2d	 --n2j --n2x --n2p --n2m
    	 --j2c --j2t --j2d --j2n       --j2x --j2p --j2m
    	 --x2c --x2t --x2d --x2n --x2j	     --x2p --x2m
    	 --p2c --p2t --p2d --p2n --p2j --p2x	   --p2m
           The letters c t d n j x p m refer to formats CSV, TSV, DKVP, NIDX, JSON, XTAB,
           PPRINT, and markdown, respectively. Note that markdown format is available for
           output only.
    
       COMPRESSED I/O
    	 --prepipe {command} This allows Miller to handle compressed inputs. You can do
    	 without this for single input files, e.g. "gunzip < myfile.csv.gz | mlr ...".
    
    	 However, when multiple input files are present, between-file separations are
    	 lost; also, the FILENAME variable doesn't iterate. Using --prepipe you can
    	 specify an action to be taken on each input file. This pre-pipe command must
    	 be able to read from standard input; it will be invoked with
    	   {command} < {filename}.
    	 Examples:
    	   mlr --prepipe 'gunzip'
    	   mlr --prepipe 'zcat -cf'
    	   mlr --prepipe 'xz -cd'
    	   mlr --prepipe cat
    	   mlr --prepipe-gunzip
    	   mlr --prepipe-zcat
    	 Note that this feature is quite general and is not limited to decompression
    	 utilities. You can use it to apply per-file filters of your choice.
    	 For output compression (or other) utilities, simply pipe the output:
    	   mlr ... | {your compression command}
    
    	 There are shorthands --prepipe-zcat and --prepipe-gunzip which are
    	 valid in .mlrrc files. The --prepipe flag is not valid in .mlrrc
    	 files since that would put execution of the prepipe command under
    	 control of the .mlrrc file.
    
       SEPARATORS
    	 --rs	  --irs     --ors	       Record separators, e.g. 'lf' or '\r\n'
    	 --fs	  --ifs     --ofs  --repifs    Field separators, e.g. comma
    	 --ps	  --ips     --ops	       Pair separators, e.g. equals sign
    
    	 Notes about line endings:
    	 * Default line endings (--irs and --ors) are "auto" which means autodetect from
    	   the input file format, as long as the input file(s) have lines ending in either
    	   LF (also known as linefeed, '\n', 0x0a, Unix-style) or CRLF (also known as
    	   carriage-return/linefeed pairs, '\r\n', 0x0d 0x0a, Windows style).
    	 * If both irs and ors are auto (which is the default) then LF input will lead to LF
    	   output and CRLF input will lead to CRLF output, regardless of the platform you're
    	   running on.
    	 * The line-ending autodetector triggers on the first line ending detected in the input
    	   stream. E.g. if you specify a CRLF-terminated file on the command line followed by an
    	   LF-terminated file then autodetected line endings will be CRLF.
    	 * If you use --ors {something else} with (default or explicitly specified) --irs auto
    	   then line endings are autodetected on input and set to what you specify on output.
    	 * If you use --irs {something else} with (default or explicitly specified) --ors auto
    	   then the output line endings used are LF on Unix/Linux/BSD/MacOSX, and CRLF on Windows.
    
    	 Notes about all other separators:
    	 * IPS/OPS are only used for DKVP and XTAB formats, since only in these formats
    	   do key-value pairs appear juxtaposed.
    	 * IRS/ORS are ignored for XTAB format. Nominally IFS and OFS are newlines;
    	   XTAB records are separated by two or more consecutive IFS/OFS -- i.e.
    	   a blank line. Everything above about --irs/--ors/--rs auto becomes --ifs/--ofs/--fs
    	   auto for XTAB format. (XTAB's default IFS/OFS are "auto".)
    	 * OFS must be single-character for PPRINT format. This is because it is used
    	   with repetition for alignment; multi-character separators would make
    	   alignment impossible.
    	 * OPS may be multi-character for XTAB format, in which case alignment is
    	   disabled.
    	 * TSV is simply CSV using tab as field separator ("--fs tab").
    	 * FS/PS are ignored for markdown format; RS is used.
    	 * All FS and PS options are ignored for JSON format, since they are not relevant
    	   to the JSON format.
    	 * You can specify separators in any of the following ways, shown by example:
    	   - Type them out, quoting as necessary for shell escapes, e.g.
    	     "--fs '|' --ips :"
    	   - C-style escape sequences, e.g. "--rs '\r\n' --fs '\t'".
    	   - To avoid backslashing, you can use any of the following names:
    	     cr crcr newline lf lflf crlf crlfcrlf tab space comma pipe slash colon semicolon equals
    	 * Default separators by format:
    	     File format  RS	   FS	    PS
    	     gen	  N/A	   (N/A)    (N/A)
    	     dkvp	  auto	   ,	    =
    	     json	  auto	   (N/A)    (N/A)
    	     nidx	  auto	   space    (N/A)
    	     csv	  auto	   ,	    (N/A)
    	     csvlite	  auto	   ,	    (N/A)
    	     markdown	  auto	   (N/A)    (N/A)
    	     pprint	  auto	   space    (N/A)
    	     xtab	  (N/A)    auto     space
    
       CSV-SPECIFIC OPTIONS
    	 --implicit-csv-header Use 1,2,3,... as field labels, rather than from line 1
    			    of input files. Tip: combine with "label" to recreate
    			    missing headers.
    	 --allow-ragged-csv-input|--ragged If a data line has fewer fields than the header line,
    			    fill remaining keys with empty string. If a data line has more
    			    fields than the header line, use integer field labels as in
    			    the implicit-header case.
    	 --headerless-csv-output   Print only CSV data lines.
    	 -N		    Keystroke-saver for --implicit-csv-header --headerless-csv-output.
    
       DOUBLE-QUOTING FOR CSV/CSVLITE OUTPUT
    	 --quote-all	    Wrap all fields in double quotes
    	 --quote-none	    Do not wrap any fields in double quotes, even if they have
    			    OFS or ORS in them
    	 --quote-minimal    Wrap fields in double quotes only if they have OFS or ORS
    			    in them (default)
    	 --quote-numeric    Wrap fields in double quotes only if they have numbers
    			    in them
    	 --quote-original   Wrap fields in double quotes if and only if they were
    			    quoted on input. This isn't sticky for computed fields:
    			    e.g. if fields a and b were quoted on input and you do
    			    "put '$c = $a . $b'" then field c won't inherit a or b's
    			    was-quoted-on-input flag.
    
       NUMERICAL FORMATTING
    	 --ofmt {format}    E.g. %.18lf, %.0lf. Please use sprintf-style codes for
    			    double-precision. Applies to verbs which compute new
    			    values, e.g. put, stats1, stats2. See also the fmtnum
    			    function within mlr put (mlr --help-all-functions).
    			    Defaults to %lf.
    
       OTHER OPTIONS
    	 --seed {n} with n of the form 12345678 or 0xcafefeed. For put/filter
    			    urand()/urandint()/urand32().
    	 --nr-progress-mod {m}, with m a positive integer: print filename and record
    			    count to stderr every m input records.
    	 --from {filename}  Use this to specify an input file before the verb(s),
    			    rather than after. May be used more than once. Example:
    			    "mlr --from a.dat --from b.dat cat" is the same as
    			    "mlr cat a.dat b.dat".
    	 -n		    Process no input files, nor standard input either. Useful
    			    for mlr put with begin/end statements only. (Same as --from
    			    /dev/null.) Also useful in "mlr -n put -v '...'" for
    			    analyzing abstract syntax trees (if that's your thing).
    	 -I		    Process files in-place. For each file name on the command
    			    line, output is written to a temp file in the same
    			    directory, which is then renamed over the original. Each
    			    file is processed in isolation: if the output format is
    			    CSV, CSV headers will be present in each output file;
    			    statistics are only over each file's own records; and so on.
    
       THEN-CHAINING
           Output of one verb may be chained as input to another using "then", e.g.
    	 mlr stats1 -a min,mean,max -f flag,u,v -g color then sort -f color
    
       AUXILIARY COMMANDS
           Miller has a few otherwise-standalone executables packaged within it.
           They do not participate in any other parts of Miller.
           Available subcommands:
    	 aux-list
    	 lecat
    	 termcvt
    	 hex
    	 unhex
    	 netbsd-strptime
           For more information, please invoke mlr {subcommand} --help
    
    MLRRC
           You can set up personal defaults via a $HOME/.mlrrc and/or ./.mlrrc.
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
           https://johnkerl.org/miller/doc/customization.html
    
    VERBS
       altkv
           Usage: mlr altkv [no options]
           Given fields with values of the form a,b,c,d,e,f emits a=b,c=d,e=f pairs.
    
       bar
           Usage: mlr bar [options]
           Replaces a numeric field with a number of asterisks, allowing for cheesy
           bar plots. These align best with --opprint or --oxtab output format.
           Options:
           -f   {a,b,c}	 Field names to convert to bars.
           -c   {character}  Fill character: default '*'.
           -x   {character}  Out-of-bounds character: default '#'.
           -b   {character}  Blank character: default '.'.
           --lo {lo}	 Lower-limit value for min-width bar: default '0.000000'.
           --hi {hi}	 Upper-limit value for max-width bar: default '100.000000'.
           -w   {n} 	 Bar-field width: default '40'.
           --auto		 Automatically computes limits, ignoring --lo and --hi.
    			 Holds all records in memory before producing any output.
    
       bootstrap
           Usage: mlr bootstrap [options]
           Emits an n-sample, with replacement, of the input records.
           Options:
           -n {number} Number of samples to output. Defaults to number of input records.
    		   Must be non-negative.
           See also mlr sample and mlr shuffle.
    
       cat
           Usage: mlr cat [options]
           Passes input records directly to output. Most useful for format conversion.
           Options:
           -n	 Prepend field "n" to each record with record-counter starting at 1
           -g {comma-separated field name(s)} When used with -n/-N, writes record-counters
    		 keyed by specified field name(s).
           -v	 Write a low-level record-structure dump to stderr.
           -N {name} Prepend field {name} to each record with record-counter starting at 1
    
       check
           Usage: mlr check
           Consumes records without printing any output.
           Useful for doing a well-formatted check on input data.
    
       clean-whitespace
           Usage: mlr clean-whitespace [options]
           For each record, for each field in the record, whitespace-cleans the keys and
           values. Whitespace-cleaning entails stripping leading and trailing whitespace,
           and replacing multiple whitespace with singles. For finer-grained control,
           please see the DSL functions lstrip, rstrip, strip, collapse_whitespace,
           and clean_whitespace.
    
           Options:
           -k|--keys-only	 Do not touch values.
           -v|--values-only  Do not touch keys.
           It is an error to specify -k as well as -v -- to clean keys and values,
           leave off -k as well as -v.
    
       count
           Usage: mlr count [options]
           Prints number of records, optionally grouped by distinct values for specified field names.
    
           Options:
           -g {a,b,c}    Field names for distinct count.
           -n	     Show only the number of distinct values. Not interesting without -g.
           -o {name}     Field name for output count. Default "count".
    
       count-distinct
           Usage: mlr count-distinct [options]
           Prints number of records having distinct values for specified field names.
           Same as uniq -c.
    
           Options:
           -f {a,b,c}    Field names for distinct count.
           -n	     Show only the number of distinct values. Not compatible with -u.
           -o {name}     Field name for output count. Default "count".
    		     Ignored with -u.
           -u	     Do unlashed counts for multiple field names. With -f a,b and
    		     without -u, computes counts for distinct combinations of a
    		     and b field values. With -f a,b and with -u, computes counts
    		     for distinct a field values and counts for distinct b field
    		     values separately.
    
       count-similar
           Usage: mlr count-similar [options]
           Ingests all records, then emits each record augmented by a count of
           the number of other records having the same group-by field values.
           Options:
           -g {d,e,f} Group-by-field names for counts.
           -o {name}  Field name for output count. Default "count".
    
       cut
           Usage: mlr cut [options]
           Passes through input records with specified fields included/excluded.
           -f {a,b,c}	Field names to include for cut.
           -o		Retain fields in the order specified here in the argument list.
    			Default is to retain them in the order found in the input data.
           -x|--complement	Exclude, rather than include, field names specified by -f.
           -r		Treat field names as regular expressions. "ab", "a.*b" will
    			match any field name containing the substring "ab" or matching
    			"a.*b", respectively; anchors of the form "^ab$", "^a.*b$" may
    			be used. The -o flag is ignored when -r is present.
           Examples:
    	 mlr cut -f hostname,status
    	 mlr cut -x -f hostname,status
    	 mlr cut -r -f '^status$,sda[0-9]'
    	 mlr cut -r -f '^status$,"sda[0-9]"'
    	 mlr cut -r -f '^status$,"sda[0-9]"i' (this is case-insensitive)
    
       decimate
           Usage: mlr decimate [options]
           -n {count}    Decimation factor; default 10
           -b	     Decimate by printing first of every n.
           -e	     Decimate by printing last of every n (default).
           -g {a,b,c}    Optional group-by-field names for decimate counts
           Passes through one of every n records, optionally by category.
    
       fill-down
           Usage: mlr fill-down [options]
           -f {a,b,c}	   Field names for fill-down
           -a|--only-if-absent Field names for fill-down
           If a given record has a missing value for a given field, fill that from
           the corresponding value from a previous record, if any.
           By default, a 'missing' field either is absent, or has the empty-string value.
           With -a, a field is 'missing' only if it is absent.
    
       filter
           Usage: mlr filter [options] {expression}
           Prints records for which {expression} evaluates to true.
           If there are multiple semicolon-delimited expressions, all of them are
           evaluated and the last one is used as the filter criterion.
    
           Conversion options:
           -S: Keeps field values as strings with no type inference to int or float.
           -F: Keeps field values as strings or floats with no inference to int.
           All field values are type-inferred to int/float/string unless this behavior is
           suppressed with -S or -F.
    
           Output/formatting options:
           --oflatsep {string}: Separator to use when flattening multi-level @-variables
    	   to output records for emit. Default ":".
           --jknquoteint: For dump output (JSON-formatted), do not quote map keys if non-string.
           --jvquoteall: For dump output (JSON-formatted), quote map values even if non-string.
           Any of the output-format command-line flags (see mlr -h). Example: using
    	 mlr --icsv --opprint ... then put --ojson 'tee > "mytap-".$a.".dat", $*' then ...
           the input is CSV, the output is pretty-print tabular, but the tee-file output
           is written in JSON format.
           --no-fflush: for emit, tee, print, and dump, don't call fflush() after every
    	   record.
    
           Expression-specification options:
           -f {filename}: the DSL expression is taken from the specified file rather
    	   than from the command line. Outer single quotes wrapping the expression
    	   should not be placed in the file. If -f is specified more than once,
    	   all input files specified using -f are concatenated to produce the expression.
    	   (For example, you can define functions in one file and call them from another.)
           -e {expression}: You can use this after -f to add an expression. Example use
    	   case: define functions/subroutines in a file you specify with -f, then call
    	   them with an expression you specify with -e.
           (If you mix -e and -f then the expressions are evaluated in the order encountered.
           Since the expression pieces are simply concatenated, please be sure to use intervening
           semicolons to separate expressions.)
    
           -s name=value: Predefines out-of-stream variable @name to have value "value".
    	   Thus mlr filter put -s foo=97 '$column += @foo' is like
    	   mlr filter put 'begin {@foo = 97} $column += @foo'.
    	   The value part is subject to type-inferencing as specified by -S/-F.
    	   May be specified more than once, e.g. -s name1=value1 -s name2=value2.
    	   Note: the value may be an environment variable, e.g. -s sequence=$SEQUENCE
    
           Tracing options:
           -v: Prints the expressions's AST (abstract syntax tree), which gives
    	   full transparency on the precedence and associativity rules of
    	   Miller's grammar, to stdout.
           -a: Prints a low-level stack-allocation trace to stdout.
           -t: Prints a low-level parser trace to stderr.
           -T: Prints a every statement to stderr as it is executed.
    
           Other options:
           -x: Prints records for which {expression} evaluates to false.
    
           Please use a dollar sign for field names and double-quotes for string
           literals. If field names have special characters such as "." then you might
           use braces, e.g. '${field.name}'. Miller built-in variables are
           NF NR FNR FILENUM FILENAME M_PI M_E, and ENV["namegoeshere"] to access environment
           variables. The environment-variable name may be an expression, e.g. a field
           value.
    
           Use # to comment to end of line.
    
           Examples:
    	 mlr filter 'log10($count) > 4.0'
    	 mlr filter 'FNR == 2'	       (second record in each file)
    	 mlr filter 'urand() < 0.001'  (subsampling)
    	 mlr filter '$color != "blue" && $value > 4.2'
    	 mlr filter '($x<.5 && $y<.5) || ($x>.5 && $y>.5)'
    	 mlr filter '($name =~ "^sys.*east$") || ($name =~ "^dev.[0-9]+"i)'
    	 mlr filter '$ab = $a+$b; $cd = $c+$d; $ab != $cd'
    	 mlr filter '
    	   NR == 1 ||
    	  #NR == 2 ||
    	   NR == 3
    	 '
    
           Please see https://miller.readthedocs.io/en/latest/reference.html for more information
           including function list. Or "mlr -f". Please also see "mlr grep" which is
           useful when you don't yet know which field name(s) you're looking for.
           Please see in particular:
    	 http://www.johnkerl.org/miller/doc/reference-verbs.html#filter
    
       format-values
           Usage: mlr format-values [options]
           Applies format strings to all field values, depending on autodetected type.
           * If a field value is detected to be integer, applies integer format.
           * Else, if a field value is detected to be float, applies float format.
           * Else, applies string format.
    
           Note: this is a low-keystroke way to apply formatting to many fields. To get
           finer control, please see the fmtnum function within the mlr put DSL.
    
           Note: this verb lets you apply arbitrary format strings, which can produce
           undefined behavior and/or program crashes.  See your system's "man printf".
    
           Options:
           -i {integer format} Defaults to "%lld".
    			   Examples: "%06lld", "%08llx".
    			   Note that Miller integers are long long so you must use
    			   formats which apply to long long, e.g. with ll in them.
    			   Undefined behavior results otherwise.
           -f {float format}   Defaults to "%lf".
    			   Examples: "%8.3lf", "%.6le".
    			   Note that Miller floats are double-precision so you must
    			   use formats which apply to double, e.g. with l[efg] in them.
    			   Undefined behavior results otherwise.
           -s {string format}  Defaults to "%s".
    			   Examples: "_%s", "%08s".
    			   Note that you must use formats which apply to string, e.g.
    			   with s in them. Undefined behavior results otherwise.
           -n		   Coerce field values autodetected as int to float, and then
    			   apply the float format.
    
       fraction
           Usage: mlr fraction [options]
           For each record's value in specified fields, computes the ratio of that
           value to the sum of values in that field over all input records.
           E.g. with input records	x=1  x=2  x=3  and  x=4, emits output records
           x=1,x_fraction=0.1  x=2,x_fraction=0.2  x=3,x_fraction=0.3  and	x=4,x_fraction=0.4
    
           Note: this is internally a two-pass algorithm: on the first pass it retains
           input records and accumulates sums; on the second pass it computes quotients
           and emits output records. This means it produces no output until all input is read.
    
           Options:
           -f {a,b,c}    Field name(s) for fraction calculation
           -g {d,e,f}    Optional group-by-field name(s) for fraction counts
           -p	     Produce percents [0..100], not fractions [0..1]. Output field names
    		     end with "_percent" rather than "_fraction"
           -c	     Produce cumulative distributions, i.e. running sums: each output
    		     value folds in the sum of the previous for the specified group
    		     E.g. with input records  x=1  x=2	x=3  and  x=4, emits output records
    		     x=1,x_cumulative_fraction=0.1  x=2,x_cumulative_fraction=0.3
    		     x=3,x_cumulative_fraction=0.6  and  x=4,x_cumulative_fraction=1.0
    
       grep
           Usage: mlr grep [options] {regular expression}
           Passes through records which match {regex}.
           Options:
           -i    Use case-insensitive search.
           -v    Invert: pass through records which do not match the regex.
           Note that "mlr filter" is more powerful, but requires you to know field names.
           By contrast, "mlr grep" allows you to regex-match the entire record. It does
           this by formatting each record in memory as DKVP, using command-line-specified
           ORS/OFS/OPS, and matching the resulting line against the regex specified
           here. In particular, the regex is not applied to the input stream: if you
           have CSV with header line "x,y,z" and data line "1,2,3" then the regex will
           be matched, not against either of these lines, but against the DKVP line
           "x=1,y=2,z=3".  Furthermore, not all the options to system grep are supported,
           and this command is intended to be merely a keystroke-saver. To get all the
           features of system grep, you can do
    	 "mlr --odkvp ... | grep ... | mlr --idkvp ..."
    
       group-by
           Usage: mlr group-by {comma-separated field names}
           Outputs records in batches having identical values at specified field names.
    
       group-like
           Usage: mlr group-like
           Outputs records in batches having identical field names.
    
       having-fields
           Usage: mlr having-fields [options]
           Conditionally passes through records depending on each record's field names.
           Options:
    	 --at-least	 {comma-separated names}
    	 --which-are	 {comma-separated names}
    	 --at-most	 {comma-separated names}
    	 --all-matching  {regular expression}
    	 --any-matching  {regular expression}
    	 --none-matching {regular expression}
           Examples:
    	 mlr having-fields --which-are amount,status,owner
    	 mlr having-fields --any-matching 'sda[0-9]'
    	 mlr having-fields --any-matching '"sda[0-9]"'
    	 mlr having-fields --any-matching '"sda[0-9]"i' (this is case-insensitive)
    
       head
           Usage: mlr head [options]
           -n {count}    Head count to print; default 10
           -g {a,b,c}    Optional group-by-field names for head counts
           Passes through the first n records, optionally by category.
           Without -g, ceases consuming more input (i.e. is fast) when n
           records have been read.
    
       histogram
           Usage: mlr histogram [options]
           -f {a,b,c}    Value-field names for histogram counts
           --lo {lo}     Histogram low value
           --hi {hi}     Histogram high value
           --nbins {n}   Number of histogram bins
           --auto	     Automatically computes limits, ignoring --lo and --hi.
    		     Holds all values in memory before producing any output.
           -o {prefix}   Prefix for output field name. Default: no prefix.
           Just a histogram. Input values < lo or > hi are not counted.
    
       join
           Usage: mlr join [options]
           Joins records from specified left file name with records from all file names
           at the end of the Miller argument list.
           Functionality is essentially the same as the system "join" command, but for
           record streams.
           Options:
    	 -f {left file name}
    	 -j {a,b,c}   Comma-separated join-field names for output
    	 -l {a,b,c}   Comma-separated join-field names for left input file;
    		      defaults to -j values if omitted.
    	 -r {a,b,c}   Comma-separated join-field names for right input file(s);
    		      defaults to -j values if omitted.
    	 --lp {text}  Additional prefix for non-join output field names from
    		      the left file
    	 --rp {text}  Additional prefix for non-join output field names from
    		      the right file(s)
    	 --np	      Do not emit paired records
    	 --ul	      Emit unpaired records from the left file
    	 --ur	      Emit unpaired records from the right file(s)
    	 -s|--sorted-input  Require sorted input: records must be sorted
    		      lexically by their join-field names, else not all records will
    		      be paired. The only likely use case for this is with a left
    		      file which is too big to fit into system memory otherwise.
    	 -u	      Enable unsorted input. (This is the default even without -u.)
    		      In this case, the entire left file will be loaded into memory.
    	 --prepipe {command} As in main input options; see mlr --help for details.
    		      If you wish to use a prepipe command for the main input as well
    		      as here, it must be specified there as well as here.
           File-format options default to those for the right file names on the Miller
           argument list, but may be overridden for the left file as follows. Please see
           the main "mlr --help" for more information on syntax for these arguments.
    	 -i {one of csv,dkvp,nidx,pprint,xtab}
    	 --irs {record-separator character}
    	 --ifs {field-separator character}
    	 --ips {pair-separator character}
    	 --repifs
    	 --repips
           Please use "mlr --usage-separator-options" for information on specifying separators.
           Please see https://miller.readthedocs.io/en/latest/reference-verbs.html#join for more information
           including examples.
    
       label
           Usage: mlr label {new1,new2,new3,...}
           Given n comma-separated names, renames the first n fields of each record to
           have the respective name. (Fields past the nth are left with their original
           names.) Particularly useful with --inidx or --implicit-csv-header, to give
           useful names to otherwise integer-indexed fields.
           Examples:
    	 "echo 'a b c d' | mlr --inidx --odkvp cat"	  gives "1=a,2=b,3=c,4=d"
    	 "echo 'a b c d' | mlr --inidx --odkvp label s,t" gives "s=a,t=b,3=c,4=d"
    
       least-frequent
           Usage: mlr least-frequent [options]
           Shows the least frequently occurring distinct values for specified field names.
           The first entry is the statistical anti-mode; the remaining are runners-up.
           Options:
           -f {one or more comma-separated field names}. Required flag.
           -n {count}. Optional flag defaulting to 10.
           -b	   Suppress counts; show only field values.
           -o {name}   Field name for output count. Default "count".
           See also "mlr most-frequent".
    
       merge-fields
           Usage: mlr merge-fields [options]
           Computes univariate statistics for each input record, accumulated across
           specified fields.
           Options:
           -a {sum,count,...}  Names of accumulators. One or more of:
    	 count	   Count instances of fields
    	 mode	   Find most-frequently-occurring values for fields; first-found wins tie
    	 antimode  Find least-frequently-occurring values for fields; first-found wins tie
    	 sum	   Compute sums of specified fields
    	 mean	   Compute averages (sample means) of specified fields
    	 stddev    Compute sample standard deviation of specified fields
    	 var	   Compute sample variance of specified fields
    	 meaneb    Estimate error bars for averages (assuming no sample autocorrelation)
    	 skewness  Compute sample skewness of specified fields
    	 kurtosis  Compute sample kurtosis of specified fields
    	 min	   Compute minimum values of specified fields
    	 max	   Compute maximum values of specified fields
           -f {a,b,c}  Value-field names on which to compute statistics. Requires -o.
           -r {a,b,c}  Regular expressions for value-field names on which to compute
    		   statistics. Requires -o.
           -c {a,b,c}  Substrings for collapse mode. All fields which have the same names
    		   after removing substrings will be accumulated together. Please see
    		   examples below.
           -i	   Use interpolated percentiles, like R's type=7; default like type=1.
    		   Not sensical for string-valued fields.
           -o {name}   Output field basename for -f/-r.
           -k	   Keep the input fields which contributed to the output statistics;
    		   the default is to omit them.
           -F	   Computes integerable things (e.g. count) in floating point.
    
           String-valued data make sense unless arithmetic on them is required,
           e.g. for sum, mean, interpolated percentiles, etc. In case of mixed data,
           numbers are less than strings.
    
           Example input data: "a_in_x=1,a_out_x=2,b_in_y=4,b_out_x=8".
           Example: mlr merge-fields -a sum,count -f a_in_x,a_out_x -o foo
    	 produces "b_in_y=4,b_out_x=8,foo_sum=3,foo_count=2" since "a_in_x,a_out_x" are
    	 summed over.
           Example: mlr merge-fields -a sum,count -r in_,out_ -o bar
    	 produces "bar_sum=15,bar_count=4" since all four fields are summed over.
           Example: mlr merge-fields -a sum,count -c in_,out_
    	 produces "a_x_sum=3,a_x_count=2,b_y_sum=4,b_y_count=1,b_x_sum=8,b_x_count=1"
    	 since "a_in_x" and "a_out_x" both collapse to "a_x", "b_in_y" collapses to
    	 "b_y", and "b_out_x" collapses to "b_x".
    
       most-frequent
           Usage: mlr most-frequent [options]
           Shows the most frequently occurring distinct values for specified field names.
           The first entry is the statistical mode; the remaining are runners-up.
           Options:
           -f {one or more comma-separated field names}. Required flag.
           -n {count}. Optional flag defaulting to 10.
           -b	   Suppress counts; show only field values.
           -o {name}   Field name for output count. Default "count".
           See also "mlr least-frequent".
    
       nest
           Usage: mlr nest [options]
           Explodes specified field values into separate fields/records, or reverses this.
           Options:
    	 --explode,--implode   One is required.
    	 --values,--pairs      One is required.
    	 --across-records,--across-fields One is required.
    	 -f {field name}       Required.
    	 --nested-fs {string}  Defaults to ";". Field separator for nested values.
    	 --nested-ps {string}  Defaults to ":". Pair separator for nested key-value pairs.
    	 --evar {string}       Shorthand for --explode --values ---across-records --nested-fs {string}
    	 --ivar {string}       Shorthand for --implode --values ---across-records --nested-fs {string}
           Please use "mlr --usage-separator-options" for information on specifying separators.
    
           Examples:
    
    	 mlr nest --explode --values --across-records -f x
    	 with input record "x=a;b;c,y=d" produces output records
    	   "x=a,y=d"
    	   "x=b,y=d"
    	   "x=c,y=d"
    	 Use --implode to do the reverse.
    
    	 mlr nest --explode --values --across-fields -f x
    	 with input record "x=a;b;c,y=d" produces output records
    	   "x_1=a,x_2=b,x_3=c,y=d"
    	 Use --implode to do the reverse.
    
    	 mlr nest --explode --pairs --across-records -f x
    	 with input record "x=a:1;b:2;c:3,y=d" produces output records
    	   "a=1,y=d"
    	   "b=2,y=d"
    	   "c=3,y=d"
    
    	 mlr nest --explode --pairs --across-fields -f x
    	 with input record "x=a:1;b:2;c:3,y=d" produces output records
    	   "a=1,b=2,c=3,y=d"
    
           Notes:
           * With --pairs, --implode doesn't make sense since the original field name has
    	 been lost.
           * The combination "--implode --values --across-records" is non-streaming:
    	 no output records are produced until all input records have been read. In
    	 particular, this means it won't work in tail -f contexts. But all other flag
    	 combinations result in streaming (tail -f friendly) data processing.
           * It's up to you to ensure that the nested-fs is distinct from your data's IFS:
    	 e.g. by default the former is semicolon and the latter is comma.
           See also mlr reshape.
    
       nothing
           Usage: mlr nothing
           Drops all input records. Useful for testing, or after tee/print/etc. have
           produced other output.
    
       put
           Usage: mlr put [options] {expression}
           Adds/updates specified field(s). Expressions are semicolon-separated and must
           either be assignments, or evaluate to boolean.  Booleans with following
           statements in curly braces control whether those statements are executed;
           booleans without following curly braces do nothing except side effects (e.g.
           regex-captures into \1, \2, etc.).
    
           Conversion options:
           -S: Keeps field values as strings with no type inference to int or float.
           -F: Keeps field values as strings or floats with no inference to int.
           All field values are type-inferred to int/float/string unless this behavior is
           suppressed with -S or -F.
    
           Output/formatting options:
           --oflatsep {string}: Separator to use when flattening multi-level @-variables
    	   to output records for emit. Default ":".
           --jknquoteint: For dump output (JSON-formatted), do not quote map keys if non-string.
           --jvquoteall: For dump output (JSON-formatted), quote map values even if non-string.
           Any of the output-format command-line flags (see mlr -h). Example: using
    	 mlr --icsv --opprint ... then put --ojson 'tee > "mytap-".$a.".dat", $*' then ...
           the input is CSV, the output is pretty-print tabular, but the tee-file output
           is written in JSON format.
           --no-fflush: for emit, tee, print, and dump, don't call fflush() after every
    	   record.
    
           Expression-specification options:
           -f {filename}: the DSL expression is taken from the specified file rather
    	   than from the command line. Outer single quotes wrapping the expression
    	   should not be placed in the file. If -f is specified more than once,
    	   all input files specified using -f are concatenated to produce the expression.
    	   (For example, you can define functions in one file and call them from another.)
           -e {expression}: You can use this after -f to add an expression. Example use
    	   case: define functions/subroutines in a file you specify with -f, then call
    	   them with an expression you specify with -e.
           (If you mix -e and -f then the expressions are evaluated in the order encountered.
           Since the expression pieces are simply concatenated, please be sure to use intervening
           semicolons to separate expressions.)
    
           -s name=value: Predefines out-of-stream variable @name to have value "value".
    	   Thus mlr put put -s foo=97 '$column += @foo' is like
    	   mlr put put 'begin {@foo = 97} $column += @foo'.
    	   The value part is subject to type-inferencing as specified by -S/-F.
    	   May be specified more than once, e.g. -s name1=value1 -s name2=value2.
    	   Note: the value may be an environment variable, e.g. -s sequence=$SEQUENCE
    
           Tracing options:
           -v: Prints the expressions's AST (abstract syntax tree), which gives
    	   full transparency on the precedence and associativity rules of
    	   Miller's grammar, to stdout.
           -a: Prints a low-level stack-allocation trace to stdout.
           -t: Prints a low-level parser trace to stderr.
           -T: Prints a every statement to stderr as it is executed.
    
           Other options:
           -q: Does not include the modified record in the output stream. Useful for when
    	   all desired output is in begin and/or end blocks.
    
           Please use a dollar sign for field names and double-quotes for string
           literals. If field names have special characters such as "." then you might
           use braces, e.g. '${field.name}'. Miller built-in variables are
           NF NR FNR FILENUM FILENAME M_PI M_E, and ENV["namegoeshere"] to access environment
           variables. The environment-variable name may be an expression, e.g. a field
           value.
    
           Use # to comment to end of line.
    
           Examples:
    	 mlr put '$y = log10($x); $z = sqrt($y)'
    	 mlr put '$x>0.0 { $y=log10($x); $z=sqrt($y) }' # does {...} only if $x > 0.0
    	 mlr put '$x>0.0;  $y=log10($x); $z=sqrt($y)'	# does all three statements
    	 mlr put '$a =~ "([a-z]+)_([0-9]+);  $b = "left_\1"; $c = "right_\2"'
    	 mlr put '$a =~ "([a-z]+)_([0-9]+) { $b = "left_\1"; $c = "right_\2" }'
    	 mlr put '$filename = FILENAME'
    	 mlr put '$colored_shape = $color . "_" . $shape'
    	 mlr put '$y = cos($theta); $z = atan2($y, $x)'
    	 mlr put '$name = sub($name, "http.*com"i, "")'
    	 mlr put -q '@sum += $x; end {emit @sum}'
    	 mlr put -q '@sum[$a] += $x; end {emit @sum, "a"}'
    	 mlr put -q '@sum[$a][$b] += $x; end {emit @sum, "a", "b"}'
    	 mlr put -q '@min=min(@min,$x);@max=max(@max,$x); end{emitf @min, @max}'
    	 mlr put -q 'is_null(@xmax) || $x > @xmax {@xmax=$x; @recmax=$*}; end {emit @recmax}'
    	 mlr put '
    	   $x = 1;
    	  #$y = 2;
    	   $z = 3
    	 '
    
           Please see also 'mlr -k' for examples using redirected output.
    
           Please see https://miller.readthedocs.io/en/latest/reference.html for more information
           including function list. Or "mlr -f".
           Please see in particular:
    	 http://www.johnkerl.org/miller/doc/reference-verbs.html#put
    
       regularize
           Usage: mlr regularize
           For records seen earlier in the data stream with same field names in
           a different order, outputs them with field names in the previously
           encountered order.
           Example: input records a=1,c=2,b=3, then e=4,d=5, then c=7,a=6,b=8
           output as	      a=1,c=2,b=3, then e=4,d=5, then a=6,c=7,b=8
    
       remove-empty-columns
           Usage: mlr remove-empty-columns
           Omits fields which are empty on every input row. Non-streaming.
    
       rename
           Usage: mlr rename [options] {old1,new1,old2,new2,...}
           Renames specified fields.
           Options:
           -r	  Treat old field  names as regular expressions. "ab", "a.*b"
    		  will match any field name containing the substring "ab" or
    		  matching "a.*b", respectively; anchors of the form "^ab$",
    		  "^a.*b$" may be used. New field names may be plain strings,
    		  or may contain capture groups of the form "\1" through
    		  "\9". Wrapping the regex in double quotes is optional, but
    		  is required if you wish to follow it with 'i' to indicate
    		  case-insensitivity.
           -g	  Do global replacement within each field name rather than
    		  first-match replacement.
           Examples:
           mlr rename old_name,new_name'
           mlr rename old_name_1,new_name_1,old_name_2,new_name_2'
           mlr rename -r 'Date_[0-9]+,Date,'  Rename all such fields to be "Date"
           mlr rename -r '"Date_[0-9]+",Date' Same
           mlr rename -r 'Date_([0-9]+).*,\1' Rename all such fields to be of the form 20151015
           mlr rename -r '"name"i,Name'	  Rename "name", "Name", "NAME", etc. to "Name"
    
       reorder
           Usage: mlr reorder [options]
           -f {a,b,c}   Field names to reorder.
           -e	    Put specified field names at record end: default is to put
    		    them at record start.
           Examples:
           mlr reorder    -f a,b sends input record "d=4,b=2,a=1,c=3" to "a=1,b=2,d=4,c=3".
           mlr reorder -e -f a,b sends input record "d=4,b=2,a=1,c=3" to "d=4,c=3,a=1,b=2".
    
       repeat
           Usage: mlr repeat [options]
           Copies input records to output records multiple times.
           Options must be exactly one of the following:
    	 -n {repeat count}  Repeat each input record this many times.
    	 -f {field name}    Same, but take the repeat count from the specified
    			    field name of each input record.
           Example:
    	 echo x=0 | mlr repeat -n 4 then put '$x=urand()'
           produces:
    	x=0.488189
    	x=0.484973
    	x=0.704983
    	x=0.147311
           Example:
    	 echo a=1,b=2,c=3 | mlr repeat -f b
           produces:
    	 a=1,b=2,c=3
    	 a=1,b=2,c=3
           Example:
    	 echo a=1,b=2,c=3 | mlr repeat -f c
           produces:
    	 a=1,b=2,c=3
    	 a=1,b=2,c=3
    	 a=1,b=2,c=3
    
       reshape
           Usage: mlr reshape [options]
           Wide-to-long options:
    	 -i {input field names}   -o {key-field name,value-field name}
    	 -r {input field regexes} -o {key-field name,value-field name}
    	 These pivot/reshape the input data such that the input fields are removed
    	 and separate records are emitted for each key/value pair.
    	 Note: this works with tail -f and produces output records for each input
    	 record seen.
           Long-to-wide options:
    	 -s {key-field name,value-field name}
    	 These pivot/reshape the input data to undo the wide-to-long operation.
    	 Note: this does not work with tail -f; it produces output records only after
    	 all input records have been read.
    
           Examples:
    
    	 Input file "wide.txt":
    	   time       X 	  Y
    	   2009-01-01 0.65473572  2.4520609
    	   2009-01-02 -0.89248112 0.2154713
    	   2009-01-03 0.98012375  1.3179287
    
    	 mlr --pprint reshape -i X,Y -o item,value wide.txt
    	   time       item value
    	   2009-01-01 X    0.65473572
    	   2009-01-01 Y    2.4520609
    	   2009-01-02 X    -0.89248112
    	   2009-01-02 Y    0.2154713
    	   2009-01-03 X    0.98012375
    	   2009-01-03 Y    1.3179287
    
    	 mlr --pprint reshape -r '[A-Z]' -o item,value wide.txt
    	   time       item value
    	   2009-01-01 X    0.65473572
    	   2009-01-01 Y    2.4520609
    	   2009-01-02 X    -0.89248112
    	   2009-01-02 Y    0.2154713
    	   2009-01-03 X    0.98012375
    	   2009-01-03 Y    1.3179287
    
    	 Input file "long.txt":
    	   time       item value
    	   2009-01-01 X    0.65473572
    	   2009-01-01 Y    2.4520609
    	   2009-01-02 X    -0.89248112
    	   2009-01-02 Y    0.2154713
    	   2009-01-03 X    0.98012375
    	   2009-01-03 Y    1.3179287
    
    	 mlr --pprint reshape -s item,value long.txt
    	   time       X 	  Y
    	   2009-01-01 0.65473572  2.4520609
    	   2009-01-02 -0.89248112 0.2154713
    	   2009-01-03 0.98012375  1.3179287
           See also mlr nest.
    
       sample
           Usage: mlr sample [options]
           Reservoir sampling (subsampling without replacement), optionally by category.
           -k {count}    Required: number of records to output, total, or by group if using -g.
           -g {a,b,c}    Optional: group-by-field names for samples.
           See also mlr bootstrap and mlr shuffle.
    
       sec2gmt
           Usage: mlr sec2gmt [options] {comma-separated list of field names}
           Replaces a numeric field representing seconds since the epoch with the
           corresponding GMT timestamp; leaves non-numbers as-is. This is nothing
           more than a keystroke-saver for the sec2gmt function:
    	 mlr sec2gmt time1,time2
           is the same as
    	 mlr put '$time1=sec2gmt($time1);$time2=sec2gmt($time2)'
           Options:
           -1 through -9: format the seconds using 1..9 decimal places, respectively.
    
       sec2gmtdate
           Usage: mlr sec2gmtdate {comma-separated list of field names}
           Replaces a numeric field representing seconds since the epoch with the
           corresponding GMT year-month-day timestamp; leaves non-numbers as-is.
           This is nothing more than a keystroke-saver for the sec2gmtdate function:
    	 mlr sec2gmtdate time1,time2
           is the same as
    	 mlr put '$time1=sec2gmtdate($time1);$time2=sec2gmtdate($time2)'
    
       seqgen
           Usage: mlr seqgen [options]
           Produces a sequence of counters.  Discards the input record stream. Produces
           output as specified by the following options:
           -f {name} Field name for counters; default "i".
           --start {number} Inclusive start value; default "1".
           --stop  {number} Inclusive stop value; default "100".
           --step  {number} Step value; default "1".
           Start, stop, and/or step may be floating-point. Output is integer if start,
           stop, and step are all integers. Step may be negative. It may not be zero
           unless start == stop.
    
       shuffle
           Usage: mlr shuffle {no options}
           Outputs records randomly permuted. No output records are produced until
           all input records are read.
           See also mlr bootstrap and mlr sample.
    
       skip-trivial-records
           Usage: mlr skip-trivial-records [options]
           Passes through all records except:
           * those with zero fields;
           * those for which all fields have empty value.
    
       sort
           Usage: mlr sort {flags}
           Flags:
    	 -f  {comma-separated field names}  Lexical ascending
    	 -n  {comma-separated field names}  Numerical ascending; nulls sort last
    	 -nf {comma-separated field names}  Same as -n
    	 -r  {comma-separated field names}  Lexical descending
    	 -nr {comma-separated field names}  Numerical descending; nulls sort first
           Sorts records primarily by the first specified field, secondarily by the second
           field, and so on.  (Any records not having all specified sort keys will appear
           at the end of the output, in the order they were encountered, regardless of the
           specified sort order.) The sort is stable: records that compare equal will sort
           in the order they were encountered in the input record stream.
    
           Example:
    	 mlr sort -f a,b -nr x,y,z
           which is the same as:
    	 mlr sort -f a -f b -nr x -nr y -nr z
    
       sort-within-records
           Usage: mlr sort-within-records [no options]
           Outputs records sorted lexically ascending by keys.
    
       stats1
           Usage: mlr stats1 [options]
           Computes univariate statistics for one or more given fields, accumulated across
           the input record stream.
           Options:
           -a {sum,count,...}  Names of accumulators: p10 p25.2 p50 p98 p100 etc. and/or
    			   one or more of:
    	  count     Count instances of fields
    	  mode	    Find most-frequently-occurring values for fields; first-found wins tie
    	  antimode  Find least-frequently-occurring values for fields; first-found wins tie
    	  sum	    Compute sums of specified fields
    	  mean	    Compute averages (sample means) of specified fields
    	  stddev    Compute sample standard deviation of specified fields
    	  var	    Compute sample variance of specified fields
    	  meaneb    Estimate error bars for averages (assuming no sample autocorrelation)
    	  skewness  Compute sample skewness of specified fields
    	  kurtosis  Compute sample kurtosis of specified fields
    	  min	    Compute minimum values of specified fields
    	  max	    Compute maximum values of specified fields
           -f {a,b,c}   Value-field names on which to compute statistics
           --fr {regex} Regex for value-field names on which to compute statistics
    		    (compute statistics on values in all field names matching regex)
           --fx {regex} Inverted regex for value-field names on which to compute statistics
    		    (compute statistics on values in all field names not matching regex)
           -g {d,e,f}   Optional group-by-field names
           --gr {regex} Regex for optional group-by-field names
    		    (group by values in field names matching regex)
           --gx {regex} Inverted regex for optional group-by-field names
    		    (group by values in field names not matching regex)
           --grfx {regex} Shorthand for --gr {regex} --fx {that same regex}
           -i	    Use interpolated percentiles, like R's type=7; default like type=1.
    		    Not sensical for string-valued fields.
           -s	    Print iterative stats. Useful in tail -f contexts (in which
    		    case please avoid pprint-format output since end of input
    		    stream will never be seen).
           -F	    Computes integerable things (e.g. count) in floating point.
           Example: mlr stats1 -a min,p10,p50,p90,max -f value -g size,shape
           Example: mlr stats1 -a count,mode -f size
           Example: mlr stats1 -a count,mode -f size -g shape
           Example: mlr stats1 -a count,mode --fr '^[a-h].*$' -gr '^k.*$'
    		This computes count and mode statistics on all field names beginning
    		with a through h, grouped by all field names starting with k.
           Notes:
           * p50 and median are synonymous.
           * min and max output the same results as p0 and p100, respectively, but use
    	 less memory.
           * String-valued data make sense unless arithmetic on them is required,
    	 e.g. for sum, mean, interpolated percentiles, etc. In case of mixed data,
    	 numbers are less than strings.
           * count and mode allow text input; the rest require numeric input.
    	 In particular, 1 and 1.0 are distinct text for count and mode.
           * When there are mode ties, the first-encountered datum wins.
    
       stats2
           Usage: mlr stats2 [options]
           Computes bivariate statistics for one or more given field-name pairs,
           accumulated across the input record stream.
           -a {linreg-ols,corr,...}  Names of accumulators: one or more of:
    	 linreg-pca   Linear regression using principal component analysis
    	 linreg-ols   Linear regression using ordinary least squares
    	 r2	      Quality metric for linreg-ols (linreg-pca emits its own)
    	 logireg      Logistic regression
    	 corr	      Sample correlation
    	 cov	      Sample covariance
    	 covx	      Sample-covariance matrix
           -f {a,b,c,d}   Value-field name-pairs on which to compute statistics.
    		      There must be an even number of names.
           -g {e,f,g}     Optional group-by-field names.
           -v	      Print additional output for linreg-pca.
           -s	      Print iterative stats. Useful in tail -f contexts (in which
    		      case please avoid pprint-format output since end of input
    		      stream will never be seen).
           --fit	      Rather than printing regression parameters, applies them to
    		      the input data to compute new fit fields. All input records are
    		      held in memory until end of input stream. Has effect only for
    		      linreg-ols, linreg-pca, and logireg.
           Only one of -s or --fit may be used.
           Example: mlr stats2 -a linreg-pca -f x,y
           Example: mlr stats2 -a linreg-ols,r2 -f x,y -g size,shape
           Example: mlr stats2 -a corr -f x,y
    
       step
           Usage: mlr step [options]
           Computes values dependent on the previous record, optionally grouped
           by category.
    
           Options:
           -a {delta,rsum,...}   Names of steppers: comma-separated, one or more of:
    	 delta	  Compute differences in field(s) between successive records
    	 shift	  Include value(s) in field(s) from previous record, if any
    	 from-first Compute differences in field(s) from first record
    	 ratio	  Compute ratios in field(s) between successive records
    	 rsum	  Compute running sums of field(s) between successive records
    	 counter  Count instances of field(s) between successive records
    	 ewma	  Exponentially weighted moving average over successive records
           -f {a,b,c} Value-field names on which to compute statistics
           -g {d,e,f} Optional group-by-field names
           -F	  Computes integerable things (e.g. counter) in floating point.
           -d {x,y,z} Weights for ewma. 1 means current sample gets all weight (no
    		  smoothing), near under under 1 is light smoothing, near over 0 is
    		  heavy smoothing. Multiple weights may be specified, e.g.
    		  "mlr step -a ewma -f sys_load -d 0.01,0.1,0.9". Default if omitted
    		  is "-d 0.5".
           -o {a,b,c} Custom suffixes for EWMA output fields. If omitted, these default to
    		  the -d values. If supplied, the number of -o values must be the same
    		  as the number of -d values.
    
           Examples:
    	 mlr step -a rsum -f request_size
    	 mlr step -a delta -f request_size -g hostname
    	 mlr step -a ewma -d 0.1,0.9 -f x,y
    	 mlr step -a ewma -d 0.1,0.9 -o smooth,rough -f x,y
    	 mlr step -a ewma -d 0.1,0.9 -o smooth,rough -f x,y -g group_name
    
           Please see https://miller.readthedocs.io/en/latest/reference-verbs.html#filter or
           https://en.wikipedia.org/wiki/Moving_average#Exponential_moving_average
           for more information on EWMA.
    
       tac
           Usage: mlr tac
           Prints records in reverse order from the order in which they were encountered.
    
       tail
           Usage: mlr tail [options]
           -n {count}    Tail count to print; default 10
           -g {a,b,c}    Optional group-by-field names for tail counts
           Passes through the last n records, optionally by category.
    
       tee
           Usage: mlr tee [options] {filename}
           Passes through input records (like mlr cat) but also writes to specified output
           file, using output-format flags from the command line (e.g. --ocsv). See also
           the "tee" keyword within mlr put, which allows data-dependent filenames.
           Options:
           -a:	    append to existing file, if any, rather than overwriting.
           --no-fflush: don't call fflush() after every record.
           Any of the output-format command-line flags (see mlr -h). Example: using
    	 mlr --icsv --opprint put '...' then tee --ojson ./mytap.dat then stats1 ...
           the input is CSV, the output is pretty-print tabular, but the tee-file output
           is written in JSON format.
    
       top
           Usage: mlr top [options]
           -f {a,b,c}    Value-field names for top counts.
           -g {d,e,f}    Optional group-by-field names for top counts.
           -n {count}    How many records to print per category; default 1.
           -a	     Print all fields for top-value records; default is
    		     to print only value and group-by fields. Requires a single
    		     value-field name only.
           --min	     Print top smallest values; default is top largest values.
           -F	     Keep top values as floats even if they look like integers.
           -o {name}     Field name for output indices. Default "top_idx".
           Prints the n records with smallest/largest values at specified fields,
           optionally by category.
    
       uniq
           Usage: mlr uniq [options]
           Prints distinct values for specified field names. With -c, same as
           count-distinct. For uniq, -f is a synonym for -g.
    
           Options:
           -g {d,e,f}    Group-by-field names for uniq counts.
           -c	     Show repeat counts in addition to unique values.
           -n	     Show only the number of distinct values.
           -o {name}     Field name for output count. Default "count".
           -a	     Output each unique record only once. Incompatible with -g.
    		     With -c, produces unique records, with repeat counts for each.
    		     With -n, produces only one record which is the unique-record count.
    		     With neither -c nor -n, produces unique records.
    
       unsparsify
           Usage: mlr unsparsify [options]
           Prints records with the union of field names over all input records.
           For field names absent in a given record but present in others, fills in a
           value. Without -f, this verb retains all input before producing any output.
    
           Options:
           --fill-with {filler string}  What to fill absent fields with. Defaults to
    				    the empty string.
           -f {a,b,c} Specify field names to be operated on. Any other fields won't be
    				    modified, and operation will be streaming.
    
           Example: if the input is two records, one being 'a=1,b=2' and the other
           being 'b=3,c=4', then the output is the two records 'a=1,b=2,c=' and
           'a=,b=3,c=4'.
    
    FUNCTIONS FOR FILTER/PUT
       +
           (class=arithmetic #args=2): Addition.
    
           + (class=arithmetic #args=1): Unary plus.
    
       -
           (class=arithmetic #args=2): Subtraction.
    
           - (class=arithmetic #args=1): Unary minus.
    
       *
           (class=arithmetic #args=2): Multiplication.
    
       /
           (class=arithmetic #args=2): Division.
    
       //
           (class=arithmetic #args=2): Integer division: rounds to negative (pythonic).
    
       .+
           (class=arithmetic #args=2): Addition, with integer-to-integer overflow
    
           .+ (class=arithmetic #args=1): Unary plus, with integer-to-integer overflow.
    
       .-
           (class=arithmetic #args=2): Subtraction, with integer-to-integer overflow.
    
           .- (class=arithmetic #args=1): Unary minus, with integer-to-integer overflow.
    
       .*
           (class=arithmetic #args=2): Multiplication, with integer-to-integer overflow.
    
       ./
           (class=arithmetic #args=2): Division, with integer-to-integer overflow.
    
       .//
           (class=arithmetic #args=2): Integer division: rounds to negative (pythonic), with integer-to-integer overflow.
    
       %
           (class=arithmetic #args=2): Remainder; never negative-valued (pythonic).
    
       **
           (class=arithmetic #args=2): Exponentiation; same as pow, but as an infix
           operator.
    
       |
           (class=arithmetic #args=2): Bitwise OR.
    
       ^
           (class=arithmetic #args=2): Bitwise XOR.
    
       &
           (class=arithmetic #args=2): Bitwise AND.
    
       ~
           (class=arithmetic #args=1): Bitwise NOT. Beware '$y=~$x' since =~ is the
           regex-match operator: try '$y = ~$x'.
    
       <<
           (class=arithmetic #args=2): Bitwise left-shift.
    
       >>
           (class=arithmetic #args=2): Bitwise right-shift.
    
       bitcount
           (class=arithmetic #args=1): Count of 1-bits
    
       ==
           (class=boolean #args=2): String/numeric equality. Mixing number and string
           results in string compare.
    
       !=
           (class=boolean #args=2): String/numeric inequality. Mixing number and string
           results in string compare.
    
       =~
           (class=boolean #args=2): String (left-hand side) matches regex (right-hand
           side), e.g. '$name =~ "^a.*b$"'.
    
       !=~
           (class=boolean #args=2): String (left-hand side) does not match regex
           (right-hand side), e.g. '$name !=~ "^a.*b$"'.
    
       >
           (class=boolean #args=2): String/numeric greater-than. Mixing number and string
           results in string compare.
    
       >=
           (class=boolean #args=2): String/numeric greater-than-or-equals. Mixing number
           and string results in string compare.
    
       <
           (class=boolean #args=2): String/numeric less-than. Mixing number and string
           results in string compare.
    
       <=
           (class=boolean #args=2): String/numeric less-than-or-equals. Mixing number
           and string results in string compare.
    
       &&
           (class=boolean #args=2): Logical AND.
    
       ||
           (class=boolean #args=2): Logical OR.
    
       ^^
           (class=boolean #args=2): Logical XOR.
    
       !
           (class=boolean #args=1): Logical negation.
    
       ? :
           (class=boolean #args=3): Ternary operator.
    
       .
           (class=string #args=2): String concatenation.
    
       gsub
           (class=string #args=3): Example: '$name=gsub($name, "old", "new")'
           (replace all).
    
       regextract
           (class=string #args=2): Example: '$name=regextract($name, "[A-Z]{3}[0-9]{2}")'
           .
    
       regextract_or_else
           (class=string #args=3): Example: '$name=regextract_or_else($name, "[A-Z]{3}[0-9]{2}", "default")'
           .
    
       strlen
           (class=string #args=1): String length.
    
       sub
           (class=string #args=3): Example: '$name=sub($name, "old", "new")'
           (replace once).
    
       ssub
           (class=string #args=3): Like sub but does no regexing. No characters are special.
    
       substr
           (class=string #args=3): substr(s,m,n) gives substring of s from 0-up position m to n
           inclusive. Negative indices -len .. -1 alias to 0 .. len-1.
    
       tolower
           (class=string #args=1): Convert string to lowercase.
    
       toupper
           (class=string #args=1): Convert string to uppercase.
    
       truncate
           (class=string #args=2): Truncates string first argument to max length of int second argument.
    
       capitalize
           (class=string #args=1): Convert string's first character to uppercase.
    
       lstrip
           (class=string #args=1): Strip leading whitespace from string.
    
       rstrip
           (class=string #args=1): Strip trailing whitespace from string.
    
       strip
           (class=string #args=1): Strip leading and trailing whitespace from string.
    
       collapse_whitespace
           (class=string #args=1): Strip repeated whitespace from string.
    
       clean_whitespace
           (class=string #args=1): Same as collapse_whitespace and strip.
    
       system
           (class=string #args=1): Run command string, yielding its stdout minus final carriage return.
    
       abs
           (class=math #args=1): Absolute value.
    
       acos
           (class=math #args=1): Inverse trigonometric cosine.
    
       acosh
           (class=math #args=1): Inverse hyperbolic cosine.
    
       asin
           (class=math #args=1): Inverse trigonometric sine.
    
       asinh
           (class=math #args=1): Inverse hyperbolic sine.
    
       atan
           (class=math #args=1): One-argument arctangent.
    
       atan2
           (class=math #args=2): Two-argument arctangent.
    
       atanh
           (class=math #args=1): Inverse hyperbolic tangent.
    
       cbrt
           (class=math #args=1): Cube root.
    
       ceil
           (class=math #args=1): Ceiling: nearest integer at or above.
    
       cos
           (class=math #args=1): Trigonometric cosine.
    
       cosh
           (class=math #args=1): Hyperbolic cosine.
    
       erf
           (class=math #args=1): Error function.
    
       erfc
           (class=math #args=1): Complementary error function.
    
       exp
           (class=math #args=1): Exponential function e**x.
    
       expm1
           (class=math #args=1): e**x - 1.
    
       floor
           (class=math #args=1): Floor: nearest integer at or below.
    
       invqnorm
           (class=math #args=1): Inverse of normal cumulative distribution
           function. Note that invqorm(urand()) is normally distributed.
    
       log
           (class=math #args=1): Natural (base-e) logarithm.
    
       log10
           (class=math #args=1): Base-10 logarithm.
    
       log1p
           (class=math #args=1): log(1-x).
    
       logifit
           (class=math #args=3): Given m and b from logistic regression, compute
           fit: $yhat=logifit($x,$m,$b).
    
       madd
           (class=math #args=3): a + b mod m (integers)
    
       max
           (class=math variadic): max of n numbers; null loses
    
       mexp
           (class=math #args=3): a ** b mod m (integers)
    
       min
           (class=math variadic): Min of n numbers; null loses
    
       mmul
           (class=math #args=3): a * b mod m (integers)
    
       msub
           (class=math #args=3): a - b mod m (integers)
    
       pow
           (class=math #args=2): Exponentiation; same as **.
    
       qnorm
           (class=math #args=1): Normal cumulative distribution function.
    
       round
           (class=math #args=1): Round to nearest integer.
    
       roundm
           (class=math #args=2): Round to nearest multiple of m: roundm($x,$m) is
           the same as round($x/$m)*$m
    
       sgn
           (class=math #args=1): +1 for positive input, 0 for zero input, -1 for
           negative input.
    
       sin
           (class=math #args=1): Trigonometric sine.
    
       sinh
           (class=math #args=1): Hyperbolic sine.
    
       sqrt
           (class=math #args=1): Square root.
    
       tan
           (class=math #args=1): Trigonometric tangent.
    
       tanh
           (class=math #args=1): Hyperbolic tangent.
    
       urand
           (class=math #args=0): Floating-point numbers uniformly distributed on the unit interval.
           Int-valued example: '$n=floor(20+urand()*11)'.
    
       urandrange
           (class=math #args=2): Floating-point numbers uniformly distributed on the interval [a, b).
    
       urand32
           (class=math #args=0): Integer uniformly distributed 0 and 2**32-1
           inclusive.
    
       urandint
           (class=math #args=2): Integer uniformly distributed between inclusive
           integer endpoints.
    
       dhms2fsec
           (class=time #args=1): Recovers floating-point seconds as in
           dhms2fsec("5d18h53m20.250000s") = 500000.250000
    
       dhms2sec
           (class=time #args=1): Recovers integer seconds as in
           dhms2sec("5d18h53m20s") = 500000
    
       fsec2dhms
           (class=time #args=1): Formats floating-point seconds as in
           fsec2dhms(500000.25) = "5d18h53m20.250000s"
    
       fsec2hms
           (class=time #args=1): Formats floating-point seconds as in
           fsec2hms(5000.25) = "01:23:20.250000"
    
       gmt2sec
           (class=time #args=1): Parses GMT timestamp as integer seconds since
           the epoch.
    
       localtime2sec
           (class=time #args=1): Parses local timestamp as integer seconds since
           the epoch. Consults $TZ environment variable.
    
       hms2fsec
           (class=time #args=1): Recovers floating-point seconds as in
           hms2fsec("01:23:20.250000") = 5000.250000
    
       hms2sec
           (class=time #args=1): Recovers integer seconds as in
           hms2sec("01:23:20") = 5000
    
       sec2dhms
           (class=time #args=1): Formats integer seconds as in sec2dhms(500000)
           = "5d18h53m20s"
    
       sec2gmt
           (class=time #args=1): Formats seconds since epoch (integer part)
           as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
           Leaves non-numbers as-is.
    
           sec2gmt (class=time #args=2): Formats seconds since epoch as GMT timestamp with n
           decimal places for seconds, e.g. sec2gmt(1440768801.7,1) = "2015-08-28T13:33:21.7Z".
           Leaves non-numbers as-is.
    
       sec2gmtdate
           (class=time #args=1): Formats seconds since epoch (integer part)
           as GMT timestamp with year-month-date, e.g. sec2gmtdate(1440768801.7) = "2015-08-28".
           Leaves non-numbers as-is.
    
       sec2localtime
           (class=time #args=1): Formats seconds since epoch (integer part)
           as local timestamp, e.g. sec2localtime(1440768801.7) = "2015-08-28T13:33:21Z".
           Consults $TZ environment variable. Leaves non-numbers as-is.
    
           sec2localtime (class=time #args=2): Formats seconds since epoch as local timestamp with n
           decimal places for seconds, e.g. sec2localtime(1440768801.7,1) = "2015-08-28T13:33:21.7Z".
           Consults $TZ environment variable. Leaves non-numbers as-is.
    
       sec2localdate
           (class=time #args=1): Formats seconds since epoch (integer part)
           as local timestamp with year-month-date, e.g. sec2localdate(1440768801.7) = "2015-08-28".
           Consults $TZ environment variable. Leaves non-numbers as-is.
    
       sec2hms
           (class=time #args=1): Formats integer seconds as in
           sec2hms(5000) = "01:23:20"
    
       strftime
           (class=time #args=2): Formats seconds since the epoch as timestamp, e.g.
           strftime(1440768801.7,"%Y-%m-%dT%H:%M:%SZ") = "2015-08-28T13:33:21Z", and
           strftime(1440768801.7,"%Y-%m-%dT%H:%M:%3SZ") = "2015-08-28T13:33:21.700Z".
           Format strings are as in the C library (please see "man strftime" on your system),
           with the Miller-specific addition of "%1S" through "%9S" which format the seconds
           with 1 through 9 decimal places, respectively. ("%S" uses no decimal places.)
           See also strftime_local.
    
       strftime_local
           (class=time #args=2): Like strftime but consults the $TZ environment variable to get local time zone.
    
       strptime
           (class=time #args=2): Parses timestamp as floating-point seconds since the epoch,
           e.g. strptime("2015-08-28T13:33:21Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.000000,
           and  strptime("2015-08-28T13:33:21.345Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.345000.
           See also strptime_local.
    
       strptime_local
           (class=time #args=2): Like strptime, but consults $TZ environment variable to find and use local timezone.
    
       systime
           (class=time #args=0): Floating-point seconds since the epoch,
           e.g. 1440768801.748936.
    
       is_absent
           (class=typing #args=1): False if field is present in input, true otherwise
    
       is_bool
           (class=typing #args=1): True if field is present with boolean value. Synonymous with is_boolean.
    
       is_boolean
           (class=typing #args=1): True if field is present with boolean value. Synonymous with is_bool.
    
       is_empty
           (class=typing #args=1): True if field is present in input with empty string value, false otherwise.
    
       is_empty_map
           (class=typing #args=1): True if argument is a map which is empty.
    
       is_float
           (class=typing #args=1): True if field is present with value inferred to be float
    
       is_int
           (class=typing #args=1): True if field is present with value inferred to be int
    
       is_map
           (class=typing #args=1): True if argument is a map.
    
       is_nonempty_map
           (class=typing #args=1): True if argument is a map which is non-empty.
    
       is_not_empty
           (class=typing #args=1): False if field is present in input with empty value, true otherwise
    
       is_not_map
           (class=typing #args=1): True if argument is not a map.
    
       is_not_null
           (class=typing #args=1): False if argument is null (empty or absent), true otherwise.
    
       is_null
           (class=typing #args=1): True if argument is null (empty or absent), false otherwise.
    
       is_numeric
           (class=typing #args=1): True if field is present with value inferred to be int or float
    
       is_present
           (class=typing #args=1): True if field is present in input, false otherwise.
    
       is_string
           (class=typing #args=1): True if field is present with string (including empty-string) value
    
       asserting_absent
           (class=typing #args=1): Returns argument if it is absent in the input data, else
           throws an error.
    
       asserting_bool
           (class=typing #args=1): Returns argument if it is present with boolean value, else
           throws an error.
    
       asserting_boolean
           (class=typing #args=1): Returns argument if it is present with boolean value, else
           throws an error.
    
       asserting_empty
           (class=typing #args=1): Returns argument if it is present in input with empty value,
           else throws an error.
    
       asserting_empty_map
           (class=typing #args=1): Returns argument if it is a map with empty value, else
           throws an error.
    
       asserting_float
           (class=typing #args=1): Returns argument if it is present with float value, else
           throws an error.
    
       asserting_int
           (class=typing #args=1): Returns argument if it is present with int value, else
           throws an error.
    
       asserting_map
           (class=typing #args=1): Returns argument if it is a map, else throws an error.
    
       asserting_nonempty_map
           (class=typing #args=1): Returns argument if it is a non-empty map, else throws
           an error.
    
       asserting_not_empty
           (class=typing #args=1): Returns argument if it is present in input with non-empty
           value, else throws an error.
    
       asserting_not_map
           (class=typing #args=1): Returns argument if it is not a map, else throws an error.
    
       asserting_not_null
           (class=typing #args=1): Returns argument if it is non-null (non-empty and non-absent),
           else throws an error.
    
       asserting_null
           (class=typing #args=1): Returns argument if it is null (empty or absent), else throws
           an error.
    
       asserting_numeric
           (class=typing #args=1): Returns argument if it is present with int or float value,
           else throws an error.
    
       asserting_present
           (class=typing #args=1): Returns argument if it is present in input, else throws
           an error.
    
       asserting_string
           (class=typing #args=1): Returns argument if it is present with string (including
           empty-string) value, else throws an error.
    
       boolean
           (class=conversion #args=1): Convert int/float/bool/string to boolean.
    
       float
           (class=conversion #args=1): Convert int/float/bool/string to float.
    
       fmtnum
           (class=conversion #args=2): Convert int/float/bool to string using
           printf-style format string, e.g. '$s = fmtnum($n, "%06lld")'. WARNING: Miller numbers
           are all long long or double. If you use formats like %d or %f, behavior is undefined.
    
       hexfmt
           (class=conversion #args=1): Convert int to string, e.g. 255 to "0xff".
    
       int
           (class=conversion #args=1): Convert int/float/bool/string to int.
    
       string
           (class=conversion #args=1): Convert int/float/bool/string to string.
    
       typeof
           (class=conversion #args=1): Convert argument to type of argument (e.g.
           MT_STRING). For debug.
    
       depth
           (class=maps #args=1): Prints maximum depth of hashmap: ''. Scalars have depth 0.
    
       haskey
           (class=maps #args=2): True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or
           'haskey(mymap, mykey)'. Error if 1st argument is not a map.
    
       joink
           (class=maps #args=2): Makes string from map keys. E.g. 'joink($*, ",")'.
    
       joinkv
           (class=maps #args=3): Makes string from map key-value pairs. E.g. 'joinkv(@v[2], "=", ",")'
    
       joinv
           (class=maps #args=2): Makes string from map values. E.g. 'joinv(mymap, ",")'.
    
       leafcount
           (class=maps #args=1): Counts total number of terminal values in hashmap. For single-level maps,
           same as length.
    
       length
           (class=maps #args=1): Counts number of top-level entries in hashmap. Scalars have length 1.
    
       mapdiff
           (class=maps variadic): With 0 args, returns empty map. With 1 arg, returns copy of arg.
           With 2 or more, returns copy of arg 1 with all keys from any of remaining argument maps removed.
    
       mapexcept
           (class=maps variadic): Returns a map with keys from remaining arguments, if any, unset.
           E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'.
    
       mapselect
           (class=maps variadic): Returns a map with only keys from remaining arguments set.
           E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'.
    
       mapsum
           (class=maps variadic): With 0 args, returns empty map. With >= 1 arg, returns a map with
           key-value pairs from all arguments. Rightmost collisions win, e.g. 'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.
    
       splitkv
           (class=maps #args=3): Splits string by separators into map with type inference.
           E.g. 'splitkv("a=1,b=2,c=3", "=", ",")' gives '{"a" : 1, "b" : 2, "c" : 3}'.
    
       splitkvx
           (class=maps #args=3): Splits string by separators into map without type inference (keys and
           values are strings). E.g. 'splitkv("a=1,b=2,c=3", "=", ",")' gives
           '{"a" : "1", "b" : "2", "c" : "3"}'.
    
       splitnv
           (class=maps #args=2): Splits string by separator into integer-indexed map with type inference.
           E.g. 'splitnv("a,b,c" , ",")' gives '{1 : "a", 2 : "b", 3 : "c"}'.
    
       splitnvx
           (class=maps #args=2): Splits string by separator into integer-indexed map without type
           inference (values are strings). E.g. 'splitnv("4,5,6" , ",")' gives '{1 : "4", 2 : "5", 3 : "6"}'.
    
    KEYWORDS FOR PUT AND FILTER
       all
           all: used in "emit", "emitp", and "unset" as a synonym for @*
    
       begin
           begin: defines a block of statements to be executed before input records
           are ingested. The body statements must be wrapped in curly braces.
           Example: 'begin { @count = 0 }'
    
       bool
           bool: declares a boolean local variable in the current curly-braced scope.
           Type-checking happens at assignment: 'bool b = 1' is an error.
    
       break
           break: causes execution to continue after the body of the current
           for/while/do-while loop.
    
       call
           call: used for invoking a user-defined subroutine.
           Example: 'subr s(k,v) { print k . " is " . v} call s("a", $a)'
    
       continue
           continue: causes execution to skip the remaining statements in the body of
           the current for/while/do-while loop. For-loop increments are still applied.
    
       do
           do: with "while", introduces a do-while loop. The body statements must be wrapped
           in curly braces.
    
       dump
           dump: prints all currently defined out-of-stream variables immediately
    	 to stdout as JSON.
    
    	 With >, >>, or |, the data do not become part of the output record stream but
    	 are instead redirected.
    
    	 The > and >> are for write and append, as in the shell, but (as with awk) the
    	 file-overwrite for > is on first write, not per record. The | is for piping to
    	 a process which will process the data. There will be one open file for each
    	 distinct file name (for > and >>) or one subordinate process for each distinct
    	 value of the piped-to command (for |). Output-formatting flags are taken from
    	 the main command line.
    
    	 Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump }'
    	 Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >  "mytap.dat"}'
    	 Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >> "mytap.dat"}'
    	 Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump | "jq .[]"}'
    
       edump
           edump: prints all currently defined out-of-stream variables immediately
    	 to stderr as JSON.
    
    	 Example: mlr --from f.dat put -q '@v[NR]=$*; end { edump }'
    
       elif
           elif: the way Miller spells "else if". The body statements must be wrapped
           in curly braces.
    
       else
           else: terminates an if/elif/elif chain. The body statements must be wrapped
           in curly braces.
    
       emit
           emit: inserts an out-of-stream variable into the output record stream. Hashmap
    	 indices present in the data but not slotted by emit arguments are not output.
    
    	 With >, >>, or |, the data do not become part of the output record stream but
    	 are instead redirected.
    
    	 The > and >> are for write and append, as in the shell, but (as with awk) the
    	 file-overwrite for > is on first write, not per record. The | is for piping to
    	 a process which will process the data. There will be one open file for each
    	 distinct file name (for > and >>) or one subordinate process for each distinct
    	 value of the piped-to command (for |). Output-formatting flags are taken from
    	 the main command line.
    
    	 You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
    	 etc., to control the format of the output if the output is redirected. See also mlr -h.
    
    	 Example: mlr --from f.dat put 'emit >	"/tmp/data-".$a, $*'
    	 Example: mlr --from f.dat put 'emit >	"/tmp/data-".$a, mapexcept($*, "a")'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums'
    	 Example: mlr --from f.dat put --ojson '@sums[$a][$b]+=$x; emit > "tap-".$a.$b.".dat", @sums'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >  "mytap.dat", @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >> "mytap.dat", @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | "gzip > mytap.dat.gz", @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit > stderr, @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | "grep somepattern", @*, "index1", "index2"'
    
    	 Please see http://johnkerl.org/miller/doc for more information.
    
       emitf
           emitf: inserts non-indexed out-of-stream variable(s) side-by-side into the
    	 output record stream.
    
    	 With >, >>, or |, the data do not become part of the output record stream but
    	 are instead redirected.
    
    	 The > and >> are for write and append, as in the shell, but (as with awk) the
    	 file-overwrite for > is on first write, not per record. The | is for piping to
    	 a process which will process the data. There will be one open file for each
    	 distinct file name (for > and >>) or one subordinate process for each distinct
    	 value of the piped-to command (for |). Output-formatting flags are taken from
    	 the main command line.
    
    	 You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
    	 etc., to control the format of the output if the output is redirected. See also mlr -h.
    
    	 Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a'
    	 Example: mlr --from f.dat put --oxtab '@a=$i;@b+=$x;@c+=$y; emitf > "tap-".$i.".dat", @a'
    	 Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a, @b, @c'
    	 Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > "mytap.dat", @a, @b, @c'
    	 Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf >> "mytap.dat", @a, @b, @c'
    	 Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > stderr, @a, @b, @c'
    	 Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | "grep somepattern", @a, @b, @c'
    	 Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | "grep somepattern > mytap.dat", @a, @b, @c'
    
    	 Please see http://johnkerl.org/miller/doc for more information.
    
       emitp
           emitp: inserts an out-of-stream variable into the output record stream.
    	 Hashmap indices present in the data but not slotted by emitp arguments are
    	 output concatenated with ":".
    
    	 With >, >>, or |, the data do not become part of the output record stream but
    	 are instead redirected.
    
    	 The > and >> are for write and append, as in the shell, but (as with awk) the
    	 file-overwrite for > is on first write, not per record. The | is for piping to
    	 a process which will process the data. There will be one open file for each
    	 distinct file name (for > and >>) or one subordinate process for each distinct
    	 value of the piped-to command (for |). Output-formatting flags are taken from
    	 the main command line.
    
    	 You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
    	 etc., to control the format of the output if the output is redirected. See also mlr -h.
    
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums'
    	 Example: mlr --from f.dat put --opprint '@sums[$a][$b]+=$x; emitp > "tap-".$a.$b.".dat", @sums'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >  "mytap.dat", @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >> "mytap.dat", @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | "gzip > mytap.dat.gz", @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp > stderr, @*, "index1", "index2"'
    	 Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | "grep somepattern", @*, "index1", "index2"'
    
    	 Please see http://johnkerl.org/miller/doc for more information.
    
       end
           end: defines a block of statements to be executed after input records
           are ingested. The body statements must be wrapped in curly braces.
           Example: 'end { emit @count }'
           Example: 'end { eprint "Final count is " . @count }'
    
       eprint
           eprint: prints expression immediately to stderr.
    	 Example: mlr --from f.dat put -q 'eprint "The sum of x and y is ".($x+$y)'
    	 Example: mlr --from f.dat put -q 'for (k, v in $*) { eprint k . " => " . v }'
    	 Example: mlr --from f.dat put	'(NR % 1000 == 0) { eprint "Checkpoint ".NR}'
    
       eprintn
           eprintn: prints expression immediately to stderr, without trailing newline.
    	 Example: mlr --from f.dat put -q 'eprintn "The sum of x and y is ".($x+$y); eprint ""'
    
       false
           false: the boolean literal value.
    
       filter
           filter: includes/excludes the record in the output record stream.
    
    	 Example: mlr --from f.dat put 'filter (NR == 2 || $x > 5.4)'
    
    	 Instead of put with 'filter false' you can simply use put -q.	The following
    	 uses the input record to accumulate data but only prints the running sum
    	 without printing the input record:
    
    	 Example: mlr --from f.dat put -q '@running_sum += $x * $y; emit @running_sum'
    
       float
           float: declares a floating-point local variable in the current curly-braced scope.
           Type-checking happens at assignment: 'float x = 0' is an error.
    
       for
           for: defines a for-loop using one of three styles. The body statements must
           be wrapped in curly braces.
           For-loop over stream record:
    	 Example:  'for (k, v in $*) { ... }'
           For-loop over out-of-stream variables:
    	 Example: 'for (k, v in @counts) { ... }'
    	 Example: 'for ((k1, k2), v in @counts) { ... }'
    	 Example: 'for ((k1, k2, k3), v in @*) { ... }'
           C-style for-loop:
    	 Example:  'for (var i = 0, var b = 1; i < 10; i += 1, b *= 2) { ... }'
    
       func
           func: used for defining a user-defined function.
           Example: 'func f(a,b) { return sqrt(a**2+b**2)} $d = f($x, $y)'
    
       if
           if: starts an if/elif/elif chain. The body statements must be wrapped
           in curly braces.
    
       in
           in: used in for-loops over stream records or out-of-stream variables.
    
       int
           int: declares an integer local variable in the current curly-braced scope.
           Type-checking happens at assignment: 'int x = 0.0' is an error.
    
       map
           map: declares an map-valued local variable in the current curly-braced scope.
           Type-checking happens at assignment: 'map b = 0' is an error. map b = {} is
           always OK. map b = a is OK or not depending on whether a is a map.
    
       num
           num: declares an int/float local variable in the current curly-braced scope.
           Type-checking happens at assignment: 'num b = true' is an error.
    
       print
           print: prints expression immediately to stdout.
    	 Example: mlr --from f.dat put -q 'print "The sum of x and y is ".($x+$y)'
    	 Example: mlr --from f.dat put -q 'for (k, v in $*) { print k . " => " . v }'
    	 Example: mlr --from f.dat put	'(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'
    
       printn
           printn: prints expression immediately to stdout, without trailing newline.
    	 Example: mlr --from f.dat put -q 'printn "."; end { print "" }'
    
       return
           return: specifies the return value from a user-defined function.
           Omitted return statements (including via if-branches) result in an absent-null
           return value, which in turns results in a skipped assignment to an LHS.
    
       stderr
           stderr: Used for tee, emit, emitf, emitp, print, and dump in place of filename
    	 to print to standard error.
    
       stdout
           stdout: Used for tee, emit, emitf, emitp, print, and dump in place of filename
    	 to print to standard output.
    
       str
           str: declares a string local variable in the current curly-braced scope.
           Type-checking happens at assignment.
    
       subr
           subr: used for defining a subroutine.
           Example: 'subr s(k,v) { print k . " is " . v} call s("a", $a)'
    
       tee
           tee: prints the current record to specified file.
    	 This is an immediate print to the specified file (except for pprint format
    	 which of course waits until the end of the input stream to format all output).
    
    	 The > and >> are for write and append, as in the shell, but (as with awk) the
    	 file-overwrite for > is on first write, not per record. The | is for piping to
    	 a process which will process the data. There will be one open file for each
    	 distinct file name (for > and >>) or one subordinate process for each distinct
    	 value of the piped-to command (for |). Output-formatting flags are taken from
    	 the main command line.
    
    	 You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
    	 etc., to control the format of the output. See also mlr -h.
    
    	 emit with redirect and tee with redirect are identical, except tee can only
    	 output $*.
    
    	 Example: mlr --from f.dat put 'tee >  "/tmp/data-".$a, $*'
    	 Example: mlr --from f.dat put 'tee >> "/tmp/data-".$a.$b, $*'
    	 Example: mlr --from f.dat put 'tee >  stderr, $*'
    	 Example: mlr --from f.dat put -q 'tee | "tr [a-z\] [A-Z\]", $*'
    	 Example: mlr --from f.dat put -q 'tee | "tr [a-z\] [A-Z\] > /tmp/data-".$a, $*'
    	 Example: mlr --from f.dat put -q 'tee | "gzip > /tmp/data-".$a.".gz", $*'
    	 Example: mlr --from f.dat put -q --ojson 'tee | "gzip > /tmp/data-".$a.".gz", $*'
    
       true
           true: the boolean literal value.
    
       unset
           unset: clears field(s) from the current record, or an out-of-stream or local variable.
    
    	 Example: mlr --from f.dat put 'unset $x'
    	 Example: mlr --from f.dat put 'unset $*'
    	 Example: mlr --from f.dat put 'for (k, v in $*) { if (k =~ "a.*") { unset $[k] } }'
    	 Example: mlr --from f.dat put '...; unset @sums'
    	 Example: mlr --from f.dat put '...; unset @sums["green"]'
    	 Example: mlr --from f.dat put '...; unset @*'
    
       var
           var: declares an untyped local variable in the current curly-braced scope.
           Examples: 'var a=1', 'var xyz=""'
    
       while
           while: introduces a while loop, or with "do", introduces a do-while loop.
           The body statements must be wrapped in curly braces.
    
       ENV
           ENV: access to environment variables by name, e.g. '$home = ENV["HOME"]'
    
       FILENAME
           FILENAME: evaluates to the name of the current file being processed.
    
       FILENUM
           FILENUM: evaluates to the number of the current file being processed,
           starting with 1.
    
       FNR
           FNR: evaluates to the number of the current record within the current file
           being processed, starting with 1. Resets at the start of each file.
    
       IFS
           IFS: evaluates to the input field separator from the command line.
    
       IPS
           IPS: evaluates to the input pair separator from the command line.
    
       IRS
           IRS: evaluates to the input record separator from the command line,
           or to LF or CRLF from the input data if in autodetect mode (which is
           the default).
    
       M_E
           M_E: the mathematical constant e.
    
       M_PI
           M_PI: the mathematical constant pi.
    
       NF
           NF: evaluates to the number of fields in the current record.
    
       NR
           NR: evaluates to the number of the current record over all files
           being processed, starting with 1. Does not reset at the start of each file.
    
       OFS
           OFS: evaluates to the output field separator from the command line.
    
       OPS
           OPS: evaluates to the output pair separator from the command line.
    
       ORS
           ORS: evaluates to the output record separator from the command line,
           or to LF or CRLF from the input data if in autodetect mode (which is
           the default).
    
    AUTHOR
           Miller is written by John Kerl <kerl.john.r@gmail.com>.
    
           This manual page has been composed from Miller's help output by Eric
           MSP Veith <eveith@veith-m.de>.
    
    SEE ALSO
           awk(1), sed(1), cut(1), join(1), sort(1), RFC 4180: Common Format and
           MIME Type for Comma-Separated Values (CSV) Files, the miller website
           http://johnkerl.org/miller/doc
    
    
    
    				  2021-03-22			     MILLER(1)
