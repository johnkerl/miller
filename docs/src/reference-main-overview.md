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
# Miller command structure

## Overview

The outline of an invocation of Miller is:

* The program name `mlr`.
* Flags controlling input/output formatting, etc. (see the [flags reference](reference-main-flag-list.md)).
* One or more verbs -- such as `cut`, `sort`, etc. (see [verbs reference](reference-verbs.md)) -- chained together using [then](reference-main-then-chaining.md). You use these to transform your data.
* Zero or more filenames, with input taken from standard input if there are no filenames present. (You can place the filenames up front using `--from` or `--mfrom` as described on the [keystroke-savers page](keystroke-savers.md#file-names-up-front-including-from).)

For example, reading from a file:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint head -n 2 then sort -f shape example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag k index quantity rate
red    square   true 2 15    79.2778  0.0130
yellow triangle true 1 11    43.6498  9.8870
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --from example.csv --icsv --opprint head -n 2 then sort -f shape</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag k index quantity rate
red    square   true 2 15    79.2778  0.0130
yellow triangle true 1 11    43.6498  9.8870
</pre>

Reading from standard input:

<pre class="pre-highlight-in-pair">
<b>cat example.csv | mlr --icsv --opprint head -n 2 then sort -f shape</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag k index quantity rate
red    square   true 2 15    79.2778  0.0130
yellow triangle true 1 11    43.6498  9.8870
</pre>

The rest of this reference section gives you full information on each of these parts of the command line.

See also the [Glossary](glossary.md) for more about terms such as
[record](glossary.md#record), [field](glossary.md#field),
[key](glossary.md#value), [streaming](glossary.md#streaming), and more.

## Verbs vs DSL

When you type `mlr {something} myfile.dat`, the `{something}` part is called a **verb**. It specifies how you want to transform your data. Most of the verbs are counterparts of built-in system tools like `cut` and `sort` -- but with file-format awareness, and giving you the ability to refer to fields by name.

The verbs `put` and `filter` are special in that they have a rich expression language (domain-specific language, or "DSL"). More information about them can be found at on the [Intro to Miller's programming language page](miller-programming-language.md); see also [DSL reference](reference-dsl.md) for more details.

Here's a comparison of verbs and `put`/`filter` DSL expressions:

Example of using a verb for data processing:

<pre class="pre-highlight-in-pair">
<b>mlr stats1 -a sum -f x -g a data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,x_sum=0.346791
a=eks,x_sum=1.140078
a=wye,x_sum=0.777891
</pre>

* Verbs are coded in Go
* They run a bit faster
* They take fewer keystrokes
* There's less to learn
* Their customization is limited to each verb's options

Example of doing the same thing using a DSL expression:

<pre class="pre-highlight-in-pair">
<b>mlr  put -q '@x_sum[$a] += $x; end{emit @x_sum, "a"}' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,x_sum=0.346791
a=eks,x_sum=1.140078
a=wye,x_sum=0.777891
</pre>

* You get to write your own expressions in Miller's programming language
* They run a bit slower
* They take more keystrokes
* There's more to learn
* They're highly customizable

## Exit codes

When you do `echo $?` immediately after running a Miller command (or any shell command), that's the _exit code_.

Miller's are as follows:

* 0 for normal exit.
* 1 if an error was encountered. There should be helpful text written to `stderr`; please [file a bug report](https://github.com/johnkerl/miller/issues/new), ideally with a reproducible scenario, if the text is either missing or unhelpful.
* 2 if there was a `panic` in the Go runtime -- you'll see lots of stack-trace lines. Please [file a bug report](https://github.com/johnkerl/miller/issues/new), ideally with a reproducible scenario, if you ever see this.
