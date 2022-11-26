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
# Questions about then-chaining

## How do I examine then-chaining?

Then-chaining found in Miller is intended to function the same as Unix pipes, but with less keystroking. You can print your data one pipeline step at a time, to see what intermediate output at one step becomes the input to the next step.

First, look at the input data:

<pre class="pre-highlight-in-pair">
<b>cat data/then-example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Status,Payment_Type,Amount
paid,cash,10.00
pending,debit,20.00
paid,cash,50.00
pending,credit,40.00
paid,debit,30.00
</pre>

Next, run the first step of your command, omitting anything from the first `then` onward:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/then-example.csv --c2p count-distinct -f Status,Payment_Type</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Status  Payment_Type count
paid    cash         2
pending debit        1
pending credit       1
paid    debit        1
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

After that, run it with the next `then` step included:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/then-example.csv --c2p count-distinct -f Status,Payment_Type \</b>
<b>  then sort -nr count</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Status  Payment_Type count
paid    cash         2
pending debit        1
pending credit       1
paid    debit        1
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

Now if you use `then` to include another verb after that, the columns `Status`, `Payment_Type`, and `count` will be the input to that verb.

Note, by the way, that you'll get the same results using pipes:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/then-example.csv --csv count-distinct -f Status,Payment_Type \</b>
<b>| mlr --c2p sort -nr count</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Status  Payment_Type count
paid    cash         2
pending debit        1
pending credit       1
paid    debit        1
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

## NR is not consecutive after then-chaining

Given this input data:

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

why don't I see `NR=1` and `NR=2` here??

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small filter '$x > 0.5' then put '$NR = NR'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=eks,b=pan,i=2,x=0.758679,y=0.522151,NR=2
a=wye,b=pan,i=5,x=0.573288,y=0.863624,NR=5
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

The reason is that `NR` is computed for the original input records and isn't dynamically updated. By contrast, `NF` is dynamically updated: it's the number of fields in the current record, and if you add/remove a field, the value of `NF` will change:

<pre class="pre-highlight-in-pair">
<b>echo x=1,y=2,z=3 | mlr put '$nf1 = NF; $u = 4; $nf2 = NF; unset $x,$y,$z; $nf3 = NF'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
nf1=3,u=4,nf2=5,nf3=3
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

`NR`, by contrast (and `FNR` as well), retains the value from the original input stream, and records may be dropped by a `filter` within a `then`-chain. To recover consecutive record numbers, you can use out-of-stream variables as follows:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --from data/small put '</b>
<b>  begin{ @nr1 = 0 }</b>
<b>  @nr1 += 1;</b>
<b>  $nr1 = @nr1</b>
<b>' \</b>
<b>then filter '$x>0.5' \</b>
<b>then put '</b>
<b>  begin{ @nr2 = 0 }</b>
<b>  @nr2 += 1;</b>
<b>  $nr2 = @nr2</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        nr1 nr2
eks pan 2 0.758679 0.522151 2   1
wye pan 5 0.573288 0.863624 5   2
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

Or, simply use `mlr cat -n`:

<pre class="pre-highlight-in-pair">
<b>mlr filter '$x > 0.5' then cat -n data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
n=1,a=eks,b=pan,i=2,x=0.758679,y=0.522151
n=2,a=wye,b=pan,i=5,x=0.573288,y=0.863624
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>
