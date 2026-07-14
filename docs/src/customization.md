<!--  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. -->
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

## Named profiles in your .mlrrc

You can group settings into INI-style named sections, called _profiles_, and select one with the
`--profile {name}` main flag, or its alias `-P {name}`. For example, given

<pre class="pre-non-highlight-non-pair">
    # Global settings, applied always:
    icsv

    [j]
    # Settings applied only with mlr --profile j (or mlr -P j):
    ojson
    jvstack

    [tsvout]
    # Settings applied only with mlr --profile tsvout (or mlr -P tsvout):
    otsv
</pre>

then `mlr cat myfile.csv` reads CSV and writes DKVP (the global setting applies, and the
sections are ignored), while `mlr -P j cat myfile.csv` reads CSV and writes vertically stacked
JSON (the global setting applies first, then the settings from the `[j]` section).

Semantics:

* Lines before any `[name]` section header are global settings, and are always applied. A
  `.mlrrc` file without any section headers behaves just as it did in older versions of Miller.

* With `--profile {name}` (or `-P {name}`), global settings are applied first, then the settings
  from the `[name]` section. It's a fatal error if no `[name]` section exists in any `.mlrrc`
  file processed -- or if no `.mlrrc` file was found at all.

* Without `--profile`, sections are ignored entirely -- their lines aren't even parsed -- so a
  typo inside an unused profile won't affect your other invocations of `mlr`.

* Section names are matched exactly (case-sensitively). Whitespace around and within the
  brackets is ignored: `[ j ]` is the same as `[j]`. Comments are allowed after section headers.

* If the same section name appears more than once, the settings from all its blocks are applied,
  in the order they appear in the file.

* If both `$HOME/.mlrrc` and `./.mlrrc` are processed, each file's global settings and matching
  section settings are applied in that per-file order, `$HOME/.mlrrc` first. The selected profile
  needs to exist in only one of them.

* Since `--profile` selects a section of your `.mlrrc`, it can't be combined with `--norc`, or
  with `MLRRC=__none__` in the environment -- that's a fatal error.

* Profiles are selected on the `mlr` command line, not from within a `.mlrrc` file: putting
  `--profile` (or `-P`) inside a `.mlrrc` file is a parse error, just as `--prepipe` is.

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
