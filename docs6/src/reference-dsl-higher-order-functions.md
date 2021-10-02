<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
<div>
<span class="quicklinks">
Quick links:
&nbsp;
<a class="quicklink" href="../reference-main-flag-list/index.html">Flag list</a>
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
# DSL higher-order functions

A _higher-order function_ is one which takes [another function](reference-dsl-user-defined-functions.md)
as an argument.
As of [Miller 6](new-in-miller-6.md) you can use
[`select`](reference-dsl-builtin-functions.md#select),
[`apply`](reference-dsl-builtin-functions.md#apply),
[`reduce`](reference-dsl-builtin-functions.md#reduce),
[`fold`](reference-dsl-builtin-functions.md#fold), and
[`sort`](reference-dsl-builtin-functions.md#sort) to express flexible,
intuitive operations on arrays and maps, as an alternative to things which
would otherwise require for-loops.

See also the [`get_keys`](reference-dsl-builtin-functions.md#get_keys) and
[`get_values`](reference-dsl-builtin-functions.md#get_values) functions which,
when given a map, return an array of its keys or an array of its values,
respectively.

## select

The [`select`](reference-dsl-builtin-functions.md#select) function takes a map
or array as its first argument and a function as second argument.  It includes
each input element in the ouptut if the function returns true.

For arrays, that function should take one argument, for array element; for
maps, it should take two, for map-element key and value. In either case it
should return a boolean.

A perhaps helpful analogy: the `select` function is to arrays and maps as the
[`filter`](reference-verbs.md#filter) is to records.

Array examples:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_array = [2, 9, 10, 3, 1, 4, 5, 8, 7, 6];</b>
<b></b>
<b>    print "Original:";</b>
<b>    print my_array;</b>
<b></b>
<b>    print;</b>
<b>    print "Evens:";</b>
<b>    print select(my_array, func (e) { return e % 2 == 0});</b>
<b></b>
<b>    print;</b>
<b>    print "Odds:";</b>
<b>    print select(my_array, func (e) { return e % 2 == 1});</b>
<b>    print;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
[2, 9, 10, 3, 1, 4, 5, 8, 7, 6]

Evens:
[2, 10, 4, 8, 6]

Odds:
[9, 3, 1, 5, 7]

</pre>

Map examples:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_map = {"cubit": 823, "dale": 13, "apple": 199, "ember": 191, "bottle": 107};</b>
<b>    print "Original:";</b>
<b>    print my_map;</b>
<b></b>
<b>    print;</b>
<b>    print "Keys with an 'o' in them:";</b>
<b>    print select(my_map, func (k,v) { return k =~ "o"});</b>
<b></b>
<b>    print;</b>
<b>    print "Values with last digit >= 5:";</b>
<b>    print select(my_map, func (k,v) { return v % 10 >= 5});</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
{
  "cubit": 823,
  "dale": 13,
  "apple": 199,
  "ember": 191,
  "bottle": 107
}

Keys with an o in them:
{
  "bottle": 107
}

Values with last digit >= 5:
{
  "apple": 199,
  "bottle": 107
}
</pre>

## apply

The [`apply`](reference-dsl-builtin-functions.md#apply) function takes a map
or array as its first argument and a function as second argument.  It applies
the function to each element of the array or map.

For arrays, the function should take one argument, for array element; it should
return a new element. For maps, it should take two, for map-element key and
value. It should return a new key-value pair (i.e. a single-entry map).

A perhaps helpful analogy: the `apply` function is to arrays and maps as the
[`put`](reference-verbs.md#put) is to records.

Array examples:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_array = [2, 9, 10, 3, 1, 4, 5, 8, 7, 6];</b>
<b>    print "Original:";</b>
<b>    print my_array;</b>
<b></b>
<b>    print;</b>
<b>    print "Squares:";</b>
<b>    print apply(my_array, func(e) { return e**2 });</b>
<b></b>
<b>    print;</b>
<b>    print "Cubes:";</b>
<b>    print apply(my_array, func(e) { return e**3 });</b>
<b></b>
<b>    print;</b>
<b>    print "Sorted cubes:";</b>
<b>    print sort(apply(my_array, func(e) { return e**3 }));</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
[2, 9, 10, 3, 1, 4, 5, 8, 7, 6]

Squares:
[4, 81, 100, 9, 1, 16, 25, 64, 49, 36]

Cubes:
[8, 729, 1000, 27, 1, 64, 125, 512, 343, 216]

Sorted cubes:
[1, 8, 27, 64, 125, 216, 343, 512, 729, 1000]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_map = {"cubit": 823, "dale": 13, "apple": 199, "ember": 191, "bottle": 107};</b>
<b>    print "Original:";</b>
<b>    print my_map;</b>
<b></b>
<b>    print;</b>
<b>    print "Squared values:";</b>
<b>    print apply(my_map, func(k,v) { return {k: v**2} });</b>
<b></b>
<b>    print;</b>
<b>    print "Cubed values, sorted by key:";</b>
<b>    print sort(apply(my_map, func(k,v) { return {k: v**3} }));</b>
<b></b>
<b>    print;</b>
<b>    print "Same, with upcased keys:";</b>
<b>    print sort(apply(my_map, func(k,v) { return {toupper(k): v**3} }));</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
{
  "cubit": 823,
  "dale": 13,
  "apple": 199,
  "ember": 191,
  "bottle": 107
}

Squared values:
{
  "cubit": 677329,
  "dale": 169,
  "apple": 39601,
  "ember": 36481,
  "bottle": 11449
}

Cubed values, sorted by key:
{
  "apple": 7880599,
  "bottle": 1225043,
  "cubit": 557441767,
  "dale": 2197,
  "ember": 6967871
}

Same, with upcased keys:
{
  "APPLE": 7880599,
  "BOTTLE": 1225043,
  "CUBIT": 557441767,
  "DALE": 2197,
  "EMBER": 6967871
}
</pre>

## reduce

The [`reduce`](reference-dsl-builtin-functions.md#reduce) function takes a map
or array as its first argument and a function as second argument. It accumulates entries into a final
output -- for example, sum or product.

For arrays, the function should take two arguments, for accumulated value and
array element; for maps, it should take four, for accumulated key and value
and map-element key and value. In either case it should return the updated
accumulator.

The start value for the accumulator is the first element for arrays, or the
first element's key-value pair for maps.

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_array = [2, 9, 10, 3, 1, 4, 5, 8, 7, 6];</b>
<b></b>
<b>    print "Original:";</b>
<b>    print my_array;</b>
<b></b>
<b>    print;</b>
<b>    print "First element:";</b>
<b>    print reduce(my_array, func (acc,e) { return acc });</b>
<b></b>
<b>    print;</b>
<b>    print "Last element:";</b>
<b>    print reduce(my_array, func (acc,e) { return e });</b>
<b></b>
<b>    print;</b>
<b>    print "Sum of values:";</b>
<b>    print reduce(my_array, func (acc,e) { return acc + e });</b>
<b></b>
<b>    print;</b>
<b>    print "Product of values:";</b>
<b>    print reduce(my_array, func (acc,e) { return acc * e });</b>
<b></b>
<b>    print;</b>
<b>    print "Concatenation of values:";</b>
<b>    print reduce(my_array, func (acc,e) { return acc. "," . e });</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
[2, 9, 10, 3, 1, 4, 5, 8, 7, 6]

First element:
2

Last element:
6

Sum of values:
55

Product of values:
3628800

Concatenation of values:
2,9,10,3,1,4,5,8,7,6
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_map = {"cubit": 823, "dale": 13, "apple": 199, "ember": 191, "bottle": 107};</b>
<b>    print "Original:";</b>
<b>    print my_map;</b>
<b></b>
<b>    print;</b>
<b>    print "First key-value pair:";</b>
<b>    print reduce(my_map, func (acck,accv,ek,ev) { return {acck: accv}});</b>
<b></b>
<b>    print;</b>
<b>    print "Last key-value pair:";</b>
<b>    print reduce(my_map, func (acck,accv,ek,ev) { return {ek: ev}});</b>
<b></b>
<b>    print;</b>
<b>    print "Concatenate keys and values:";</b>
<b>    print reduce(my_map, func (acck,accv,ek,ev) { return {acck . "," . ek: accv . "," . ev}});</b>
<b></b>
<b>    print;</b>
<b>    print "Sum of values:";</b>
<b>    print reduce(my_map, func (acck,accv,ek,ev) { return {"sum": accv + ev }});</b>
<b></b>
<b>    print;</b>
<b>    print "Product of values:";</b>
<b>    print reduce(my_map, func (acck,accv,ek,ev) { return {"product": accv * ev }});</b>
<b></b>
<b>    print;</b>
<b>    print "String-join of values:";</b>
<b>    print reduce(my_map, func (acck,accv,ek,ev) { return {"joined": accv . "," . ev }});</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
{
  "cubit": 823,
  "dale": 13,
  "apple": 199,
  "ember": 191,
  "bottle": 107
}

First key-value pair:
{
  "cubit": 823
}

Last key-value pair:
{
  "bottle": 107
}

Concatenate keys and values:
{
  "cubit,dale,apple,ember,bottle": "823,13,199,191,107"
}

Sum of values:
{
  "sum": 1333
}

Product of values:
{
  "product": 43512437137
}

String-join of values:
{
  "joined": "823,13,199,191,107"
}
</pre>

## fold

The [`fold`](reference-dsl-builtin-functions.md#fold) function is the same as
`reduce`, except that instead of the starting value for the accumulation being
taken from the first entry of the array/map, you specify it as the third
argument.

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_array = [2, 9, 10, 3, 1, 4, 5, 8, 7, 6];</b>
<b></b>
<b>    print "Original:";</b>
<b>    print my_array;</b>
<b></b>
<b>    print;</b>
<b>    print "Sum with reduce:";</b>
<b>    print reduce(my_array, func (acc,e) { return acc + e });</b>
<b></b>
<b>    print;</b>
<b>    print "Sum with fold and 0 initial value:";</b>
<b>    print fold(my_array, func (acc,e) { return acc + e }, 0);</b>
<b></b>
<b>    print;</b>
<b>    print "Sum with fold and 1000000 initial value:";</b>
<b>    print fold(my_array, func (acc,e) { return acc + e }, 1000000);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
[2, 9, 10, 3, 1, 4, 5, 8, 7, 6]

Sum with reduce:
55

Sum with fold and 0 initial value:
55

Sum with fold and 1000000 initial value:
1000055
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_map = {"cubit": 823, "dale": 13, "apple": 199, "ember": 191, "bottle": 107};</b>
<b>    print "Original:";</b>
<b>    print my_map;</b>
<b></b>
<b>    print;</b>
<b>    print "First key-value pair -- note this is the starting accumulator:";</b>
<b>    print fold(my_map, func (acck,accv,ek,ev) { return {acck: accv}}, {"start": 999});</b>
<b></b>
<b>    print;</b>
<b>    print "Last key-value pair:";</b>
<b>    print fold(my_map, func (acck,accv,ek,ev) { return {ek: ev}}, {"start": 999});</b>
<b></b>
<b>    print;</b>
<b>    print "Sum of values with fold and 0 initial value:";</b>
<b>    print fold(my_map, func (acck,accv,ek,ev) { return {"sum": accv + ev} }, {"sum": 0});</b>
<b></b>
<b>    print;</b>
<b>    print "Sum of values with fold and 1000000 initial value:";</b>
<b>    print fold(my_map, func (acck,accv,ek,ev) { return {"sum": accv + ev} }, {"sum": 1000000});</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
{
  "cubit": 823,
  "dale": 13,
  "apple": 199,
  "ember": 191,
  "bottle": 107
}

First key-value pair -- note this is the starting accumulator:
{
  "start": 999
}

Last key-value pair:
{
  "bottle": 107
}

Sum of values with fold and 0 initial value:
{
  "sum": 1333
}

Sum of values with fold and 1000000 initial value:
{
  "sum": 1001333
}
</pre>

## sort

The [`sort`](reference-dsl-builtin-functions.md#sort) function takes a map or
array as its first argument, and it can take a function as second argument.
Unlike the other higher-order functions, the second argument can be omitted
when the natural ordering is desired -- ordered by array element for arrays, or by
key for maps.

As a second option, character flags such as `r` for reverse or `c` for
case-folded lexical sort can be supplied as the second argument.

As a third option, a function can be supplied as the second argument.

For arrays, that function should take two arguments `a` and `b`, returning a
negative, zero, or positive number as `a<b`, `a==b`, or `a>b` respectively.
For maps, the function should take four arguments `ak`, `av`, `bk`, and `bv`,
again returning negative, zero, or positive, using `a` and `b`'s keys and
values.

Array examples:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_array = [2, 9, 10, 3, 1, 4, 5, 8, 7, 6];</b>
<b></b>
<b>    print "Original:";</b>
<b>    print my_array;</b>
<b></b>
<b>    print;</b>
<b>    print "Ascending:";</b>
<b>    print sort(my_array);</b>
<b>    print sort(my_array, func (a,b) { return a <=> b });</b>
<b></b>
<b>    print;</b>
<b>    print "Descending:";</b>
<b>    print sort(my_array, "r");</b>
<b>    print sort(my_array, func (a,b) { return b <=> a });</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
[2, 9, 10, 3, 1, 4, 5, 8, 7, 6]

Ascending:
[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]

Descending:
[10, 9, 8, 7, 6, 5, 4, 3, 2, 1]
[10, 9, 8, 7, 6, 5, 4, 3, 2, 1]
</pre>

Map examples:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    my_map = {"cubit": 823, "dale": 13, "apple": 199, "ember": 191, "bottle": 107};</b>
<b></b>
<b>    print "Original:";</b>
<b>    print my_map;</b>
<b></b>
<b>    print;</b>
<b>    print "Ascending by key:";</b>
<b>    print sort(my_map);</b>
<b>    print sort(my_map, func(ak,av,bk,bv) { return ak <=> bk });</b>
<b></b>
<b>    print;</b>
<b>    print "Descending by key:";</b>
<b>    print sort(my_map, "r");</b>
<b>    print sort(my_map, func(ak,av,bk,bv) { return bk <=> ak });</b>
<b></b>
<b>    print;</b>
<b>    print "Ascending by value:";</b>
<b>    print sort(my_map, func(ak,av,bk,bv) { return av <=> bv });</b>
<b></b>
<b>    print;</b>
<b>    print "Descending by value:";</b>
<b>    print sort(my_map, func(ak,av,bk,bv) { return bv <=> av });</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
{
  "cubit": 823,
  "dale": 13,
  "apple": 199,
  "ember": 191,
  "bottle": 107
}

Ascending by key:
{
  "apple": 199,
  "bottle": 107,
  "cubit": 823,
  "dale": 13,
  "ember": 191
}
{
  "apple": 199,
  "bottle": 107,
  "cubit": 823,
  "dale": 13,
  "ember": 191
}

Descending by key:
{
  "ember": 191,
  "dale": 13,
  "cubit": 823,
  "bottle": 107,
  "apple": 199
}
{
  "ember": 191,
  "dale": 13,
  "cubit": 823,
  "bottle": 107,
  "apple": 199
}

Ascending by value:
{
  "dale": 13,
  "bottle": 107,
  "ember": 191,
  "apple": 199,
  "cubit": 823
}

Descending by value:
{
  "cubit": 823,
  "apple": 199,
  "ember": 191,
  "bottle": 107,
  "dale": 13
}
</pre>

Please see the [sorting page](sorting.md) for more examples.

## Combined examples

Using a paradigm from the [page on operating on all
records](operating-on-all-records.md), we can retain a column from the input
data as an array, then apply some higher-order functions to it:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p cat example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
yellow triangle true  1  11    43.6498  9.8870
red    square   true  2  15    79.2778  0.0130
red    circle   true  3  16    13.8103  2.9010
red    square   false 4  48    77.5542  7.4670
purple triangle false 5  51    81.2290  8.5910
red    square   false 6  64    77.1991  9.5310
purple triangle false 7  65    80.1405  5.8240
yellow circle   true  8  73    63.9785  4.2370
yellow circle   true  9  87    63.5058  8.3350
purple square   false 10 91    72.3735  8.2430
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv put -q '</b>
<b>  begin {</b>
<b>    @indexes = [] # So auto-extend will make an array, not a map</b>
<b>  }</b>
<b>  @indexes[NR] = $index;</b>
<b>  end {</b>
<b></b>
<b>    print "Original:";</b>
<b>    print @indexes;</b>
<b></b>
<b>    print;</b>
<b>    print "Sorted:";</b>
<b>    print sort(@indexes, "r");</b>
<b></b>
<b>    print;</b>
<b>    print "Sorted, then cubed:";</b>
<b>    print apply(</b>
<b>      sort(@indexes, "r"),</b>
<b>      func(e) { return e**3 },</b>
<b>    );</b>
<b></b>
<b>    print;</b>
<b>    print "Sorted, then cubed, then summed:";</b>
<b>    print reduce(</b>
<b>      apply(</b>
<b>        sort(@indexes, "r"),</b>
<b>        func(e) { return e**3 },</b>
<b>      ),</b>
<b>      func(acc, e) { return acc + e },</b>
<b>    )</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Original:
[11, 15, 16, 48, 51, 64, 65, 73, 87, 91]

Sorted:
[91, 87, 73, 65, 64, 51, 48, 16, 15, 11]

Sorted, then cubed:
[753571, 658503, 389017, 274625, 262144, 132651, 110592, 4096, 3375, 1331]

Sorted, then cubed, then summed:
2589905
</pre>

## Caveats

### Remember return

From other languages it's easy to accidentially write

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end { print select([1,2,3,4,5], func (e) { e >= 3 })}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
mlr: select: selector function returned non-boolean "(absent)".
</pre>

instead of

<pre class="pre-highlight-in-pair">
<b>mlr -n put 'end { print select([1,2,3,4,5], func (e) { return e >= 3 })}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[3, 4, 5]
</pre>

### No IIFEs

As of September 2021, immediately invoked function expressions (IIFEs) are not part of the Miller DSL's grammar. For example, this doesn't work yet:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = 3;</b>
<b>    y = (func (e) { return e**7 })(x);</b>
<b>    print y;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
mlr: cannot parse DSL expression.
Parse error on token "(" at line 4 columnn 35.
Please check for missing semicolon.
Expected one of:
  ; } > >> | ? || ^^ && ?? ??? =~ !=~ == != <=> >= < <= ^ & << >>> + - .+
  .- * / // % .* ./ .// . **

</pre>

but this does:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = 3;</b>
<b>    f = func (e) { return e**7 };</b>
<b>    y = f(x);</b>
<b>    print y;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
2187
</pre>

### Built-in functions currently unsupported as arguments

[Built-in functions](reference-dsl-user-defined-functions.md) are, as of
September 2021, a bit separate from [user-defined
functions](reference-dsl-builtin-functions.md) internally to Miller, and can't
be used directly as arguments to higher-order functions.

For example, this doesn't work yet:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    notches = [0,1,2,3];</b>
<b>    radians = apply(notches, func (e) { return e * M_PI / 8 });</b>
<b>    cosines = apply(radians, cos);</b>
<b>    print cosines;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
mlr: apply: second argument must be a function; got absent.
</pre>

but this does:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    notches = [0,1,2,3];</b>
<b>    radians = apply(notches, func (e) { return e * M_PI / 8 });</b>
<b>    # cosines = apply(radians, cos);</b>
<b>    cosines = apply(radians, func (e) { return cos(e) });</b>
<b>    print cosines;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[1, 0.9238795325112867, 0.7071067811865476, 0.38268343236508984]
</pre>
