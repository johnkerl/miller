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
# DSL control structures

## Pattern-action blocks

These are reminiscent of `awk` syntax.  They can be used to allow assignments to be done only when appropriate -- e.g. for math-function domain restrictions, regex-matching, and so on:

<pre class="pre-highlight-in-pair">
<b>mlr cat data/put-gating-example-1.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=-1
x=0
x=1
x=2
x=3
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$x > 0.0 { $y = log10($x); $z = sqrt($y) }' data/put-gating-example-1.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=-1
x=0
x=1,y=0,z=0
x=2,y=0.3010299956639812,z=0.5486620049392715
x=3,y=0.4771212547196624,z=0.6907396432228734
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr cat data/put-gating-example-2.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=abc_123
a=some other name
a=xyz_789
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '</b>
<b>  $a =~ "([a-z]+)_([0-9]+)" {</b>
<b>    $b = "left_\1"; $c = "right_\2"</b>
<b>  }' \</b>
<b>  data/put-gating-example-2.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=abc_123,b=left_abc,c=right_123
a=some other name
a=xyz_789,b=left_xyz,c=right_789
</pre>

This produces heteregenous output which Miller, of course, has no problems with (see [Record Heterogeneity](record-heterogeneity.md)).  But if you want homogeneous output, the curly braces can be replaced with a semicolon between the expression and the body statements.  This causes `put` to evaluate the boolean expression (along with any side effects, namely, regex-captures `\1`, `\2`, etc.) but doesn't use it as a criterion for whether subsequent assignments should be executed. Instead, subsequent assignments are done unconditionally:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '</b>
<b>  $a =~ "([a-z]+)_([0-9]+)";</b>
<b>  $b = "left_\1";</b>
<b>  $c = "right_\2"</b>
<b>' data/put-gating-example-2.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a               b        c
abc_123         left_abc right_123
some other name left_    right_
xyz_789         left_xyz right_789
</pre>

Note that pattern-action blocks are just a syntactic variation of if-statements. The following do the same thing:

<pre class="pre-non-highlight-non-pair">
  boolean_condition {
    body
  }
</pre>

<pre class="pre-non-highlight-non-pair">
  if (boolean_condition) {
    body
  }
</pre>

## If-statements

These are again reminiscent of `awk`. Pattern-action blocks are a special case of `if` with no `elif` or `else` blocks, no `if` keyword, and parentheses optional around the boolean expression:

<pre class="pre-highlight-non-pair">
<b>mlr put 'NR == 4 {$foo = "bar"}'</b>
</pre>

<pre class="pre-highlight-non-pair">
<b>mlr put 'if (NR == 4) {$foo = "bar"}'</b>
</pre>

Compound statements use `elif` (rather than `elsif` or `else if`):

<pre class="pre-highlight-non-pair">
<b>mlr put '</b>
<b>  if (NR == 2) {</b>
<b>    ...</b>
<b>  } elif (NR ==4) {</b>
<b>    ...</b>
<b>  } elif (NR ==6) {</b>
<b>    ...</b>
<b>  } else {</b>
<b>    ...</b>
<b>  }</b>
<b>'</b>
</pre>

## While and do-while loops

Miller's `while` and `do-while` are unsurprising in comparison to various languages, as are `break` and `continue`:

<pre class="pre-highlight-in-pair">
<b>echo x=1,y=2 | mlr put '</b>
<b>  while (NF < 10) {</b>
<b>    $[NF+1] = ""</b>
<b>  }</b>
<b>  $foo = "bar"</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=1,y=2,3=,4=,5=,6=,7=,8=,9=,10=,foo=bar
</pre>

<pre class="pre-highlight-in-pair">
<b>echo x=1,y=2 | mlr put '</b>
<b>  do {</b>
<b>    $[NF+1] = "";</b>
<b>    if (NF == 5) {</b>
<b>      break</b>
<b>    }</b>
<b>  } while (NF < 10);</b>
<b>  $foo = "bar"</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=1,y=2,3=,4=,5=,foo=bar
</pre>

A `break` or `continue` within nested conditional blocks or if-statements will,
of course, propagate to the innermost loop enclosing them, if any. A `break` or
`continue` outside a loop is a syntax error that will be flagged as soon as the
expression is parsed, before any input records are ingested.

The existence of `while`, `do-while`, and `for` loops in Miller's DSL means
that you can create infinite-loop scenarios inadvertently.  In particular,
please recall that DSL statements are executed once if in `begin` or `end`
blocks, and once *per record* otherwise. For example, **while (NR < 10) will
never terminate**. The [`NR`
variable](reference-dsl-variables.md#built-in-variables) is only incremented
between records, and each DSL expression is invoked once per record: so, once
for `NR=1`, once for `NR=2`, etc.

If you do want to loop over records, see [Operating on all
records](operating-on-all-records.md) for some options.

## For-loops

While Miller's `while` and `do-while` statements are much as in many other languages, `for` loops are more idiosyncratic to Miller. They are loops over key-value pairs, whether in stream records, out-of-stream variables, local variables, or map-literals: more reminiscent of `foreach`, as in (for example) PHP. There are **for-loops over map keys** and **for-loops over key-value tuples**.  Additionally, Miller has a **C-style triple-for loop** with initialize, test, and update statements. Each is described below.

As with `while` and `do-while`, a `break` or `continue` within nested control structures will propagate to the innermost loop enclosing them, if any, and a `break` or `continue` outside a loop is a syntax error that will be flagged as soon as the expression is parsed, before any input records are ingested.

### Single-variable for-loops

For [maps](reference-main-maps.md), the single variable is always bound to the *key* of key-value pairs:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small put -q '</b>
<b>  print "NR = ".NR;</b>
<b>  for (e in $*) {</b>
<b>    print "  key:", e, "value:", $[e];</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
NR = 1
  key: a value: pan
  key: b value: pan
  key: i value: 1
  key: x value: 0.346791
  key: y value: 0.726802
NR = 2
  key: a value: eks
  key: b value: pan
  key: i value: 2
  key: x value: 0.758679
  key: y value: 0.522151
NR = 3
  key: a value: wye
  key: b value: wye
  key: i value: 3
  key: x value: 0.204603
  key: y value: 0.338318
NR = 4
  key: a value: eks
  key: b value: wye
  key: i value: 4
  key: x value: 0.381399
  key: y value: 0.134188
NR = 5
  key: a value: wye
  key: b value: pan
  key: i value: 5
  key: x value: 0.573288
  key: y value: 0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put -q '</b>
<b>  end {</b>
<b>    o = {"a":1, "b":{"c":3}};</b>
<b>    for (e in o) {</b>
<b>      print "key:", e, "valuetype:", typeof(o[e]);</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key: a valuetype: int
key: b valuetype: map
</pre>

Note that the value corresponding to a given key may be gotten as through a **computed field name** using square brackets as in `$[e]` for stream records, or by indexing the looped-over variable using square brackets.

For [arrays](reference-main-arrays.md), the single variable is always bound to the *value* (not the array index):

<pre class="pre-highlight-in-pair">
<b>mlr -n put -q '</b>
<b>  end {</b>
<b>    o = [10, "20", {}, "four", true];</b>
<b>    for (e in o) {</b>
<b>      print "value:", e, "valuetype:", typeof(e);</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
value: 10 valuetype: int
value: 20 valuetype: string
value: {} valuetype: map
value: four valuetype: string
value: true valuetype: bool
</pre>

### Key-value for-loops

For [maps](reference-main-maps.md), the first loop variable is the key and the
second is the value; for [arrays](reference-main-arrays.md), the first loop
variable is the (1-up) array index and the second is the value.

Single-level keys may be gotten at using either `for(k,v)` or `for((k),v)`; multi-level keys may be gotten at using `for((k1,k2,k3),v)` and so on.  The `v` variable will be bound to to a scalar value (non-array/non-map) if the map stops at that level, or to a map-valued or array-valued variable if the map goes deeper. If the map isn't deep enough then the loop body won't be executed.

<pre class="pre-highlight-in-pair">
<b>cat data/for-srec-example.tbl</b>
</pre>
<pre class="pre-non-highlight-in-pair">
label1 label2 f1  f2  f3
blue   green  100 240 350
red    green  120 11  195
yellow blue   140 0   240
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --pprint --from data/for-srec-example.tbl put '</b>
<b>  $sum1 = $f1 + $f2 + $f3;</b>
<b>  $sum2 = 0;</b>
<b>  $sum3 = 0;</b>
<b>  for (key, value in $*) {</b>
<b>    if (key =~ "^f[0-9]+") {</b>
<b>      $sum2 += value;</b>
<b>      $sum3 += $[key];</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
label1 label2 f1  f2  f3  sum1 sum2 sum3
blue   green  100 240 350 690  690  690
red    green  120 11  195 326  326  326
yellow blue   140 0   240 380  380  380
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put 'for (k,v in $*) { $[k."_type"] = typeof(v) }'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        a_type b_type i_type x_type y_type
pan pan 1 0.346791 0.726802 string string int    float  float
eks pan 2 0.758679 0.522151 string string int    float  float
wye wye 3 0.204603 0.338318 string string int    float  float
eks wye 4 0.381399 0.134188 string string int    float  float
wye pan 5 0.573288 0.863624 string string int    float  float
</pre>

Note that the value of the current field in the for-loop can be gotten either using the bound variable `value`, or through a **computed field name** using square brackets as in `$[key]`.

Important note: to avoid inconsistent looping behavior in case you're setting new fields (and/or unsetting existing ones) while looping over the record, **Miller makes a copy of the record before the loop: loop variables are bound from the copy and all other reads/writes involve the record itself**:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put '</b>
<b>  $sum1 = 0;</b>
<b>  $sum2 = 0;</b>
<b>  for (k,v in $*) {</b>
<b>    if (is_numeric(v)) {</b>
<b>      $sum1 +=v;</b>
<b>      $sum2 += $[k];</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        sum1               sum2
pan pan 1 0.346791 0.726802 2.073593           8.294372
eks pan 2 0.758679 0.522151 3.28083            13.12332
wye wye 3 0.204603 0.338318 3.542921           14.171684
eks wye 4 0.381399 0.134188 4.515587           18.062348
wye pan 5 0.573288 0.863624 6.4369119999999995 25.747647999999998
</pre>

It can be confusing to modify the stream record while iterating over a copy of it, so instead you might find it simpler to use a local variable in the loop and only update the stream record after the loop:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put '</b>
<b>  sum = 0;</b>
<b>  for (k,v in $*) {</b>
<b>    if (is_numeric(v)) {</b>
<b>      sum += $[k];</b>
<b>    }</b>
<b>  }</b>
<b>  $sum = sum</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        sum
pan pan 1 0.346791 0.726802 2.073593
eks pan 2 0.758679 0.522151 3.28083
wye wye 3 0.204603 0.338318 3.542921
eks wye 4 0.381399 0.134188 4.515587
wye pan 5 0.573288 0.863624 6.4369119999999995
</pre>

You can also start iterating on sub-maps of an out-of-stream or local variable; you can loop over nested keys; you can loop over all out-of-stream variables.  The bound variables are bound to a copy of the sub-map as it was before the loop started.  The sub-map is specified by square-bracketed indices after `in`, and additional deeper indices are bound to loop key-variables. The terminal values are bound to the loop value-variable whenever the keys are not too shallow. The value-variable may refer to a terminal (string, number) or it may be map-valued if the map goes deeper. Example indexing is as follows:

<pre class="pre-non-highlight-non-pair">
# Parentheses are optional for single key:
for (k1,           v in @a["b"]["c"]) { ... }
for ((k1),         v in @a["b"]["c"]) { ... }
# Parentheses are required for multiple keys:
for ((k1, k2),     v in @a["b"]["c"]) { ... } # Loop over subhashmap of a variable
for ((k1, k2, k3), v in @a["b"]["c"]) { ... } # Ditto
for ((k1, k2, k3), v in @a { ... }            # Loop over variable starting from basename
for ((k1, k2, k3), v in @* { ... }            # Loop over all variables (k1 is bound to basename)
</pre>

That's confusing in the abstract, so a concrete example is in order. Suppose the out-of-stream variable `@myvar` is populated as follows:

<pre class="pre-highlight-in-pair">
<b>mlr -n put --jknquoteint -q '</b>
<b>  begin {</b>
<b>    @myvar = {</b>
<b>      1: 2,</b>
<b>      3: { 4 : 5 },</b>
<b>      6: { 7: { 8: 9 } }</b>
<b>    }</b>
<b>  }</b>
<b>  end { dump }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "myvar": {
    "1": 2,
    "3": {
      "4": 5
    },
    "6": {
      "7": {
        "8": 9
      }
    }
  }
}
</pre>

Then we can get at various values as follows:

<pre class="pre-highlight-in-pair">
<b>mlr -n put --jknquoteint -q '</b>
<b>  begin {</b>
<b>    @myvar = {</b>
<b>      1: 2,</b>
<b>      3: { 4 : 5 },</b>
<b>      6: { 7: { 8: 9 } }</b>
<b>    }</b>
<b>  }</b>
<b>  end {</b>
<b>    for (k, v in @myvar) {</b>
<b>      print</b>
<b>        "key=" . k .</b>
<b>        ",valuetype=" . typeof(v);</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key=1,valuetype=int
key=3,valuetype=map
key=6,valuetype=map
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put --jknquoteint -q '</b>
<b>  begin {</b>
<b>    @myvar = {</b>
<b>      1: 2,</b>
<b>      3: { 4 : 5 },</b>
<b>      6: { 7: { 8: 9 } }</b>
<b>    }</b>
<b>  }</b>
<b>  end {</b>
<b>    for ((k1, k2), v in @myvar) {</b>
<b>      print</b>
<b>        "key1=" . k1 .</b>
<b>        ",key2=" . k2 .</b>
<b>        ",valuetype=" . typeof(v);</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key1=3,key2=4,valuetype=int
key1=6,key2=7,valuetype=map
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put --jknquoteint -q '</b>
<b>  begin {</b>
<b>    @myvar = {</b>
<b>      1: 2,</b>
<b>      3: { 4 : 5 },</b>
<b>      6: { 7: { 8: 9 } }</b>
<b>    }</b>
<b>  }</b>
<b>  end {</b>
<b>    for ((k1, k2), v in @myvar[6]) {</b>
<b>      print</b>
<b>        "key1=" . k1 .</b>
<b>        ",key2=" . k2 .</b>
<b>        ",valuetype=" . typeof(v);</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key1=7,key2=8,valuetype=int
</pre>

### C-style triple-for loops

These are supported as follows:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put '</b>
<b>  num suma = 0;</b>
<b>  for (a = 1; a <= NR; a += 1) {</b>
<b>    suma += a;</b>
<b>  }</b>
<b>  $suma = suma;</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        suma
pan pan 1 0.346791 0.726802 1
eks pan 2 0.758679 0.522151 3
wye wye 3 0.204603 0.338318 6
eks wye 4 0.381399 0.134188 10
wye pan 5 0.573288 0.863624 15
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small --opprint put '</b>
<b>  num suma = 0;</b>
<b>  num sumb = 0;</b>
<b>  for (num a = 1, num b = 1; a <= NR; a += 1, b *= 2) {</b>
<b>    suma += a;</b>
<b>    sumb += b;</b>
<b>  }</b>
<b>  $suma = suma;</b>
<b>  $sumb = sumb;</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        suma sumb
pan pan 1 0.346791 0.726802 1    1
eks pan 2 0.758679 0.522151 3    3
wye wye 3 0.204603 0.338318 6    7
eks wye 4 0.381399 0.134188 10   15
wye pan 5 0.573288 0.863624 15   31
</pre>

Notes:

* In `for (start; continuation; update) { body }`, the start, continuation, and update statements may be empty, single statements, or multiple comma-separated statements. If the continuation is empty (e.g. `for(i=1;;i+=1)`) it defaults to true.

* In particular, you may use `$`-variables and/or `@`-variables in the start, continuation, and/or update steps (as well as the body, of course).

* The typedecls such as `int` or `num` are optional.  If a typedecl is provided (for a local variable), it binds a variable scoped to the for-loop regardless of whether a same-name variable is present in outer scope. If a typedecl is not provided, then the variable is scoped to the for-loop if no same-name variable is present in outer scope, or if a same-name variable is present in outer scope then it is modified.

* Miller has no `++` or `--` operators.

* As with all `for`/`if`/`while` statements in Miller, the curly braces are required even if the body is a single statement, or empty.

## Begin/end blocks

Miller supports an `awk`-like `begin/end` syntax.  The statements in the `begin` block are executed before any input records are read; the statements in the `end` block are executed after the last input record is read.  (If you want to execute some statement at the start of each file, not at the start of the first file as with `begin`, you might use a pattern/action block of the form `FNR == 1 { ... }`.) All statements outside of `begin` or `end` are, of course, executed on every input record. Semicolons separate statements inside or outside of begin/end blocks; semicolons are required between begin/end block bodies and any subsequent statement.  For example:

<pre class="pre-highlight-in-pair">
<b>mlr put '</b>
<b>  begin { @sum = 0 };</b>
<b>  @x_sum += $x;</b>
<b>  end { emit @x_sum }</b>
<b>' ./data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
x_sum=2.26476
</pre>

Since uninitialized out-of-stream variables default to 0 for addition/subtraction and 1 for multiplication when they appear on expression right-hand sides (not quite as in `awk`, where they'd default to 0 either way), the above can be written more succinctly as

<pre class="pre-highlight-in-pair">
<b>mlr put '</b>
<b>  @x_sum += $x;</b>
<b>  end { emit @x_sum }</b>
<b>' ./data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
x_sum=2.26476
</pre>

The **put -q** option suppresses printing of each output record, with only `emit` statements being output. So to get only summary outputs, you could write

<pre class="pre-highlight-in-pair">
<b>mlr put -q '</b>
<b>  @x_sum += $x;</b>
<b>  end { emit @x_sum }</b>
<b>' ./data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_sum=2.26476
</pre>

We can do similarly with multiple out-of-stream variables:

<pre class="pre-highlight-in-pair">
<b>mlr put -q '</b>
<b>  @x_count += 1;</b>
<b>  @x_sum += $x;</b>
<b>  end {</b>
<b>    emit @x_count;</b>
<b>    emit @x_sum;</b>
<b>  }</b>
<b>' ./data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_count=5
x_sum=2.26476
</pre>

This is of course (see also [here](reference-dsl.md#verbs-compared-to-dsl)) not much different than

<pre class="pre-highlight-in-pair">
<b>mlr stats1 -a count,sum -f x ./data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x_count=5,x_sum=2.26476
</pre>

Note that it's a syntax error for begin/end blocks to refer to field names (beginning with `$`), since begin/end blocks execute outside the context of input records.

