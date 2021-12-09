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
# DSL user-defined functions

As of Miller 5.0.0 you can define your own functions, as well as subroutines.

## User-defined functions

Here's the obligatory example of a recursive function to compute the factorial function:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --from data/small put '</b>
<b>    func f(n) {</b>
<b>        if (is_numeric(n)) {</b>
<b>            if (n > 0) {</b>
<b>                return n * f(n-1);</b>
<b>            } else {</b>
<b>                return 1;</b>
<b>            }</b>
<b>        }</b>
<b>        # implicitly return absent-null if non-numeric</b>
<b>    }</b>
<b>    $ox = f($x + NR);</b>
<b>    $oi = f($i);</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a   b   i x        y        ox                 oi
pan pan 1 0.346791 0.726802 0.4670549976810001 1
eks pan 2 0.758679 0.522151 3.6808304227112796 2
wye wye 3 0.204603 0.338318 1.7412477437471126 6
eks wye 4 0.381399 0.134188 18.588317372151177 24
wye pan 5 0.573288 0.863624 211.38663947090302 120
</pre>

Properties of user-defined functions:

* Function bodies start with `func` and a parameter list, defined outside of `begin`, `end`, or other `func` or `subr` blocks. (I.e. the Miller DSL has no nested functions.)

* A function (uniqified by its name) may not be redefined: either by redefining a user-defined function, or by redefining a built-in function. However, functions and subroutines have separate namespaces: you can define a subroutine `log` (for logging messages to stderr, say) which does not clash with the mathematical `log` (logarithm) function.

* Functions may be defined either before or after use -- there is an object-binding/linkage step at startup.  More specifically, functions may be either recursive or mutually recursive.

* Functions may be defined and called either within `mlr filter` or `mlr put`.

* Argument values may be reassigned: they are not read-only.

* When a return value is not implicitly returned, this results in a return value of [absent-null](reference-main-null-data.md). (In the example above, if there were records for which the argument to `f` is non-numeric, the assignments would be skipped.) See also the [null-data reference page](reference-main-null-data.md).

* See the section on [Local variables](reference-dsl-variables.md#local-variables) for information on scope and extent of arguments, as well as for information on the use of local variables within functions.

* See the section on [Expressions from files](reference-dsl-syntax.md#expressions-from-files) for information on the use of `-f` and `-e` flags.

## User-defined subroutines

Example:

<pre class="pre-highlight-in-pair">
<b>mlr --opprint --from data/small put -q '</b>
<b>  begin {</b>
<b>    @call_count = 0;</b>
<b>  }</b>
<b>  subr s(n) {</b>
<b>    @call_count += 1;</b>
<b>    if (is_numeric(n)) {</b>
<b>      if (n > 1) {</b>
<b>        call s(n-1);</b>
<b>      } else {</b>
<b>        print "numcalls=" . @call_count;</b>
<b>      }</b>
<b>    }</b>
<b>  }</b>
<b>  print "NR=" . NR;</b>
<b>  call s(NR);</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
NR=1
numcalls=1
NR=2
numcalls=3
NR=3
numcalls=6
NR=4
numcalls=10
NR=5
numcalls=15
</pre>

Properties of user-defined subroutines:

* Subroutine bodies start with `subr` and a parameter list, defined outside of `begin`, `end`, or other `func` or `subr` blocks. (I.e. the Miller DSL has no nested subroutines.)

* A subroutine (uniqified by its name) may not be redefined. However, functions and subroutines have separate namespaces: you can define a subroutine `log` which does not clash with the mathematical `log` function.

* Subroutines may be defined either before or after use -- there is an object-binding/linkage step at startup.  More specifically, subroutines may be either recursive or mutually recursive. Subroutines may call functions.

* Subroutines may be defined and called either within `mlr put` or `mlr put`.

* Subroutines have read/write access to `$`-variables and `@`-variables.

* Argument values may be reassigned: they are not read-only.

* See the section on [local variables](reference-dsl-variables.md#local-variables) for information on scope and extent of arguments, as well as for information on the use of local variables within functions.

* See the section on [Expressions from files](reference-dsl-syntax.md#expressions-from-files) for information on the use of `-f` and `-e` flags.

## Differences between functions and subroutines

Subroutines cannot return values, and they are invoked by the keyword `call`.

In hindsight, subroutines needn't have been invented. If `foo` is a function
then you can write `foo(1,2,3)` while ignoring its return value, and that plays
the role of subroutine quite well.

## Loading a library of functions

If you have a file with UDFs you use frequently, say `my-udfs.mlr`, you can use
`--load` or `--mload` to define them for your Miller scripts. For example, in
your shell, 

<pre class="pre-highlight-non-pair">
<b>alias mlr='mlr --load ~/my-functions.mlr'</b>
</pre>

or

<pre class="pre-highlight-non-pair">
<b>alias mlr='mlr --load /u/miller-udfs/'</b>
</pre>

See the [miscellaneous-flags page](reference-main-flag-list.md#miscellaneous-flags) for more information.

## Function literals

You can define unmnamed functions and assign the to variables, or pass them to functions.

See also the [page on higher-order functions](reference-dsl-higher-order-functions.md)
for more information on
[`select`](reference-dsl-builtin-functions.md#select),
[`apply`](reference-dsl-builtin-functions.md#apply),
[`reduce`](reference-dsl-builtin-functions.md#reduce),
[`fold`](reference-dsl-builtin-functions.md#fold), and sort.
[`sort`](reference-dsl-builtin-functions.md#sort),

For example:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv put '</b>
<b>  f = func(s, t) {</b>
<b>    return s . ":" . t;</b>
<b>  };</b>
<b>  $z = f($color, $shape);</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   z
yellow triangle true  1  11    43.6498  9.8870 yellow:triangle
red    square   true  2  15    79.2778  0.0130 red:square
red    circle   true  3  16    13.8103  2.9010 red:circle
red    square   false 4  48    77.5542  7.4670 red:square
purple triangle false 5  51    81.2290  8.5910 purple:triangle
red    square   false 6  64    77.1991  9.5310 red:square
purple triangle false 7  65    80.1405  5.8240 purple:triangle
yellow circle   true  8  73    63.9785  4.2370 yellow:circle
yellow circle   true  9  87    63.5058  8.3350 yellow:circle
purple square   false 10 91    72.3735  8.2430 purple:square
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv put '</b>
<b>  a = func(s, t) {</b>
<b>    return s . ":" . t . " above";</b>
<b>  };</b>
<b>  b = func(s, t) {</b>
<b>    return s . ":" . t . " below";</b>
<b>  };</b>
<b>  f = $index >= 50 ? a : b;</b>
<b>  $z = f($color, $shape);</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   z
yellow triangle true  1  11    43.6498  9.8870 yellow:triangle below
red    square   true  2  15    79.2778  0.0130 red:square below
red    circle   true  3  16    13.8103  2.9010 red:circle below
red    square   false 4  48    77.5542  7.4670 red:square below
purple triangle false 5  51    81.2290  8.5910 purple:triangle above
red    square   false 6  64    77.1991  9.5310 red:square above
purple triangle false 7  65    80.1405  5.8240 purple:triangle above
yellow circle   true  8  73    63.9785  4.2370 yellow:circle above
yellow circle   true  9  87    63.5058  8.3350 yellow:circle above
purple square   false 10 91    72.3735  8.2430 purple:square above
</pre>

Note that you need a semicolon after the closing curly brace of the function literal.

Unlike named functions, function literals (also known as unnamed functions)
have access to local variables defined in their enclosing scope. That's
so you can do things like this:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv put '</b>
<b>  f = func(s, t, i) {</b>
<b>    if (i >= cap) {</b>
<b>      return s . ":" . t . " above";</b>
<b>    } else {</b>
<b>      return s . ":" . t . " below";</b>
<b>    }</b>
<b>  };</b>
<b>  cap = 10;</b>
<b>  $z = f($color, $shape, $index);</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   z
yellow triangle true  1  11    43.6498  9.8870 yellow:triangle above
red    square   true  2  15    79.2778  0.0130 red:square above
red    circle   true  3  16    13.8103  2.9010 red:circle above
red    square   false 4  48    77.5542  7.4670 red:square above
purple triangle false 5  51    81.2290  8.5910 purple:triangle above
red    square   false 6  64    77.1991  9.5310 red:square above
purple triangle false 7  65    80.1405  5.8240 purple:triangle above
yellow circle   true  8  73    63.9785  4.2370 yellow:circle above
yellow circle   true  9  87    63.5058  8.3350 yellow:circle above
purple square   false 10 91    72.3735  8.2430 purple:square above
</pre>

See the [page on higher-order functions](reference-dsl-higher-order-functions.md) for more.
