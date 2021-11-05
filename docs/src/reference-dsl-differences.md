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
  $ ; > >> | ? || ^^ && ?? ??? =~ !=~ == != <=> >= < <= ^ & << >>> + - .+
  .- * / // % .* ./ .// . ** [ [[ [[[

</pre>

## elif

Miller has [`elif`](reference-dsl-control-structures.md#if-statements), not `else if` or `elsif`.

## If-statement variable scoping

Miller is simple-minded about scoping [local variables](reference-dsl-variables.md#local-variables) to blocks.
If you have

<pre class="pre-non-highlight-non-pair">
  if (something) {
    x = 1
  } else {
    x = 2
  }
</pre>

then there are two `x` variable, each confined only to their enclosing curly
braces; there is no hoisting out of the `if` and `else` blocks.

A suggestion is

<pre class="pre-non-highlight-non-pair">
  var x
  if (something) {
    x = 1
  } else {
    x = 2
  }
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

I.e. ints/strings/etc are neither "truthy" nor "falsy".
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

## String literals with double quotes only

In some languages, like Ruby and Bash, string literals can be in single quotes or double quotes,
where single quotes suppress `\n` converting to a newline character and double quotes allowing it:
`'a\nb'` prints as the four characters `a`, `\`, `n`, and `b` on one line; `"a\nb"` prints as an
`a` on one line and a `b` on another.

In others, like Python and JavaScript, string literals can be in single quotes or double quotes,
interchangeably -- so you can have `"don't"` or `'the "right" thing'` as you wish.

In yet others, such as C/C++ and Java, string literals are in double auotes, like `"abc"`,
while single quotes are for character literals like `'a'` or `'\n'`. In these, if `s` is a non-empty string,
then `s[0]` is its first character.

In the [Miller programming language](miller-programming-language.md):

* String literals are always in double quotes, like `"abc"`.
* String-indexing/slicing always results in strings (even of length 1): `"abc"[1:1]` is the string `"a"`, and there is no notion in the Miller programming language of a character type.
* The single-quote character plays no role whatsoever in the grammar of the Miller programming language.
* Single quotes are reserved for wrapping expressions at the system command line. For example, in `mlr put '$message = "hello"' ...`, the [`put` verb](reference-dsl.md) gets the string `$message = "hello"`; the shell has consumed the outer single quotes by the time the Miller parser receives it.
* Things are a little different on Windows, where `"""` sequences are sometimes necessary: see the [Miller on Windows page](miller-on-windows.md).

## Absent-null

Miller has a somewhat novel flavor of null data called _absent_: if a record
has a field `x` then `$y=$x` creates a field `y`, but if it doesn't then the assignment
is skipped. See the [null-data page](reference-main-null-data.md) for more
information.

## Maps

See the [maps page](reference-main-maps.md).

## Arrays, including 1-up array indices

Arrays and strings are indexed starting with 1, not 0. This is discussed in
detail on the [arrays page](reference-main-arrays.md) and the [strings
page](reference-main-strings.md).

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

Also, slices for arrays and strings are _doubly inclusive_: `x[3:5]` gets you
elements 3, 4, and 5 of the array or string named `x`.

See the [arrays page](reference-main-arrays.md) for more about arrays; see the
[strings page](reference-main-strings.md) for more about strings.

## Two-variable for-loops

Miller has a [key-value loop flavor](reference-dsl-control-structures.md#key-value-for-loops): whether `x` is a map or array, in `for (k,v in x) { ... }` the `k` will be bound to successive map keys (for maps) or 1-up array indices (for arrays), and the `v` will be bound to successive map values.

## Semantics for one-variable for-loops

Miller also has a [single-variable loop flavor](reference-dsl-control-structures.md#single-variable-for-loops). If `x` is a map then `for (e in x) { ... }` binds `e` to successive map _keys_ (not values as in PHP). But if `x` is an array then `for e in x) { ... }` binds `e` to successive array _values_ (not indices).

## JSON parse, stringify, decode, and encode

Miller has the verbs
[`json-parse`](reference-verbs.md#json-parse) and
[`json-stringify`](reference-verbs.md#json-stringify), and the DSL functions
[`json_parse`](reference-dsl-builtin-functions.md#json_parse) and
[`json_stringify`](reference-dsl-builtin-functions.md#json_stringify).
In some other lannguages these are called `json_decode` and `json_encode`.
