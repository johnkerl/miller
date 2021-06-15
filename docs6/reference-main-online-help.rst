..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

Reference: online help
================================================================

TODO: expand this section

Examples:

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr --help
    Usage: mlr [I/O options] {verb} [verb-dependent options ...] {zero or more file names}
    
    COMMAND-LINE-SYNTAX EXAMPLES:
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
      mlr --from infile.dat put  '(NR % 1000 == 0) { print > os.Stderr, "Checkpoint ".NR}'
    
    DATA-FORMAT EXAMPLES:
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
    
    HELP OPTIONS:
      -h or --help                 Show this message.
      --version                    Show the software version.
      {verb name} --help           Show verb-specific help.
      --help-all-verbs             Show help on all verbs.
      -l or --list-all-verbs       List only verb names.
      -L                           List only verb names, one per line.
      -f or --help-all-functions   Show help on all built-in functions.
      -F                           Show a bare listing of built-in functions by name.
      -k or --help-all-keywords    Show help on all keywords.
      -K                           Show a bare listing of keywords by name.
    
    CUSTOMIZATION VIA .MLRRC:
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
    https://miller.readthedocs.io/en/latest/customization.html
    
    VERBS:
       altkv bar bootstrap cat check clean-whitespace count-distinct count
       count-similar cut decimate fill-down fill-empty filter flatten format-values
       fraction gap grep group-by group-like having-fields head histogram
       json-parse json-stringify join label least-frequent merge-fields
       most-frequent nest nothing put regularize remove-empty-columns rename
       reorder repeat reshape sample sec2gmtdate sec2gmt seqgen shuffle
       skip-trivial-records sort sort-within-records stats1 stats2 step tac tail
       tee top unflatten uniq unsparsify
    
    FUNCTIONS FOR THE FILTER AND PUT VERBS:
    +
    -
    *
    /
    //
    **
    pow
    .+
    .-
    .*
    ./
    %
    ~
    &
    |
    ^
    <<
    >>
    >>>
    bitcount
    madd
    msub
    mmul
    mexp
    !
    ==
    !=
    >
    >=
    <
    <=
    =~
    !=~
    &&
    ||
    ^^
    ??
    ???
    ?:
    .
    capitalize
    clean_whitespace
    collapse_whitespace
    gsub
    lstrip
    regextract
    regextract_or_else
    rstrip
    strip
    strlen
    ssub
    sub
    substr0
    substr1
    substr
    tolower
    toupper
    truncate
    md5
    sha1
    sha256
    sha512
    abs
    acos
    acosh
    asin
    asinh
    atan
    atan2
    atanh
    cbrt
    ceil
    cos
    cosh
    erf
    erfc
    exp
    expm1
    floor
    invqnorm
    log
    log10
    log1p
    logifit
    max
    min
    qnorm
    round
    sgn
    sin
    sinh
    sqrt
    tan
    tanh
    roundm
    urand
    urandint
    urandrange
    urand32
    gmt2sec
    sec2gmt
    sec2gmtdate
    systime
    systimeint
    uptime
    strftime
    strptime
    dhms2fsec
    dhms2sec
    fsec2dhms
    fsec2hms
    hms2fsec
    hms2sec
    sec2dhms
    sec2hms
    is_absent
    is_array
    is_bool
    is_boolean
    is_empty
    is_empty_map
    is_error
    is_float
    is_int
    is_map
    is_nonempty_map
    is_not_empty
    is_not_map
    is_not_array
    is_not_null
    is_null
    is_numeric
    is_present
    is_string
    asserting_absent
    asserting_array
    asserting_bool
    asserting_boolean
    asserting_error
    asserting_empty
    asserting_empty_map
    asserting_float
    asserting_int
    asserting_map
    asserting_nonempty_map
    asserting_not_empty
    asserting_not_map
    asserting_not_array
    asserting_not_null
    asserting_null
    asserting_numeric
    asserting_present
    asserting_string
    typeof
    boolean
    float
    fmtnum
    hexfmt
    int
    joink
    joinv
    joinkv
    splita
    splitax
    splitkv
    splitkvx
    splitnv
    splitnvx
    string
    append
    arrayify
    depth
    flatten
    get_keys
    get_values
    haskey
    json_parse
    json_stringify
    leafcount
    length
    mapdiff
    mapexcept
    mapselect
    mapsum
    unflatten
    hostname
    os
    system
    version
    Please use "mlr --help-function {function name}" for function-specific help.
    
    DATA-FORMAT OPTIONS, FOR INPUT, OUTPUT, OR BOTH:
    
    	  --idkvp   --odkvp   --dkvp      Delimited key-value pairs, e.g "a=1,b=2"
    	                                  (this is Miller's default format).
    
    	  --inidx   --onidx   --nidx      Implicitly-integer-indexed fields
    	                                  (Unix-toolkit style).
    	  -T                              Synonymous with "--nidx --fs tab".
    
    	  --icsv    --ocsv    --csv       Comma-separated value (or tab-separated
    	                                  with --fs tab, etc.)
    
    	  --itsv    --otsv    --tsv       Keystroke-savers for "--icsv --ifs tab",
    	                                  "--ocsv --ofs tab", "--csv --fs tab".
    	  --iasv    --oasv    --asv       Similar but using ASCII FS 0x1f and RS 0x1e\n",
    	  --iusv    --ousv    --usv       Similar but using Unicode FS U+241F (UTF-8 0xe2909f)\n",
    	                                  and RS U+241E (UTF-8 0xe2909e)\n",
    
    	  --icsvlite --ocsvlite --csvlite Comma-separated value (or tab-separated
    	                                  with --fs tab, etc.). The 'lite' CSV does not handle
    	                                  RFC-CSV double-quoting rules; is slightly faster;
    	                                  and handles heterogeneity in the input stream via
    	                                  empty newline followed by new header line. See also
    	                                  http://johnkerl.org/miller/doc/file-formats.html#CSV/TSV/etc.
    
    	  --itsvlite --otsvlite --tsvlite Keystroke-savers for "--icsvlite --ifs tab",
    	                                  "--ocsvlite --ofs tab", "--csvlite --fs tab".
    	  -t                              Synonymous with --tsvlite.
    	  --iasvlite --oasvlite --asvlite Similar to --itsvlite et al. but using ASCII FS 0x1f and RS 0x1e\n",
    	  --iusvlite --ousvlite --usvlite Similar to --itsvlite et al. but using Unicode FS U+241F (UTF-8 0xe2909f)\n",
    	                                  and RS U+241E (UTF-8 0xe2909e)\n",
    
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
    	                      --jvstack   Put one key-value pair per line for JSON output.
    	                   --no-jvstack   Put objects/arrays all on one line for JSON output.
    	                --jsonx --ojsonx  Keystroke-savers for --json --jvstack
    	                --jsonx --ojsonx  and --ojson --jvstack, respectively.
    	                      --jlistwrap Wrap JSON output in outermost [ ].
    	                    --jknquoteint Do not quote non-string map keys in JSON output.
    	                     --jvquoteall Quote map values in JSON output, even if they're
    	                                  numeric.
    	              --oflatsep {string} Separator for flattening multi-level JSON keys,
    	                                  e.g. '{"a":{"b":3}}' becomes a:b => 3 for
    	                                  non-JSON formats. Defaults to ..\n",
    
    	  -p is a keystroke-saver for --nidx --fs space --repifs
    
    	  Examples: --csv for CSV-formatted input and output; --idkvp --opprint for
    	  DKVP-formatted input and pretty-printed output.
    
    	  Please use --iformat1 --oformat2 rather than --format1 --oformat2.
    	  The latter sets up input and output flags for format1, not all of which
    	  are overridden in all cases by setting output format to format2.
    
    
    COMMENTS IN DATA:
      --skip-comments                 Ignore commented lines (prefixed by "#")
                                      within the input.
      --skip-comments-with {string}   Ignore commented lines within input, with
                                      specified prefix.
      --pass-comments                 Immediately print commented lines (prefixed by "#")
                                      within the input.
      --pass-comments-with {string}   Immediately print commented lines within input, with
                                      specified prefix.
    Notes:
    * Comments are only honored at the start of a line.
    * In the absence of any of the above four options, comments are data like
      any other text.
    * When pass-comments is used, comment lines are written to standard output
      immediately upon being read; they are not part of the record stream.
      Results may be counterintuitive. A suggestion is to place comments at the
      start of data files.
    
    FORMAT-CONVERSION KEYSTROKE-SAVER OPTIONS:
    As keystroke-savers for format-conversion you may use the following:
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
    
    COMPRESSED-DATA OPTIONS:
      Decompression done within the Miller process itself:
      --gzin  Uncompress gzip within the Miller process. Done by default if file ends in ".gz".
      --bz2in Uncompress bz2ip within the Miller process. Done by default if file ends in ".bz2".
      --zin   Uncompress zlib within the Miller process. Done by default if file ends in ".z".
      Decompression done outside the Miller processn  --prepipe {command} You can, of course, already do without this for single input files,
      e.g. "gunzip < myfile.csv.gz | mlr ...".
      However, when multiple input files are present, between-file separations are
      lost; also, the FILENAME variable doesn't iterate. Using --prepipe you can
      specify an action to be taken on each input file. This prepipe command must
      be able to read from standard input; it will be invoked with
        {command} < {filename}.
      --prepipex {command} Like --prepipe with one exception: doesn't insert '<' between
      command and filename at runtime. Useful for some commands like 'unzip -qc' which don't
      read standard input.
      Examples:
        mlr --prepipe 'gunzip'
        mlr --prepipe 'zcat -cf'
        mlr --prepipe 'xz -cd'
        mlr --prepipe cat
      Note that this feature is quite general and is not limited to decompression
      utilities. You can use it to apply per-file filters of your choice.
      For output compression (or other) utilities, simply pipe the output:
        mlr ... | {your compression command}
      Lastly, note that if --prepipe is specified, it replaces any decisions that might
      have been made based on the file suffix. Also, --gzin/--bz2in/--zin are ignored
      if --prepipe is also specified.
    
    RELEVANT TO CSV/CSV-LITE INPUT ONLY:
      --implicit-csv-header Use 1,2,3,... as field labels, rather than from line 1
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
    
    NUMERICAL FORMATTING:
      --ofmt {format}    E.g. %.18f, %.0f, %9.6e. Please use sprintf-style codes for
                         floating-point nummbers. If not specified, default formatting is used.
                         See also the fmtnum function within mlr put (mlr --help-all-functions);
                         see also the format-values function.
    
    OUTPUT COLORIZATION:
    Things having colors:
    * Keys in CSV header lines, JSON keys, etc
    * Values in CSV data lines, JSON scalar values, etc
    * "PASS" and "FAIL" in regression-test output
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
      * Command-line flags ``--no-color`` or ``-M``, ``--always-color`` or ``-C``.
    * Color choices can be specified by using environment variables, or command-line flags,
      with values 0..255:
      * export MLR_KEY_COLOR=208, MLR_VALUE_COLOR-33, etc.
      * Command-line flags --key-color 208, --value-color 33, etc.
      * This is particularly useful if your terminal's background color clashes with current settings.
    * If environment-variable settings and command-line flags are both provided,the latter take precedence.
    * Please do mlr --list-colors to see the available color codes.
    
    OTHER OPTIONS:
      --seed {n} with n of the form 12345678 or 0xcafefeed. For put/filter
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
                         CSV, CSV headers will be present in each output file;
                         statistics are only over each file's own records; and so on.
    
    THEN-CHAINING:
    Output of one verb may be chained as input to another using "then", e.g.
      mlr stats1 -a min,mean,max -f flag,u,v -g color then sort -f color
    
    AUXILIARY COMMANDS:
    Miller has a few otherwise-standalone executables packaged within it.
    They do not participate in any other parts of Miller.
    Please use "mlr aux-list" for more information.
    
    SEE ALSO:
    For more information please see http://johnkerl.org/miller/doc and/or
    http://github.com/johnkerl/miller. This is Miller version v6.0.0-dev.

.. code-block:: none
   :emphasize-lines: 1-1

    $ mlr sort --help
    Usage: mlr sort {flags}
    Sorts records primarily by the first specified field, secondarily by the second
    field, and so on.  (Any records not having all specified sort keys will appear
    at the end of the output, in the order they were encountered, regardless of the
    specified sort order.) The sort is stable: records that compare equal will sort
    in the order they were encountered in the input record stream.
    
    Options:
    -f  {comma-separated field names}  Lexical ascending
    -n  {comma-separated field names}  Numerical ascending; nulls sort last
    -nf {comma-separated field names}  Same as -n
    -r  {comma-separated field names}  Lexical descending
    -nr {comma-separated field names}  Numerical descending; nulls sort first
    -h|--help Show this message.
    
    Example:
      mlr sort -f a,b -nr x,y,z
    which is the same as:
      mlr sort -f a -f b -nr x -nr y -nr z
