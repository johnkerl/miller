<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# DSL reference: filter statements

You can use ``filter`` within ``put``. In fact, the following two are synonymous:

<pre>
<b>mlr filter 'NR==2 || NR==3' data/small</b>
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
</pre>

<pre>
<b>mlr put 'filter NR==2 || NR==3' data/small</b>
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
</pre>

The former, of course, is much easier to type. But the latter allows you to define more complex expressions for the filter, and/or do other things in addition to the filter:

<pre>
<b>mlr put '@running_sum += $x; filter @running_sum > 1.3' data/small</b>
a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729
</pre>

<pre>
<b>mlr put '$z = $x * $y; filter $z > 0.3' data/small</b>
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,z=0.3961455844854848
a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,z=0.4951063394654227
</pre>
