<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Two-pass algorithms

## Overview

Miller is a streaming record processor; commands are performed once per record. This makes Miller particularly suitable for single-pass algorithms, allowing many of its verbs to process files that are (much) larger than the amount of RAM present in your system. (Of course, Miller verbs such as `sort`, `tac`, etc. all must ingest and retain all input records before emitting any output records.) You can also use out-of-stream variables to perform multi-pass computations, at the price of retaining all input records in memory.

One of Miller's strengths is its compact notation: for example, given input of the form

<pre class="pre-highlight-in-pair">
<b>head -n 5 ./data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776
a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729
</pre>

you can simply do

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab stats1 -a sum -f x ./data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_sum 4986.019681679581
</pre>

or

<pre class="pre-highlight-in-pair">
<b>mlr --opprint stats1 -a sum -f x -g b ./data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
b   x_sum
pan 965.7636699425815
wye 1023.5484702619565
zee 979.7420161495838
eks 1016.7728571314786
hat 1000.192668193983
</pre>

rather than the more tedious

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab put -q '</b>
<b>  @x_sum += $x;</b>
<b>  end {</b>
<b>    emit @x_sum</b>
<b>  }</b>
<b>' data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_sum 4986.019681679581
</pre>

or

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put -q '</b>
<b>  @x_sum[$b] += $x;</b>
<b>  end {</b>
<b>    emit @x_sum, "b"</b>
<b>  }</b>
<b>' data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
b   x_sum
pan 965.7636699425815
wye 1023.5484702619565
zee 979.7420161495838
eks 1016.7728571314786
hat 1000.192668193983
</pre>

The former (`mlr stats1` et al.) has the advantages of being easier to type, being less error-prone to type, and running faster.

Nonetheless, out-of-stream variables (which I whimsically call *oosvars*), begin/end blocks, and emit statements give you the ability to implement logic -- if you wish to do so -- which isn't present in other Miller verbs.  (If you find yourself often using the same out-of-stream-variable logic over and over, please file a request at [https://github.com/johnkerl/miller/issues](https://github.com/johnkerl/miller/issues) to get it implemented directly in Go as a Miller verb of its own.)

The following examples compute some things using oosvars which are already computable using Miller verbs, by way of providing food for thought.

## Computation of percentages

For example, mapping numeric values down a column to the percentage between their min and max values is two-pass: on the first pass you find the min and max values, then on the second, map each record's value to a percentage.

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put -q '</b>
<b>  # These are executed once per record, which is the first pass.</b>
<b>  # The key is to use NR to index an out-of-stream variable to</b>
<b>  # retain all the x-field values.</b>
<b>  @x_min = min($x, @x_min);</b>
<b>  @x_max = max($x, @x_max);</b>
<b>  @x[NR] = $x;</b>
<b></b>
<b>  # The second pass is in a for-loop in an end-block.</b>
<b>  end {</b>
<b>    for (nr, x in @x) {</b>
<b>      @x_pct[nr] = 100 * (x - @x_min) / (@x_max - @x_min);</b>
<b>    }</b>
<b>    emit (@x, @x_pct), "NR"</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
NR x                   x_pct
1  0.3467901443380824  25.66194338926441
2  0.7586799647899636  100
3  0.20460330576630303 0
4  0.38139939387114097 31.90823602213647
5  0.5732889198020006  66.54054236562845
</pre>

## Line-number ratios

Similarly, finding the total record count requires first reading through all the data:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --from data/small put -q '</b>
<b>  @records[NR] = $*;</b>
<b>  end {</b>
<b>    for((I,k),v in @records) {</b>
<b>      @records[I]["I"] = I;</b>
<b>      @records[I]["N"] = NR;</b>
<b>      @records[I]["PCT"] = 100*I/NR</b>
<b>    }</b>
<b>    emit @records,"I"</b>
<b>  }</b>
<b>' then reorder -f I,N,PCT</b>
</pre>
<pre class="pre-non-highlight-in-pair">
I N PCT     a   b   i x                   y
1 5 (error) pan pan 1 0.3467901443380824  0.7268028627434533
2 5 (error) eks pan 2 0.7586799647899636  0.5221511083334797
3 5 (error) wye wye 3 0.20460330576630303 0.33831852551664776
4 5 (error) eks wye 4 0.38139939387114097 0.13418874328430463
5 5 (error) wye pan 5 0.5732889198020006  0.8636244699032729
</pre>

## Records having max value

The idea is to retain records having the largest value of `n` in the following data:

<pre class="pre-highlight-in-pair">
<b>mlr --itsv --opprint cat data/maxrows.tsv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a      b      n score
purple red    5 0.743231
blue   purple 2 0.093710
red    purple 2 0.802103
purple red    5 0.389055
red    purple 2 0.880457
orange red    2 0.540349
purple purple 1 0.634451
orange purple 5 0.257223
orange purple 5 0.693499
red    red    4 0.981355
blue   purple 5 0.157052
purple purple 1 0.441784
red    purple 1 0.124912
orange blue   1 0.921944
blue   purple 4 0.490909
purple red    5 0.454779
green  purple 4 0.198278
orange blue   5 0.705700
red    red    3 0.940705
purple red    5 0.072936
orange blue   3 0.389463
orange purple 2 0.664985
blue   purple 1 0.371813
red    purple 4 0.984571
green  purple 5 0.203577
green  purple 3 0.900873
purple purple 0 0.965677
blue   purple 2 0.208785
purple purple 1 0.455077
red    purple 4 0.477187
blue   red    4 0.007487
</pre>

Of course, the largest value of `n` isn't known until after all data have been read. Using an out-of-stream variable we can retain all records as they are read, then filter them at the end:

<pre class="pre-highlight-in-pair">
<b>cat data/maxrows.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
# Retain all records
@records[NR] = $*;
# Track max value of n
@maxn = max(@maxn, $n);

# After all records have been read, loop through retained records
# and print those with the max n value.
end {
  for (nr in @records) {
    map record = @records[nr];
    if (record["n"] == @maxn) {
      emit record;
    }
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --itsv --opprint put -q -f data/maxrows.mlr data/maxrows.tsv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a      b      n score
purple red    5 0.743231
purple red    5 0.389055
orange purple 5 0.257223
orange purple 5 0.693499
blue   purple 5 0.157052
purple red    5 0.454779
orange blue   5 0.705700
purple red    5 0.072936
green  purple 5 0.203577
</pre>

## Feature-counting

Suppose you have some heterogeneous data like this:

<pre class="pre-non-highlight-non-pair">
{ "qoh": 29874, "rate": 1.68, "latency": 0.02 }
{ "name": "alice", "uid": 572 }
{ "qoh": 1227, "rate": 1.01, "latency": 0.07 }
{ "qoh": 13458, "rate": 1.72, "latency": 0.04 }
{ "qoh": 56782, "rate": 1.64 }
{ "qoh": 23512, "rate": 1.71, "latency": 0.03 }
{ "qoh": 9876, "rate": 1.89, "latency": 0.08 }
{ "name": "bill", "uid": 684 }
{ "name": "chuck", "uid2": 908 }
{ "name": "dottie", "uid": 440 }
{ "qoh": 0, "rate": 0.40, "latency": 0.01 }
{ "qoh": 5438, "rate": 1.56, "latency": 0.17 }
</pre>

A reasonable question to ask is, how many occurrences of each field are there? And, what percentage of total row count has each of them? Since the denominator of the percentage is not known until the end, this is a two-pass algorithm:

<pre class="pre-non-highlight-non-pair">
for (key in $*) {
  @key_counts[key] += 1;
}
@record_count += 1;

end {
  for (key in @key_counts) {
      @key_fraction[key] = @key_counts[key] / @record_count
  }
  emit @record_count;
  emit @key_counts, "key";
  emit @key_fraction,"key"
}
</pre>

Then

<pre class="pre-highlight-in-pair">
<b>mlr --json put -q -f data/feature-count.mlr data/features.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "record_count": 12
}
{
  "key": "qoh",
  "key_counts": 8
}
{
  "key": "rate",
  "key_counts": 8
}
{
  "key": "latency",
  "key_counts": 7
}
{
  "key": "name",
  "key_counts": 4
}
{
  "key": "uid",
  "key_counts": 3
}
{
  "key": "uid2",
  "key_counts": 1
}
{
  "key": "qoh",
  "key_fraction": 0.6666666666666666
}
{
  "key": "rate",
  "key_fraction": 0.6666666666666666
}
{
  "key": "latency",
  "key_fraction": 0.5833333333333334
}
{
  "key": "name",
  "key_fraction": 0.3333333333333333
}
{
  "key": "uid",
  "key_fraction": 0.25
}
{
  "key": "uid2",
  "key_fraction": 0.08333333333333333
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint put -q -f data/feature-count.mlr data/features.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
record_count
12

key     key_counts
qoh     8
rate    8
latency 7
name    4
uid     3
uid2    1

key     key_fraction
qoh     0.6666666666666666
rate    0.6666666666666666
latency 0.5833333333333334
name    0.3333333333333333
uid     0.25
uid2    0.08333333333333333
</pre>

## Unsparsing

The previous section discussed how to fill out missing data fields within CSV with full header line -- so the list of all field names is present within the header line. Next, let's look at a related problem: we have data where each record has various key names but we want to produce rectangular output having the union of all key names.

For example, suppose you have JSON input like this:

<pre class="pre-highlight-in-pair">
<b>cat data/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{"a":1,"b":2,"v":3}
{"u":1,"b":2}
{"a":1,"v":2,"x":3}
{"v":1,"w":2}
</pre>

There are field names `a`, `b`, `v`, `u`, `x`, `w` in the data -- but not all in every record.  Since we don't know the names of all the keys until we've read them all, this needs to be a two-pass algorithm. On the first pass, remember all the unique key names and all the records; on the second pass, loop through the records filling in absent values, then producing output. Use `put -q` since we don't want to produce per-record output, only emitting output in the `end` block:

<pre class="pre-highlight-in-pair">
<b>cat data/unsparsify.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
# First pass:
# Remember all unique key names:
for (k in $*) {
  @all_keys[k] = 1;
}
# Remember all input records:
@records[NR] = $*;

# Second pass:
end {
  for (nr in @records) {
    # Get the sparsely keyed input record:
    irecord = @records[nr];
    # Fill in missing keys with empty string:
    map orecord = {};
    for (k in @all_keys) {
      if (haskey(irecord, k)) {
        orecord[k] = irecord[k];
      } else {
        orecord[k] = "";
      }
    }
    # Produce the output:
    emit orecord;
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --json put -q -f data/unsparsify.mlr data/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": 2,
  "v": 3,
  "u": "",
  "x": "",
  "w": ""
}
{
  "a": "",
  "b": 2,
  "v": "",
  "u": 1,
  "x": "",
  "w": ""
}
{
  "a": 1,
  "b": "",
  "v": 2,
  "u": "",
  "x": 3,
  "w": ""
}
{
  "a": "",
  "b": "",
  "v": 1,
  "u": "",
  "x": "",
  "w": 2
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --ocsv put -q -f data/unsparsify.mlr data/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,v,u,x,w
1,2,3,,,
,2,,1,,
1,,2,,3,
,,1,,,2
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint put -q -f data/unsparsify.mlr data/sparse.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b v u x w
1 2 3 - - -
- 2 - 1 - -
1 - 2 - 3 -
- - 1 - - 2
</pre>

There is a keystroke-saving verb for this: [unsparsify](reference-verbs.md#unsparsify).

## Mean without/with oosvars

<pre class="pre-highlight-in-pair">
<b>mlr --opprint stats1 -a mean -f x data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_mean
0.49860196816795804
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put -q '</b>
<b>  @x_sum += $x;</b>
<b>  @x_count += 1;</b>
<b>  end {</b>
<b>    @x_mean = @x_sum / @x_count;</b>
<b>    emit @x_mean</b>
<b>  }</b>
<b>' data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_mean
0.49860196816795804
</pre>

## Keyed mean without/with oosvars

<pre class="pre-highlight-in-pair">
<b>mlr --opprint stats1 -a mean -f x -g a,b data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   x_mean
pan pan 0.5133141190437597
eks pan 0.48507555383425127
wye wye 0.49150092785839306
eks wye 0.4838950517724162
wye pan 0.4996119901034838
zee pan 0.5198298297816007
eks zee 0.49546320772681596
zee wye 0.5142667998230479
hat wye 0.49381326184632596
pan wye 0.5023618498923658
zee eks 0.4883932942792647
hat zee 0.5099985721987774
hat eks 0.48587864619953547
wye hat 0.4977304763723314
pan eks 0.5036718595143479
eks eks 0.5227992666570941
hat hat 0.47993053101017374
hat pan 0.4643355557376876
zee zee 0.5127559183726382
pan hat 0.492140950155604
pan zee 0.4966041598627583
zee hat 0.46772617655014515
wye zee 0.5059066170573692
eks hat 0.5006790659966355
wye eks 0.5306035254809106
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put -q '</b>
<b>  @x_sum[$a][$b] += $x;</b>
<b>  @x_count[$a][$b] += 1;</b>
<b>  end{</b>
<b>    for ((a, b), v in @x_sum) {</b>
<b>      @x_mean[a][b] = @x_sum[a][b] / @x_count[a][b];</b>
<b>    }</b>
<b>    emit @x_mean, "a", "b"</b>
<b>  }</b>
<b>' data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   x_mean
pan pan 0.5133141190437597
pan wye 0.5023618498923658
pan eks 0.5036718595143479
pan hat 0.492140950155604
pan zee 0.4966041598627583
eks pan 0.48507555383425127
eks wye 0.4838950517724162
eks zee 0.49546320772681596
eks eks 0.5227992666570941
eks hat 0.5006790659966355
wye wye 0.49150092785839306
wye pan 0.4996119901034838
wye hat 0.4977304763723314
wye zee 0.5059066170573692
wye eks 0.5306035254809106
zee pan 0.5198298297816007
zee wye 0.5142667998230479
zee eks 0.4883932942792647
zee zee 0.5127559183726382
zee hat 0.46772617655014515
hat wye 0.49381326184632596
hat zee 0.5099985721987774
hat eks 0.48587864619953547
hat hat 0.47993053101017374
hat pan 0.4643355557376876
</pre>

## Variance and standard deviation without/with oosvars

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab stats1 -a count,sum,mean,var,stddev -f x data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_count  10000
x_sum    4986.019681679581
x_mean   0.49860196816795804
x_var    0.08426974433144456
x_stddev 0.2902925151144007
</pre>

<pre class="pre-highlight-in-pair">
<b>cat variance.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
@n += 1;
@sumx += $x;
@sumx2 += $x**2;
end {
  @mean = @sumx / @n;
  @var = (@sumx2 - @mean * (2 * @sumx - @n * @mean)) / (@n - 1);
  @stddev = sqrt(@var);
  emitf @n, @sumx, @sumx2, @mean, @var, @stddev
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab put -q -f variance.mlr data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
n      10000
sumx   4986.019681679581
sumx2  3328.652400179729
mean   0.49860196816795804
var    0.08426974433144456
stddev 0.2902925151144007
</pre>

You can also do this keyed, of course, imitating the keyed-mean example above.

## Min/max without/with oosvars

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab stats1 -a min,max -f x data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_min 0.00004509679127584487
x_max 0.999952670371898
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab put -q '</b>
<b>  @x_min = min(@x_min, $x);</b>
<b>  @x_max = max(@x_max, $x);</b>
<b>  end{emitf @x_min, @x_max}</b>
<b>' data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_min 0.00004509679127584487
x_max 0.999952670371898
</pre>

## Keyed min/max without/with oosvars

<pre class="pre-highlight-in-pair">
<b>mlr --opprint stats1 -a min,max -f x -g a data/medium</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   x_min                  x_max
pan 0.00020390740306253097 0.9994029107062516
eks 0.0006917972627396018  0.9988110946859143
wye 0.0001874794831505655  0.9998228522652893
zee 0.0005486114815762555  0.9994904324789629
hat 0.00004509679127584487 0.999952670371898
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --from data/medium put -q '</b>
<b>  @min[$a] = min(@min[$a], $x);</b>
<b>  @max[$a] = max(@max[$a], $x);</b>
<b>  end{</b>
<b>    emit (@min, @max), "a";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   min                    max
pan 0.00020390740306253097 0.9994029107062516
eks 0.0006917972627396018  0.9988110946859143
wye 0.0001874794831505655  0.9998228522652893
zee 0.0005486114815762555  0.9994904324789629
hat 0.00004509679127584487 0.999952670371898
</pre>

## Delta without/with oosvars

<pre class="pre-highlight-in-pair">
<b>mlr --opprint step -a delta -f x data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x                   y                   x_delta
pan pan 1 0.3467901443380824  0.7268028627434533  0
eks pan 2 0.7586799647899636  0.5221511083334797  0.41188982045188116
wye wye 3 0.20460330576630303 0.33831852551664776 -0.5540766590236605
eks wye 4 0.38139939387114097 0.13418874328430463 0.17679608810483793
wye pan 5 0.5732889198020006  0.8636244699032729  0.19188952593085962
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '</b>
<b>  $x_delta = is_present(@last) ? $x - @last : 0;</b>
<b>  @last = $x</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x                   y                   x_delta
pan pan 1 0.3467901443380824  0.7268028627434533  0
eks pan 2 0.7586799647899636  0.5221511083334797  0.41188982045188116
wye wye 3 0.20460330576630303 0.33831852551664776 -0.5540766590236605
eks wye 4 0.38139939387114097 0.13418874328430463 0.17679608810483793
wye pan 5 0.5732889198020006  0.8636244699032729  0.19188952593085962
</pre>

## Keyed delta without/with oosvars

<pre class="pre-highlight-in-pair">
<b>mlr --opprint step -a delta -f x -g a data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x                   y                   x_delta
pan pan 1 0.3467901443380824  0.7268028627434533  0
eks pan 2 0.7586799647899636  0.5221511083334797  0
wye wye 3 0.20460330576630303 0.33831852551664776 0
eks wye 4 0.38139939387114097 0.13418874328430463 -0.3772805709188226
wye pan 5 0.5732889198020006  0.8636244699032729  0.36868561403569755
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '</b>
<b>  $x_delta = is_present(@last[$a]) ? $x - @last[$a] : 0;</b>
<b>  @last[$a]=$x</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x                   y                   x_delta
pan pan 1 0.3467901443380824  0.7268028627434533  0
eks pan 2 0.7586799647899636  0.5221511083334797  0
wye wye 3 0.20460330576630303 0.33831852551664776 0
eks wye 4 0.38139939387114097 0.13418874328430463 -0.3772805709188226
wye pan 5 0.5732889198020006  0.8636244699032729  0.36868561403569755
</pre>

## Exponentially weighted moving averages without/with oosvars

<pre class="pre-highlight-in-pair">
<b>mlr --opprint step -a ewma -d 0.1 -f x data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x                   y                   x_ewma_0.1
pan pan 1 0.3467901443380824  0.7268028627434533  0.3467901443380824
eks pan 2 0.7586799647899636  0.5221511083334797  0.3879791263832706
wye wye 3 0.20460330576630303 0.33831852551664776 0.36964154432157387
eks wye 4 0.38139939387114097 0.13418874328430463 0.37081732927653055
wye pan 5 0.5732889198020006  0.8636244699032729  0.3910644883290776
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '</b>
<b>  begin{ @a=0.1 };</b>
<b>  $e = NR==1 ? $x : @a * $x + (1 - @a) * @e;</b>
<b>  @e=$e</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x                   y                   e
pan pan 1 0.3467901443380824  0.7268028627434533  0.3467901443380824
eks pan 2 0.7586799647899636  0.5221511083334797  0.3879791263832706
wye wye 3 0.20460330576630303 0.33831852551664776 0.36964154432157387
eks wye 4 0.38139939387114097 0.13418874328430463 0.37081732927653055
wye pan 5 0.5732889198020006  0.8636244699032729  0.3910644883290776
</pre>
