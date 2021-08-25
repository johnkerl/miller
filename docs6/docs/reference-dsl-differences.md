<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# Differences from other programming languages

The Miller programming language is intended to be straightforward and familiar,
as well as [not overly complex](reference-dsl-complexity.md). It doesn't try to
break new ground in terms of syntax; there are no classes or closures, and so
on.

While the [Principle of Least
Surprise](https://en.wikipedia.org/wiki/Principle_of_least_astonishment) is
often held to, nonetheless the following may be surprising.

## No ++ or --

There is no `++` or `--` [operator](reference-dsl-operators.md). To increment
`x`, use `x = x+1` or `x += 1`, and similarly for decrement.

## Semicolons as delimiters

You don't need a semicolon to end expressions, only to separate them. This
was done intentionally from the very start of Miller: you should be able to do
simple things like `mlr put '$z = $x * $y' myfile.dat` without needing a
semicolon.

Note that since you also don't need a semicolon before or after closing curly
braces (such as `begin`/`end` blocks, `if`-statements, `for`-loops, etc.) it's
easy to key in a few semicolon-free statements, and then to forget a
semicolon where one is needed . The parser tries to remind you about semicolons
whenever there's a chance a missing semicolon might be involved in a parse
error.

<pre class="pre-highlight-non-pair">
<b>mlr --csv --from example.csv put -q '</b>
<b>  begin {</b>
<b>    @count = 0 # No semicolon required -- before closing curly brace</b>
<b>  }</b>
<b>  $x=1         # No semicolon required -- at end of expression</b>
<b>'</b>
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from example.csv put -q '</b>
<b>  begin {</b>
<b>    @count = 0 # No semicolon required -- before closing curly brace</b>
<b>  }</b>
<b>  $x=1         # Needs a semicolon after it</b>
<b>  $y=2         # No semicolon required -- at end of expression</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
mlr: cannot parse DSL expression.
Parse error on token "$y" at line 6 columnn 3.
Please check for missing semicolon.
Expected one of:
  $ ; > >> | ? || ^^ && ?? ??? =~ !=~ == != >= < <= ^ & << >>> + - .+ .- .
  * / // % .* ./ .// ** [ [[ [[[
</pre>

## Required curly braces

Bodies for all compound statements must be enclosed in curly braces, even if the body is a single statement:

<pre class="pre-highlight-non-pair">
<b>mlr ... put 'if ($x == 1) $y = 2' # Syntax error</b>
</pre>

<pre class="pre-highlight-non-pair">
<b>mlr ... put 'if ($x == 1) { $y = 2 }' # This is OK</b>
</pre>

## No autoconvert to boolean

Boolean tests in `if`/`while`/`for`/etc must always take a boolean expression:
`if (1) {...}` results in the parse error
`Miller: conditional expression did not evaluate to boolean.`,
Likewise `if (x) {...}`, unless `x` is a variable of boolean type.
Please use `if (x != 0) {...}`, etc.

## Integer-preserving arithmetic

As discussed on the [arithmetic page](reference-main-arithmetic.md) the sum, difference, and product of two integers is again an integer, unless overflow occurs -- in which case Miller tries to convert to float in the least obtrusive way possible.

Likewise, while quotient and remainder are generally pythonic, the quotient and exponentiation of two integers is an integer when possible.

<pre class="pre-highlight-in-pair">
<b>$ mlr repl -q</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[mlr] 6/2
3

[mlr] typeof(6/2)
int

[mlr] 6/5
1.2

[mlr] typeof(6/5)
float

[mlr] typeof(7**8)
int

[mlr] typeof(7**80)
float
</pre>

## 1-up array indices

Arrays are indexed starting with 1, not 0. This is discussed in detail on the [arrays page](reference-dsl-arrays.md).

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data/short.csv cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
word,value
apple,37
ball,28
cat,54
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from data/short.csv put -q '</b>
<b>  @records[NR] = $*;</b>
<b>  end {</b>
<b>    for (i = 1; i <= NR; i += 1) {</b>
<b>      print "Record", i, "has word", @records[i]["word"];</b>
<b>    }</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
Record 1 has word apple
Record 2 has word ball
Record 3 has word cat
</pre>

## Print adds spaces around multiple arguments

As seen in the previous example,
[`print`](reference-dsl-output-statements.md#print-statements) with multiple
comma-delimited arguments fills in intervening spaces for you. If you want to
avoid this, use the dot operator for string-concatenation instead.

<pre class="pre-highlight-in-pair">
<b>mlr -n put -q '</b>
<b>  end {</b>
<b>    print "[", "a", "b", "c", "]";</b>
<b>    print "[" . "a" . "b" . "c" . "]";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[ a b c ]
[abc]
</pre>

Similarly, a final newline is printed for you; use [`printn`](reference-dsl-output-statements.md#print-statements) to avoid this.

## Insertion-order-preserving hashmaps

Miller's [maps](reference-dsl-maps.md) (as in many modern languages) preserve insertion order. If you set `x["foo"]=1` and then `x["bar"]=2`, then you are guaranteed that any looping over `x` will retrieve the `"foo"` key-value pair first, and the `"bar"` key-value pair second.

<pre class="pre-highlight-in-pair">
<b>mlr -n put -q 'end {</b>
<b>  x["foo"] = 1;</b>
<b>  x["bar"] = 2;</b>
<b>  dump x;</b>
<b>  for (k,v in x) {</b>
<b>    print "key", k, "value", v</b>
<b>  }</b>
<b>}'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "foo": 1,
  "bar": 2
}
key foo value 1
key bar value 2
</pre>

## Two-variable for-loops

Miller has a [key-value loop flavor](reference-dsl-control-structures.md#key-value-for-loops): whether `x` is a map or array, in `for (k,v in x) { ... }` the `k` will be bound to successive map keys (for maps) or 1-up array indices (for arrays), and the `v` will be bound to successive map values.

## Semantics for one-variable for-loops

Miller also has a [single-variable loop flavor](reference-dsl-control-structures.md#single-variable-for-loops). If `x` is a map then `for (e in x) { ... }` binds `e` to successive map _keys_ (not values as in PHP). But if `x` is an array then `for e in x) { ... }` binds `e` to successive array _values_ (not indices).

## Absent-null

Miller has a somewhat novel flavor of null data called _absent_: if a record
has a field `x` then `$y=$x` creates a field `y`, but if it doesn't then the assignment
is skipped. See the [null-data page](reference-main-null-data.md) for more
information.
