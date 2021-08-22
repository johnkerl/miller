<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# DSL overview

## Verbs compared to DSL

Here's comparison of verbs and `put`/`filter` DSL expressions:

Example:

<pre class="pre-highlight-in-pair">
<b>mlr stats1 -a sum -f x -g a data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,x_sum=0.3467901443380824
a=eks,x_sum=1.1400793586611044
a=wye,x_sum=0.7778922255683036
</pre>

* Verbs are coded in Go
* They run a bit faster
* They take fewer keystrokes
* There is less to learn
* Their customization is limited to each verb's options

Example:

<pre class="pre-highlight-in-pair">
<b>mlr  put -q '@x_sum[$a] += $x; end{emit @x_sum, "a"}' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,x_sum=0.3467901443380824
a=eks,x_sum=1.1400793586611044
a=wye,x_sum=0.7778922255683036
</pre>

* You get to write your own DSL expressions
* They run a bit slower
* They take more keystrokes
* There is more to learn
* They are highly customizable

Please see [Verbs Reference](reference-verbs.md) for information on verbs other than `put` and `filter`.

## Implicit loop over records for main statements

The most important point about the Miller DSL is that it is designed for _streaming operation over records_.

DSL statements include:

* `func` and `subr` for user-defined functions and subroutines, which we'll look at later in the [separate page about them](reference-dsl-user-defined-functions.md);
* `begin` and `end` blocks, for statements you want to run before the first record, or after the last one;
* everything else, which collectively are called _main statements_.

The feature of _streaming operation over records_ is implemented by the main
statements getting invoked once per record. You don't explicitly loop over
records, as you would in some dataframes contexts; rather, _Miller loops over
records for you_, and it lets you specify what to do on each record: you write
the body of the loop.

(You can, if you like, use the per-record statements to grow a list of records,
then loop over them all in an `end` block. This is described in the page on
[operating over all records](operating-over-all-records.md)).

To see this in action, let's take a look at the [data/short.csv](./data/short.csv) file:

<pre class="pre-highlight-in-pair">
<b>cat data/short.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
word,value
apple,37
ball,28
cat,54
</pre>

There are three records in this file, with `word=apple`, `word=ball`, and
`word=cat`, respectively. Let's print something in a `begin` statement, add a
field in a main statement, and print something else in an `end` statement:

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data/short.csv put '</b>
<b>  begin {</b>
<b>    print "begin";</b>
<b>  }</b>
<b>  $nr = NR;</b>
<b>  end {</b>
<b>    print "end";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
begin
word,value,nr
apple,37,1
ball,28,2
cat,54,3
end
</pre>

The `print` statements for `begin` and `end` went out before the first record
was seen and after the last was seen; the field-creation statement `$nr = NR`
was invoked three times, once for each record. We didn't explicitly loop over
records, since Miller was already looping over records, and invoked our main
statement on each loop iteration.

For almost all simple uses of the Miller programming language, this implicit
looping over records is probably all you will need. (For more involved cases you
can see the pages on [operating over all records](operating-on-all-records.md),
[out-of-stream variables](reference-dsl-variables.md#out-of-stream-variables),
and [two-pass algorithms](two-pass-algorithms.md).)

## Essential use: record-selection and record-updating

The essential usages of `mlr filter` and `mlr put` are for record-selection and
record-updating expressions, respectively. For example, given the following
input data:

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

you might retain only the records whose `a` field has value `eks`:

<pre class="pre-highlight-in-pair">
<b>mlr filter '$a == "eks"' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797
a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463
</pre>

or you might add a new field which is a function of existing fields:

<pre class="pre-highlight-in-pair">
<b>mlr put '$ab = $a . "_" . $b ' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.3467901443380824,y=0.7268028627434533,ab=pan_pan
a=eks,b=pan,i=2,x=0.7586799647899636,y=0.5221511083334797,ab=eks_pan
a=wye,b=wye,i=3,x=0.20460330576630303,y=0.33831852551664776,ab=wye_wye
a=eks,b=wye,i=4,x=0.38139939387114097,y=0.13418874328430463,ab=eks_wye
a=wye,b=pan,i=5,x=0.5732889198020006,y=0.8636244699032729,ab=wye_pan
</pre>

## Differences between put and filter

The two verbs `mlr filter` and `mlr put` are essentially the same. The only differences are:

* Expressions sent to `mlr filter` should contain a boolean expression, which is the filtering criterion. (If not, all records pass through.)

* `mlr filter` expressions may not reference the `filter` keyword within them.

## Location of boolean expression for filter

You can define and invoke functions and subroutines to help produce the bare-boolean statement, and record fields may be assigned in the statements before or after the bare-boolean statement. For example:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv filter '</b>
<b>  # Bare-boolean filter expression: only records matching this pass through:</b>
<b>  $quantity >= 70;</b>
<b>  # For records that do pass through, set these:</b>
<b>  if ($rate > 8) {</b>
<b>    $description = "high rate";</b>
<b>  } else {</b>
<b>    $description = "low rate";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   description
red    square   true  2  15    79.2778  0.0130 low rate
red    square   false 4  48    77.5542  7.4670 low rate
purple triangle false 5  51    81.2290  8.5910 high rate
red    square   false 6  64    77.1991  9.5310 high rate
purple triangle false 7  65    80.1405  5.8240 low rate
purple square   false 10 91    72.3735  8.2430 high rate
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv filter '</b>
<b>  # Bare-boolean filter expression: only records matching this pass through:</b>
<b>  $shape =~ "^(...)(...)$";</b>
<b>  # For records that do pass through, capture the first "(...)" into $left and</b>
<b>  # the second "(...)" into $right</b>
<b>  $left = "\1";</b>
<b>  $right = "\2";</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape  flag  k  index quantity rate   left right
red    square true  2  15    79.2778  0.0130 squ  are
red    circle true  3  16    13.8103  2.9010 cir  cle
red    square false 4  48    77.5542  7.4670 squ  are
red    square false 6  64    77.1991  9.5310 squ  are
yellow circle true  8  73    63.9785  4.2370 cir  cle
yellow circle true  9  87    63.5058  8.3350 cir  cle
purple square false 10 91    72.3735  8.2430 squ  are
</pre>


There are more details and more choices, of course, as detailed in the following sections.

