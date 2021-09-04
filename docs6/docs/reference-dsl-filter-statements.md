<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
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
# DSL filter statements

You can use the `filter` DSL keyword within the `put` verb. In fact, the following two are synonymous:

<pre class="pre-highlight-in-pair">
<b>mlr --csv filter 'NR==2 || NR==3' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv put 'filter NR==2 || NR==3' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
red,square,true,2,15,79.2778,0.0130
red,circle,true,3,16,13.8103,2.9010
</pre>

The former, of course, is a little easier to type. For another example:

<pre class="pre-highlight-in-pair">
<b>mlr --csv put '@running_sum += $quantity; filter @running_sum > 500' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
purple,square,false,10,91,72.3735,8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv filter '@running_sum += $quantity; @running_sum > 500' example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color,shape,flag,k,index,quantity,rate
yellow,circle,true,8,73,63.9785,4.2370
yellow,circle,true,9,87,63.5058,8.3350
purple,square,false,10,91,72.3735,8.2430
</pre>
