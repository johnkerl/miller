<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
DSL reference: unset statements
# 
You can clear a map key by assigning the empty string as its value: `$x=""` or `@x=""`. Using `unset` you can remove the key entirely. Examples:

<pre class="pre-highlight-in-pair">
<b>cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put 'unset $x, $a' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
b=pan,i=1,y=0.7268028627434533
b=pan,i=2,y=0.5221511083334797
b=wye,i=3,y=0.33831852551664776
b=wye,i=4,y=0.13418874328430463
b=pan,i=5,y=0.8636244699032729
</pre>

This can also be done, of course, using `mlr cut -x`. You can also clear out-of-stream or local variables, at the base name level, or at an indexed sublevel:

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { dump; unset @sum; dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": {
    "pan": {
      "pan": 0.3467901443380824
    },
    "eks": {
      "pan": 0.7586799647899636,
      "wye": 0.38139939387114097
    },
    "wye": {
      "wye": 0.20460330576630303,
      "pan": 0.5732889198020006
    }
  }
}
{}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put -q '@sum[$a][$b] += $x; end { dump; unset @sum["eks"]; dump }' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "sum": {
    "pan": {
      "pan": 0.3467901443380824
    },
    "eks": {
      "pan": 0.7586799647899636,
      "wye": 0.38139939387114097
    },
    "wye": {
      "wye": 0.20460330576630303,
      "pan": 0.5732889198020006
    }
  }
}
{
  "sum": {
    "pan": {
      "pan": 0.3467901443380824
    },
    "wye": {
      "wye": 0.20460330576630303,
      "pan": 0.5732889198020006
    }
  }
}
</pre>

If you use `unset all` (or `unset @*` which is synonymous), that will unset all out-of-stream variables which have been defined up to that point.
