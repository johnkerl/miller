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
# Auxiliary commands

There are a few nearly-standalone programs which have a little to do with the rest of Miller, do not participate in record streams, and do not deal with file formats. They might as well be little standalone executables, but instead they're delivered within the main Miller executable for convenience.

<pre class="pre-highlight-in-pair">
<b>mlr aux-list</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Available subcommands:
  aux-list
  hex
  lecat
  termcvt
  unhex
  help
  regtest
  repl
For more information, please invoke mlr {subcommand} --help.
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr lecat --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr lecat [options] {zero or more file names}
Simply echoes input, but flags CR characters in red and LF characters in green.
If zero file names are supplied, standard input is read.
Options:
--mono: don't try to colorize the output
-h or --help: print this message
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr termcvt --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr termcvt [option] {zero or more file names}
Option (exactly one is required):
--cr2crlf
--lf2crlf
--crlf2cr
--crlf2lf
--cr2lf
--lf2cr
-I in-place processing (default is to write to stdout)
-h or --help: print this message
Zero file names means read from standard input.
Output is always to standard output; files are not written in-place.
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr hex --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr hex [options] {zero or more file names}
Simple hex-dump.
If zero file names are supplied, standard input is read.
Options:
-r: print only raw hex without leading offset indicators or trailing ASCII dump.
-h or --help: print this message
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr unhex --help</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Usage: mlr unhex [option] {zero or more file names}
Options:
-h or --help: print this message
Zero file names means read from standard input.
Output is always to standard output; files are not written in-place.
</pre>

Examples:

<pre class="pre-highlight-in-pair">
<b>echo 'Hello, world!' | mlr lecat --mono</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Hello, world![LF]
</pre>

<pre class="pre-highlight-in-pair">
<b>echo 'Hello, world!' | mlr termcvt --lf2crlf | mlr lecat --mono</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Hello, world![CR][LF]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr hex data/budget.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
00000000: 23 20 41 73  61 6e 61 20  2d 2d 20 68  65 72 65 20 |# Asana -- here |
00000010: 61 72 65 20  74 68 65 20  62 75 64 67  65 74 20 66 |are the budget f|
00000020: 69 67 75 72  65 73 20 79  6f 75 20 61  73 6b 65 64 |igures you asked|
00000030: 20 66 6f 72  21 0a 74 79  70 65 2c 71  75 61 6e 74 | for!.type,quant|
00000040: 69 74 79 0a  70 75 72 70  6c 65 2c 34  35 36 2e 37 |ity.purple,456.7|
00000050: 38 0a 67 72  65 65 6e 2c  36 37 38 2e  31 32 0a 6f |8.green,678.12.o|
00000060: 72 61 6e 67  65 2c 31 32  33 2e 34 35  0a          |range,123.45.|
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr hex -r data/budget.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
23 20 41 73  61 6e 61 20  2d 2d 20 68  65 72 65 20 
61 72 65 20  74 68 65 20  62 75 64 67  65 74 20 66 
69 67 75 72  65 73 20 79  6f 75 20 61  73 6b 65 64 
20 66 6f 72  21 0a 74 79  70 65 2c 71  75 61 6e 74 
69 74 79 0a  70 75 72 70  6c 65 2c 34  35 36 2e 37 
38 0a 67 72  65 65 6e 2c  36 37 38 2e  31 32 0a 6f 
72 61 6e 67  65 2c 31 32  33 2e 34 35  0a          
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr hex -r data/budget.csv | sed 's/20/2a/g' | mlr unhex</b>
</pre>
<pre class="pre-non-highlight-in-pair">
#*Asana*--*here*are*the*budget*figures*you*asked*for!
type,quantity
purple,456.78
green,678.12
orange,123.45
</pre>

Additionally, [`mlr help`](online-help.md), [`mlr repl`](repl.md), and [`mlr regtest`](https://github.com/johnkerl/miller/blob/main/test/README.md) are implemented here.
