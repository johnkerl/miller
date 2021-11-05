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
# Arithmetic

## Input scanning

Numbers in Miller are double-precision float or 64-bit signed integers. Anything scannable as int, e.g `123` or `0xabcd`, is treated as an integer; otherwise, input scannable as float (`4.56` or `8e9`) is treated as float; everything else is a string.

If you want all numbers to be treated as floats, then you may use `float()` in your filter/put expressions (e.g. replacing `$c = $a * $b` with `$c = float($a) * float($b)`).

<!--- TODO: probably remove this entirely for Miller6.
Or, more simply, use `mlr filter -F` and `mlr put -F` which forces all numeric input, whether from expression literals or field values, to float. Likewise `mlr stats1 -F` and `mlr step -F` force integerable accumulators (such as `count`) to be done in floating-point.
-->

## Conversion by math routines

For most math functions, integers are cast to float on input, and produce float output: e.g. `exp(0) = 1.0` rather than `1`.  The following, however, produce integer output if their inputs are integers: `+` `-` `*` `/` `//` `%` `abs` `ceil` `floor` `max` `min` `round` `roundm` `sgn`. As well, `stats1 -a min`, `stats1 -a max`, `stats1 -a sum`, `step -a delta`, and `step -a rsum` produce integer output if their inputs are integers.

## Conversion by arithmetic operators

The sum, difference, and product of integers is again integer, except for when that would overflow a 64-bit integer at which point Miller converts the result to float.

The short of it is that Miller does this transparently for you so you needn't think about it.

Implementation details of this, for the interested: integer adds and subtracts overflow by at most one bit so it suffices to check sign-changes. Thus, Miller allows you to add and subtract arbitrary 64-bit signed integers, converting only to float precisely when the result is less than -2\*\*63 or greater than 2\*\*63 - 1.  Multiplies, on the other hand, can overflow by a word size and a sign-change technique does not suffice to detect overflow. Instead, Miller tests whether the floating-point product exceeds the representable integer range. Now, 64-bit integers have 64-bit precision while IEEE-doubles have only 52-bit mantissas -- so, there are 53 bits including implicit leading one.  The following experiment explicitly demonstrates the resolution at this range:

<pre class="pre-non-highlight-non-pair">
64-bit integer     64-bit integer     Casted to double           Back to 64-bit
in hex             in decimal                                    integer
0x7ffffffffffff9ff 9223372036854774271 9223372036854773760.000000 0x7ffffffffffff800
0x7ffffffffffffa00 9223372036854774272 9223372036854773760.000000 0x7ffffffffffff800
0x7ffffffffffffbff 9223372036854774783 9223372036854774784.000000 0x7ffffffffffffc00
0x7ffffffffffffc00 9223372036854774784 9223372036854774784.000000 0x7ffffffffffffc00
0x7ffffffffffffdff 9223372036854775295 9223372036854774784.000000 0x7ffffffffffffc00
0x7ffffffffffffe00 9223372036854775296 9223372036854775808.000000 0x8000000000000000
0x7ffffffffffffffe 9223372036854775806 9223372036854775808.000000 0x8000000000000000
0x7fffffffffffffff 9223372036854775807 9223372036854775808.000000 0x8000000000000000
</pre>

That is, one cannot check an integer product to see if it is precisely greater than 2\*\*63 - 1 or less than -2\*\*63 using either integer arithmetic (it may have already overflowed) or using double-precision (due to granularity).  Instead, Miller checks for overflow in 64-bit integer multiplication by seeing whether the absolute value of the double-precision product exceeds the largest representable IEEE double less than 2\*\*63, which we see from the listing above is 9223372036854774784. (An alternative would be to do all integer multiplies using handcrafted multi-word 128-bit arithmetic.  This approach is not taken.)

## Pythonic division

Division and remainder are [pythonic](http://python-history.blogspot.com/2010/08/why-pythons-integer-division-floors.html):

* Quotient of integers is floating-point unless (unlike Python) exactly representable as integer: `7/2` is `3.5` but `6/2` is `3` (not 3.0).
* Integer division is done with `//`: `7//2` is `3`.  This rounds toward the negative.
* Remainders are non-negative.
