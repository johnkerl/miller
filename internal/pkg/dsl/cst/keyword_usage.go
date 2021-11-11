package cst

import (
	"fmt"

	"mlr/internal/pkg/colorizer"
	"mlr/internal/pkg/lib"
)

// ----------------------------------------------------------------
type tKeywordUsageFunc func()

type tKeywordUsageEntry struct {
	name      string
	usageFunc tKeywordUsageFunc
}

var KEYWORD_USAGE_TABLE = []tKeywordUsageEntry{
	{"all", allKeywordUsage},
	{"begin", beginKeywordUsage},
	{"bool", boolKeywordUsage},
	{"break", breakKeywordUsage},
	{"call", callKeywordUsage},
	{"continue", continueKeywordUsage},
	{"do", doKeywordUsage},
	{"dump", dumpKeywordUsage},
	{"edump", edumpKeywordUsage},
	{"elif", elifKeywordUsage},
	{"else", elseKeywordUsage},
	{"emit1", emit1KeywordUsage},
	{"emit", emitKeywordUsage},
	{"emitf", emitfKeywordUsage},
	{"emitp", emitpKeywordUsage},
	{"end", endKeywordUsage},
	{"eprint", eprintKeywordUsage},
	{"eprintn", eprintnKeywordUsage},
	{"false", falseKeywordUsage},
	{"filter", filterKeywordUsage},
	{"float", floatKeywordUsage},
	{"for", forKeywordUsage},
	{"func", funcKeywordUsage},
	{"if", ifKeywordUsage},
	{"in", inKeywordUsage},
	{"int", intKeywordUsage},
	{"map", mapKeywordUsage},
	{"num", numKeywordUsage},
	{"print", printKeywordUsage},
	{"printn", printnKeywordUsage},
	{"return", returnKeywordUsage},
	{"stderr", stderrKeywordUsage},
	{"stdout", stdoutKeywordUsage},
	{"str", strKeywordUsage},
	{"subr", subrKeywordUsage},
	{"tee", teeKeywordUsage},
	{"true", trueKeywordUsage},
	{"unset", unsetKeywordUsage},
	{"var", varKeywordUsage},
	{"while", whileKeywordUsage},
	{"ENV", ENVKeywordUsage},
	{"FILENAME", FILENAMEKeywordUsage},
	{"FILENUM", FILENUMKeywordUsage},
	{"FNR", FNRKeywordUsage},
	{"IFS", IFSKeywordUsage},
	{"IPS", IPSKeywordUsage},
	{"IRS", IRSKeywordUsage},
	{"M_E", M_EKeywordUsage},
	{"M_PI", M_PIKeywordUsage},
	{"NF", NFKeywordUsage},
	{"NR", NRKeywordUsage},
	{"OFS", OFSKeywordUsage},
	{"OPS", OPSKeywordUsage},
	{"ORS", ORSKeywordUsage},
}

// ----------------------------------------------------------------

// Pass function_name == NULL to get usage for all keywords.
func UsageKeywords() {
	for i, entry := range KEYWORD_USAGE_TABLE {
		if i > 0 {
			fmt.Println()
		}
		fmt.Printf("%s: ", colorizer.MaybeColorizeHelp(entry.name, true))
		entry.usageFunc()
	}
}

func UsageForKeyword(name string) {
	if !TryUsageForKeyword(name) {
		fmt.Printf("mlr: unrecognized keyword \"%s\".\n", name)
	}
}

func TryUsageForKeyword(name string) bool {
	for _, entry := range KEYWORD_USAGE_TABLE {
		if entry.name == name {
			fmt.Printf("%s: ", colorizer.MaybeColorizeHelp(entry.name, true))
			entry.usageFunc()
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------
func ListKeywordsVertically() {
	for _, entry := range KEYWORD_USAGE_TABLE {
		fmt.Println(entry.name)
	}
}

// ----------------------------------------------------------------
func ListKeywordsAsParagraph() {
	keywords := make([]string, len(KEYWORD_USAGE_TABLE))
	for i, entry := range KEYWORD_USAGE_TABLE {
		keywords[i] = entry.name
	}
	lib.PrintWordsAsParagraph(keywords)
}

// ----------------------------------------------------------------
func allKeywordUsage() {
	fmt.Println(
		`used in "emit1", "emit", "emitp", and "unset" as a synonym for @*`,
	)
}

func beginKeywordUsage() {
	fmt.Println(
		`defines a block of statements to be executed before input records
are ingested. The body statements must be wrapped in curly braces.

  Example: 'begin { @count = 0 }'`)
}

func boolKeywordUsage() {
	fmt.Println(
		`declares a boolean local variable in the current curly-braced scope.
Type-checking happens at assignment: 'bool b = 1' is an error.`)
}

func breakKeywordUsage() {
	fmt.Println(
		`causes execution to continue after the body of the current for/while/do-while loop.`)
}

func callKeywordUsage() {
	fmt.Println(
		`used for invoking a user-defined subroutine.

  Example: 'subr s(k,v) { print k . " is " . v} call s("a", $a)'`)
}

func continueKeywordUsage() {
	fmt.Println(
		`causes execution to skip the remaining statements in the body of
the current for/while/do-while loop. For-loop increments are still applied.`)
}

func doKeywordUsage() {
	fmt.Println(
		`with "while", introduces a do-while loop. The body statements must be wrapped
in curly braces.`)
}

func dumpKeywordUsage() {
	fmt.Println(
		`prints all currently defined out-of-stream variables immediately
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
  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump | "jq .[]"}'`)
}

func edumpKeywordUsage() {
	fmt.Println(
		`prints all currently defined out-of-stream variables immediately
to stderr as JSON.

  Example: mlr --from f.dat put -q '@v[NR]=$*; end { edump }'`)
}

func elifKeywordUsage() {
	fmt.Println(
		`the way Miller spells "else if". The body statements must be wrapped
in curly braces.`)
}

func elseKeywordUsage() {
	fmt.Println(
		`terminates an if/elif/elif chain. The body statements must be wrapped
in curly braces.`)
}

func emit1KeywordUsage() {
	fmt.Printf(
		`inserts an out-of-stream variable into the output record stream. Unlike
the other map variants, side-by-sides, indexing, and redirection are not supported,
but you can emit any map-valued expression.

  Example: mlr --from f.dat put 'emit1 $*'
  Example: mlr --from f.dat put 'emit1 mapsum({"id": NR}, $*)'

Please see %s://johnkerl.org/miller/doc for more information.
`, lib.DOC_URL)
}

func emitKeywordUsage() {
	fmt.Printf(
		`inserts an out-of-stream variable into the output record stream. Hashmap
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

  Example: mlr --from f.dat put 'emit >  "/tmp/data-".$a, $*'
  Example: mlr --from f.dat put 'emit >  "/tmp/data-".$a, mapexcept($*, "a")'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums'
  Example: mlr --from f.dat put --ojson '@sums[$a][$b]+=$x; emit > "tap-".$a.$b.".dat", @sums'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >  "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >> "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | "gzip > mytap.dat.gz", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit > stderr, @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | "grep somepattern", @*, "index1", "index2"'

Please see %s://johnkerl.org/miller/doc for more information.
`, lib.DOC_URL)
}

func emitfKeywordUsage() {
	fmt.Printf(
		`inserts non-indexed out-of-stream variable(s) side-by-side into the
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

Please see %s://johnkerl.org/miller/doc for more information.
`, lib.DOC_URL)
}

func emitpKeywordUsage() {
	fmt.Printf(
		`inserts an out-of-stream variable into the output record stream.
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

Please see %s://johnkerl.org/miller/doc for more information.
`, lib.DOC_URL)
}

func endKeywordUsage() {
	fmt.Println(
		`defines a block of statements to be executed after input records
are ingested. The body statements must be wrapped in curly braces.

  Example: 'end { emit @count }'
  Example: 'end { eprint "Final count is " . @count }'`)
}

func eprintKeywordUsage() {
	fmt.Println(
		`prints expression immediately to stderr.

  Example: mlr --from f.dat put -q 'eprint "The sum of x and y is ".($x+$y)'
  Example: mlr --from f.dat put -q 'for (k, v in $*) { eprint k . " => " . v }'
  Example: mlr --from f.dat put  '(NR % 1000 == 0) { eprint "Checkpoint ".NR}'`)
}

func eprintnKeywordUsage() {
	fmt.Println(
		`prints expression immediately to stderr, without trailing newline.

  Example: mlr --from f.dat put -q 'eprintn "The sum of x and y is ".($x+$y); eprint ""'`)
}

func falseKeywordUsage() {
	fmt.Println(`the boolean literal value.`)
}

func filterKeywordUsage() {
	fmt.Println(
		`includes/excludes the record in the output record stream.

  Example: mlr --from f.dat put 'filter (NR == 2 || $x > 5.4)'

Instead of put with 'filter false' you can simply use put -q.  The following
uses the input record to accumulate data but only prints the running sum
without printing the input record:

  Example: mlr --from f.dat put -q '@running_sum += $x * $y; emit @running_sum'`)
}

func floatKeywordUsage() {
	fmt.Println(
		`declares a floating-point local variable in the current curly-braced scope.
Type-checking happens at assignment: 'float x = 0' is an error.`)
}

func forKeywordUsage() {
	fmt.Println(
		`defines a for-loop using one of three styles. The body statements must
be wrapped in curly braces.
For-loop over stream record:

  Example:  'for (k, v in $*) { ... }'

For-loop over out-of-stream variables:

  Example: 'for (k, v in @counts) { ... }'
  Example: 'for ((k1, k2), v in @counts) { ... }'
  Example: 'for ((k1, k2, k3), v in @*) { ... }'

C-style for-loop:

  Example:  'for (var i = 0, var b = 1; i < 10; i += 1, b *= 2) { ... }'`)
}

func funcKeywordUsage() {
	fmt.Println(
		`used for defining a user-defined function.

  Example: 'func f(a,b) { return sqrt(a**2+b**2)} $d = f($x, $y)'`)
}

func ifKeywordUsage() {
	fmt.Println(
		`starts an if/elif/elif chain. The body statements must be wrapped
in curly braces.`)
}

func inKeywordUsage() {
	fmt.Println(`used in for-loops over stream records or out-of-stream variables.`)
}

func intKeywordUsage() {
	fmt.Println(
		`declares an integer local variable in the current curly-braced scope.
Type-checking happens at assignment: 'int x = 0.0' is an error.`)
}

func mapKeywordUsage() {
	fmt.Println(
		`declares an map-valued local variable in the current curly-braced scope.
Type-checking happens at assignment: 'map b = 0' is an error. map b = {} is
always OK. map b = a is OK or not depending on whether a is a map.`)
}

func numKeywordUsage() {
	fmt.Println(
		`declares an int/float local variable in the current curly-braced scope.
Type-checking happens at assignment: 'num b = true' is an error.`)
}

func printKeywordUsage() {
	fmt.Println(
		`prints expression immediately to stdout.

  Example: mlr --from f.dat put -q 'print "The sum of x and y is ".($x+$y)'
  Example: mlr --from f.dat put -q 'for (k, v in $*) { print k . " => " . v }'
  Example: mlr --from f.dat put  '(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'`)
}

func printnKeywordUsage() {
	fmt.Println(
		`prints expression immediately to stdout, without trailing newline.

  Example: mlr --from f.dat put -q 'printn "."; end { print "" }'`)
}

func returnKeywordUsage() {
	fmt.Println(
		`specifies the return value from a user-defined function.
Omitted return statements (including via if-branches) result in an absent-null
return value, which in turns results in a skipped assignment to an LHS.`)
}

func stderrKeywordUsage() {
	fmt.Println(
		`Used for tee, emit, emitf, emitp, print, and dump in place of filename
to print to standard error.`)
}

func stdoutKeywordUsage() {
	fmt.Println(
		`Used for tee, emit, emitf, emitp, print, and dump in place of filename
to print to standard output.`)
}

func strKeywordUsage() {
	fmt.Println(
		`declares a string local variable in the current curly-braced scope.
Type-checking happens at assignment.`)
}

func subrKeywordUsage() {
	fmt.Println(
		`used for defining a subroutine.

  Example: 'subr s(k,v) { print k . " is " . v} call s("a", $a)'`)
}

func teeKeywordUsage() {
	fmt.Println(
		`prints the current record to specified file.
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
  Example: mlr --from f.dat put -q 'tee | "tr \[a-z\\] \[A-Z\\]", $*'
  Example: mlr --from f.dat put -q 'tee | "tr \[a-z\\] \[A-Z\\] > /tmp/data-".$a, $*'
  Example: mlr --from f.dat put -q 'tee | "gzip > /tmp/data-".$a.".gz", $*'
  Example: mlr --from f.dat put -q --ojson 'tee | "gzip > /tmp/data-".$a.".gz", $*'`)
}

func trueKeywordUsage() {
	fmt.Println(`the boolean literal value.`)
}

func unsetKeywordUsage() {
	fmt.Println(
		`clears field(s) from the current record, or an out-of-stream or local variable.

  Example: mlr --from f.dat put 'unset $x'
  Example: mlr --from f.dat put 'unset $*'
  Example: mlr --from f.dat put 'for (k, v in $*) { if (k =~ "a.*") { unset $[k] } }'
  Example: mlr --from f.dat put '...; unset @sums'
  Example: mlr --from f.dat put '...; unset @sums["green"]'
  Example: mlr --from f.dat put '...; unset @*'`)
}

func varKeywordUsage() {
	fmt.Println(
		`declares an untyped local variable in the current curly-braced scope.

  Examples: 'var a=1', 'var xyz=""'`)
}

func whileKeywordUsage() {
	fmt.Println(
		`introduces a while loop, or with "do", introduces a do-while loop.
The body statements must be wrapped in curly braces.`)
}

func ENVKeywordUsage() {
	fmt.Println(`access to environment variables by name, e.g. '$home = ENV["HOME"]'`)
}

func FILENAMEKeywordUsage() {
	fmt.Println(`evaluates to the name of the current file being processed.`)
}

func FILENUMKeywordUsage() {
	fmt.Println(
		`evaluates to the number of the current file being processed,
starting with 1.`)
}

func FNRKeywordUsage() {
	fmt.Println(
		`evaluates to the number of the current record within the current file
being processed, starting with 1. Resets at the start of each file.`)
}

func IFSKeywordUsage() {
	fmt.Println(`evaluates to the input field separator from the command line.`)
}

func IPSKeywordUsage() {
	fmt.Println(`evaluates to the input pair separator from the command line.`)
}

func IRSKeywordUsage() {
	fmt.Println(
		`evaluates to the input record separator from the command line,
or to LF or CRLF from the input data if in autodetect mode (which is
the default).`)
}

func M_EKeywordUsage() {
	fmt.Println(`the mathematical constant e.`)
}

func M_PIKeywordUsage() {
	fmt.Println(`the mathematical constant pi.`)
}

func NFKeywordUsage() {
	fmt.Println(`evaluates to the number of fields in the current record.`)
}

func NRKeywordUsage() {
	fmt.Println(
		`evaluates to the number of the current record over all files
being processed, starting with 1. Does not reset at the start of each file.`)
}

func OFSKeywordUsage() {
	fmt.Println(`evaluates to the output field separator from the command line.`)
}

func OPSKeywordUsage() {
	fmt.Println(`evaluates to the output pair separator from the command line.`)
}

func ORSKeywordUsage() {
	fmt.Println(
		`evaluates to the output record separator from the command line,
or to LF or CRLF from the input data if in autodetect mode (which is
the default).`)
}
