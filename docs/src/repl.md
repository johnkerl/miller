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
# The REPL

The Miller REPL (read-evaluate-print loop) is an interactive counterpart to record-processing using the `put`/`filter` language. (A REPL is anything that evaluates what you type into it -- like `python` with no arguments, or Ruby's `irb`, or `node` with no arguments, etc.)

Miller's REPL isn't a source-level debugger which lets you execute one source-code *statement* at a time -- however, it does let you operate on one *record* at a time. Further, it lets you use "immediate expressions", namely, you can interact with the [Miller programming language](miller-programming-language.md) without having to provide data from an input file.

<pre class="pre-highlight-in-pair">
<b>mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">

[mlr] 1 + 2
3
</pre>

## Using Miller without the REPL

Using `put` and `filter`, you can do the following as we've seen above:

* Specify input format (e.g. `--icsv`), output format (e.g. `--ojson`), etc. using command-line flags.
* Specify filenames on the command line.
* Define `begin {...}` blocks which are executed before the first record is read.
* Define `end {...}` blocks which are executed after the last record is read.
* Define user-defined functions/subroutines using `func` and `subr`.
* Specify statements to be executed on each record -- which are anything outside of `begin`/`end`/`func`/`subr`.
* Example:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from example.csv head -n 2 \</b>
<b>  then put 'begin {print "HELLO"} $qr = $quantity / $rate; end {print "GOODBYE"}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
HELLO
{
  "color": "yellow",
  "shape": "triangle",
  "flag": "true",
  "k": 1,
  "index": 11,
  "quantity": 43.6498,
  "rate": 9.8870,
  "qr": 4.414868008496004
}
{
  "color": "red",
  "shape": "square",
  "flag": "true",
  "k": 2,
  "index": 15,
  "quantity": 79.2778,
  "rate": 0.0130,
  "qr": 6098.292307692308
}
GOODBYE
</pre>

## Using Miller with the REPL

Using the REPL, by contrast, you get interactive control over those same steps:

* Specify input format (e.g. `--icsv`), output format (e.g. `--ojson`), etc. using command-line flags.
* REPL-only statements (non-DSL statements) start with `:`, such as `:help` or `:quit`
  or `:open`.
* Specify filenames either on the command line or via `:open` at the Miller REPL.
* Read records one at a time using `:read`.
* Write the current record (maybe after you've modified it with things like `$z = $x + $y`)
  using `:write`. This goes to the terminal; you can use `:> {filename}` to make writes
  go to a file, or `:>> {filename}` to append.
* You can type `:reopen` to go back to the start of the same file(s) you specified
  with `:open`.
* Skip ahead using statements `:skip 10` or `:skip until NR == 100` or
  `:skip until $status_code != 200`.
* Similarly, but processing records rather than skipping past them, using
  `:process` rather than `:skip`. Like `:write`, these go to the screen;
  use `:> {filename}` or `:>> {filename}` to log to a file instead.
* Define `begin {...}` blocks; invoke them at will using `:begin`.
* Define `end {...}` blocks; invoke them at will using `:end`.
* Define user-defined functions/subroutines using `func`/`subr`; call them from other statements.
* Interactively specify statements to be executed immediately on the current record.
* Load any of the above from Miller-script files using `:load`.

The input "record" by default is the empty map but you can do things like
`$x=3`, or `unset $y`, or `$* = {"x": 3, "y": 4}` to populate it. Or, `:open
foo.dat` followed by `:read` to populate it from a data file.

Non-assignment expressions, such as `7` or `true`, operate as filter conditions
in the `put` DSL: they can be used to specify whether a record will or won't be
included in the output-record stream.  But here in the REPL, they are simply
printed to the terminal, e.g. if you type `1+2`, you will see `3`.

## Entering multi-line statements

* To enter multi-line statements, enter `<` on a line by itself, then the code (taking care
  for semicolons), then `>` on a line by itself. These will be executed immediately.
* If you enter `<<` on a line by itself, then the code, then `>>` on a line by
  itself, the statements will be remembered for executing on records with
  `:main`, as if you had done `:load` to load statements from a file.

## Examples

Use the REPL to look at arithmetic:

<pre class="pre-highlight-in-pair">
<b>mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">

[mlr] 6/3
2

[mlr] 6/5
1.2

[mlr] typeof(6/3)
int

[mlr] typeof(6/5)
float
</pre>

Read the first record from a small file:

<pre class="pre-highlight-in-pair">
<b>mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">

[mlr] :open foo.dat

[mlr] :read

[mlr] :context
FILENAME="foo.dat",FILENUM=1,NR=1,FNR=1

[mlr] $*
{
  "a": "eks",
  "b": "wye",
  "i": 4,
  "x": 0.38139939387114097,
  "y": 0.13418874328430463
}

[mlr] $z = $x + $i

[mlr] :write
a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,z=4.381399393871141
</pre>

Skip until deep into a larger file, then inspect a record:

<pre class="pre-highlight-in-pair">
<b>mlr repl --csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">

[mlr] :open data/colored-shapes.csv
[mlr] :skip until NR == 10000
[mlr] :r
[mlr] $*
{
  "color": "yellow",
  "shape": "circle",
  "flag": 1,
  "i": 496422,
  "u": 0.6530503199545348,
  "v": 0.23908588907834516,
  "w": 0.4799125551304738,
  "x": 6.379888206335166
}
</pre>

## History-editing

No command-line-history-editing feature is built in but **rlwrap mlr repl** is a
delight. You may need `brew install rlwrap`, `sudo apt-get install rlwrap`,
etc. depending on your platform.

Suggestion: `alias mrpl='rlwrap mlr repl'` in your shell's startup file.

## Online help

After `mlr repl`, type `:help` to see more about your options. In particular, `:help examples`.
