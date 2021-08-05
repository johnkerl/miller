<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Reference: Miller commands

## Overview

The outline of an invocation of Miller is

* `mlr`
* Options controlling input/output formatting, etc. ([Reference: I/O options](reference-main-io-options.md)).
* One or more verbs (such as `cut`, `sort`, etc.) ([Verbs Reference](reference-verbs.md)) -- chained together using [then](reference-main-then-chaining.md)). You use these to transform your data.
* Zero or more filenames, with input taken from standard input if there are no filenames present.

For example, reading from a file:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint head -n 2 then sort -f shape example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag index quantity rate
red    square   true 15    79.2778  0.0130
yellow triangle true 11    43.6498  9.8870
</pre>

Reading from standard input:

<pre class="pre-highlight-in-pair">
<b>cat example.csv | mlr --icsv --opprint head -n 2 then sort -f shape</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag index quantity rate
red    square   true 15    79.2778  0.0130
yellow triangle true 11    43.6498  9.8870
</pre>

The rest of this reference section gives you full information on each of these parts of the command line.

## Verbs vs DSL

When you type `mlr {something} myfile.dat`, the `{something}` part is called a **verb**. It specifies how you want to transform your data. Most of the verbs are counterparts of built-in system tools like `cut` and `sort` -- but with file-format awareness, and giving you the ability to refer to fields by name.

The verbs `put` and `filter` are special in that they have a rich expression language (domain-specific language, or "DSL"). More information about them can be found at [DSL reference](reference-dsl.md).

Here's a comparison of verbs and `put`/`filter` DSL expressions:

Example of using a verb for data processing:

<pre class="pre-highlight-in-pair">
<b>mlr stats1 -a sum -f x -g a data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,x_sum=0.3467901443380824
a=eks,x_sum=1.1400793586611044
a=wye,x_sum=0.7778922255683036
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
a=pan,x_sum=0.3467901443380824
a=eks,x_sum=1.1400793586611044
a=wye,x_sum=0.7778922255683036
</pre>

* You get to write your own expressions in Miller's programming language
* They run a bit slower
* They take more keystrokes
* There's more to learn
* They're highly customizable
