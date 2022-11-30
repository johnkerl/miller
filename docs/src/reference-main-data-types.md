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
# Data types

## List of types

Miller's types are:

* Scalars:
    * **string**: such as `"abcdefg"`, supporting concatenation, one-up indexing and slicing, and [library functions](reference-dsl-builtin-functions.md#string-functions). See the pages on [strings](reference-main-strings.md) and [regular expressions](reference-main-regular-expressions.md).
    * **float** and **int**: such as `1.2` and `3`: double-precision and 64-bit signed, respectively. See the section on [arithmetic operators and math-related library functions](reference-dsl-builtin-functions.md#math-functions) as well as the [Arithmetic](reference-main-arithmetic.md) page.
    * dates/times are not a separate data type; Miller uses ints for [seconds since the epoch](https://en.wikipedia.org/wiki/Unix_time) and strings for formatted date/times. See the [DSL datetime/timezone functions page](reference-dsl-time.md) for more information.
    * **boolean**: literals `true` and `false`; results of `==`, `<`, `>`, etc. See the section on [boolean operators](reference-dsl-builtin-functions.md#boolean-functions).
* Collections:
    * **map**: such as `{"a":1,"b":[2,3,4]}`, supporting key-indexing, preservation of insertion order, [library functions](reference-dsl-builtin-functions.md#collections-functions), etc. See the [Maps](reference-main-maps.md) page.
    * **array**: such as `["a", 2, true]`, supporting one-up indexing and slicing, [library functions](reference-dsl-builtin-functions.md#collections-functions), etc. See the [Arrays](reference-main-arrays.md) page.
* Nulls and error:
    * **absent-null**: Such as on reads of unset right-hand sides, or fall-through non-explicit return values from user-defined functions. See the [null-data page](reference-main-null-data.md).
    * **JSON-null**: For `null` in JSON files; also used in [gapped auto-extend of arrays](reference-main-arrays.md#auto-extend-and-null-gaps). See the [null-data page](reference-main-null-data.md).
    * **error** -- for various results which cannot be computed, often when the input to a [built-in function](reference-dsl-builtin-functions.md) is of the wrong type. For example, doing [strlen](reference-dsl-builtin-functions.md#strlen) or [substr](reference-dsl-builtin-functions.md#substr) on a non-string, [sec2gmt](reference-dsl-builtin-functions.md#sec2gmt) on a non-integer, etc.
* Functions:
    * As described in the [page on function literals](reference-dsl-user-defined-functions.md#function-literals), you can define unnamed functions and assign them to variables, or pass them to functions.
    * These can also be (named) [user-defined functions](reference-dsl-user-defined-functions.md).
    * Use-cases include [custom sorting](reference-dsl-builtin-functions.md#sort), along with higher-order-functions such as [`select`](reference-dsl-builtin-functions.md#select), [`apply`](reference-dsl-builtin-functions.md#apply), [`reduce`](reference-dsl-builtin-functions.md#reduce), and [`fold`](reference-dsl-builtin-functions.md#fold).

See also the list of
[type-checking functions](reference-dsl-builtin-functions.md#type-checkin -functions) for the
[Miller programming language](miller-programming-language.md).

See also [Differences from other programming languages](reference-dsl-differences.md).

## Type inference for literal and record data

Miller's input and output are all text-oriented: all the
[file formats supported by Miller](file-formats.md) are human-readable text,
such as CSV, TSV, and JSON; binary formats such as
[BSON](https://bsonspec.org/) and [Parquet](https://parquet.apache.org/) are
not supported (as of mid-2021). In this sense, everything is a string in and out of
Miller -- be it in data files, or in DSL expressions you key in.

In the [DSL](miller-programming-language.md), `7` is an `int` and `8.9` is a float, as
one would expect.  Likewise, on input from [data files](file-formats.md),
string values representable as numbers, e.g. `1.2` or `3`, are treated as int
or float, respectively. If a record has `x=1,y=2` then `mlr put '$z=$x+$y'`
will produce `x=1,y=2,z=3`.

Numbers retain their original string representation, so if `x` is `1.2` on one
record and `1.200` on another, they'll print out that way on output (unless of
course they've been modified during processing, e.g. `mlr put '$x = $x + 10`).

Generally strings, numbers, and booleans don't mix; use type-casting like
`string($x)` to convert. However, the dot (string-concatenation) operator has
been special-cased: `mlr put '$z=$x.$y'` does not give an error, because the
dot operator has been generalized to stringify non-strings

Examples:

<pre class="pre-highlight-in-pair">
<b>mlr --csv cat data/type-infer.csv</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a,b,c
1.2,3,true
4,5.6,buongiorno
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --oxtab --from data/type-infer.csv put '</b>
<b>  $d = $a . $c;</b>
<b>  $e = 7;</b>
<b>  $f = 8.9;</b>
<b>  $g = $e + $f;</b>
<b>  $ta = typeof($a);</b>
<b>  $tb = typeof($b);</b>
<b>  $tc = typeof($c);</b>
<b>  $td = typeof($d);</b>
<b>  $te = typeof($e);</b>
<b>  $tf = typeof($f);</b>
<b>  $tg = typeof($g);</b>
<b>' then reorder -f a,ta,b,tb,c,tc,d,td,e,te,f,tf,g,tg</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a  1.2
ta float
b  3
tb int
c  true
tc string
d  1.2true
td string
e  7
te int
f  8.9
tf float
g  15.9
tg float

a  4
ta int
b  5.6
tb float
c  buongiorno
tc string
d  4buongiorno
td string
e  7
te int
f  8.9
tf float
g  15.9
tg float
</pre>

On input, string values representable as boolean  (e.g. `"true"`, `"false"`)
are *not* automatically treated as boolean.  This is because `"true"` and
`"false"` are ordinary words, and auto string-to-boolean on a column consisting
of words would result in some strings mixed with some booleans. Use the
`boolean` function to coerce: e.g. giving the record `x=1,y=2,w=false` to `mlr
filter '$z=($x<$y) || boolean($w)'`.

The same is true for `inf`, `+inf`, `-inf`, `infinity`, `+infinity`,
`-infinity`, `NaN`, and all upper-cased/lower-cased/mixed-case variants of
those. These are valid IEEE floating-point numbers, but Miller treats these as
strings. You can explicit force conversion: if `x=infinity` in a data file,
then `typeof($x)` is `string` but `typeof(float($x))` is `float`.

## JSON parse and stringify

If you have, say, a CSV file whose columns contain strings which are well-formatted JSON,
they will not be auto-converted, but you can use the
[`json-parse` verb](reference-verbs.md#json-parse)
or the
[`json_parse` DSL function](reference-dsl-builtin-functions.md#json_parse):

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data/json-in-csv.csv cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
id,blob
100,"{""a"":1,""b"":[2,3,4]}"
105,"{""a"":6,""b"":[7,8,9]}"
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/json-in-csv.csv cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "id": 100,
  "blob": "{\"a\":1,\"b\":[2,3,4]}"
},
{
  "id": 105,
  "blob": "{\"a\":6,\"b\":[7,8,9]}"
}
]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/json-in-csv.csv json-parse -f blob</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "id": 100,
  "blob": {
    "a": 1,
    "b": [2, 3, 4]
  }
},
{
  "id": 105,
  "blob": {
    "a": 6,
    "b": [7, 8, 9]
  }
}
]
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --ojson --from data/json-in-csv.csv put '$blob = json_parse($blob)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
{
  "id": 100,
  "blob": {
    "a": 1,
    "b": [2, 3, 4]
  }
},
{
  "id": 105,
  "blob": {
    "a": 6,
    "b": [7, 8, 9]
  }
}
]
</pre>

These have their respective operations to convert back to string: the
[`json-stringify` verb](reference-verbs.md#json-stringify)
and
[`json_stringify` DSL function](reference-dsl-builtin-functions.md#json_stringify).
