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
# Scripting with Miller

Suppose you are often doing something like

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint \</b>
<b>filter '$quantity != 20' \</b>
<b>then count-distinct -f shape \</b>
<b>then fraction -f count \</b>
<b>example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    count count_fraction
triangle 3     0.3
square   4     0.4
circle   3     0.3
</pre>

Typing this out can get a bit old, if the only thing that changes for you is the filename. Some options include:

* On Linux/Mac/etc you can make a script with `#!/bin/sh` which invokes Miller as part of the shell-script body.
* On Linux/Mac/etc you can make a script with `#!/usr/bin/env mlr -s` which invokes Miller.
* On any platform you can put the reusable part of your command line into a text file (say `myflags.txt`), then `mlr -s myflags-txt filename-which-varies.csv`.

Let's look at examples of each.

## Shell scripts

A shell-script option:

<pre class="pre-highlight-in-pair">
<b>cat example-shell-script</b>
</pre>
<pre class="pre-non-highlight-in-pair">
#!/bin/bash
mlr --c2p \
  filter '$quantity != 20' \
  then count-distinct -f shape \
  then fraction -f count \
  -- "$@"
</pre>

Key points here:

* Use `--` before `"$@"` at the end so that main-flags like `--json` won't be confused for options to the `fraction` verb.
* Use `"$@"` at the end which means "all remaining arguments to the script".
* Use `chmod +x example-shell-script` after you create one of these.

Then you can do

<pre class="pre-highlight-in-pair">
<b>example-shell-script example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    count count_fraction
triangle 3     0.3
square   4     0.4
circle   3     0.3
</pre>

<pre class="pre-highlight-in-pair">
<b>cat example.csv | example-shell-script</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    count count_fraction
triangle 3     0.3
square   4     0.4
circle   3     0.3
</pre>

<pre class="pre-highlight-in-pair">
<b>example-shell-script --ojson example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "shape": "triangle",
  "count": 3,
  "count_fraction": 0.3
}
{
  "shape": "square",
  "count": 4,
  "count_fraction": 0.4
}
{
  "shape": "circle",
  "count": 3,
  "count_fraction": 0.3
}
</pre>

<pre class="pre-highlight-in-pair">
<b>example-shell-script --ojson then filter '$count == 3' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "shape": "triangle",
  "count": 3,
  "count_fraction": 0.3
}
{
  "shape": "circle",
  "count": 3,
  "count_fraction": 0.3
}
</pre>

etc.

## Miller scripts

Here instead of putting `#!/bin/bash` on the first line, we can put `mlr` directly:

<pre class="pre-highlight-in-pair">
<b>cat example-mlr-s-script</b>
</pre>
<pre class="pre-non-highlight-in-pair">
#!/usr/bin/env mlr -s
--c2p
filter '$quantity != 20'
then count-distinct -f shape
then fraction -f count
</pre>

Points:

* This is largely the same as a shell script.
* Use `chmod +x example-mlr-s-script` after you create one of these.
* You leave off the initial `mlr` since that's present on line 1.
* You don't need all the backslashing for line-continuations.
* You don't need the explicit `--` or `"$@"`.

Then you can do

<pre class="pre-highlight-in-pair">
<b>example-mlr-s-script example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    count count_fraction
triangle 3     0.3
square   4     0.4
circle   3     0.3
</pre>

<pre class="pre-highlight-in-pair">
<b>cat example.csv | example-mlr-s-script</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    count count_fraction
triangle 3     0.3
square   4     0.4
circle   3     0.3
</pre>

<pre class="pre-highlight-in-pair">
<b>example-mlr-s-script --ojson example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "shape": "triangle",
  "count": 3,
  "count_fraction": 0.3
}
{
  "shape": "square",
  "count": 4,
  "count_fraction": 0.4
}
{
  "shape": "circle",
  "count": 3,
  "count_fraction": 0.3
}
</pre>

<pre class="pre-highlight-in-pair">
<b>example-mlr-s-script --ojson then filter '$count == 3' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "shape": "triangle",
  "count": 3,
  "count_fraction": 0.3
}
{
  "shape": "circle",
  "count": 3,
  "count_fraction": 0.3
}
</pre>

## Miller scripts on Windows

Both the previous options require executable mode with `chmod`, and a _shebang
line_ with `#!...`, which are unixisms.

One of the nice features of `mlr -s` is it can be done without a shebang line,
and this works fine on Windows. For example:

<pre class="pre-highlight-in-pair">
<b>cat example-mlr-s-script-no-shebang</b>
</pre>
<pre class="pre-non-highlight-in-pair">
--c2p
filter '$quantity != 20'
then count-distinct -f shape
then fraction -f count
</pre>

Points:

* Same as above, where the `#!` line isn't needed. (But you can include a `#!` line; `mlr -s` will simply see it as a comment line.).
* As above, you don't need all the backslashing for line-continuations.
* As above, you don't need the explicit `--` or `"$@"`.

Then you can do

<pre class="pre-highlight-in-pair">
<b>mlr -s example-mlr-s-script-no-shebang example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
shape    count count_fraction
triangle 3     0.3
square   4     0.4
circle   3     0.3
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -s example-mlr-s-script-no-shebang --ojson example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "shape": "triangle",
  "count": 3,
  "count_fraction": 0.3
}
{
  "shape": "square",
  "count": 4,
  "count_fraction": 0.4
}
{
  "shape": "circle",
  "count": 3,
  "count_fraction": 0.3
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -s example-mlr-s-script-no-shebang --ojson then filter '$count == 3' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "shape": "triangle",
  "count": 3,
  "count_fraction": 0.3
}
{
  "shape": "circle",
  "count": 3,
  "count_fraction": 0.3
}
</pre>

and so on. See also [Miller on Windows](miller-on-windows.md).
