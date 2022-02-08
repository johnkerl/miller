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
# DSL variables

Miller has the following kinds of variables:

**Fields of stream records**, accessed using the `$` prefix. These refer to fields of the current data-stream record. For example, in `echo x=1,y=2 | mlr put '$z = $x + $y'`, `$x` and `$y` refer to input fields, and `$z` refers to a new, computed output field. In a few contexts, presented below, you can refer to the entire record as `$*`.

**Out-of-stream variables** accessed using the `@` prefix. These refer to data which persist from one record to the next, including in `begin` and `end` blocks (which execute before/after the record stream is consumed, respectively). You use them to remember values across records, such as sums, differences, counters, and so on.  In a few contexts, presented below, you can refer to the entire out-of-stream-variables collection as `@*`.

**Local variables** are limited in scope and extent to the current statements being executed: these include function arguments, bound variables in for loops, and local variables.

**Built-in variables** such as `NF`, `NR`, `FILENAME`, `M_PI`, and `M_E`.  These are all capital letters and are read-only (although some of them change value from one record to another).

**Keywords** are not variables, but since their names are reserved, you cannot use these names for local variables.

## Field names

Names of fields within stream records must be specified using a `$` in [filter and put expressions](reference-dsl.md), even though the dollar signs don't appear in the data stream itself. For integer-indexed data, this looks like `awk`'s `$1,$2,$3`, except that Miller allows non-numeric names such as `$quantity` or `$hostname`.  Likewise, enclose string literals in double quotes in `filter` expressions even though they don't appear in file data.  In particular, `mlr filter '$x=="abc"'` passes through the record `x=abc`.

If field names have **special characters** such as `.` then you can use braces, e.g. `'${field.name}'`.

You may also use a **computed field name** in square brackets, e.g.

<pre class="pre-highlight-non-pair">
<b>echo a=3,b=4 | mlr filter '$["x"] < 0.5'</b>
</pre>

<pre class="pre-highlight-in-pair">
<b>echo s=green,t=blue,a=3,b=4 | mlr put '$[$s."_".$t] = $a * $b'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
s=green,t=blue,a=3,b=4,green_blue=12
</pre>

Notes:

The names of record fields depend on the contents of your input data stream, and their values change from one record to the next as Miller scans through your input data stream.

Their **extent** is limited to the current record; their **scope** is the `filter` or `put` command in which they appear.

These are **read-write**: you can do `$y=2*$x`, `$x=$x+1`, etc.

Records are Miller's output: field names present in the input stream are passed through to output (written to standard output) unless fields are removed with `cut`, or records are excluded with `filter` or `put -q`, etc. Simply assign a value to a field and it will be output.

## Positional field names

Even though Miller's main selling point is name-indexing, sometimes you really want to refer to a field name by its positional index (starting from 1).

Use `$[[3]]` to access the name of field 3.  More generally, any expression evaluating to an integer can go between `$[[` and `]]`.

Then using a computed field name, `$[ $[[3]] ]` is the value in the third field. This has the shorter equivalent notation `$[[[3]]]`.

<pre class="pre-highlight-in-pair">
<b>mlr cat data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$[[3]] = "NEW"' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,NEW=1,x=0.346791,y=0.726802
a=eks,b=pan,NEW=2,x=0.758679,y=0.522151
a=wye,b=wye,NEW=3,x=0.204603,y=0.338318
a=eks,b=wye,NEW=4,x=0.381399,y=0.134188
a=wye,b=pan,NEW=5,x=0.573288,y=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$[[[3]]] = "NEW"' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=NEW,x=0.346791,y=0.726802
a=eks,b=pan,i=NEW,x=0.758679,y=0.522151
a=wye,b=wye,i=NEW,x=0.204603,y=0.338318
a=eks,b=wye,i=NEW,x=0.381399,y=0.134188
a=wye,b=pan,i=NEW,x=0.573288,y=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$NEW = $[[NR]]' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802,NEW=a
a=eks,b=pan,i=2,x=0.758679,y=0.522151,NEW=b
a=wye,b=wye,i=3,x=0.204603,y=0.338318,NEW=i
a=eks,b=wye,i=4,x=0.381399,y=0.134188,NEW=x
a=wye,b=pan,i=5,x=0.573288,y=0.863624,NEW=y
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$NEW = $[[[NR]]]' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802,NEW=pan
a=eks,b=pan,i=2,x=0.758679,y=0.522151,NEW=pan
a=wye,b=wye,i=3,x=0.204603,y=0.338318,NEW=3
a=eks,b=wye,i=4,x=0.381399,y=0.134188,NEW=0.381399
a=wye,b=pan,i=5,x=0.573288,y=0.863624,NEW=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$[[[NR]]] = "NEW"' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=NEW,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=NEW,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=NEW,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=NEW,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=NEW
</pre>

Right-hand side accesses to non-existent fields -- i.e. with index less than 1 or greater than `NF` -- return an absent value. Likewise, left-hand side accesses only refer to fields which already exist. For example, if a field has 5 records then assigning the name or value of the 6th (or 600th) field results in a no-op.

<pre class="pre-highlight-in-pair">
<b>mlr put '$[[6]] = "NEW"' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$[[[6]]] = "NEW"' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=wye,i=4,x=0.381399,y=0.134188
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

## Out-of-stream variables

These are prefixed with an at-sign, e.g. `@sum`.  Furthermore, unlike built-in variables and stream-record fields, they are maintained in an arbitrarily nested map: you can do `@sum += $quantity`, or `@sum[$color] += $quantity`, or `@sum[$color][$shape] += $quantity`. The keys for the multi-level map can be any expression which evaluates to string or integer: e.g.  `@sum[NR] = $a + $b`, `@sum[$a."-".$b] = $x`, etc.

Their names and their values are entirely under your control; they change only when you assign to them.

Just as for field names in stream records, if you want to define out-of-stream variables with **special characters** such as `.` then you can use braces, e.g. `'@{variable.name}["index"]'`.

You may use a **computed key** in square brackets, e.g.

<pre class="pre-highlight-in-pair">
<b>echo s=green,t=blue,a=3,b=4 | mlr put -q '@[$s."_".$t] = $a * $b; emit all'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
green_blue=12
</pre>

Out-of-stream variables are **scoped** to the `put` command in which they appear.  In particular, if you have two or more `put` commands separated by `then`, each put will have its own set of out-of-stream variables:

<pre class="pre-highlight-in-pair">
<b>cat data/a.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=1,b=2,c=3
a=4,b=5,c=6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '@sum += $a; end {emit @sum}' \</b>
<b>  then put 'is_present($a) {$a=10*$a; @sum += $a}; end {emit @sum}' \</b>
<b>  data/a.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=10,b=2,c=3
a=40,b=5,c=6
sum=5
sum=50
</pre>

Out-of-stream variables' **extent** is from the start to the end of the record stream, i.e. every time the `put` or `filter` statement referring to them is executed.

Out-of-stream variables are **read-write**: you can do `$sum=@sum`, `@sum=$sum`, etc.

## Indexed out-of-stream variables

Using an index on the `@count` and `@sum` variables, we get the benefit of the `-g` (group-by) option which `mlr stats1` and various other Miller commands have:

<pre class="pre-highlight-in-pair">
<b>mlr put -q '</b>
<b>  @x_count[$a] += 1;</b>
<b>  @x_sum[$a] += $x;</b>
<b>  end {</b>
<b>    emit @x_count, "a";</b>
<b>    emit @x_sum, "a";</b>
<b>  }</b>
<b>' ./data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,x_count=1
a=eks,x_count=2
a=wye,x_count=2
a=pan,x_sum=0.346791
a=eks,x_sum=1.140078
a=wye,x_sum=0.777891
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr stats1 -a count,sum -f x -g a ./data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,x_count=1,x_sum=0.346791
a=eks,x_count=2,x_sum=1.140078
a=wye,x_count=2,x_sum=0.777891
</pre>

Indices can be arbitrarily deep -- here there are two or more of them:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/medium put -q '</b>
<b>  @x_count[$a][$b] += 1;</b>
<b>  @x_sum[$a][$b] += $x;</b>
<b>  end {</b>
<b>    emit (@x_count, @x_sum), "a", "b";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,x_count=427,x_sum=219.1851288316854
a=pan,b=wye,x_count=395,x_sum=198.43293070748447
a=pan,b=eks,x_count=429,x_sum=216.07522773165525
a=pan,b=hat,x_count=417,x_sum=205.22277621488686
a=pan,b=zee,x_count=413,x_sum=205.09751802331917
a=eks,b=pan,x_count=371,x_sum=179.96303047250723
a=eks,b=wye,x_count=407,x_sum=196.9452860713734
a=eks,b=zee,x_count=357,x_sum=176.8803651584733
a=eks,b=eks,x_count=413,x_sum=215.91609712937984
a=eks,b=hat,x_count=417,x_sum=208.783170520597
a=wye,b=wye,x_count=377,x_sum=185.29584980261419
a=wye,b=pan,x_count=392,x_sum=195.84790012056564
a=wye,b=hat,x_count=426,x_sum=212.0331829346132
a=wye,b=zee,x_count=385,x_sum=194.77404756708714
a=wye,b=eks,x_count=386,x_sum=204.8129608356315
a=zee,b=pan,x_count=389,x_sum=202.21380378504267
a=zee,b=wye,x_count=455,x_sum=233.9913939194868
a=zee,b=eks,x_count=391,x_sum=190.9617780631925
a=zee,b=zee,x_count=403,x_sum=206.64063510417319
a=zee,b=hat,x_count=409,x_sum=191.30000620900935
a=hat,b=wye,x_count=423,x_sum=208.8830097609959
a=hat,b=zee,x_count=385,x_sum=196.3494502965293
a=hat,b=eks,x_count=389,x_sum=189.0067933716193
a=hat,b=hat,x_count=381,x_sum=182.8535323148762
a=hat,b=pan,x_count=363,x_sum=168.5538067327806
</pre>

The idea is that `stats1`, and other Miller verbs, encapsulate frequently-used patterns with a minimum of keystroking (and run a little faster), whereas using out-of-stream variables you have more flexibility and control in what you do.

Begin/end blocks can be mixed with pattern/action blocks. For example:

<pre class="pre-highlight-in-pair">
<b>mlr put '</b>
<b>  begin {</b>
<b>    @num_total = 0;</b>
<b>    @num_positive = 0;</b>
<b>  };</b>
<b>  @num_total += 1;</b>
<b>  $x > 0.0 {</b>
<b>    @num_positive += 1;</b>
<b>    $y = log10($x); $z = sqrt($y)</b>
<b>  };</b>
<b>  end {</b>
<b>    emitf @num_total, @num_positive</b>
<b>  }</b>
<b>' data/put-gating-example-1.dkvp</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=-1
x=0
x=1,y=0,z=0
x=2,y=0.3010299956639812,z=0.5486620049392715
x=3,y=0.4771212547196624,z=0.6907396432228734
num_total=5,num_positive=3
</pre>

## Local variables

Local variables are similar to out-of-stream variables, except that their extent is limited to the expressions in which they appear (and their basenames can't be computed using square brackets). There are three kinds of local variables: **arguments** to functions/subroutines, **variables bound within for-loops**, and **locals** defined within control blocks. They may be untyped using `var`, or typed using `num`, `int`, `float`, `str`, `bool`, and `map`.

For example:

<pre class="pre-highlight-in-pair">
<b># Here I'm using a specified random-number seed so this example always</b>
<b># produces the same output for this web document: in everyday practice we</b>
<b># would leave off the --seed 12345 part.</b>
<b>mlr --seed 12345 seqgen --start 1 --stop 10 then put '</b>
<b>  func f(a, b) {                          # function arguments a and b</b>
<b>      r = 0.0;                            # local r scoped to the function</b>
<b>      for (int i = 0; i < 6; i += 1) {    # local i scoped to the for-loop</b>
<b>          num u = urand();                # local u scoped to the for-loop</b>
<b>          r += u;                         # updates r from the enclosing scope</b>
<b>      }</b>
<b>      r /= 6;</b>
<b>      return a + (b - a) * r;</b>
<b>  }</b>
<b>  num o = f(10, 20);                      # local to the top-level scope</b>
<b>  $o = o;</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
i=1,o=15.952526011537227
i=2,o=12.782237754999116
i=3,o=15.126606630220966
i=4,o=14.794357488895775
i=5,o=15.168665974047421
i=6,o=16.20662783079942
i=7,o=13.966128063060479
i=8,o=13.99248245928659
i=9,o=15.784270485515197
i=10,o=15.37686787628025
</pre>

Things which are completely unsurprising, resembling many other languages:

* Parameter names are bound to their arguments but can be reassigned, e.g. if there is a parameter named `a` then you can reassign the value of `a` to be something else within the function if you like.

* However, you cannot redeclare the *type* of an argument or a local: `var a=1; var a=2` is an error but `var a=1;  a=2` is OK.

* All argument-passing is positional rather than by name; arguments are passed by value, not by reference. (This is also true for map-valued variables: they are not, and cannot be, passed by reference.)

* You can define locals (using `var`, `num`, etc.) at any scope (if-statements, else-statements, while-loops, for-loops, or the top-level scope), and nested scopes will have access (more details on scope in the next section).  If you define a local variable with the same name inside an inner scope, then a new variable is created with the narrower scope.

* If you assign to a local variable for the first time in a scope without declaring it as `var`, `num`, etc. then: if it exists in an outer scope, that outer-scope variable will be updated; if not, it will be defined in the current scope as if `var` had been used. (See also [Type-checking](reference-dsl-variables.md#type-checking) for an example.) I recommend always declaring variables explicitly to make the intended scoping clear.

* Functions and subroutines never have access to locals from their callee (unless passed by value as arguments).

Things which are perhaps surprising compared to other languages:

* Type declarations using `var`, or typed using `num`, `int`, `float`, `str`, and `bool` are not necessary to declare local variables.  Function arguments and variables bound in for-loops over stream records and out-of-stream variables are *implicitly* declared using `var`. (Some examples are shown below.)

* Type-checking is done at assignment time. For example, `float f = 0` is an error (since `0` is an integer), as is `float f = 0.0; f = 1`. For this reason I prefer to use `num` over `float` in most contexts since `num` encompasses integer and floating-point values. More information is at [Type-checking](reference-dsl-variables.md#type-checking).

* Bound variables in for-loops over stream records and out-of-stream variables are implicitly local to that block. E.g. in `for (k, v in $*) { ... }` `for ((k1, k2), v in @*) { ... }` if there are `k`, `v`, etc. in the enclosing scope then those will be masked by the loop-local bound variables in the loop, and moreover the values of the loop-local bound variables are not available after the end of the loop.

* For C-style triple-for loops, if a for-loop variable is defined using `var`, `int`, etc. then it is scoped to that for-loop. E.g. `for (i = 0; i < 10; i += 1) { ... }` and `for (int i = 0; i < 10; i += 1) { ... }`. (This is unsurprising.). If there is no typedecl and an outer-scope variable of that name exists, then it is used. (This is also unsurprising.) But if there is no outer-scope variable of that name, then the variable is scoped to the for-loop only.

The following example demonstrates the scope rules:

<pre class="pre-highlight-in-pair">
<b>cat data/scope-example.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
func f(a) {      # argument is local to the function
  var b = 100;   # local to the function
  c = 100;       # local to the function; does not overwrite outer c
  return a + 1;
}
var a = 10;      # local at top level
var b = 20;      # local at top level
c = 30;          # local at top level; there is no more-outer-scope c
if (NR == 3) {
  var a = 40;    # scoped to the if-statement; doesn't overwrite outer a
  b = 50;        # not scoped to the if-statement; overwrites outer b
  c = 60;        # not scoped to the if-statement; overwrites outer c
  d = 70;        # there is no outer d so a local d is created here

  $inner_a = a;
  $inner_b = b;
  $inner_c = c;
  $inner_d = d;
}
$outer_a = a;
$outer_b = b;
$outer_c = c;
$outer_d = d;    # there is no outer d defined so no assignment happens
</pre>

<pre class="pre-highlight-in-pair">
<b>cat data/scope-example.dat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
n=1,x=123
n=2,x=456
n=3,x=789
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --oxtab --from data/scope-example.dat put -f data/scope-example.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
n       1
x       123
outer_a 10
outer_b 20
outer_c 30

n       2
x       456
outer_a 10
outer_b 20
outer_c 30

n       3
x       789
inner_a 40
inner_b 50
inner_c 60
inner_d 70
outer_a 10
outer_b 50
outer_c 60
</pre>

And this example demonstrates the type-declaration rules:

<pre class="pre-highlight-in-pair">
<b>cat data/type-decl-example.mlr</b>
</pre>
<pre class="pre-non-highlight-in-pair">
subr s(a, str b, int c) {   # a is implicitly var (untyped).
                            # b is explicitly str.
                            # c is explicitly int.
                            # The type-checking is done at the callsite
                            # when arguments are bound to parameters.
                            #
    var b = 100;   # error  # Re-declaration in the same scope is disallowed.
    int n = 10;             # Declaration of variable local to the subroutine.
    n = 20;                 # Assignment is OK.
    int n = 30;    # error  # Re-declaration in the same scope is disallowed.
    str n = "abc"; # error  # Re-declaration in the same scope is disallowed.
                            #
    float f1 = 1;  # error  # 1 is an int, not a float.
    float f2 = 2.0;         # 2.0 is a float.
    num f3 = 3;             # 3 is a num.
    num f4 = 4.0;           # 4.0 is a num.
}                           #
                            #
call s(1, 2, 3);            # Type-assertion '3 is int' is done here at the callsite.
                            #
k = "def";                  # Top-level variable k.
                            #
for (str k, v in $*) {      # k and v are bound here, masking outer k.
  print k . ":" . v;        # k is explicitly str; v is implicitly var.
}                           #
                            #
print "k is".k;             # k at this scope level is still "def".
print "v is".v;             # v is undefined in this scope.
                            #
i = -1;                     #
for (i = 1, int j = 2; i <= 10; i += 1, j *= 2) {
                            # C-style triple-for variables use enclosing scope,
                            # unless declared local: i is outer, j is local to the loop.
  print "inner i =", i;     #
  print "inner j =", j;     #
}                           #
print "outer i =", i;       # i has been modified by the loop.
print "outer j =", j;       # j is undefined in this scope.
</pre>

## Map literals

Miller's `put`/`filter` DSL has four kinds of maps. **Stream records** are (single-level) maps from name to value. **Out-of-stream variables** and **local variables** can also be maps, although they can be multi-level maps (e.g. `@sum[$x][$y]`).  The fourth kind is **map literals**. These cannot be on the left-hand side of assignment expressions. Syntactically they look like JSON, although Miller allows string and integer keys in its map literals while JSON allows only string keys (e.g. `"3"` rather than `3`). Note though that integer keys become stringified in Miller: `@mymap[3]=4` results in `@mymap` being `{"3":4}`.

For example, the following swaps the input stream's `a` and `i` fields, modifies `y`, and drops the rest:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put '</b>
<b>  $* = {</b>
<b>    "a": $i,</b>
<b>    "i": $a,</b>
<b>    "y": $y * 10,</b>
<b>  }</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a i   y
1 pan 7.26802
2 eks 5.22151
3 wye 3.3831800000000003
4 eks 1.34188
5 wye 8.636239999999999
</pre>

Likewise, you can assign map literals to out-of-stream variables or local variables; pass them as arguments to user-defined functions, return them from functions, and so on:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small put '</b>
<b>  func f(map m): map {</b>
<b>    m["x"] *= 200;</b>
<b>    return m;</b>
<b>  }</b>
<b>  $* = f({"a": $a, "x": $x});</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,x=69.3582
a=eks,x=151.7358
a=wye,x=40.9206
a=eks,x=76.2798
a=wye,x=114.6576
</pre>

Like out-of-stream and local variables, map literals can be multi-level:

<pre class="pre-highlight-in-pair">
<b>mlr --from data/small put -q '</b>
<b>  begin {</b>
<b>    @o = {</b>
<b>      "nrec": 0,</b>
<b>      "nkey": {"numeric":0, "non-numeric":0},</b>
<b>    };</b>
<b>  }</b>
<b>  @o["nrec"] += 1;</b>
<b>  for (k, v in $*) {</b>
<b>    if (is_numeric(v)) {</b>
<b>      @o["nkey"]["numeric"] += 1;</b>
<b>    } else {</b>
<b>      @o["nkey"]["non-numeric"] += 1;</b>
<b>    }</b>
<b>  }</b>
<b>  end {</b>
<b>    dump @o;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "nrec": 5,
  "nkey": {
    "numeric": 15,
    "non-numeric": 10
  }
}
</pre>

See also the [Maps page](reference-main-maps.md).

## Built-in variables

These are written all in capital letters, and only a small, specific set of them is defined by Miller.

Namely, Miller supports the following five built-in variables for [filter and
put](reference-dsl.md), all `awk`-inspired: `NF`, `NR`, `FNR`, `FILENUM`, and
`FILENAME`, as well as the mathematical constants `M_PI` and `M_E`.  As well,
there are the read-only separator variables `IRS`, `ORS`, `IFS`, `OFS`, `IPS`,
and `OPS` as discussed on the [separators page](reference-main-separators.md),
and the flatten/unflatten separator `FLATSEP` discussed on the
[flatten/unflatten page](flatten-unflatten.md).  Lastly, the `ENV` map allows
read/write access to environment variables, e.g.  `ENV["HOME"]` or
`ENV["foo_".$hostname]` or `ENV["VERSION"]="1.2.3"`.

<!--- TODO: FLATSEP IFLATSEP OFLATSEP --->

<pre class="pre-highlight-in-pair">
<b>mlr filter 'FNR == 2' data/small*</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=eks,b=pan,i=2,x=0.758679,y=0.522151
1=pan,2=pan,3=1,4=0.3467901443380824,5=0.7268028627434533
a=wye,b=eks,i=10000,x=0.734806020620654365,y=0.884788571337605134
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr put '$fnr = FNR' data/small*</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802,fnr=1
a=eks,b=pan,i=2,x=0.758679,y=0.522151,fnr=2
a=wye,b=wye,i=3,x=0.204603,y=0.338318,fnr=3
a=eks,b=wye,i=4,x=0.381399,y=0.134188,fnr=4
a=wye,b=pan,i=5,x=0.573288,y=0.863624,fnr=5
1=a,2=b,3=i,4=x,5=y,fnr=1
1=pan,2=pan,3=1,4=0.3467901443380824,5=0.7268028627434533,fnr=2
1=eks,2=pan,3=2,4=0.7586799647899636,5=0.5221511083334797,fnr=3
1=wye,2=wye,3=3,4=0.20460330576630303,5=0.33831852551664776,fnr=4
1=eks,2=wye,3=4,4=0.38139939387114097,5=0.13418874328430463,fnr=5
1=wye,2=pan,3=5,4=0.5732889198020006,5=0.8636244699032729,fnr=6
a=pan,b=eks,i=9999,x=0.267481232652199086,y=0.557077185510228001,fnr=1
a=wye,b=eks,i=10000,x=0.734806020620654365,y=0.884788571337605134,fnr=2
a=pan,b=wye,i=10001,x=0.870530722602517626,y=0.009854780514656930,fnr=3
a=hat,b=wye,i=10002,x=0.321507044286237609,y=0.568893318795083758,fnr=4
a=pan,b=zee,i=10003,x=0.272054845593895200,y=0.425789896597056627,fnr=5
</pre>

Their values of `NF`, `NR`, `FNR`, `FILENUM`, and `FILENAME` change from one
record to the next as Miller scans through your input data stream. The
mathematical constants, of course, do not change; `ENV` is populated from the
system environment variables at the time Miller starts. Any changes made to
`ENV` by assigning to it will affect any subprocesses, such as using
[piped tee](reference-dsl-output-statements.md#redirected-output-statements).

Their **scope is global**: you can refer to them in any `filter` or `put` statement. Their values are assigned by the input-record reader:

<pre class="pre-highlight-in-pair">
<b>mlr --csv put '$nr = NR' data/a.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c,nr
1,2,3,1
4,5,6,2
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv repeat -n 3 then put '$nr = NR' data/a.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c,nr
1,2,3,1
1,2,3,1
1,2,3,1
4,5,6,2
4,5,6,2
4,5,6,2
</pre>

The **extent** is for the duration of the put/filter: in a `begin` statement (which executes before the first input record is consumed) you will find `NR=1` and in an `end` statement (which is executed after the last input record is consumed) you will find `NR` to be the total number of records ingested.

These are all **read-only** for the `mlr put` and `mlr filter` DSL: they may be assigned from, e.g. `$nr=NR`, but they may not be assigned to: `NR=100` is a syntax error.

## Type-checking

Miller's `put`/`filter` DSL supports two optional kinds of type-checking.  One is inline **type-tests** and **type-assertions** within expressions.  The other is **type declarations** for assignments to local variables, binding of arguments to user-defined functions, and return values from user-defined functions, These are discussed in the following subsections.

Use of type-checking is entirely up to you: omit it if you want flexibility with heterogeneous data; use it if you want to help catch misspellings in your DSL code or unexpected irregularities in your input data.

### Type-test and type-assertion expressions

The following `is_...` functions take a value and return a boolean indicating whether the argument is of the indicated type. The `assert_...` functions return their argument if it is of the specified type, and cause a fatal error otherwise:

<pre class="pre-highlight-in-pair">
<b>mlr -f | grep ^is</b>
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
is_nan
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

<pre class="pre-highlight-in-pair">
<b>mlr -f | grep ^assert</b>
</pre>
<pre class="pre-non-highlight-in-pair">
asserting_absent
asserting_array
asserting_bool
asserting_boolean
asserting_empty
asserting_empty_map
asserting_error
asserting_float
asserting_int
asserting_map
asserting_nonempty_map
asserting_not_array
asserting_not_empty
asserting_not_map
asserting_not_null
asserting_null
asserting_numeric
asserting_present
asserting_string
</pre>

See [Data-cleaning Examples](data-cleaning-examples.md) for examples of how to use these.

### Type-declarations for local variables, function parameter, and function return values

Local variables can be defined either untyped as in `x = 1`, or typed as in `int x = 1`. Types include **var** (explicitly untyped), **int**, **float**, **num** (int or float), **str**, **bool**, and **map**. These optional type declarations are enforced at the time values are assigned to variables: whether at the initial value assignment as in `int x = 1` or in any subsequent assignments to the same variable farther down in the scope.

The reason for `num` is that `int` and `float` typedecls are very precise:

<pre class="pre-non-highlight-non-pair">
float a = 0;   # Runtime error since 0 is int not float
int   b = 1.0; # Runtime error since 1.0 is float not int
num   c = 0;   # OK
num   d = 1.0; # OK
</pre>

A suggestion is to use `num` for general use when you want numeric content, and use `int` when you genuinely want integer-only values, e.g. in loop indices or map keys (since Miller map keys can only be strings or ints).

The `var` type declaration indicates no type restrictions, e.g. `var x = 1` has the same type restrictions on `x` as `x = 1`. The difference is in intentional shadowing: if you have `x = 1` in outer scope and `x = 2` in inner scope (e.g. within a for-loop or an if-statement) then outer-scope `x` has value 2 after the second assignment.  But if you have `var x = 2` in the inner scope, then you are declaring a variable scoped to the inner block.) For example:

<pre class="pre-non-highlight-non-pair">
x = 1;
if (NR == 4) {
  x = 2; # Refers to outer-scope x: value changes from 1 to 2.
}
print x; # Value of x is now two
</pre>

<pre class="pre-non-highlight-non-pair">
x = 1;
if (NR == 4) {
  var x = 2; # Defines a new inner-scope x with value 2
}
print x;     # Value of this x is still 1
</pre>

Likewise function arguments can optionally be typed, with type enforced when the function is called:

<pre class="pre-non-highlight-non-pair">
func f(map m, int i) {
  ...
}
$a = f({1:2, 3:4}, 5);     # OK
$b = f({1:2, 3:4}, "abc"); # Runtime error
$c = f({1:2, 3:4}, $x);    # Runtime error for records with non-integer field named x
if (NR == 4) {
  var x = 2; # Defines a new inner-scope x with value 2
}
print x;     # Value of this x is still 1
</pre>

Thirdly, function return values can be type-checked at the point of `return` using `:` and a typedecl after the parameter list:

<pre class="pre-non-highlight-non-pair">
func f(map m, int i): bool {
  ...
  ...
  if (...) {
    return "false"; # Runtime error if this branch is taken
  }
  ...
  ...
  if (...) {
    return retval; # Runtime error if this function doesn't have an in-scope
    # boolean-valued variable named retval
  }
  ...
  ...
  # In Miller if your functions don't explicitly return a value, they return absent-null.
  # So it would also be a runtime error on reaching the end of this function without
  # an explicit return statement.
}
</pre>

## Aggregate variable assignments

There are three remaining kinds of variable assignment using out-of-stream variables, the last two of which use the `$*` syntax:

* Recursive copy of out-of-stream variables
* Out-of-stream variable assigned to full stream record
* Full stream record assigned to an out-of-stream variable

Example recursive copy of out-of-stream variables:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --from data/small put -q '</b>
<b>  @v["sum"] += $x;</b>
<b>  @v["count"] += 1;</b>
<b>  end{</b>
<b>    dump;</b>
<b>    @w = @v;</b>
<b>    dump</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "v": {
    "sum": 2.26476,
    "count": 5
  }
}
{
  "v": {
    "sum": 2.26476,
    "count": 5
  },
  "w": {
    "sum": 2.26476,
    "count": 5
  }
}
</pre>

Example of out-of-stream variable assigned to full stream record, where the 2nd record is stashed, and the 4th record is overwritten with that:

<pre class="pre-highlight-in-pair">
<b>mlr put 'NR == 2 {@keep = $*}; NR == 4 {$* = @keep}' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a=pan,b=pan,i=1,x=0.346791,y=0.726802
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=wye,i=3,x=0.204603,y=0.338318
a=eks,b=pan,i=2,x=0.758679,y=0.522151
a=wye,b=pan,i=5,x=0.573288,y=0.863624
</pre>

Example of full stream record assigned to an out-of-stream variable, finding the record for which the `x` field has the largest value in the input stream:

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

<pre class="pre-highlight-in-pair">
<b>mlr --opprint put -q '</b>
<b>  is_null(@xmax) || $x > @xmax {@xmax = $x; @recmax = $*};</b>
<b>  end {emit @recmax}</b>
<b>' data/small</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y
eks pan 2 0.758679 0.522151
</pre>

## Keywords for filter and put

<pre class="pre-highlight-in-pair">
<b>mlr help list-keywords # you can also use mlr -k</b>
</pre>
<pre class="pre-non-highlight-in-pair">
all
begin
bool
break
call
continue
do
dump
edump
elif
else
emit1
emit
emitf
emitp
end
eprint
eprintn
false
filter
float
for
func
funct
if
in
int
map
num
print
printn
return
stderr
stdout
str
subr
tee
true
unset
var
while
ENV
FILENAME
FILENUM
FNR
IFS
IPS
IRS
M_E
M_PI
NF
NR
OFS
OPS
ORS
</pre>


<pre class="pre-highlight-in-pair">
<b>mlr help usage-keywords # you can also use mlr -K</b>
</pre>
<pre class="pre-non-highlight-in-pair">
all: used in "emit1", "emit", "emitp", and "unset" as a synonym for @*

begin: defines a block of statements to be executed before input records
are ingested. The body statements must be wrapped in curly braces.

  Example: 'begin { @count = 0 }'

bool: declares a boolean local variable in the current curly-braced scope.
Type-checking happens at assignment: 'bool b = 1' is an error.

break: causes execution to continue after the body of the current for/while/do-while loop.

call: used for invoking a user-defined subroutine.

  Example: 'subr s(k,v) { print k . " is " . v} call s("a", $a)'

continue: causes execution to skip the remaining statements in the body of
the current for/while/do-while loop. For-loop increments are still applied.

do: with "while", introduces a do-while loop. The body statements must be wrapped
in curly braces.

dump: prints all currently defined out-of-stream variables immediately
to stdout as JSON.

With >, >>, or |, the data do not become part of the output record stream but
are instead redirected.

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump }'
  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >  "mytap.dat"}'
  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump >> "mytap.dat"}'
  Example: mlr --from f.dat put -q '@v[NR]=$*; end { dump | "jq .[]"}'

edump: prints all currently defined out-of-stream variables immediately
to stderr as JSON.

  Example: mlr --from f.dat put -q '@v[NR]=$*; end { edump }'

elif: the way Miller spells "else if". The body statements must be wrapped
in curly braces.

else: terminates an if/elif/elif chain. The body statements must be wrapped
in curly braces.

emit1: inserts an out-of-stream variable into the output record stream. Unlike
the other map variants, side-by-sides, indexing, and redirection are not supported,
but you can emit any map-valued expression.

  Example: mlr --from f.dat put 'emit1 $*'
  Example: mlr --from f.dat put 'emit1 mapsum({"id": NR}, $*)'

Please see https://miller.readthedocs.io://johnkerl.org/miller/doc for more information.

emit: inserts an out-of-stream variable into the output record stream. Hashmap
indices present in the data but not slotted by emit arguments are not output.

With >, >>, or |, the data do not become part of the output record stream but
are instead redirected.

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
etc., to control the format of the output if the output is redirected. See also mlr -h.

  Example: mlr --from f.dat put 'emit >  "/tmp/data-".$a, $*'
  Example: mlr --from f.dat put 'emit >  "/tmp/data-".$a, mapexcept($*, "a")'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums'
  Example: mlr --from f.dat put --ojson '@sums[$a][$b]+=$x; emit > "tap-".$a.$b.".dat", @sums'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @sums, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >  "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit >> "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | "gzip > mytap.dat.gz", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit > stderr, @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emit | "grep somepattern", @*, "index1", "index2"'

Please see https://miller.readthedocs.io://johnkerl.org/miller/doc for more information.

emitf: inserts non-indexed out-of-stream variable(s) side-by-side into the
output record stream.

With >, >>, or |, the data do not become part of the output record stream but
are instead redirected.

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
etc., to control the format of the output if the output is redirected. See also mlr -h.

  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a'
  Example: mlr --from f.dat put --oxtab '@a=$i;@b+=$x;@c+=$y; emitf > "tap-".$i.".dat", @a'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > "mytap.dat", @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf >> "mytap.dat", @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf > stderr, @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | "grep somepattern", @a, @b, @c'
  Example: mlr --from f.dat put '@a=$i;@b+=$x;@c+=$y; emitf | "grep somepattern > mytap.dat", @a, @b, @c'

Please see https://miller.readthedocs.io://johnkerl.org/miller/doc for more information.

emitp: inserts an out-of-stream variable into the output record stream.
Hashmap indices present in the data but not slotted by emitp arguments are
output concatenated with ":".

With >, >>, or |, the data do not become part of the output record stream but
are instead redirected.

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
etc., to control the format of the output if the output is redirected. See also mlr -h.

  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums'
  Example: mlr --from f.dat put --opprint '@sums[$a][$b]+=$x; emitp > "tap-".$a.$b.".dat", @sums'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @sums, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >  "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp >> "mytap.dat", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | "gzip > mytap.dat.gz", @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp > stderr, @*, "index1", "index2"'
  Example: mlr --from f.dat put '@sums[$a][$b]+=$x; emitp | "grep somepattern", @*, "index1", "index2"'

Please see https://miller.readthedocs.io://johnkerl.org/miller/doc for more information.

end: defines a block of statements to be executed after input records
are ingested. The body statements must be wrapped in curly braces.

  Example: 'end { emit @count }'
  Example: 'end { eprint "Final count is " . @count }'

eprint: prints expression immediately to stderr.

  Example: mlr --from f.dat put -q 'eprint "The sum of x and y is ".($x+$y)'
  Example: mlr --from f.dat put -q 'for (k, v in $*) { eprint k . " => " . v }'
  Example: mlr --from f.dat put  '(NR % 1000 == 0) { eprint "Checkpoint ".NR}'

eprintn: prints expression immediately to stderr, without trailing newline.

  Example: mlr --from f.dat put -q 'eprintn "The sum of x and y is ".($x+$y); eprint ""'

false: the boolean literal value.

filter: includes/excludes the record in the output record stream.

  Example: mlr --from f.dat put 'filter (NR == 2 || $x > 5.4)'

Instead of put with 'filter false' you can simply use put -q.  The following
uses the input record to accumulate data but only prints the running sum
without printing the input record:

  Example: mlr --from f.dat put -q '@running_sum += $x * $y; emit @running_sum'

float: declares a floating-point local variable in the current curly-braced scope.
Type-checking happens at assignment: 'float x = 0' is an error.

for: defines a for-loop using one of three styles. The body statements must
be wrapped in curly braces.
For-loop over stream record:

  Example:  'for (k, v in $*) { ... }'

For-loop over out-of-stream variables:

  Example: 'for (k, v in @counts) { ... }'
  Example: 'for ((k1, k2), v in @counts) { ... }'
  Example: 'for ((k1, k2, k3), v in @*) { ... }'

C-style for-loop:

  Example:  'for (var i = 0, var b = 1; i < 10; i += 1, b *= 2) { ... }'

func: used for defining a user-defined function.

  Example: 'func f(a,b) { return sqrt(a**2+b**2)} $d = f($x, $y)'

funct: used for saying that a function argument is a user-defined function.

  Example: 'func g(num a, num b, funct f) :num { return f(a**2+b**2) }'

if: starts an if/elif/elif chain. The body statements must be wrapped
in curly braces.

in: used in for-loops over stream records or out-of-stream variables.

int: declares an integer local variable in the current curly-braced scope.
Type-checking happens at assignment: 'int x = 0.0' is an error.

map: declares an map-valued local variable in the current curly-braced scope.
Type-checking happens at assignment: 'map b = 0' is an error. map b = {} is
always OK. map b = a is OK or not depending on whether a is a map.

num: declares an int/float local variable in the current curly-braced scope.
Type-checking happens at assignment: 'num b = true' is an error.

print: prints expression immediately to stdout.

  Example: mlr --from f.dat put -q 'print "The sum of x and y is ".($x+$y)'
  Example: mlr --from f.dat put -q 'for (k, v in $*) { print k . " => " . v }'
  Example: mlr --from f.dat put  '(NR % 1000 == 0) { print > stderr, "Checkpoint ".NR}'

printn: prints expression immediately to stdout, without trailing newline.

  Example: mlr --from f.dat put -q 'printn "."; end { print "" }'

return: specifies the return value from a user-defined function.
Omitted return statements (including via if-branches) result in an absent-null
return value, which in turns results in a skipped assignment to an LHS.

stderr: Used for tee, emit, emitf, emitp, print, and dump in place of filename
to print to standard error.

stdout: Used for tee, emit, emitf, emitp, print, and dump in place of filename
to print to standard output.

str: declares a string local variable in the current curly-braced scope.
Type-checking happens at assignment.

subr: used for defining a subroutine.

  Example: 'subr s(k,v) { print k . " is " . v} call s("a", $a)'

tee: prints the current record to specified file.
This is an immediate print to the specified file (except for pprint format
which of course waits until the end of the input stream to format all output).

The > and >> are for write and append, as in the shell, but (as with awk) the
file-overwrite for > is on first write, not per record. The | is for piping to
a process which will process the data. There will be one open file for each
distinct file name (for > and >>) or one subordinate process for each distinct
value of the piped-to command (for |). Output-formatting flags are taken from
the main command line.

You can use any of the output-format command-line flags, e.g. --ocsv, --ofs,
etc., to control the format of the output. See also mlr -h.

emit with redirect and tee with redirect are identical, except tee can only
output $*.

  Example: mlr --from f.dat put 'tee >  "/tmp/data-".$a, $*'
  Example: mlr --from f.dat put 'tee >> "/tmp/data-".$a.$b, $*'
  Example: mlr --from f.dat put 'tee >  stderr, $*'
  Example: mlr --from f.dat put -q 'tee | "tr \[a-z\\] \[A-Z\\]", $*'
  Example: mlr --from f.dat put -q 'tee | "tr \[a-z\\] \[A-Z\\] > /tmp/data-".$a, $*'
  Example: mlr --from f.dat put -q 'tee | "gzip > /tmp/data-".$a.".gz", $*'
  Example: mlr --from f.dat put -q --ojson 'tee | "gzip > /tmp/data-".$a.".gz", $*'

true: the boolean literal value.

unset: clears field(s) from the current record, or an out-of-stream or local variable.

  Example: mlr --from f.dat put 'unset $x'
  Example: mlr --from f.dat put 'unset $*'
  Example: mlr --from f.dat put 'for (k, v in $*) { if (k =~ "a.*") { unset $[k] } }'
  Example: mlr --from f.dat put '...; unset @sums'
  Example: mlr --from f.dat put '...; unset @sums["green"]'
  Example: mlr --from f.dat put '...; unset @*'

var: declares an untyped local variable in the current curly-braced scope.

  Examples: 'var a=1', 'var xyz=""'

while: introduces a while loop, or with "do", introduces a do-while loop.
The body statements must be wrapped in curly braces.

ENV: access to environment variables by name, e.g. '$home = ENV["HOME"]'

FILENAME: evaluates to the name of the current file being processed.

FILENUM: evaluates to the number of the current file being processed,
starting with 1.

FNR: evaluates to the number of the current record within the current file
being processed, starting with 1. Resets at the start of each file.

IFS: evaluates to the input field separator from the command line.

IPS: evaluates to the input pair separator from the command line.

IRS: evaluates to the input record separator from the command line,
or to LF or CRLF from the input data if in autodetect mode (which is
the default).

M_E: the mathematical constant e.

M_PI: the mathematical constant pi.

NF: evaluates to the number of fields in the current record.

NR: evaluates to the number of the current record over all files
being processed, starting with 1. Does not reset at the start of each file.

OFS: evaluates to the output field separator from the command line.

OPS: evaluates to the output pair separator from the command line.

ORS: evaluates to the output record separator from the command line,
or to LF or CRLF from the input data if in autodetect mode (which is
the default).
</pre>

