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
# What's new in Miller 6

See also the [list of issues tagged with go-port](https://github.com/johnkerl/miller/issues?q=label%3Ago-port).

## User experience

### Documentation improvements

Documentation (what you're reading here) and online help (`mlr --help`) have been completely reworked.

In the initial release, the focus was convincing users already familiar with
`awk`/`grep`/`cut` that Miller was a viable alternative -- but over time it's
become clear that many Miller users aren't expert with those tools. The focus
has shifted toward a higher quantity of more introductory/accessible material
for command-line data processing.

Similarly, the FAQ/recipe material has been expanded to include more, and
simpler, use-cases including resolved questions from
[Miller Issues](https://github.com/johnkerl/miller/issues)
and
[Miller Discussions](https://github.com/johnkerl/miller/discussions);
more complex/niche material has been pushed farther down. The long reference
pages have been split up into separate pages. (See also
[Structure of these documents](structure-of-these-documents.md).)

Since CSV is overwhelmingly the most popular data format for Miller, it is
now discussed first, and more examples use CSV.

### Improved Windows experience

Stronger support for Windows (with or without MSYS2), with a couple of
exceptions.  See [Miller on Windows](miller-on-windows.md) for more information.

Binaries are reliably available using GitHub Actions: see also [Installation](installing-miller.md).

### Output colorization

Miller uses separate, customizable colors for keys and values whenever the output is to a terminal. See [Output Colorization](output-colorization.md).

### Improved command-line parsing

Miller 6 has getoptish command-line parsing ([pull request 467](https://github.com/johnkerl/miller/pull/467)):

* `-xyz` expands automatically to `-x -y -z`, so (for example) `mlr cut -of shape,flag` is the same as `mlr cut -o -f shape,flag`.
* `--foo=bar` expands automatically to  `--foo bar`, so (for example) `mlr --ifs=comma` is the same as `mlr --ifs comma`.
* `--mfrom`, `--load`, `--mload` as described in the [flags reference](reference-main-flag-list.md#miscellaneous-flags).

A small but nice item: since **mlr --csv** and **mlr --json** are so common, you can now use alternate shorthands **mlr -c** and **mlr -j**, respectively.

### Improved error messages for DSL parsing

For `mlr put` and `mlr filter`, parse-error messages now include location information:

<pre class="pre-non-highlight-non-pair">
mlr: cannot parse DSL expression.
Parse error on token ">" at line 63 columnn 7.
</pre>

### REPL

Miller now has a read-evaluate-print-loop ([REPL](repl.md)) where you can single-step through your data-file record, express arbitrary statements to converse with the data, etc.

<pre class="pre-highlight-in-pair">
<b>mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">

[mlr] 1 + 2
3

[mlr] apply([1,2,3,4,5], func(e) {return e ** 3})
[1, 8, 27, 64, 125]

[mlr] :open example.csv

[mlr] :read

[mlr] $*
{
  "color": "yellow",
  "shape": "triangle",
  "flag": "true",
  "k": 1,
  "index": 11,
  "quantity": 43.6498,
  "rate": 9.8870
}

</pre>

## Localization and internationalization

### Improved internationalization support

You can now write field names, local variables, etc. all in UTF-8, e.g.  `mlr
--c2p filter '$σχήμα == "κύκλος"' παράδειγμα.csv`. See the
[internationalization page](internationalization.md) for examples.

### Improved datetime/timezone support

Including support for specifying timezone via function arguments, as an alternative to
the `TZ` environment variable. Please see [DSL datetime/timezone functions](reference-dsl-time.md).

## Data ingestion

### In-process support for compressed input

In addition to `--prepipe gunzip`, you can now use the `--gzin` flag. In fact, if your files end in `.gz` you don't even need to do that -- Miller will autodetect by file extension and automatically uncompress `mlr --csv cat foo.csv.gz`. Similarly for `.z` and `.bz2` files.  Please see the page on [Compressed data](reference-main-compressed-data.md) for more information.

### Support for reading web URLs

You can read input with prefixes `https://`, `http://`, and `file://`:

<pre class="pre-highlight-in-pair">
<b>mlr --csv sort -f shape \</b>
<b>  https://raw.githubusercontent.com/johnkerl/miller/main/docs/src/gz-example.csv.gz</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
red,circle,true,3,16,13.8103,2.9010
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
red,square,true,2,15,79.2778,0.0130
red,square,false,4,48,77.5542,7.4670
red,square,false,6,64,77.1991,9.5310
purple,square,false,10,91,72.3735,8.2430
yellow,triangle,true,1,11,43.6498,9.8870
purple,triangle,false,5,51,81.2290,8.5910
purple,triangle,false,7,65,80.1405,5.8240
</pre>

## Data processing

### Improved JSON support, and arrays

Arrays are now supported in Miller's `put`/`filter` programming language, as
described in the [Arrays reference](reference-main-arrays.md). (Also, `array` is
now a keyword so this is no longer usable as a local-variable or UDF name.)

JSON support is improved:

* Direct support for arrays means that you can now use Miller to process more JSON files.
* Streamable JSON parsing: Miller's internal record-processing pipeline starts as soon as the first record is read (which was already the case for other file formats). This means that, unless records are wrapped with outermost `[...]`, Miller now handles JSON in `tail -f` contexts like it does for other file formats.
* Flatten/unflatten -- conversion of JSON nested data structures (arrays and/or maps in record values) to/from non-JSON formats is a powerful new feature, discussed in the page [Flatten/unflatten: JSON vs. tabular formats](flatten-unflatten.md).
* Since types are better handled now, the workaround flags `--jvquoteall` and `--jknquoteint` no longer have meaning -- although they're accepted as no-ops at the command line for backward compatibility.
* Multi-line JSON is now the default. Use `--no-jvstack` for Miller-5 style, which required `--jvstack` to get multiline output.

See also the [Arrays reference](reference-main-arrays.md) for more information.

### Improved numeric conversion

The most central part of Miller 6 is a deep refactor of how data values are parsed
from file contents, how types are inferred, and how they're converted back to
text into output files.

This was all initiated by [https://github.com/johnkerl/miller/issues/151](https://github.com/johnkerl/miller/issues/151).

In Miller 5 and below, all values were stored as strings, then only converted
to int/float as-needed, for example when a particular field was referenced in
the `stats1` or `put` verbs. This led to awkwardnesses such as the `-S`
and `-F` flags for `put` and `filter`.

In Miller 6, things parseable as int/float are treated as such from the moment
the input data is read, and these are passed along through the verb chain.  All
values are typed from when they're read, and their types are passed along.
Meanwhile the original string representation of each value is also retained. If
a numeric field isn't modified during the processing chain, it's printed out
the way it arrived. Also, quoted values in JSON strings are flagged as being
strings throughout the processing chain.

For example (see [https://github.com/johnkerl/miller/issues/178](https://github.com/johnkerl/miller/issues/178)) you can now do

<pre class="pre-highlight-in-pair">
<b>echo '{ "a": "0123" }' | mlr --json cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": "0123"
}
</pre>

<pre class="pre-highlight-in-pair">
<b>echo '{ "x": 1.230, "y": 1.230000000 }' | mlr --json cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "x": 1.230,
  "y": 1.230000000
}
</pre>

### Regex support for IFS and IPS

You can now split fields on whitespace when whitespace is a mix of tabs and
spaces.  As well, you can use regular expressions for the input field-separator
and the input pair-separator.  Please see the section on
[multi-character and regular-expression separators](reference-main-separators.md#multi-character-and-regular-expression-separators).

In particular, for NIDX format, the default IFS now allows splitting on one or more of space or tab.

### Case-folded sorting options

The [sort](reference-verbs.md#sort) verb now accepts `-c` and `-cr` options for case-folded ascending/descending sort, respetively.

### New DSL functions / operators

* Higher-order functions [`select`](reference-dsl-builtin-functions.md#select), [`apply`](reference-dsl-builtin-functions.md#apply), [`reduce`](reference-dsl-builtin-functions.md#reduce), [`fold`](reference-dsl-builtin-functions.md#fold), and [`sort`](reference-dsl-builtin-functions.md#sort).  See the [sorting page](sorting.md) and the [higher-order-functions page](reference-dsl-higher-order-functions.md) for more information.

* Absent-coalesce operator [`??`](reference-dsl-builtin-functions.md#absent-coalesce) along with `??=`;
absent-empty-coalesce operator [`???`](reference-dsl-builtin-functions.md#absent-empty-coalesce) along with `???=`.

* The dot operator is not new, but it has a new role: in addition to its existing use for string-concatenation like `"a"."b" = "ab"`, you can now also use it for keying maps. For example, `$req.headers.host` is the same as `$req["headers"]["host"]`. See the [dot-operator reference](reference-dsl-operators.md#the-double-purpose-dot-operator) for more information.

* String-hashing functions [md5](reference-dsl-builtin-functions.md#md5), [sha1](reference-dsl-builtin-functions.md#sha1), [sha256](reference-dsl-builtin-functions.md#sha256), and [sha512](reference-dsl-builtin-functions.md#sha512).

* Platform-property functions [hostname](reference-dsl-builtin-functions.md#hostname), [os](reference-dsl-builtin-functions.md#os), and [version](reference-dsl-builtin-functions.md#version).

* Unsigned right-shift [`>>>`](reference-dsl-builtin-functions.md#ursh) along with `>>>=`.

## Developer-specific aspects

* Miller has been ported from C to Go. Developer notes: [https://github.com/johnkerl/miller/blob/main/README-go-port.md](https://github.com/johnkerl/miller/blob/main/README-go-port.md).
* Regression testing has been completely reworked, including regression-testing now running fully on Windows (alongside Linux and Mac) [on each GitHub commit](https://github.com/johnkerl/miller/actions).

## Changes from Miller 5

The following differences are rather technical. If they don't sound familiar to you, not to worry! Most users won't be affected by the (relatively minor) changes between Miller 5 and Miller 6.

* Line endings: The `--auto` flag is now ignored. Before, if a file had CR/LF (Windows-style) line endings on input (on any platform), it would have the same on output; likewise, LF (Unix-style) line endings. Now, files with CR/LF or LF line endings are processed on any platform, but the output line-ending is for the platform. E.g. reading CR/LF files on Linux will now produce LF output.
* Since IFS and IPS can be regular expressions now, if your `IFS` or `IPS` is a `.` (which is a regular-expression metacharacter) you can either use `--ifs '\.'` or `--ifs . --no-ifs-regex` (and similarly for `IPS`).  Please see the section on [multi-character and regular-expression separators](reference-main-separators.md#multi-character-and-regular-expression-separators).
* JSON formatting:
    * `--jvstack` and `--jsonx` (multi-line JSON output) is now the default for JSON; use `--no-jvstack` to suppress it.
    * `--jknquoteint` and `jquoteall` are ignored; they were workarounds for the (now much-improved) tyoe-inference and type-tracking in Miller 6.
    * `--json-fatal-arrays-on-input`, `--json-map-arrays-on-input`, and `--json-skip-arrays-on-input` are ignored; Miller 6 now supports arrays fully.
    * See also `mlr help legacy-flags` or the [legacy-flags reference](reference-main-flag-list.md#legacy-flags).
* Type-inference:
    * The `-S` and `-F` flags to `mlr put` and `mlr filter` are ignored, since type-inference is no longer done in `mlr put` and `mlr filter`, but rather, when records are first read. You can use `mlr -S` and `mlr -A`, respectively, instead to control type-inference within the record-readers.
    * Similarly, use `mlr -O` to force octal-looking strings to remain strings like `"0123"`, not ints like `0123` which is 83 in decimal.
    * See also the [miscellaneous-flags reference](reference-main-flag-list.md#miscellaneous-flags).
* Emitting a map-valued expression now requires either a temporary variable or the new `emit1` keyword. Please see the
[page on emit statements](reference-dsl-output-statements.md#emit1-and-emitemitpemitf) for more information.
