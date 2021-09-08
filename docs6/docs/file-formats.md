<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flag list</a>
&nbsp;
<a class="quicklink" href="../reference-verbs/index.html">Verb list</a>
&nbsp;
<a class="quicklink" href="../reference-dsl-builtin-functions/index.html">Function list</a>
&nbsp;
<a class="quicklink" href="../glossary/index.html">Glossary</a>
&nbsp;
<a class="quicklink" href="https://github.com/johnkerl/miller" target="_blank">Repository ↗</a>
</span>
</div>
# File formats

Miller handles name-indexed data using several formats: some you probably know by name, such as CSV, TSV, and JSON -- and other formats you're likely already seeing and using in your structured data.

Additionally, Miller gives you the option of including comments within your data.

## Examples

<pre class="pre-highlight-in-pair">
<b>mlr help data-formats</b>
</pre>
<pre class="pre-non-highlight-in-pair">
CSV/CSV-lite: comma-separated values with separate header line
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
</pre>

## CSV/TSV/ASV/USV/etc.

When `mlr` is invoked with the `--csv` or `--csvlite` option, key names are found on the first record and values are taken from subsequent records.  This includes the case of CSV-formatted files.  See [Record Heterogeneity](record-heterogeneity.md) for how Miller handles changes of field names within a single data stream.

Miller has record separator `RS` and field separator `FS`, just as `awk` does.  For TSV, use `--fs tab`; to convert TSV to CSV, use `--ifs tab --ofs comma`, etc.  (See also [Record/field/pair separators](reference-main-io-options.md#recordfieldpair-separators).)

**TSV (tab-separated values):** the following are synonymous pairs:

* `--tsv` and `--csv --fs tab`
* `--itsv` and `--icsv --ifs tab`
* `--otsv` and `--ocsv --ofs tab`
* `--tsvlite` and `--csvlite --fs tab`
* `--itsvlite` and `--icsvlite --ifs tab`
* `--otsvlite` and `--ocsvlite --ofs tab`

**ASV (ASCII-separated values):** the flags `--asv`, `--iasv`, `--oasv`, `--asvlite`, `--iasvlite`, and `--oasvlite` are analogous except they use ASCII FS and RS 0x1f and 0x1e, respectively.

**USV (Unicode-separated values):** likewise, the flags `--usv`, `--iusv`, `--ousv`, `--usvlite`, `--iusvlite`, and `--ousvlite` use Unicode FS and RS U+241F (UTF-8 0x0xe2909f) and U+241E (UTF-8 0xe2909e), respectively.

Miller's `--csv` flag supports [RFC-4180 CSV](https://tools.ietf.org/html/rfc4180). This includes CRLF line-terminators by default, regardless of platform.

Here are the differences between CSV and CSV-lite:

* CSV supports [RFC-4180](https://tools.ietf.org/html/rfc4180)-style double-quoting, including the ability to have commas and/or LF/CRLF line-endings contained within an input field; CSV-lite does not.

* CSV does not allow heterogeneous data; CSV-lite does (see also [Record Heterogeneity](record-heterogeneity.md)).

* The CSV-lite input-reading code is fractionally more efficient than the CSV input-reader.

Here are things they have in common:

* The ability to specify record/field separators other than the default, e.g. CR-LF vs. LF, or tab instead of comma for TSV, and so on.

* The `--implicit-csv-header` flag for input and the `--headerless-csv-output` flag for output.

## JSON

JSON is a format which supports scalars (numbers, strings, boolean, etc.) as
well as "objects" (maps) and "arrays" (lists), while Miller is a tool for
handling **tabular data** only.  By *tabular JSON* I mean the data is either a
sequence of one or more objects, or an array consisting of one or more objects.
Miller treats JSON objects as name-indexed records.

This means Miller cannot (and should not) handle arbitrary JSON.  In practice,
though, Miller can handle single JSON objects as well as list of them. The only
kinds of JSON that are unmillerable are single scalars (e.g. file contents `3`)
and arrays of non-object (e.g. file contents `[1,2,3,4,5]`).  Check out
[jq](https://stedolan.github.io/jq/) for a tool which handles all valid JSON.

In short, if you have tabular data represented in JSON -- lists of objects,
either with or without outermost `[...]` -- [then Miller can handle that for
you.

### Single-level JSON objects

An **array of single-level objects** is, quite simply, **a table**:

<pre class="pre-highlight-in-pair">
<b>mlr --json head -n 2 then cut -f color,shape data/json-example-1.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "color": "yellow",
  "shape": "triangle"
}
{
  "color": "red",
  "shape": "square"
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --json head -n 2 then cut -f color,u,v data/json-example-1.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "color": "yellow",
  "u": 0.632170,
  "v": 0.988721
}
{
  "color": "red",
  "u": 0.219668,
  "v": 0.001257
}
</pre>

Single-level JSON data goes back and forth between JSON and tabular formats
in the direct way:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint head -n 2 then cut -f color,u,v data/json-example-1.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  u        v
yellow 0.632170 0.988721
red    0.219668 0.001257
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint cat data/json-example-1.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag i  u        v        w        x
yellow triangle 1    11 0.632170 0.988721 0.436498 5.798188
red    square   1    15 0.219668 0.001257 0.792778 2.944117
red    circle   1    16 0.209017 0.290052 0.138103 5.065034
red    square   0    48 0.956274 0.746720 0.775542 7.117831
purple triangle 0    51 0.435535 0.859129 0.812290 5.753095
red    square   0    64 0.201551 0.953110 0.771991 5.612050
purple triangle 0    65 0.684281 0.582372 0.801405 5.805148
yellow circle   1    73 0.603365 0.423708 0.639785 7.006414
yellow circle   1    87 0.285656 0.833516 0.635058 6.350036
purple square   0    91 0.259926 0.824322 0.723735 6.854221
</pre>

### Nested JSON objects

Additionally, Miller can **tabularize nested objects by concatentating keys**. If your processing has
input as well as output in JSON format, JSON structure is preserved throughout the processing:

<pre class="pre-highlight-in-pair">
<b>mlr --json --jvstack head -n 2 data/json-example-2.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "flag": 1,
  "i": 11,
  "attributes": {
    "color": "yellow",
    "shape": "triangle"
  },
  "values": {
    "u": 0.632170,
    "v": 0.988721,
    "w": 0.436498,
    "x": 5.798188
  }
}
{
  "flag": 1,
  "i": 15,
  "attributes": {
    "color": "red",
    "shape": "square"
  },
  "values": {
    "u": 0.219668,
    "v": 0.001257,
    "w": 0.792778,
    "x": 2.944117
  }
}
</pre>

But if the input format is JSON and the output format is not (or vice versa) then key-concatenation applies:

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint head -n 4 data/json-example-2.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
flag i  attributes.color attributes.shape values.u values.v values.w values.x
1    11 yellow           triangle         0.632170 0.988721 0.436498 5.798188
1    15 red              square           0.219668 0.001257 0.792778 2.944117
1    16 red              circle           0.209017 0.290052 0.138103 5.065034
0    48 red              square           0.956274 0.746720 0.775542 7.117831
</pre>

This is discussed in more detail on the page [Flatten/unflatten: JSON vs. tabular formats](flatten-unflatten.md).

### JSON-formatting options

JSON isn't a parameterized format, so `RS`, `FS`, `PS` aren't specifiable. Nonetheless, you can do the following:

* Use `--no-jvstack` to print JSON objects one per line.  By default, each Miller record (JSON object) is pretty-printed in multi-line format.

* Use `--jlistwrap` to print the sequence of JSON objects wrapped in an outermost `[` and `]`. By default, these aren't printed.

<!---
TODO: probably remove entirely
* Use `--jquoteall` to double-quote all object values. By default, integers, floating-point numbers, and booleans `true` and `false` are not double-quoted when they appear as JSON-object keys.
-->

* Use `--jflatsep yourseparatorhere` to specify the string used for key concatenation: this defaults to a single dot.

### JSON-in-CSV

It's quite common to have CSV data which contains stringified JSON as a column.
See the [JSON parse and stringify
section](reference-main-data-types.md#json-parse-and-stringify) for ways to
decode these in Miller.

## PPRINT: Pretty-printed tabular

Miller's pretty-print format is like CSV, but column-aligned.  For example, compare

<pre class="pre-highlight-in-pair">
<b>mlr --ocsv cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,i,x,y
pan,pan,1,0.346791,0.726802
eks,pan,2,0.758679,0.522151
wye,wye,3,0.204603,0.338318
eks,wye,4,0.381399,0.134188
wye,pan,5,0.573288,0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y
pan pan 1 0.346791 0.726802
eks pan 2 0.758679 0.522151
wye wye 3 0.204603 0.338318
eks wye 4 0.381399 0.134188
wye pan 5 0.573288 0.863624
</pre>

Note that while Miller is a line-at-a-time processor and retains input lines in memory only where necessary (e.g. for sort), pretty-print output requires it to accumulate all input lines (so that it can compute maximum column widths) before producing any output. This has two consequences: (a) pretty-print output won't work on `tail -f` contexts, where Miller will be waiting for an end-of-file marker which never arrives; (b) pretty-print output for large files is constrained by available machine memory.

See [Record Heterogeneity](record-heterogeneity.md) for how Miller handles changes of field names within a single data stream.

For output only (this isn't supported in the input-scanner as of 5.0.0) you can use `--barred` with pprint output format:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --barred cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
+-----+-----+---+----------+----------+
| a   | b   | i | x        | y        |
+-----+-----+---+----------+----------+
| pan | pan | 1 | 0.346791 | 0.726802 |
| eks | pan | 2 | 0.758679 | 0.522151 |
| wye | wye | 3 | 0.204603 | 0.338318 |
| eks | wye | 4 | 0.381399 | 0.134188 |
| wye | pan | 5 | 0.573288 | 0.863624 |
+-----+-----+---+----------+----------+
</pre>

## Markdown tabular

Markdown format looks like this:

<pre class="pre-highlight-in-pair">
<b>mlr --omd cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
| a | b | i | x | y |
| --- | --- | --- | --- | --- |
| pan | pan | 1 | 0.346791 | 0.726802 |
| eks | pan | 2 | 0.758679 | 0.522151 |
| wye | wye | 3 | 0.204603 | 0.338318 |
| eks | wye | 4 | 0.381399 | 0.134188 |
| wye | pan | 5 | 0.573288 | 0.863624 |
</pre>

which renders like this when dropped into various web tools (e.g. github comments):

![pix/omd.png](pix/omd.png)

As of Miller 4.3.0, markdown format is supported only for output, not input.

## XTAB: Vertical tabular

This is perhaps most useful for looking a very wide and/or multi-column data which causes line-wraps on the screen (but see also
[ngrid](https://github.com/twosigma/ngrid/) for an entirely different, very powerful option). Namely:

<pre class="pre-highlight-in-pair">
<b>$ grep -v '^#' /etc/passwd | head -n 6 | mlr --nidx --fs : --opprint cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1          2 3  4  5                          6               7
nobody     * -2 -2 Unprivileged User          /var/empty      /usr/bin/false
root       * 0  0  System Administrator       /var/root       /bin/sh
daemon     * 1  1  System Services            /var/root       /usr/bin/false
_uucp      * 4  4  Unix to Unix Copy Protocol /var/spool/uucp /usr/sbin/uucico
_taskgated * 13 13 Task Gate Daemon           /var/empty      /usr/bin/false
_networkd  * 24 24 Network Services           /var/networkd   /usr/bin/false
</pre>

<pre class="pre-highlight-in-pair">
<b>$ grep -v '^#' /etc/passwd | head -n 2 | mlr --nidx --fs : --oxtab cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1 nobody
2 *
3 -2
4 -2
5 Unprivileged User
6 /var/empty
7 /usr/bin/false

1 root
2 *
3 0
4 0
5 System Administrator
6 /var/root
7 /bin/sh
</pre>

<pre class="pre-highlight-in-pair">
<b>$ grep -v '^#' /etc/passwd | head -n 2 | \</b>
<b>  mlr --nidx --fs : --ojson --jvstack --jlistwrap \</b>
<b>    label name,password,uid,gid,gecos,home_dir,shell</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "name": "nobody",
  "password": "*",
  "uid": -2,
  "gid": -2,
  "gecos": "Unprivileged User",
  "home_dir": "/var/empty",
  "shell": "/usr/bin/false"
},
{
  "name": "root",
  "password": "*",
  "uid": 0,
  "gid": 0,
  "gecos": "System Administrator",
  "home_dir": "/var/root",
  "shell": "/bin/sh"
}
]
</pre>

## DKVP: Key-value pairs

Miller's default file format is DKVP, for **delimited key-value pairs**. Example:

<pre class="pre-highlight-in-pair">
<b>mlr cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

Such data are easy to generate, e.g. in Ruby with

<pre class="pre-non-highlight-non-pair">
puts "host=#{hostname},seconds=#{t2-t1},message=#{msg}"
</pre>

<pre class="pre-non-highlight-non-pair">
puts mymap.collect{|k,v| "#{k}=#{v}"}.join(',')
</pre>

or `print` statements in various languages, e.g.

<pre class="pre-non-highlight-non-pair">
echo "type=3,user=$USER,date=$date\n";
</pre>

<pre class="pre-non-highlight-non-pair">
logger.log("type=3,user=$USER,date=$date\n");
</pre>

Fields lacking an IPS will have positional index (starting at 1) used as the key, as in NIDX format. For example, `dish=7,egg=8,flint` is parsed as `"dish" => "7", "egg" => "8", "3" => "flint"` and `dish,egg,flint` is parsed as `"1" => "dish", "2" => "egg", "3" => "flint"`.

As discussed in [Record Heterogeneity](record-heterogeneity.md), Miller handles changes of field names within the same data stream. But using DKVP format this is particularly natural. One of my favorite use-cases for Miller is in application/server logs, where I log all sorts of lines such as

<pre class="pre-non-highlight-non-pair">
resource=/path/to/file,loadsec=0.45,ok=true
record_count=100, resource=/path/to/file
resource=/some/other/path,loadsec=0.97,ok=false
</pre>

etc. and I just log them as needed. Then later, I can use `grep`, `mlr --opprint group-like`, etc.
to analyze my logs.

See the [I/O options reference](reference-main-io-options.md) regarding how to specify separators other than the default equals-sign and comma.

## NIDX: Index-numbered (toolkit style)

With `--inidx --ifs ' ' --repifs`, Miller splits lines on whitespace and assigns integer field names starting with 1.

This recapitulates Unix-toolkit behavior.

Example with index-numbered output:

<pre class="pre-highlight-in-pair">
<b>cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --onidx --ofs ' ' cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
pan pan 1 0.346791 0.726802
eks pan 2 0.758679 0.522151
wye wye 3 0.204603 0.338318
eks wye 4 0.381399 0.134188
wye pan 5 0.573288 0.863624
</pre>

Example with index-numbered input:

<pre class="pre-highlight-in-pair">
<b>cat data/mydata.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
oh say can you see
by the dawn's
early light
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --inidx --ifs ' ' --odkvp cat data/mydata.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
1=oh,2=say,3=can,4=you,5=see
1=by,2=the,3=dawn's
1=early,2=light
</pre>

Example with index-numbered input and output:

<pre class="pre-highlight-in-pair">
<b>cat data/mydata.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
oh say can you see
by the dawn's
early light
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --nidx --fs ' ' --repifs cut -f 2,3 data/mydata.txt</b>
</pre>
<pre class="pre-non-highlight-in-pair">
say can
the dawn's
light
</pre>

## Data-conversion keystroke-savers

While you can do format conversion using `mlr --icsv --ojson cat myfile.csv`, there are also keystroke-savers for this purpose, such as `mlr --c2j cat myfile.csv`.  For a complete list:

<pre class="pre-highlight-in-pair">
<b>mlr help format-conversion</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Type 'mlr help {topic}' for any of the following:
Essentials:
  mlr help topics
  mlr help basic-examples
  mlr help data-formats
Flags:
  mlr help flags
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
</pre>

<!---
TODO: probably entirely unsupport this feature in Miller6.

## Autodetect of line endings

Default line endings (`--irs` and `--ors`) are `'auto'` which means **autodetect from the input file format**, as long as the input file(s) have lines ending in either LF (also known as linefeed, `'\n'`, `0x0a`, Unix-style) or CRLF (also known as carriage-return/linefeed pairs, `'\r\n'`, `0x0d 0x0a`, Windows style).

**If both IRS and ORS are auto (which is the default) then LF input will lead to LF output and CRLF input will lead to CRLF output, regardless of the platform you're running on.**

The line-ending autodetector triggers on the first line ending detected in the input stream. E.g. if you specify a CRLF-terminated file on the command line followed by an LF-terminated file then autodetected line endings will be CRLF.

If you use `--ors {something else}` with (default or explicitly specified) `--irs auto` then line endings are autodetected on input and set to what you specify on output.

If you use `--irs {something else}` with (default or explicitly specified) `--ors auto` then the output line endings used are LF on Unix/Linux/BSD/MacOS X, and CRLF on Windows.

See also [Record/field/pair separators](reference-main-io-options.md#recordfieldpair-separators) for more information about record/field/pair separators.
--->

## Comments in data

You can include comments within your data files, and either have them ignored, or passed directly through to the standard output as soon as they are encountered:

<pre class="pre-highlight-in-pair">
<b>mlr help comments-in-data</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Type 'mlr help {topic}' for any of the following:
Essentials:
  mlr help topics
  mlr help basic-examples
  mlr help data-formats
Flags:
  mlr help flags
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
</pre>

Examples:

<pre class="pre-highlight-in-pair">
<b>cat data/budget.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
# Asana -- here are the budget figures you asked for!
type,quantity
purple,456.78
green,678.12
orange,123.45
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --skip-comments --icsv --opprint sort -nr quantity data/budget.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
type   quantity
green  678.12
purple 456.78
orange 123.45
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --pass-comments --icsv --opprint sort -nr quantity data/budget.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
# Asana -- here are the budget figures you asked for!
type   quantity
green  678.12
purple 456.78
orange 123.45
</pre>
