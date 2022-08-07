<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flags</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verbs</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Functions</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="../release-docs/index.html">Release docs</a>
</span>
</div>
# Online help

Miller has several online help mechanisms built in.

## Main help

The front door is `mlr --help` or its synonym `mlr -h`. This leads you to `mlr help topics` with its list of specific areas:

<pre class="pre-highlight-in-pair">
<b>mlr --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr [flags] {verb} [verb-dependent options ...] {zero or more file names}

If zero file names are provided, standard input is read, e.g.
  mlr --csv sort -f shape example.csv

Output of one verb may be chained as input to another using "then", e.g.
  mlr --csv stats1 -a min,mean,max -f quantity then sort -f color example.csv

Please see 'mlr help topics' for more information.
Please also see https://miller.readthedocs.io
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help topics</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Type 'mlr help {topic}' for any of the following:
Essentials:
  mlr help topics
  mlr help basic-examples
  mlr help file-formats
Flags:
  mlr help flags
  mlr help list-separator-aliases
  mlr help list-separator-regex-aliases
  mlr help comments-in-data-flags
  mlr help compressed-data-flags
  mlr help csv/tsv-only-flags
  mlr help file-format-flags
  mlr help flatten-unflatten-flags
  mlr help format-conversion-keystroke-saver-flags
  mlr help legacy-flags
  mlr help miscellaneous-flags
  mlr help output-colorization-flags
  mlr help pprint-only-flags
  mlr help profiling-flags
  mlr help separator-flags
Verbs:
  mlr help list-verbs
  mlr help usage-verbs
  mlr help verb
Functions:
  mlr help list-functions
  mlr help list-function-classes
  mlr help list-functions-in-class
  mlr help usage-functions
  mlr help usage-functions-by-class
  mlr help function
Keywords:
  mlr help list-keywords
  mlr help usage-keywords
  mlr help keyword
Other:
  mlr help auxents
  mlr help mlrrc
  mlr help output-colorization
  mlr help type-arithmetic-info
Shorthands:
  mlr -g = mlr help flags
  mlr -l = mlr help list-verbs
  mlr -L = mlr help usage-verbs
  mlr -f = mlr help list-functions
  mlr -F = mlr help usage-functions
  mlr -k = mlr help list-keywords
  mlr -K = mlr help usage-keywords
Lastly, 'mlr help ...' will search for your exact text '...' using the sources of
'mlr help flag', 'mlr help verb', 'mlr help function', and 'mlr help keyword'.
Use 'mlr help find ...' for approximate (substring) matches, e.g. 'mlr help find map'
for all things with "map" in their names.
</pre>

If you know the name of the thing you're looking for, use `mlr help`:

<pre class="pre-highlight-in-pair">
<b>mlr help map</b>
</pre>
<pre class="pre-non-highlight-in-pair">
map: declares an map-valued local variable in the current curly-braced scope.
Type-checking happens at assignment: 'map b = 0' is an error. map b = {} is
always OK. map b = a is OK or not depending on whether a is a map.
</pre>

To search by substring, use `mlr help find`:

<pre class="pre-highlight-in-pair">
<b>mlr help find gmt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
sec2gmtdate
Usage: ../c/mlr sec2gmtdate {comma-separated list of field names}
Replaces a numeric field representing seconds since the epoch with the
corresponding GMT year-month-day timestamp; leaves non-numbers as-is.
This is nothing more than a keystroke-saver for the sec2gmtdate function:
  ../c/mlr sec2gmtdate time1,time2
is the same as
  ../c/mlr put '$time1=sec2gmtdate($time1);$time2=sec2gmtdate($time2)'
sec2gmt
Usage: mlr sec2gmt [options] {comma-separated list of field names}
Replaces a numeric field representing seconds since the epoch with the
corresponding GMT timestamp; leaves non-numbers as-is. This is nothing
more than a keystroke-saver for the sec2gmt function:
  mlr sec2gmt time1,time2
is the same as
  mlr put '$time1 = sec2gmt($time1); $time2 = sec2gmt($time2)'
Options:
-1 through -9: format the seconds using 1..9 decimal places, respectively.
--millis Input numbers are treated as milliseconds since the epoch.
--micros Input numbers are treated as microseconds since the epoch.
--nanos  Input numbers are treated as nanoseconds since the epoch.
-h|--help Show this message.
gmt2localtime  (class=time #args=1,2) Convert from a GMT-time string to a local-time string. Consulting $TZ unless second argument is supplied.
Examples:
gmt2localtime("1999-12-31T22:00:00Z") = "2000-01-01 00:00:00" with TZ="Asia/Istanbul"
gmt2localtime("1999-12-31T22:00:00Z", "Asia/Istanbul") = "2000-01-01 00:00:00"
gmt2sec  (class=time #args=1) Parses GMT timestamp as integer seconds since the epoch.
Example:
gmt2sec("2001-02-03T04:05:06Z") = 981173106
localtime2gmt  (class=time #args=1,2) Convert from a local-time string to a GMT-time string. Consults $TZ unless second argument is supplied.
Examples:
localtime2gmt("2000-01-01 00:00:00") = "1999-12-31T22:00:00Z" with TZ="Asia/Istanbul"
localtime2gmt("2000-01-01 00:00:00", "Asia/Istanbul") = "1999-12-31T22:00:00Z"
sec2gmt  (class=time #args=1,2) Formats seconds since epoch as GMT timestamp. Leaves non-numbers as-is. With second integer argument n, includes n decimal places for the seconds part.
Examples:
sec2gmt(1234567890)           = "2009-02-13T23:31:30Z"
sec2gmt(1234567890.123456)    = "2009-02-13T23:31:30Z"
sec2gmt(1234567890.123456, 6) = "2009-02-13T23:31:30.123456Z"
sec2gmtdate  (class=time #args=1) Formats seconds since epoch (integer part) as GMT timestamp with year-month-date. Leaves non-numbers as-is.
Example:
sec2gmtdate(1440768801.7) = "2015-08-28".
</pre>

Etc.

## Command-line flags

This is a command-line version of the [List of command-line flags](reference-main-flag-list.md) page.
See `mlr help flags` for a full listing.

## Per-verb help

This is a command-line version of the [List of verbs](reference-verbs.md) page.
Given the name of a verb (from `mlr -l`) you can invoke it with `--help` or `-h` -- or, use `mlr help verb`:

<pre class="pre-highlight-in-pair">
<b>mlr cat --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr cat [options]
Passes input records directly to output. Most useful for format conversion.
Options:
-n         Prepend field "n" to each record with record-counter starting at 1.
-N {name}  Prepend field {name} to each record with record-counter starting at 1.
-g {a,b,c} Optional group-by-field names for counters, e.g. a,b,c
-h|--help Show this message.
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr group-like -h</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr group-like [options]
Outputs records in batches having identical field names.
Options:
-h|--help Show this message.
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help verb sort</b>
</pre>
<pre class="pre-non-highlight-in-pair">
sort
Usage: mlr sort {flags}
Sorts records primarily by the first specified field, secondarily by the second
field, and so on.  (Any records not having all specified sort keys will appear
at the end of the output, in the order they were encountered, regardless of the
specified sort order.) The sort is stable: records that compare equal will sort
in the order they were encountered in the input record stream.

Options:
-f  {comma-separated field names}  Lexical ascending
-r  {comma-separated field names}  Lexical descending
-c  {comma-separated field names}  Case-folded lexical ascending
-cr {comma-separated field names}  Case-folded lexical descending
-n  {comma-separated field names}  Numerical ascending; nulls sort last
-nf {comma-separated field names}  Same as -n
-nr {comma-separated field names}  Numerical descending; nulls sort first
-t  {comma-separated field names}  Natural ascending
-tr|-rt {comma-separated field names}  Natural descending
-h|--help Show this message.

Example:
  mlr sort -f a,b -nr x,y,z
which is the same as:
  mlr sort -f a -f b -nr x -nr y -nr z
</pre>

Etc.

## Per-function help

This is a command-line version of the [DSL built-in functions](reference-dsl-builtin-functions.md) page.
Given the name of a DSL function (from `mlr -f`) you can use `mlr help function` for details:

<pre class="pre-highlight-in-pair">
<b>mlr help function append</b>
</pre>
<pre class="pre-non-highlight-in-pair">
append  (class=collections #args=2) Appends second argument to end of first argument, which must be an array.
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help function split</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Function "split" not found.
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help function splita</b>
</pre>
<pre class="pre-non-highlight-in-pair">
splita  (class=conversion #args=2) Splits string into array with type inference. First argument is string to split; second is the separator to split on.
Example:
splita("3,4,5", ",") = [3,4,5]
</pre>

Etc.

## REPL help

You can use `:h` or `:help` inside the [REPL](repl.md):

<!--- TODO: repl-executor genmd function -->
<pre class="pre-highlight-in-pair">
<b>$ mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Miller v6.0.0-dev REPL for darwin:amd64:go1.16.5
Docs: https://miller.readthedocs.io
Type ':h' or ':help' for on-line help; ':q' or ':quit' to quit.
[mlr] :h
Options:
:help intro
:help examples
:help repl-list
:help repl-details
:help prompt
:help function-names
:help function-details
:help {function name}, e.g. :help sec2gmt
:help {function name}, e.g. :help sec2gmt
[mlr]
</pre>

## Manual page

If you've gotten Miller from a package installer, you should have `man mlr` producing a traditional manual page.
If not, no worries -- the manual page is a concatenated listing of the same information also available by each of the topics in `mlr help topics`. See also the [Manual page](manpage.md) which is an online copy.
