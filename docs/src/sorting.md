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
# Sorting

Miller gives you three ways to sort your data:

* The [`sort`](reference-verbs.md#sort) verb lets you sort records (rows) by various fields (columns).
* The [`sort-within-records`](reference-verbs.md#sort-within-records) verb lets you sort fields within records.
* The [`sort`](reference-dsl-builtin-functions.md#sort) DSL function gives you more customizable options for sorting data either within fields or across records. (See also the [higher-order-functions page](reference-dsl-higher-order-functions.md) for related information.)

## Sorting records: the sort verb

The `sort` verb (see [its documentation](reference-verbs.md#sort) for more
information) reorders entire records within the data stream. You can sort
lexically (with or without case-folding) or numerically, ascending or
descending; and you can sort primary by one column, then secondarily by
another, etc.

Input data:

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

Sorted numerically ascending by rate:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p sort -n rate example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
red    square   true  2  15    79.2778  0.0130
red    circle   true  3  16    13.8103  2.9010
yellow circle   true  8  73    63.9785  4.2370
purple triangle false 7  65    80.1405  5.8240
red    square   false 4  48    77.5542  7.4670
purple square   false 10 91    72.3735  8.2430
yellow circle   true  9  87    63.5058  8.3350
purple triangle false 5  51    81.2290  8.5910
red    square   false 6  64    77.1991  9.5310
yellow triangle true  1  11    43.6498  9.8870
</pre>

Sorted lexically ascending by color; then, within each color, numerically descending by quantity:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p sort -f color -nr quantity example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
purple triangle false 5  51    81.2290  8.5910
purple triangle false 7  65    80.1405  5.8240
purple square   false 10 91    72.3735  8.2430
red    square   true  2  15    79.2778  0.0130
red    square   false 4  48    77.5542  7.4670
red    square   false 6  64    77.1991  9.5310
red    circle   true  3  16    13.8103  2.9010
yellow circle   true  8  73    63.9785  4.2370
yellow circle   true  9  87    63.5058  8.3350
yellow triangle true  1  11    43.6498  9.8870
</pre>

## Sorting fields within records: the sort-within-records verb

The `sort-within-records` verb (see [its
documentation](reference-verbs.md#sort-within-records) for more information)
leaves records in their original order in the data stream, but reorders fields
within each record. A typical use-case is for given all records the same column-ordering,
in particular for converting JSON to CSV (or other tabular formats):

<pre class="pre-highlight-in-pair">
<b>cat data/sort-within-records.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 1,
  "b": 2,
  "c": 3
}
{
  "b": 4,
  "a": 5,
  "c": 6
}
{
  "c": 7,
  "b": 8,
  "a": 9
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint cat data/sort-within-records.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b c
1 2 3

b a c
4 5 6

c b a
7 8 9
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --ijson --opprint sort-within-records data/sort-within-records.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a b c
1 2 3
5 4 6
9 8 7
</pre>

## The sort function by example

* It returns a sorted copy of an input array or map.
* Without second argument, uses the natural ordering.
* With second which is string, takes sorting flags from it: `"f"` for lexical or `"c"` for case-folded lexical, and/or `"r"` for reverse/descending.

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    # Sort array with natural ordering</b>
<b>    print sort([5,2,3,1,4]);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[1, 2, 3, 4, 5]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    # Sort array with reverse-natural ordering</b>
<b>    print sort([5,2,3,1,4], "r");</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[5, 4, 3, 2, 1]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    # Sort array with custom function: natural ordering</b>
<b>    print sort([5,2,3,1,4], func(a,b) { return a <=> b});</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[1, 2, 3, 4, 5]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    # Sort array with custom function: reverse-natural ordering</b>
<b>    print sort([5,2,3,1,4], func(a,b) { return b <=> a});</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[5, 4, 3, 2, 1]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    # Sort map with natural ordering on keys</b>
<b>    print sort({"c":2, "a": 3, "b": 1});</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 3,
  "b": 1,
  "c": 2
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    # Sort map with reverse-natural ordering on keys</b>
<b>    print sort({"c":2, "a": 3, "b": 1}, "r");</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "c": 2,
  "b": 1,
  "a": 3
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    # Sort map with custom function: natural ordering on values</b>
<b>    print sort({"c":2, "a": 3, "b": 1}, func(ak,av,bk,bv){return av <=> bv});</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "b": 1,
  "c": 2,
  "a": 3
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    # Sort map with custom function: reverse-natural ordering on values</b>
<b>    print sort({"c":2, "a": 3, "b": 1}, func(ak,av,bk,bv){return bv <=> av});</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 3,
  "c": 2,
  "b": 1
}
</pre>

In the rest of this page we'll look more closely at these variants.

## Simple sorting of arrays

Using the [`sort`](reference-dsl-builtin-functions.md#sort) function, you can
get a copy of an array, sorted by its values -- optionally, with reversed
order, and/or lexical/case-folded sorting. The first argument is an array to be
sorted. The optional second argument is a string containing any of the
characters `n` for numeric (the default anyway), `f` for lexical, or `c` for
case-folded lexical, and `r` for reverse.  Note that `sort` does not modify
its argument; it returns a sorted copy.

Also note that all the flags to `sort` allow you to operate on arrays which
contain strings, floats, and booleans; if you need to sort an array whose
values are themselves maps or arrays, you'll need `sort` with function argument
as described further down in this page.

<pre class="pre-highlight-in-pair">
<b>cat data/sorta-example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key,values
alpha,4;6;1;5
beta,7;9;9;8
gamma,11;2;1;12
</pre>

Default sort is numerical ascending:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/sorta-example.csv put '</b>
<b>  $values = splita($values, ";");</b>
<b>  $values = sort($values);        # default flags</b>
<b>  $values = joinv($values, ";");</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key   values
alpha 1;4;5;6
beta  7;8;9;9
gamma 1;2;11;12
</pre>

Use the `"r"` flag for reverse, which is numerical descending:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/sorta-example.csv put '</b>
<b>  $values = splita($values, ";");</b>
<b>  $values = sort($values, "r");   # 'r' flag for reverse sort</b>
<b>  $values = joinv($values, ";");</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key   values
alpha 6;5;4;1
beta  9;9;8;7
gamma 12;11;2;1
</pre>

Use the `"f"` flag for lexical ascending sort (and `"fr"` would lexical descending):

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/sorta-example.csv put '</b>
<b>  $values = splita($values, ";");</b>
<b>  $values = sort($values, "f");   # 'f' flag for lexical sort</b>
<b>  $values = joinv($values, ";");</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key   values
alpha 1;4;5;6
beta  7;8;9;9
gamma 1;11;12;2
</pre>

Without and with case-folding:

<pre class="pre-highlight-in-pair">
<b>cat data/sorta-example-text.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key,values
alpha,cat;bat;Australia;Bavaria;apple;Colombia
alpha,cat;bat;Australia;Bavaria;apple;Colombia
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from data/sorta-example-text.csv put '</b>
<b>  $values = splita($values, ";");</b>
<b>  if (NR == 1) {</b>
<b>    $values = sort($values, "f"); # 'f' flag for (non-folded) lexical sort</b>
<b>  } else {</b>
<b>    $values = sort($values, "c"); # 'c' flag for case-folded lexical sort</b>
<b>  }</b>
<b>  $values = joinv($values, ";");</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key   values
alpha Australia;Bavaria;Colombia;apple;bat;cat
alpha apple;Australia;bat;Bavaria;cat;Colombia
</pre>

## Simple sorting of maps within records

Using the [`sort`](reference-dsl-builtin-functions.md#sort) function, you
can sort a map by its keys.

Since `sort` only gives you options for sorting a map by its keys, if you want
to sort a map by its values you'll need `sort` with function argument as
described further down in this page.

Also note that, unlike the `sort-within-record` verb with its `-r` flag,
`sort` doesn't recurse into submaps and sort those.

<pre class="pre-highlight-in-pair">
<b>cat data/server-log.json</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "hostname": "localhost",
  "pid": 12345,
  "req": {
    "id": 6789,
    "method": "GET",
    "path": "api/check",
    "host": "foo.bar",
    "headers": {
      "host": "bar.baz",
      "user-agent": "browser"
    }
  },
  "res": {
    "status_code": 200,
    "header": {
      "content-type": "text",
      "content-encoding": "plain"
    }
  }
}
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --json --from data/server-log.json put '</b>
<b>  $req = sort($req);      # Ascending here</b>
<b>  $res = sort($res, "r"); # Descending here</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "hostname": "localhost",
  "pid": 12345,
  "req": {
    "headers": {
      "host": "bar.baz",
      "user-agent": "browser"
    },
    "host": "foo.bar",
    "id": 6789,
    "method": "GET",
    "path": "api/check"
  },
  "res": {
    "status_code": 200,
    "header": {
      "content-type": "text",
      "content-encoding": "plain"
    }
  }
}
</pre>

## Simple sorting of maps across records

As discussed in the page on
[operating on all records](operating-on-all-records.md), while Miller is normally
[streaming](streaming-and-memory.md) (we operate on one record at a time), we
can accumulate records in an array-valued or map-valued
[out-of-stream variable](reference-dsl-variables.md#out-of-stream-variables),
then operate on that record-list in an `end` block. This includes the possibility
of accumulating records in a map, then sorting the map.

Using the `f` flag we're sorting the map keys (1-up NR) lexically, so we
have 1, then 10, then 2:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv put -q '</b>
<b>  begin {</b>
<b>    @records = {};  # Define as a map</b>
<b>  }</b>
<b>  $nr = NR;</b>
<b>  @records[NR] = $*; # Accumulate</b>
<b>  end {</b>
<b>    @records = sort(@records, "f");</b>
<b>    for (_, record in @records) {</b>
<b>      emit1 record;</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   nr
yellow triangle true  1  11    43.6498  9.8870 1
purple square   false 10 91    72.3735  8.2430 10
red    square   true  2  15    79.2778  0.0130 2
red    circle   true  3  16    13.8103  2.9010 3
red    square   false 4  48    77.5542  7.4670 4
purple triangle false 5  51    81.2290  8.5910 5
red    square   false 6  64    77.1991  9.5310 6
purple triangle false 7  65    80.1405  5.8240 7
yellow circle   true  8  73    63.9785  4.2370 8
yellow circle   true  9  87    63.5058  8.3350 9
</pre>

## Custom sorting of arrays within records

Using the [`sort`](reference-dsl-builtin-functions.md#sort) function, you
can sort an array by its values, using another function (which you specify --
see the [page on user-defined functions](reference-dsl-user-defined-functions.md))
for comparing elements.

* Your function must take two arguments, which will range over various pairs of values in your array;
* It must return a number which is negative, zero, or positive depending on whether you want the first argument to sort less than, equal to, or greater than the second, respectively.

For example, let's use the following input data. Instead of having an array, it
has some semicolon-delimited data in a field which we can split and sort:

<pre class="pre-highlight-in-pair">
<b>cat data/sortaf-example.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
key,values
alpha,5;2;8;6;1;4;9;10;3;7
</pre>

In the following example we sort data in several ways -- the first two just
recaptiulate (for reference) what `sort` with default flags already does; the third is novel:

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/sortaf-example.csv put '</b>
<b></b>
<b>  # Same as sort($values)</b>
<b>  func forward(a,b) {</b>
<b>    return a <=> b</b>
<b>  }</b>
<b></b>
<b>  # Same as sort($values, "r")</b>
<b>  func reverse(a,b) {</b>
<b>    return b <=> a</b>
<b>  }</b>
<b></b>
<b>  # Custom sort</b>
<b>  func even_then_odd(a,b) {</b>
<b>    ax = a % 2;</b>
<b>    bx = b % 2;</b>
<b>    if (ax == bx) {</b>
<b>      return a <=> b</b>
<b>    } elif (bx == 1) {</b>
<b>      return -1</b>
<b>    } else {</b>
<b>      return 1</b>
<b>    }</b>
<b>  }</b>
<b></b>
<b>  split_values = splita($values, ";");</b>
<b>  $forward = sort(split_values, forward);</b>
<b>  $reverse = sort(split_values, reverse);</b>
<b>  $even_then_odd = sort(split_values, even_then_odd);</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "key": "alpha",
  "values": "5;2;8;6;1;4;9;10;3;7",
  "forward": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
  "reverse": [10, 9, 8, 7, 6, 5, 4, 3, 2, 1],
  "even_then_odd": [2, 4, 6, 8, 10, 1, 3, 5, 7, 9]
}
</pre>

## Custom sorting of arrays across records

As noted above, we can use the
[operating-on-all-records](operating-on-all-records.md) paradigm
to accumulate records in an array-valued or map-valued
[out-of-stream variable](reference-dsl-variables.md#out-of-stream-variables),
then operate on that record-list in an `end` block. This includes the possibility
of accumulating records in an array, then sorting the array.

Note that here the array elements are maps, so the `a` and `b` arguments to our
functions are maps -- and we have to access the `index` field using either
`a["index"]` and `b["index"]`, or (using the [dot operator for
indexing](reference-dsl-operators.md#the-double-purpose-dot-operator))
`a.index` and `b.index`.

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv put -q '</b>
<b>  # Sort primarily ascending on the shape field, then secondarily</b>
<b>  # descending numeric on the index field.</b>
<b>  func cmp(a, b) {</b>
<b>    cmp1 = a.shape <=> b.shape;</b>
<b>    if (cmp1 != 0) {</b>
<b>      return cmp1</b>
<b>    } else {</b>
<b>      return b.index <=> a.index;</b>
<b>    }</b>
<b>  }</b>
<b>  begin {</b>
<b>    @records = [];  # Define as an array, else auto-create will make a map</b>
<b>  }</b>
<b>  @records[NR] = $*; # Accumulate</b>
<b>  end {</b>
<b>    @records = sort(@records, cmp);</b>
<b>    for (record in @records) {</b>
<b>      emit1 record;</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
yellow circle   true  9  87    63.5058  8.3350
yellow circle   true  8  73    63.9785  4.2370
red    circle   true  3  16    13.8103  2.9010
purple square   false 10 91    72.3735  8.2430
red    square   false 6  64    77.1991  9.5310
red    square   false 4  48    77.5542  7.4670
red    square   true  2  15    79.2778  0.0130
purple triangle false 7  65    80.1405  5.8240
purple triangle false 5  51    81.2290  8.5910
yellow triangle true  1  11    43.6498  9.8870
</pre>

## Custom sorting of maps within records

Using the [`sort`](reference-dsl-builtin-functions.md#sort) function, you
can sort a map using a function which you specify (see the [page on
user-defined functions](reference-dsl-user-defined-functions.md)) for comparing
keys and/or values.

* Your function must take four arguments, which will range over various pairs of key-value pairs in your map;
* It must return a number which is negative, zero, or positive depending on whether you want the first argument to sort less than, equal to, or greater than the second, respectively.

For example, we can sort ascending or descending by map key or map value:

<pre class="pre-highlight-in-pair">
<b>mlr -n put -q '</b>
<b>  func f1(ak, av, bk, bv) {</b>
<b>    return ak <=> bk</b>
<b>  }</b>
<b>  func f2(ak, av, bk, bv) {</b>
<b>    return bk <=> ak</b>
<b>  }</b>
<b>  func f3(ak, av, bk, bv) {</b>
<b>    return av <=> bv</b>
<b>  }</b>
<b>  func f4(ak, av, bk, bv) {</b>
<b>    return bv <=> av</b>
<b>  }</b>
<b>  end {</b>
<b>    x = {</b>
<b>      "c":1,</b>
<b>      "a":3,</b>
<b>      "b":2,</b>
<b>    };</b>
<b></b>
<b>    print sort(x, f1);</b>
<b>    print sort(x, f2);</b>
<b>    print sort(x, f3);</b>
<b>    print sort(x, f4);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "a": 3,
  "b": 2,
  "c": 1
}
{
  "c": 1,
  "b": 2,
  "a": 3
}
{
  "c": 1,
  "b": 2,
  "a": 3
}
{
  "a": 3,
  "b": 2,
  "c": 1
}
</pre>

## Custom sorting of maps across records

We can modify our above example just a bit, where we accumulate records in a map rather than
an array. Here the map keys will be `NR` values `"1"`, `"2"`, etc.

Why would we do this? When we're operating across all records and keeping all
of them -- densely -- accumulating them in an array is fine. If we're only
taking a subset -- sparsely -- and we want to retain the original `NR` as keys,
using a map is handy, since we don't need continguous keys.

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv put -q '</b>
<b>  # Sort descending numeric on the index field</b>
<b>  func cmp(ak, av, bk, bv) {</b>
<b>    return bv.index <=> av.index</b>
<b>  }</b>
<b>  begin {</b>
<b>    @records = {};  # Define as a map</b>
<b>  }</b>
<b>  @records[NR] = $*; # Accumulate</b>
<b>  end {</b>
<b>    @records = sort(@records, cmp);</b>
<b>    for (_, record in @records) {</b>
<b>      emit1 record;</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
purple square   false 10 91    72.3735  8.2430
yellow circle   true  9  87    63.5058  8.3350
yellow circle   true  8  73    63.9785  4.2370
purple triangle false 7  65    80.1405  5.8240
red    square   false 6  64    77.1991  9.5310
purple triangle false 5  51    81.2290  8.5910
red    square   false 4  48    77.5542  7.4670
red    circle   true  3  16    13.8103  2.9010
red    square   true  2  15    79.2778  0.0130
yellow triangle true  1  11    43.6498  9.8870
</pre>
