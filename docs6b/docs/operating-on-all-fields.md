<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Operating on all fields

## Bulk rename of fields

Suppose you want to replace spaces with underscores in your column names:

<pre class="pre-highlight">
<b>cat data/spaces.csv</b>
</pre>
<pre class="pre-non-highlight">
a b c,def,g h i
123,4567,890
2468,1357,3579
9987,3312,4543
</pre>

The simplest way is to use `mlr rename` with `-g` (for global replace, not just first occurrence of space within each field) and `-r` for pattern-matching (rather than explicit single-column renames):

<pre class="pre-highlight">
<b>mlr --csv rename -g -r ' ,_'  data/spaces.csv</b>
</pre>
<pre class="pre-non-highlight">
a_b_c,def,g_h_i
123,4567,890
2468,1357,3579
9987,3312,4543
</pre>

<pre class="pre-highlight">
<b>mlr --csv --opprint rename -g -r ' ,_'  data/spaces.csv</b>
</pre>
<pre class="pre-non-highlight">
a_b_c def  g_h_i
123   4567 890
2468  1357 3579
9987  3312 4543
</pre>

You can also do this with a for-loop:

<pre class="pre-highlight">
<b>cat data/bulk-rename-for-loop.mlr</b>
</pre>
<pre class="pre-non-highlight">
map newrec = {};
for (oldk, v in $*) {
    newrec[gsub(oldk, " ", "_")] = v;
}
$* = newrec
</pre>

<pre class="pre-highlight">
<b>mlr --icsv --opprint put -f data/bulk-rename-for-loop.mlr data/spaces.csv</b>
</pre>
<pre class="pre-non-highlight">
a_b_c def  g_h_i
123   4567 890
2468  1357 3579
9987  3312 4543
</pre>

## Search-and-replace over all fields

How to do `$name = gsub($name, "old", "new")` for all fields?

<pre class="pre-highlight">
<b>cat data/sar.csv</b>
</pre>
<pre class="pre-non-highlight">
a,b,c
the quick,brown fox,jumped
over,the,lazy dogs
</pre>

<pre class="pre-highlight">
<b>cat data/sar.mlr</b>
</pre>
<pre class="pre-non-highlight">
  for (k in $*) {
    $[k] = gsub($[k], "e", "X");
  }
</pre>

<pre class="pre-highlight">
<b>mlr --csv put -f data/sar.mlr data/sar.csv</b>
</pre>
<pre class="pre-non-highlight">
a,b,c
thX quick,brown fox,jumpXd
ovXr,thX,lazy dogs
</pre>

## Full field renames and reassigns

Using Miller 5.0.0's map literals and assigning to `$*`, you can fully generalize [rename](reference-verbs.md#rename), [reorder](reference-verbs.md#reorder), etc.

<pre class="pre-highlight">
<b>cat data/small</b>
</pre>
<pre class="pre-non-highlight">
a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729
</pre>

<pre class="pre-highlight">
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
<pre class="pre-non-highlight">
z=0.3467901443380824,KEYFIELD=pan,i=1,b=pan,y=0.3467901443380824,x=0.7268028627434533
z=0.7586799647899636,KEYFIELD=eks,i=3,b=pan,y=0.7586799647899636,x=0.5221511083334797
z=0.20460330576630303,KEYFIELD=wye,i=6,b=wye,y=0.20460330576630303,x=0.33831852551664776
z=0.38139939387114097,KEYFIELD=eks,i=10,b=wye,y=0.38139939387114097,x=0.13418874328430463
z=0.5732889198020006,KEYFIELD=wye,i=15,b=pan,y=0.5732889198020006,x=0.8636244699032729
</pre>
