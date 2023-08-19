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
# Customization: .mlrrc

## How to use .mlrrc

Suppose you always use CSV files. Then instead of always having to type `--csv` as in

<pre class="pre-highlight-non-pair">
<b>mlr --csv cut -x -f extra mydata.csv</b>
</pre>

<pre class="pre-highlight-non-pair">
<b>mlr --csv sort -n id mydata.csv</b>
</pre>

and so on, you can instead put the following into your `$HOME/.mlrrc`:

<pre class="pre-non-highlight-non-pair">
    --csv
</pre>

Then you can just type things like

<pre class="pre-highlight-non-pair">
<b>mlr cut -x -f extra mydata.csv</b>
</pre>

<pre class="pre-highlight-non-pair">
<b>mlr sort -n id mydata.csv</b>
</pre>

and the `--csv` part will automatically be understood. If you do want to process, say, a JSON file then `mlr --json ...` at the command line will still override the defaults you've placed in your `.mlrrc`.

## What you can put in your .mlrrc

* You can include any command-line flags, except the "terminal" ones such as `--help`.

* The `--prepipe`, `--load`, and `--mload` flags aren't allowed in `.mlrrc` as they control code execution, and could result in your scripts running things you don't expect if you receive data from someone with a `./.mlrrc` in it. You can use `--prepipe-bz2`, `--prepipe-gunzip`, `--prepipe-zcat`, and `--prepipe-zstdcat` in `.mlrrc`, though.

* The formatting rule is you need to put one flag beginning with `--` per line: for example, `--csv` on one line and `--nr-progress-mod 1000` on a separate line.

* Since every line starts with a `--` option, you can leave off the initial `--` if you want. For example, `ojson` is the same as `--ojson`, and `nr-progress-mod 1000` is the same as `--nr-progress-mod 1000`.

* Comments are from a `#` to the end of the line.

* Empty lines are ignored -- including lines which are empty after comments are removed.

Here is an example `.mlrrc` file:

<pre class="pre-non-highlight-non-pair">
# Input and output formats are CSV by default (unless otherwise specified
# on the mlr command line):
csv

# If a data line has fewer fields than the header line, instead of erroring
# (which is the default), just insert empty values for the missing ones:
allow-ragged-csv-input

# These are no-ops for CSV, but when I do use JSON output, I want these
# pretty-printing options to be used:
jvstack
jlistwrap

# Use "@", rather than "#", for comments within data files:
skip-comments-with @
</pre>

## Where to put your .mlrrc

If the environment variable `MLRRC` is set:

* If its value is `__none__` then no `.mlrrc` files are processed.  (This is nice for things like regression testing.)

* Otherwise, its value (as a filename) is loaded and processed. If there are syntax errors, they abort `mlr` with a usage message (as if you had mistyped something on the command line). If the file can't be loaded at all, though, it is silently skipped.

* Any `.mlrrc` in your home directory or current directory is ignored whenever `MLRRC` is set in the environment.

* Example line in your shell's rc file: `export MLRRC=/path/to/my/mlrrc`

Otherwise:

* If `$HOME/.mlrrc` exists, it's processed as above.

* If `./.mlrrc` exists, it's then also processed as above.

* The idea is you can have all your settings in your `$HOME/.mlrrc`, then maybe more project-specific settings for your current directory if you like.
