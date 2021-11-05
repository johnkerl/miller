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
# Operating on all fields

## Bulk rename of fields

Suppose you want to replace spaces with underscores in your column names:

<pre class="pre-highlight-in-pair">
<b>cat data/spaces.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b c,def,g h i
123,4567,890
2468,1357,3579
9987,3312,4543
</pre>

The simplest way is to use `mlr rename` with `-g` (for global replace, not just first occurrence of space within each field) and `-r` for pattern-matching (rather than explicit single-column renames):

<pre class="pre-highlight-in-pair">
<b>mlr --csv rename -g -r ' ,_'  data/spaces.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a_b_c,def,g_h_i
123,4567,890
2468,1357,3579
9987,3312,4543
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv --opprint rename -g -r ' ,_'  data/spaces.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a_b_c def  g_h_i
123   4567 890
2468  1357 3579
9987  3312 4543
</pre>

You can also do this with a for-loop:

<pre class="pre-highlight-in-pair">
<b>cat data/bulk-rename-for-loop.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
map newrec = {};
for (oldk, v in $*) {
    newrec[gsub(oldk, " ", "_")] = v;
}
$* = newrec
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint put -f data/bulk-rename-for-loop.mlr data/spaces.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a_b_c def  g_h_i
123   4567 890
2468  1357 3579
9987  3312 4543
</pre>

## Search-and-replace over all fields

How to do `$name = gsub($name, "old", "new")` for all fields?

<pre class="pre-highlight-in-pair">
<b>cat data/sar.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
the quick,brown fox,jumped
over,the,lazy dogs
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/sar.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
for (k in $*) {
  $[k] = gsub($[k], "e", "X");
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv put -f data/sar.mlr data/sar.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
thX quick,brown fox,jumpXd
ovXr,thX,lazy dogs
</pre>

## Full field renames and reassigns

Using Miller 5.0.0's map literals and assigning to `$*`, you can fully generalize [rename](reference-verbs.md#rename), [reorder](reference-verbs.md#reorder), etc.

<pre class="pre-highlight-in-pair">
<b>cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '</b>
<b>  begin {</b>
<b>    @i_cumu = 0;</b>
<b>  }</b>
<b></b>
<b>  @i_cumu += $i;</b>
<b>  $* = {</b>
<b>    "z": $x + y,</b>
<b>    "KEYFIELD": $a,</b>
<b>    "i": @i_cumu,</b>
<b>    "b": $b,</b>
<b>    "y": $x,</b>
<b>    "x": $y,</b>
<b>  };</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
z=0.346791,KEYFIELD=pan,i=1,b=pan,y=0.346791,x=0.726802
z=0.758679,KEYFIELD=eks,i=3,b=pan,y=0.758679,x=0.522151
z=0.204603,KEYFIELD=wye,i=6,b=wye,y=0.204603,x=0.338318
z=0.381399,KEYFIELD=eks,i=10,b=wye,y=0.381399,x=0.134188
z=0.573288,KEYFIELD=wye,i=15,b=pan,y=0.573288,x=0.863624
</pre>
