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
# Null/empty/absent data

One of Miller's key features is its support for **heterogeneous** data.  For example, take `mlr sort`: if you try to sort on field `hostname` when not all records in the data stream *have* a field named `hostname`, it is not an error (although you could pre-filter the data stream using `mlr having-fields --at-least hostname then sort ...`).  Rather, records lacking one or more sort keys are simply output contiguously by `mlr sort`.

## The three types

Miller has three kinds of null data:

* **Empty (key present, value empty)**: a field name is present in a record (or in an out-of-stream variable) with empty value: e.g. `x=,y=2` in the data input stream, or assignment `$x=""` or `@x=""` in `mlr put`.

* **Absent (key not present)**: a field name is not present, e.g. input record is `x=1,y=2` and a `put` or `filter` expression refers to `$z`. Or, reading an out-of-stream variable which hasn't been assigned a value yet, e.g.  `mlr put -q '@sum += $x; end{emit @sum}'` or `mlr put -q '@sum[$a][$b] += $x; end{emit @sum, "a", "b"}'`.

* **JSON null**: The main purpose of this is to support reading the `null` type in JSON files. The [Miller programming language](miller-programming-language.md) has a `null` keyword as well, so you can also write the null type using `$x = null`. Additionally, though, when you write past the end of an array, leaving gaps -- e.g. writing `a[12]` when the array `a` has length 10 -- JSON-null is used to fill the gaps. See also the [arrays page](reference-main-arrays.md#auto-extend-and-null-gaps).

You can test these programmatically using the functions `is_empty`/`is_not_empty`, `is_absent`/`is_present`, and `is_null`/`is_not_null`. For the last pair, note that null means either empty or absent. Here is a full list of such functions:

<pre class="pre-highlight-in-pair">
<b>mlr -f | grep is_</b>
</pre>
<pre class="pre-non-highlight-in-pair">
is_absent
is_array
is_bool
is_boolean
is_empty
is_empty_map
is_error
is_float
is_int
is_map
is_nonempty_map
is_not_array
is_not_empty
is_not_map
is_not_null
is_null
is_numeric
is_present
is_string
</pre>

## Rules for null-handling

* Records with one or more empty sort-field values sort after records with all sort-field values present:

<pre class="pre-highlight-in-pair">
<b>mlr cat data/sort-null.dat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=3,b=2
a=1,b=8
a=,b=4
x=9,b=10
a=5,b=7
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr sort -n a data/sort-null.dat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=1,b=8
a=3,b=2
a=5,b=7
a=,b=4
x=9,b=10
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr sort -nr a data/sort-null.dat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=,b=4
a=5,b=7
a=3,b=2
a=1,b=8
x=9,b=10
</pre>

* Functions/operators which have one or more *empty* arguments produce empty output: e.g.

<pre class="pre-highlight-in-pair">
<b>echo 'x=2,y=3' | mlr put '$a=$x+$y'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=2,y=3,a=5
</pre>

<pre class="pre-highlight-in-pair">
<b>echo 'x=,y=3' | mlr put '$a=$x+$y'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=,y=3,a=
</pre>

<pre class="pre-highlight-in-pair">
<b>echo 'x=,y=3' | mlr put '$a=log($x);$b=log($y)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=,y=3,a=,b=1.0986122886681096
</pre>

with the exception that the `min` and `max` functions are special: if one argument is non-null, it wins:

<pre class="pre-highlight-in-pair">
<b>echo 'x=,y=3' | mlr put '$a=min($x,$y);$b=max($x,$y)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=,y=3,a=3,b=
</pre>

* Functions of *absent* variables (e.g. `mlr put '$y = log10($nonesuch)'`) evaluate to absent, and arithmetic/bitwise/boolean operators with both operands being absent evaluate to absent. Arithmetic operators with one absent operand return the other operand. More specifically, absent values act like zero for addition/subtraction, and one for multiplication: Furthermore, **any expression which evaluates to absent is not stored in the left-hand side of an assignment statement**:

<pre class="pre-highlight-in-pair">
<b>echo 'x=2,y=3' | mlr put '$a=$u+$v; $b=$u+$y; $c=$x+$y'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=2,y=3,b=3,c=5
</pre>

<pre class="pre-highlight-in-pair">
<b>echo 'x=2,y=3' | mlr put '$a=min($x,$v);$b=max($u,$y);$c=min($u,$v)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=2,y=3,a=2,b=3
</pre>

* Likewise, for assignment to maps, **absent-valued keys or values result in a skipped assignment**.

The reasoning is as follows:

* Empty values are explicit in the data so they should explicitly affect accumulations: `mlr put '@sum += $x'` should accumulate numeric `x` values into the sum but an empty `x`, when encountered in the input data stream, should make the sum non-numeric. To work around this you can use the `is_not_null` function as follows: `mlr put 'is_not_null($x) { @sum += $x }'`

* Absent stream-record values should not break accumulations, since Miller by design handles heterogeneous data: the running `@sum` in `mlr put '@sum += $x'` should not be invalidated for records which have no `x`.

* Absent out-of-stream-variable values are precisely what allow you to write `mlr put '@sum += $x'`. Otherwise you would have to write `mlr put 'begin{@sum = 0}; @sum += $x'` -- which is tolerable -- but for `mlr put 'begin{...}; @sum[$a][$b] += $x'` you'd have to pre-initialize `@sum` for all values of `$a` and `$b` in your input data stream, which is intolerable.

* The penalty for the absent feature is that misspelled variables can be hard to find: e.g. in `mlr put 'begin{@sumx = 10}; ...; update @sumx somehow per-record; ...; end {@something = @sum * 2}'` the accumulator is spelt `@sumx` in the begin-block but `@sum` in the end-block, where since it is absent, `@sum*2` evaluates to 2. See also the section on [DSL errors and transparency](reference-dsl-errors.md).

## Absent-test functions

Since absent plus absent is absent (and likewise for other operators), accumulations such as `@sum += $x` work correctly on heterogeneous data, as do within-record formulas if both operands are absent. If one operand is present, you may get behavior you don't desire.  To work around this -- namely, to set an output field only for records which have all the inputs present -- you can use a pattern-action block with `is_present`:

<pre class="pre-highlight-in-pair">
<b>mlr cat data/het.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
resource=/path/to/file,loadsec=0.45,ok=true
record_count=100,resource=/path/to/file
resource=/path/to/second/file,loadsec=0.32,ok=true
record_count=150,resource=/path/to/second/file
resource=/some/other/path,loadsec=0.97,ok=false
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put 'is_present($loadsec) { $loadmillis = $loadsec * 1000 }' data/het.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
resource=/path/to/file,loadsec=0.45,ok=true,loadmillis=450
record_count=100,resource=/path/to/file
resource=/path/to/second/file,loadsec=0.32,ok=true,loadmillis=320
record_count=150,resource=/path/to/second/file
resource=/some/other/path,loadsec=0.97,ok=false,loadmillis=970
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$loadmillis = (is_present($loadsec) ? $loadsec : 0.0) * 1000' data/het.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
resource=/path/to/file,loadsec=0.45,ok=true,loadmillis=450
record_count=100,resource=/path/to/file,loadmillis=0
resource=/path/to/second/file,loadsec=0.32,ok=true,loadmillis=320
record_count=150,resource=/path/to/second/file,loadmillis=0
resource=/some/other/path,loadsec=0.97,ok=false,loadmillis=970
</pre>

## Arithmetic rules

If you're interested in a formal description of how empty and absent fields participate in arithmetic, here's a table for plus (other arithmetic/boolean/bitwise operators are similar):

<pre class="pre-highlight-in-pair">
<b>mlr help type-arithmetic-info</b>
</pre>
<pre class="pre-non-highlight-in-pair">
(+)        | 1          2.5        (absent)   (error)   
------     + ------     ------     ------     ------    
1          | 2          3.5        1          (error)   
2.5        | 3.5        5          2.5        (error)   
(absent)   | 1          2.5        (absent)   (error)   
(error)    | (error)    (error)    (error)    (error)   
</pre>
