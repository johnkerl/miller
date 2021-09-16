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
# Data-cleaning examples

Here are some ways to use the type-checking options as described in the [Type-checking page](reference-dsl-variables.md#type-checking).  Suppose you have the following data file, with inconsistent typing for boolean. (Also imagine that, for the sake of discussion, we have a million-line file rather than a four-line file, so we can't see it all at once and some automation is called for.)

<pre class="pre-highlight-in-pair">
<b>cat data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name,reachable
barney,false
betty,true
fred,true
wilma,1
</pre>

One option is to coerce everything to boolean, or integer:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint put '$reachable = boolean($reachable)' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name   reachable
barney false
betty  true
fred   true
wilma  true
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint put '$reachable = int(boolean($reachable))' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name   reachable
barney 0
betty  1
fred   1
wilma  1
</pre>

A second option is to flag badly formatted data within the output stream:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint put '$format_ok = is_string($reachable)' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
name   reachable format_ok
barney false     true
betty  true      true
fred   true      true
wilma  1         false
</pre>

Or perhaps to flag badly formatted data outside the output stream:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint put '</b>
<b>  if (!is_string($reachable)) {eprint "Malformed at NR=".NR}</b>
<b>' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Malformed at NR=4
name   reachable
barney false
betty  true
fred   true
wilma  1
</pre>

A third way is to abort the process on first instance of bad data:

<pre class="pre-highlight-in-pair">
<b>mlr --csv put '$reachable = asserting_string($reachable)' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Miller: is_string type-assertion failed at NR=4 FNR=4 FILENAME=data/het-bool.csv
name,reachable
barney,false
</pre>
