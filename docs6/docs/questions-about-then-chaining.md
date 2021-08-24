<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
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
</pre>

## NR is not consecutive after then-chaining

Given this input data:

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

why don't I see `NR=1` and `NR=2` here??

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small filter '$x > 0.5' then put '$NR = NR'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,NR=2
a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,NR=5
</pre>

The reason is that `NR` is computed for the original input records and isn't dynamically updated. By contrast, `NF` is dynamically updated: it's the number of fields in the current record, and if you add/remove a field, the value of `NF` will change:

<pre class="pre-highlight-in-pair">
<b>echo x=1,y=2,z=3 | mlr put '$nf1 = NF; $u = 4; $nf2 = NF; unset $x,$y,$z; $nf3 = NF'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
nf1=3,u=4,nf2=5,nf3=3
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
a   b   i x                  y                  nr1 nr2
eks pan 2 0.7586799647899636 0.5221511083334797 2   1
wye pan 5 0.5732889198020006 0.8636244699032729 5   2
</pre>

Or, simply use `mlr cat -n`:

<pre class="pre-highlight-in-pair">
<b>mlr filter '$x > 0.5' then cat -n data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
n=1,a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
n=2,a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729
</pre>
