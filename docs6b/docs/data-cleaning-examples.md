<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Data-cleaning examples

Here are some ways to use the type-checking options as described in [Type-checking](reference-dsl-variables.md#type-checking).  Suppose you have the following data file, with inconsistent typing for boolean. (Also imagine that, for the sake of discussion, we have a million-line file rather than a four-line file, so we can't see it all at once and some automation is called for.)

<pre class="pre-highlight">
<b>cat data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight">
name,reachable
barney,false
betty,true
fred,true
wilma,1
</pre>

One option is to coerce everything to boolean, or integer:

<pre class="pre-highlight">
<b>mlr --icsv --opprint put '$reachable = boolean($reachable)' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight">
name   reachable
barney false
betty  true
fred   true
wilma  true
</pre>

<pre class="pre-highlight">
<b>mlr --icsv --opprint put '$reachable = int(boolean($reachable))' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight">
name   reachable
barney 0
betty  1
fred   1
wilma  1
</pre>

A second option is to flag badly formatted data within the output stream:

<pre class="pre-highlight">
<b>mlr --icsv --opprint put '$format_ok = is_string($reachable)' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight">
name   reachable format_ok
barney false     false
betty  true      false
fred   true      false
wilma  1         false
</pre>

Or perhaps to flag badly formatted data outside the output stream:

<pre class="pre-highlight">
<b>mlr --icsv --opprint put '</b>
<b>  if (!is_string($reachable)) {eprint "Malformed at NR=".NR}</b>
<b>' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight">
Malformed at NR=1
Malformed at NR=2
Malformed at NR=3
Malformed at NR=4
name   reachable
barney false
betty  true
fred   true
wilma  1
</pre>

A third way is to abort the process on fimd.instance of bad data:

<pre class="pre-highlight">
<b>mlr --csv put '$reachable = asserting_string($reachable)' data/het-bool.csv</b>
</pre>
<pre class="pre-non-highlight">
Miller: is_string type-assertion failed at NR=1 FNR=1 FILENAME=data/het-bool.csv
</pre>
