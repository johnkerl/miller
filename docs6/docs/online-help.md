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
# Online help

Miller has several online help mechanisms built in.

## Main help

The front door is `mlr --help` or its synonym `mlr -h`. This leads you to `mlr help topics` with its list of specific areas:

<pre class="pre-highlight-in-pair">
<b>mlr --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr [flags] {verb} [verb-dependent options ...] {zero or more file names}
Output of one verb may be chained as input to another using "then", e.g.
  mlr stats1 -a min,mean,max -f flag,u,v -g color then sort -f color
Please see 'mlr help topics' for more information.
Please also see https://johnkerl.org/miller6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help topics</b>
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

<pre class="pre-highlight-in-pair">
<b>mlr help functions</b>
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
No exact match for "split". Inexact matches:
  splita
  splitax
  splitkv
  splitkvx
  splitnv
  splitnvx
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr help function splita</b>
</pre>
<pre class="pre-non-highlight-in-pair">
splita  (class=conversion #args=2) Splits string into array with type inference. Example:
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
Pre-release docs for Miller 6: https://johnkerl.org/miller6
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
