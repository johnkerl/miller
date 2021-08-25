<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Running shell commands

TODO: while-read example from issues

The [system](reference-dsl.md#system) DSL function allows you to run a specific shell command and put its output -- minus the final newline -- into a record field. The command itself is any string, either a literal string, or a concatenation of strings, perhaps including other field values or what have you.

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '$o = system("echo hello world")' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        o
pan pan 1 0.346791 0.726802 hello world
eks pan 2 0.758679 0.522151 hello world
wye wye 3 0.204603 0.338318 hello world
eks wye 4 0.381399 0.134188 hello world
wye pan 5 0.573288 0.863624 hello world
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '$o = system("echo {" . NR . "}")' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        o
pan pan 1 0.346791 0.726802 {1}
eks pan 2 0.758679 0.522151 {2}
wye wye 3 0.204603 0.338318 {3}
eks wye 4 0.381399 0.134188 {4}
wye pan 5 0.573288 0.863624 {5}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '$o = system("echo -n ".$a."| md5")' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        o
pan pan 1 0.346791 0.726802 ccc62bbd08bdc21905f4909463ccdf7c
eks pan 2 0.758679 0.522151 585d25a8ff04840f77779eeff61167dc
wye wye 3 0.204603 0.338318 fb6361a373147c163e65ada94719fa16
eks wye 4 0.381399 0.134188 585d25a8ff04840f77779eeff61167dc
wye pan 5 0.573288 0.863624 fb6361a373147c163e65ada94719fa16
</pre>

Note that running a subprocess on every record takes a non-trivial amount of time. Comparing asking the system `date` command for the current time in nanoseconds versus computing it in process:

<!--- hard-coded, not live-code, since %N doesn't exist on all platforms -->

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '$t=system("date +%s.%N")' then step -a delta -f t data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x                   y                   t                    t_delta
pan pan 1 0.3467901443380824  0.7268028627434533  1568774318.513903817 0
eks pan 2 0.7586799647899636  0.5221511083334797  1568774318.514722876 0.000819
wye wye 3 0.20460330576630303 0.33831852551664776 1568774318.515618046 0.000895
eks wye 4 0.38139939387114097 0.13418874328430463 1568774318.516547441 0.000929
wye pan 5 0.5732889198020006  0.8636244699032729  1568774318.517518828 0.000971
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '$t=systime()' then step -a delta -f t data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x                   y                   t                 t_delta
pan pan 1 0.3467901443380824  0.7268028627434533  1568774318.518699 0
eks pan 2 0.7586799647899636  0.5221511083334797  1568774318.518717 0.000018
wye wye 3 0.20460330576630303 0.33831852551664776 1568774318.518723 0.000006
eks wye 4 0.38139939387114097 0.13418874328430463 1568774318.518727 0.000004
wye pan 5 0.5732889198020006  0.8636244699032729  1568774318.518730 0.000003
</pre>
