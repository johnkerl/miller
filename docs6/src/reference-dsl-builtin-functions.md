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
# DSL built-in functions

These are functions in the [Miller programming language](miller-programming-language.md)
that you can call when you use `mlr put` and `mlr filter`. For example, when you type

<pre class="pre-highlight-in-pair">
<b>mlr --icsv --opprint --from example.csv put '</b>
<b>  $color = toupper($color);</b>
<b>  $shape = gsub($shape, "[aeiou]", "*");</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate
YELLOW tr**ngl* true  1  11    43.6498  9.8870
RED    sq**r*   true  2  15    79.2778  0.0130
RED    c*rcl*   true  3  16    13.8103  2.9010
RED    sq**r*   false 4  48    77.5542  7.4670
PURPLE tr**ngl* false 5  51    81.2290  8.5910
RED    sq**r*   false 6  64    77.1991  9.5310
PURPLE tr**ngl* false 7  65    80.1405  5.8240
YELLOW c*rcl*   true  8  73    63.9785  4.2370
YELLOW c*rcl*   true  9  87    63.5058  8.3350
PURPLE sq**r*   false 10 91    72.3735  8.2430
</pre>

the `toupper` and `gsub` bits are _functions_.

## Overview

At the command line, you can use `mlr -f` and `mlr -F` for information much
like what's on this page.

Each function takes a specific number of arguments, as shown below, except for
functions marked as variadic such as `min` and `max`. (The latter compute min
and max of any number of arguments.) There is no notion of optional or
default-on-absent arguments. All argument-passing is positional rather than by
name; arguments are passed by value, not by reference.

At the command line, you can get a list of all functions using `mlr -f`, with
details using `mlr -F`.  (Or, `mlr help usage-functions-by-class` to get
details in the order shown on this page.) You can get detail for a given
function using `mlr help function namegoeshere`, e.g.  `mlr help function
gsub`.

Operators are listed here along with functions. In this case, the
argument-count is the number of items involved in the infix operator, e.g. we
say `x+y` so the details for the `+` operator say that its number of arguments
is 2. Unary operators such as `!` and `~` show argument-count of 1; the ternary
`? :` operator shows an argument-count of 3.


## Functions by class

* [**Arithmetic functions**](#arithmetic-functions):  [bitcount](#bitcount),  [madd](#madd),  [mexp](#mexp),  [mmul](#mmul),  [msub](#msub),  [pow](#pow),  [%](#percent),  [&](#bitwise-and),  [\*](#times),  [\**](#exponentiation),  [\+](#plus),  [\-](#minus),  [\.\*](#dot-times),  [\.\+](#dot-plus),  [\.\-](#dot-minus),  [\./](#dot-slash),  [/](#slash),  [//](#slash-slash),  [<<](#lsh),  [>>](#srsh),  [>>>](#ursh),  [^](#bitwise-xor),  [\|](#bitwise-or),  [~](#bitwise-not).
* [**Boolean functions**](#boolean-functions):  [\!](#exclamation-point),  [\!=](#exclamation-point-equals),  [!=~](#regnotmatch),  [&&](#logical-and),  [<](#less-than),  [<=](#less-than-or-equals),  [<=>](#<=>),  [==](#double-equals),  [=~](#regmatch),  [>](#greater-than),  [>=](#greater-than-or-equals),  [?:](#question-mark-colon),  [??](#absent-coalesce),  [???](#absent-empty-coalesce),  [^^](#logical-xor),  [\|\|](#logical-or).
* [**Collections functions**](#collections-functions):  [append](#append),  [arrayify](#arrayify),  [depth](#depth),  [flatten](#flatten),  [get_keys](#get_keys),  [get_values](#get_values),  [haskey](#haskey),  [json_parse](#json_parse),  [json_stringify](#json_stringify),  [leafcount](#leafcount),  [length](#length),  [mapdiff](#mapdiff),  [mapexcept](#mapexcept),  [mapselect](#mapselect),  [mapsum](#mapsum),  [sorta](#sorta),  [sortaf](#sortaf),  [sortmf](#sortmf),  [sortmk](#sortmk),  [unflatten](#unflatten).
* [**Conversion functions**](#conversion-functions):  [boolean](#boolean),  [float](#float),  [fmtnum](#fmtnum),  [hexfmt](#hexfmt),  [int](#int),  [joink](#joink),  [joinkv](#joinkv),  [joinv](#joinv),  [splita](#splita),  [splitax](#splitax),  [splitkv](#splitkv),  [splitkvx](#splitkvx),  [splitnv](#splitnv),  [splitnvx](#splitnvx),  [string](#string).
* [**Hashing functions**](#hashing-functions):  [md5](#md5),  [sha1](#sha1),  [sha256](#sha256),  [sha512](#sha512).
* [**Math functions**](#math-functions):  [abs](#abs),  [acos](#acos),  [acosh](#acosh),  [asin](#asin),  [asinh](#asinh),  [atan](#atan),  [atan2](#atan2),  [atanh](#atanh),  [cbrt](#cbrt),  [ceil](#ceil),  [cos](#cos),  [cosh](#cosh),  [erf](#erf),  [erfc](#erfc),  [exp](#exp),  [expm1](#expm1),  [floor](#floor),  [invqnorm](#invqnorm),  [log](#log),  [log10](#log10),  [log1p](#log1p),  [logifit](#logifit),  [max](#max),  [min](#min),  [qnorm](#qnorm),  [round](#round),  [roundm](#roundm),  [sgn](#sgn),  [sin](#sin),  [sinh](#sinh),  [sqrt](#sqrt),  [tan](#tan),  [tanh](#tanh),  [urand](#urand),  [urand32](#urand32),  [urandint](#urandint),  [urandrange](#urandrange).
* [**String functions**](#string-functions):  [capitalize](#capitalize),  [clean_whitespace](#clean_whitespace),  [collapse_whitespace](#collapse_whitespace),  [gsub](#gsub),  [lstrip](#lstrip),  [regextract](#regextract),  [regextract_or_else](#regextract_or_else),  [rstrip](#rstrip),  [ssub](#ssub),  [strip](#strip),  [strlen](#strlen),  [sub](#sub),  [substr](#substr),  [substr0](#substr0),  [substr1](#substr1),  [tolower](#tolower),  [toupper](#toupper),  [truncate](#truncate),  [\.](#dot).
* [**System functions**](#system-functions):  [hostname](#hostname),  [os](#os),  [system](#system),  [version](#version).
* [**Time functions**](#time-functions):  [dhms2fsec](#dhms2fsec),  [dhms2sec](#dhms2sec),  [fsec2dhms](#fsec2dhms),  [fsec2hms](#fsec2hms),  [gmt2sec](#gmt2sec),  [hms2fsec](#hms2fsec),  [hms2sec](#hms2sec),  [sec2dhms](#sec2dhms),  [sec2gmt](#sec2gmt),  [sec2gmtdate](#sec2gmtdate),  [sec2hms](#sec2hms),  [strftime](#strftime),  [strptime](#strptime),  [systime](#systime),  [systimeint](#systimeint),  [uptime](#uptime).
* [**Typing functions**](#typing-functions):  [asserting_absent](#asserting_absent),  [asserting_array](#asserting_array),  [asserting_bool](#asserting_bool),  [asserting_boolean](#asserting_boolean),  [asserting_empty](#asserting_empty),  [asserting_empty_map](#asserting_empty_map),  [asserting_error](#asserting_error),  [asserting_float](#asserting_float),  [asserting_int](#asserting_int),  [asserting_map](#asserting_map),  [asserting_nonempty_map](#asserting_nonempty_map),  [asserting_not_array](#asserting_not_array),  [asserting_not_empty](#asserting_not_empty),  [asserting_not_map](#asserting_not_map),  [asserting_not_null](#asserting_not_null),  [asserting_null](#asserting_null),  [asserting_numeric](#asserting_numeric),  [asserting_present](#asserting_present),  [asserting_string](#asserting_string),  [is_absent](#is_absent),  [is_array](#is_array),  [is_bool](#is_bool),  [is_boolean](#is_boolean),  [is_empty](#is_empty),  [is_empty_map](#is_empty_map),  [is_error](#is_error),  [is_float](#is_float),  [is_int](#is_int),  [is_map](#is_map),  [is_nonempty_map](#is_nonempty_map),  [is_not_array](#is_not_array),  [is_not_empty](#is_not_empty),  [is_not_map](#is_not_map),  [is_not_null](#is_not_null),  [is_null](#is_null),  [is_numeric](#is_numeric),  [is_present](#is_present),  [is_string](#is_string),  [typeof](#typeof).

## Arithmetic functions


### bitcount
<pre class="pre-non-highlight-non-pair">
bitcount  (class=arithmetic #args=1) Count of 1-bits.
</pre>


### madd
<pre class="pre-non-highlight-non-pair">
madd  (class=arithmetic #args=3) a + b mod m (integers)
</pre>


### mexp
<pre class="pre-non-highlight-non-pair">
mexp  (class=arithmetic #args=3) a ** b mod m (integers)
</pre>


### mmul
<pre class="pre-non-highlight-non-pair">
mmul  (class=arithmetic #args=3) a * b mod m (integers)
</pre>


### msub
<pre class="pre-non-highlight-non-pair">
msub  (class=arithmetic #args=3) a - b mod m (integers)
</pre>


### pow
<pre class="pre-non-highlight-non-pair">
pow  (class=arithmetic #args=2) Exponentiation. Same as **, but as a function.
</pre>


<a id=percent />
### %
<pre class="pre-non-highlight-non-pair">
%  (class=arithmetic #args=2) Remainder; never negative-valued (pythonic).
</pre>


<a id=bitwise-and />
### &
<pre class="pre-non-highlight-non-pair">
&  (class=arithmetic #args=2) Bitwise AND.
</pre>


<a id=times />
### \*
<pre class="pre-non-highlight-non-pair">
*  (class=arithmetic #args=2) Multiplication, with integer*integer overflow to float.
</pre>


<a id=exponentiation />
### \**
<pre class="pre-non-highlight-non-pair">
**  (class=arithmetic #args=2) Exponentiation. Same as pow, but as an infix operator.
</pre>


<a id=plus />
### \+
<pre class="pre-non-highlight-non-pair">
+  (class=arithmetic #args=1,2) Addition as binary operator; unary plus operator.
</pre>


<a id=minus />
### \-
<pre class="pre-non-highlight-non-pair">
-  (class=arithmetic #args=1,2) Subtraction as binary operator; unary negation operator.
</pre>


<a id=dot-times />
### \.\*
<pre class="pre-non-highlight-non-pair">
.*  (class=arithmetic #args=2) Multiplication, with integer-to-integer overflow.
</pre>


<a id=dot-plus />
### \.\+
<pre class="pre-non-highlight-non-pair">
.+  (class=arithmetic #args=2) Addition, with integer-to-integer overflow.
</pre>


<a id=dot-minus />
### \.\-
<pre class="pre-non-highlight-non-pair">
.-  (class=arithmetic #args=2) Subtraction, with integer-to-integer overflow.
</pre>


<a id=dot-slash />
### \./
<pre class="pre-non-highlight-non-pair">
./  (class=arithmetic #args=2) Integer division; not pythonic.
</pre>


<a id=slash />
### /
<pre class="pre-non-highlight-non-pair">
/  (class=arithmetic #args=2) Division. Integer / integer is floating-point.
</pre>


<a id=slash-slash />
### //
<pre class="pre-non-highlight-non-pair">
//  (class=arithmetic #args=2) Pythonic integer division, rounding toward negative.
</pre>


<a id=lsh />
### <<
<pre class="pre-non-highlight-non-pair">
<<  (class=arithmetic #args=2) Bitwise left-shift.
</pre>


<a id=srsh />
### >>
<pre class="pre-non-highlight-non-pair">
>>  (class=arithmetic #args=2) Bitwise signed right-shift.
</pre>


<a id=ursh />
### >>>
<pre class="pre-non-highlight-non-pair">
>>>  (class=arithmetic #args=2) Bitwise unsigned right-shift.
</pre>


<a id=bitwise-xor />
### ^
<pre class="pre-non-highlight-non-pair">
^  (class=arithmetic #args=2) Bitwise XOR.
</pre>


<a id=bitwise-or />
### \|
<pre class="pre-non-highlight-non-pair">
|  (class=arithmetic #args=2) Bitwise OR.
</pre>


<a id=bitwise-not />
### ~
<pre class="pre-non-highlight-non-pair">
~  (class=arithmetic #args=1) Bitwise NOT. Beware '$y=~$x' since =~ is the
regex-match operator: try '$y = ~$x'.
</pre>

## Boolean functions


<a id=exclamation-point />
### \!
<pre class="pre-non-highlight-non-pair">
!  (class=boolean #args=1) Logical negation.
</pre>


<a id=exclamation-point-equals />
### \!=
<pre class="pre-non-highlight-non-pair">
!=  (class=boolean #args=2) String/numeric inequality. Mixing number and string results in string compare.
</pre>


<a id=regnotmatch />
### !=~
<pre class="pre-non-highlight-non-pair">
!=~  (class=boolean #args=2) String (left-hand side) does not match regex (right-hand side), e.g. '$name !=~ "^a.*b$"'.
</pre>


<a id=logical-and />
### &&
<pre class="pre-non-highlight-non-pair">
&&  (class=boolean #args=2) Logical AND.
</pre>


<a id=less-than />
### <
<pre class="pre-non-highlight-non-pair">
<  (class=boolean #args=2) String/numeric less-than. Mixing number and string results in string compare.
</pre>


<a id=less-than-or-equals />
### <=
<pre class="pre-non-highlight-non-pair">
<=  (class=boolean #args=2) String/numeric less-than-or-equals. Mixing number and string results in string compare.
</pre>


### <=>
<pre class="pre-non-highlight-non-pair">
<=>  (class=boolean #args=2) Comparator, nominally for sorting. Given a <=> b, returns <0, 0, >0 as a < b, a == b, or a > b, respectively.
</pre>


<a id=double-equals />
### ==
<pre class="pre-non-highlight-non-pair">
==  (class=boolean #args=2) String/numeric equality. Mixing number and string results in string compare.
</pre>


<a id=regmatch />
### =~
<pre class="pre-non-highlight-non-pair">
=~  (class=boolean #args=2) String (left-hand side) matches regex (right-hand side), e.g. '$name =~ "^a.*b$"'.
</pre>


<a id=greater-than />
### >
<pre class="pre-non-highlight-non-pair">
>  (class=boolean #args=2) String/numeric greater-than. Mixing number and string results in string compare.
</pre>


<a id=greater-than-or-equals />
### >=
<pre class="pre-non-highlight-non-pair">
>=  (class=boolean #args=2) String/numeric greater-than-or-equals. Mixing number and string results in string compare.
</pre>


<a id=question-mark-colon />
### ?:
<pre class="pre-non-highlight-non-pair">
?:  (class=boolean #args=3) Standard ternary operator.
</pre>


<a id=absent-coalesce />
### ??
<pre class="pre-non-highlight-non-pair">
??  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record.
</pre>


<a id=absent-empty-coalesce />
### ???
<pre class="pre-non-highlight-non-pair">
???  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record, or has empty value.
</pre>


<a id=logical-xor />
### ^^
<pre class="pre-non-highlight-non-pair">
^^  (class=boolean #args=2) Logical XOR.
</pre>


<a id=logical-or />
### \|\|
<pre class="pre-non-highlight-non-pair">
||  (class=boolean #args=2) Logical OR.
</pre>

## Collections functions


### append
<pre class="pre-non-highlight-non-pair">
append  (class=collections #args=2) Appends second argument to end of first argument, which must be an array.
</pre>


### arrayify
<pre class="pre-non-highlight-non-pair">
arrayify  (class=collections #args=1) Walks through a nested map/array, converting any map with consecutive keys
"1", "2", ... into an array. Useful to wrap the output of unflatten.
</pre>


### depth
<pre class="pre-non-highlight-non-pair">
depth  (class=collections #args=1) Prints maximum depth of map/array. Scalars have depth 0.
</pre>


### flatten
<pre class="pre-non-highlight-non-pair">
flatten  (class=collections #args=3) Flattens multi-level maps to single-level ones. Examples:
flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
Two-argument version: flatten($*, ".") is the same as flatten("", ".", $*).
Useful for nested JSON-like structures for non-JSON file formats like CSV.
</pre>


### get_keys
<pre class="pre-non-highlight-non-pair">
get_keys  (class=collections #args=1) Returns array of keys of map or array
</pre>


### get_values
<pre class="pre-non-highlight-non-pair">
get_values  (class=collections #args=1) Returns array of keys of map or array -- in the latter case, returns a copy of the array
</pre>


### haskey
<pre class="pre-non-highlight-non-pair">
haskey  (class=collections #args=2) True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or
'haskey(mymap, mykey)', or true/false if array index is in bounds / out of bounds.
Error if 1st argument is not a map or array. Note -n..-1 alias to 1..n in Miller arrays.
</pre>


### json_parse
<pre class="pre-non-highlight-non-pair">
json_parse  (class=collections #args=1) Converts value from JSON-formatted string.
</pre>


### json_stringify
<pre class="pre-non-highlight-non-pair">
json_stringify  (class=collections #args=1,2) Converts value to JSON-formatted string. Default output is single-line.
With optional second boolean argument set to true, produces multiline output.
</pre>


### leafcount
<pre class="pre-non-highlight-non-pair">
leafcount  (class=collections #args=1) Counts total number of terminal values in map/array. For single-level
map/array, same as length.
</pre>


### length
<pre class="pre-non-highlight-non-pair">
length  (class=collections #args=1) Counts number of top-level entries in array/map. Scalars have length 1.
</pre>


### mapdiff
<pre class="pre-non-highlight-non-pair">
mapdiff  (class=collections #args=variadic) With 0 args, returns empty map. With 1 arg, returns copy of arg.
With 2 or more, returns copy of arg 1 with all keys from any of remaining
argument maps removed.
</pre>


### mapexcept
<pre class="pre-non-highlight-non-pair">
mapexcept  (class=collections #args=variadic) Returns a map with keys from remaining arguments, if any, unset.
Remaining arguments can be strings or arrays of string.
E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'
and  'mapexcept({1:2,3:4,5:6}, [1, 5, 7])' is '{3:4}'.
</pre>


### mapselect
<pre class="pre-non-highlight-non-pair">
mapselect  (class=collections #args=variadic) Returns a map with only keys from remaining arguments set.
Remaining arguments can be strings or arrays of string.
E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'
and  'mapselect({1:2,3:4,5:6}, [1, 5, 7])' is '{1:2,5:6}'.
</pre>


### mapsum
<pre class="pre-non-highlight-non-pair">
mapsum  (class=collections #args=variadic) With 0 args, returns empty map. With >= 1 arg, returns a map with
key-value pairs from all arguments. Rightmost collisions win, e.g.
'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.
</pre>


### sorta
<pre class="pre-non-highlight-non-pair">
sorta  (class=collections #args=1-2) Returns a copy of an array, sorted ascending. Coming soon: other sort options.
</pre>


### sortaf
<pre class="pre-non-highlight-non-pair">
sortaf  (class=collections #args=2) Sorts an array (1st argument) using a comparator function you specify by name (2nd argument). The function is a comparator: it should take two arguments, returning a number <0, ==0, >0 as a<b, a==b, or a>b respectively. Example: 'sortaf([5,2,3,1,4], "f")'. Forward sort: 'func f(a,b) {return a <=> b}'. Reverse sort: 'func f(a,b) {return b <=> a}'. And so on -- you can implement logic you choose.
</pre>


### sortmf
<pre class="pre-non-highlight-non-pair">
sortmf  (class=collections #args=2) Sorts an array (1st argument) using a comparator function you specify by name (2nd argument). The function is a comparator: it should take four arguments, for one keyk, one value, other key, other value. It should return a number <0, ==0, >0 as a<b, a==b, or a>b respectively. Example: 'sortaf({"c":1,"b":3,"a":1}, "f")'. Forward sort by key: 'func f(ak,av,bk,bv) {return ak <=> bk}'. Reverse sort by key: 'func f(ak,av,bk,bv) {return bk <=> ak}'. And so on -- you can implement logic you choose.
</pre>


### sortmk
<pre class="pre-non-highlight-non-pair">
sortmk  (class=collections #args=1-2) Returns a copy of a map, sorted ascending by map key. Coming soon: other sort options.
</pre>


### unflatten
<pre class="pre-non-highlight-non-pair">
unflatten  (class=collections #args=2) Reverses flatten. Example:
unflatten({"a.b.c" : 4}, ".") is {"a": "b": { "c": 4 }}.
Useful for nested JSON-like structures for non-JSON file formats like CSV.
See also arrayify.
</pre>

## Conversion functions


### boolean
<pre class="pre-non-highlight-non-pair">
boolean  (class=conversion #args=1) Convert int/float/bool/string to boolean.
</pre>


### float
<pre class="pre-non-highlight-non-pair">
float  (class=conversion #args=1) Convert int/float/bool/string to float.
</pre>


### fmtnum
<pre class="pre-non-highlight-non-pair">
fmtnum  (class=conversion #args=2) Convert int/float/bool to string using
printf-style format string, e.g. '$s = fmtnum($n, "%06lld")'.
</pre>


### hexfmt
<pre class="pre-non-highlight-non-pair">
hexfmt  (class=conversion #args=1) Convert int to hex string, e.g. 255 to "0xff".
</pre>


### int
<pre class="pre-non-highlight-non-pair">
int  (class=conversion #args=1) Convert int/float/bool/string to int.
</pre>


### joink
<pre class="pre-non-highlight-non-pair">
joink  (class=conversion #args=2) Makes string from map/array keys. Examples:
joink({"a":3,"b":4,"c":5}, ",") = "a,b,c"
joink([1,2,3], ",") = "1,2,3".
</pre>


### joinkv
<pre class="pre-non-highlight-non-pair">
joinkv  (class=conversion #args=3) Makes string from map/array key-value pairs. Examples:
joinkv([3,4,5], "=", ",") = "1=3,2=4,3=5"
joinkv({"a":3,"b":4,"c":5}, "=", ",") = "a=3,b=4,c=5"
</pre>


### joinv
<pre class="pre-non-highlight-non-pair">
joinv  (class=conversion #args=2) Makes string from map/array values.
joinv([3,4,5], ",") = "3,4,5"
joinv({"a":3,"b":4,"c":5}, ",") = "3,4,5"
</pre>


### splita
<pre class="pre-non-highlight-non-pair">
splita  (class=conversion #args=2) Splits string into array with type inference. Example:
splita("3,4,5", ",") = [3,4,5]
</pre>


### splitax
<pre class="pre-non-highlight-non-pair">
splitax  (class=conversion #args=2) Splits string into array without type inference. Example:
splita("3,4,5", ",") = ["3","4","5"]
</pre>


### splitkv
<pre class="pre-non-highlight-non-pair">
splitkv  (class=conversion #args=3) Splits string by separators into map with type inference. Example:
splitkv("a=3,b=4,c=5", "=", ",") = {"a":3,"b":4,"c":5}
</pre>


### splitkvx
<pre class="pre-non-highlight-non-pair">
splitkvx  (class=conversion #args=3) Splits string by separators into map without type inference (keys and
values are strings). Example:
splitkvx("a=3,b=4,c=5", "=", ",") = {"a":"3","b":"4","c":"5"}
</pre>


### splitnv
<pre class="pre-non-highlight-non-pair">
splitnv  (class=conversion #args=2) Splits string by separator into integer-indexed map with type inference. Example:
splitnv("a,b,c", ",") = {"1":"a","2":"b","3":"c"}
</pre>


### splitnvx
<pre class="pre-non-highlight-non-pair">
splitnvx  (class=conversion #args=2) Splits string by separator into integer-indexed map without type
inference (values are strings). Example:
splitnvx("3,4,5", ",") = {"1":"3","2":"4","3":"5"}
</pre>


### string
<pre class="pre-non-highlight-non-pair">
string  (class=conversion #args=1) Convert int/float/bool/string/array/map to string.
</pre>

## Hashing functions


### md5
<pre class="pre-non-highlight-non-pair">
md5  (class=hashing #args=1) MD5 hash.
</pre>


### sha1
<pre class="pre-non-highlight-non-pair">
sha1  (class=hashing #args=1) SHA1 hash.
</pre>


### sha256
<pre class="pre-non-highlight-non-pair">
sha256  (class=hashing #args=1) SHA256 hash.
</pre>


### sha512
<pre class="pre-non-highlight-non-pair">
sha512  (class=hashing #args=1) SHA512 hash.
</pre>

## Math functions


### abs
<pre class="pre-non-highlight-non-pair">
abs  (class=math #args=1) Absolute value.
</pre>


### acos
<pre class="pre-non-highlight-non-pair">
acos  (class=math #args=1) Inverse trigonometric cosine.
</pre>


### acosh
<pre class="pre-non-highlight-non-pair">
acosh  (class=math #args=1) Inverse hyperbolic cosine.
</pre>


### asin
<pre class="pre-non-highlight-non-pair">
asin  (class=math #args=1) Inverse trigonometric sine.
</pre>


### asinh
<pre class="pre-non-highlight-non-pair">
asinh  (class=math #args=1) Inverse hyperbolic sine.
</pre>


### atan
<pre class="pre-non-highlight-non-pair">
atan  (class=math #args=1) One-argument arctangent.
</pre>


### atan2
<pre class="pre-non-highlight-non-pair">
atan2  (class=math #args=2) Two-argument arctangent.
</pre>


### atanh
<pre class="pre-non-highlight-non-pair">
atanh  (class=math #args=1) Inverse hyperbolic tangent.
</pre>


### cbrt
<pre class="pre-non-highlight-non-pair">
cbrt  (class=math #args=1) Cube root.
</pre>


### ceil
<pre class="pre-non-highlight-non-pair">
ceil  (class=math #args=1) Ceiling: nearest integer at or above.
</pre>


### cos
<pre class="pre-non-highlight-non-pair">
cos  (class=math #args=1) Trigonometric cosine.
</pre>


### cosh
<pre class="pre-non-highlight-non-pair">
cosh  (class=math #args=1) Hyperbolic cosine.
</pre>


### erf
<pre class="pre-non-highlight-non-pair">
erf  (class=math #args=1) Error function.
</pre>


### erfc
<pre class="pre-non-highlight-non-pair">
erfc  (class=math #args=1) Complementary error function.
</pre>


### exp
<pre class="pre-non-highlight-non-pair">
exp  (class=math #args=1) Exponential function e**x.
</pre>


### expm1
<pre class="pre-non-highlight-non-pair">
expm1  (class=math #args=1) e**x - 1.
</pre>


### floor
<pre class="pre-non-highlight-non-pair">
floor  (class=math #args=1) Floor: nearest integer at or below.
</pre>


### invqnorm
<pre class="pre-non-highlight-non-pair">
invqnorm  (class=math #args=1) Inverse of normal cumulative distribution function.
Note that invqorm(urand()) is normally distributed.
</pre>


### log
<pre class="pre-non-highlight-non-pair">
log  (class=math #args=1) Natural (base-e) logarithm.
</pre>


### log10
<pre class="pre-non-highlight-non-pair">
log10  (class=math #args=1) Base-10 logarithm.
</pre>


### log1p
<pre class="pre-non-highlight-non-pair">
log1p  (class=math #args=1) log(1-x).
</pre>


### logifit
<pre class="pre-non-highlight-non-pair">
logifit  (class=math #args=3)  Given m and b from logistic regression, compute fit:
$yhat=logifit($x,$m,$b).
</pre>


### max
<pre class="pre-non-highlight-non-pair">
max  (class=math #args=variadic) Max of n numbers; null loses.
</pre>


### min
<pre class="pre-non-highlight-non-pair">
min  (class=math #args=variadic) Min of n numbers; null loses.
</pre>


### qnorm
<pre class="pre-non-highlight-non-pair">
qnorm  (class=math #args=1) Normal cumulative distribution function.
</pre>


### round
<pre class="pre-non-highlight-non-pair">
round  (class=math #args=1) Round to nearest integer.
</pre>


### roundm
<pre class="pre-non-highlight-non-pair">
roundm  (class=math #args=2) Round to nearest multiple of m: roundm($x,$m) is
the same as round($x/$m)*$m.
</pre>


### sgn
<pre class="pre-non-highlight-non-pair">
sgn  (class=math #args=1)  +1, 0, -1 for positive, zero, negative input respectively.
</pre>


### sin
<pre class="pre-non-highlight-non-pair">
sin  (class=math #args=1) Trigonometric sine.
</pre>


### sinh
<pre class="pre-non-highlight-non-pair">
sinh  (class=math #args=1) Hyperbolic sine.
</pre>


### sqrt
<pre class="pre-non-highlight-non-pair">
sqrt  (class=math #args=1) Square root.
</pre>


### tan
<pre class="pre-non-highlight-non-pair">
tan  (class=math #args=1) Trigonometric tangent.
</pre>


### tanh
<pre class="pre-non-highlight-non-pair">
tanh  (class=math #args=1) Hyperbolic tangent.
</pre>


### urand
<pre class="pre-non-highlight-non-pair">
urand  (class=math #args=0) Floating-point numbers uniformly distributed on the unit interval.
Int-valued example: '$n=floor(20+urand()*11)'.
</pre>


### urand32
<pre class="pre-non-highlight-non-pair">
urand32  (class=math #args=0) Integer uniformly distributed 0 and 2**32-1 inclusive.
</pre>


### urandint
<pre class="pre-non-highlight-non-pair">
urandint  (class=math #args=2) Integer uniformly distributed between inclusive integer endpoints.
</pre>


### urandrange
<pre class="pre-non-highlight-non-pair">
urandrange  (class=math #args=2) Floating-point numbers uniformly distributed on the interval [a, b).
</pre>

## String functions


### capitalize
<pre class="pre-non-highlight-non-pair">
capitalize  (class=string #args=1) Convert string's first character to uppercase.
</pre>


### clean_whitespace
<pre class="pre-non-highlight-non-pair">
clean_whitespace  (class=string #args=1) Same as collapse_whitespace and strip.
</pre>


### collapse_whitespace
<pre class="pre-non-highlight-non-pair">
collapse_whitespace  (class=string #args=1) Strip repeated whitespace from string.
</pre>


### gsub
<pre class="pre-non-highlight-non-pair">
gsub  (class=string #args=3) Example: '$name=gsub($name, "old", "new")' (replace all).
</pre>


### lstrip
<pre class="pre-non-highlight-non-pair">
lstrip  (class=string #args=1) Strip leading whitespace from string.
</pre>


### regextract
<pre class="pre-non-highlight-non-pair">
regextract  (class=string #args=2) Example: '$name=regextract($name, "[A-Z]{3}[0-9]{2}")'
</pre>


### regextract_or_else
<pre class="pre-non-highlight-non-pair">
regextract_or_else  (class=string #args=3) Example: '$name=regextract_or_else($name, "[A-Z]{3}[0-9]{2}", "default")'
</pre>


### rstrip
<pre class="pre-non-highlight-non-pair">
rstrip  (class=string #args=1) Strip trailing whitespace from string.
</pre>


### ssub
<pre class="pre-non-highlight-non-pair">
ssub  (class=string #args=3) Like sub but does no regexing. No characters are special.
</pre>


### strip
<pre class="pre-non-highlight-non-pair">
strip  (class=string #args=1) Strip leading and trailing whitespace from string.
</pre>


### strlen
<pre class="pre-non-highlight-non-pair">
strlen  (class=string #args=1) String length.
</pre>


### sub
<pre class="pre-non-highlight-non-pair">
sub  (class=string #args=3) Example: '$name=sub($name, "old", "new")' (replace once).
</pre>


### substr
<pre class="pre-non-highlight-non-pair">
substr  (class=string #args=3) substr is an alias for substr0. See also substr1. Miller is generally 1-up
with all array and string indices, but, this is a backward-compatibility issue with Miller 5
and below. Arrays are new in Miller 6; the substr function is older.
</pre>


### substr0
<pre class="pre-non-highlight-non-pair">
substr0  (class=string #args=3) substr0(s,m,n) gives substring of s from 0-up position m to n
inclusive. Negative indices -len .. -1 alias to 0 .. len-1. See also substr and substr1.
</pre>


### substr1
<pre class="pre-non-highlight-non-pair">
substr1  (class=string #args=3) substr1(s,m,n) gives substring of s from 1-up position m to n
inclusive. Negative indices -len .. -1 alias to 1 .. len. See also substr and substr0.
</pre>


### tolower
<pre class="pre-non-highlight-non-pair">
tolower  (class=string #args=1) Convert string to lowercase.
</pre>


### toupper
<pre class="pre-non-highlight-non-pair">
toupper  (class=string #args=1) Convert string to uppercase.
</pre>


### truncate
<pre class="pre-non-highlight-non-pair">
truncate  (class=string #args=2) Truncates string first argument to max length of int second argument.
</pre>


<a id=dot />
### \.
<pre class="pre-non-highlight-non-pair">
.  (class=string #args=2) String concatenation.
</pre>

## System functions


### hostname
<pre class="pre-non-highlight-non-pair">
hostname  (class=system #args=0) Returns the hostname as a string.
</pre>


### os
<pre class="pre-non-highlight-non-pair">
os  (class=system #args=0) Returns the operating-system name as a string.
</pre>


### system
<pre class="pre-non-highlight-non-pair">
system  (class=system #args=1) Run command string, yielding its stdout minus final carriage return.
</pre>


### version
<pre class="pre-non-highlight-non-pair">
version  (class=system #args=0) Returns the Miller version as a string.
</pre>

## Time functions


### dhms2fsec
<pre class="pre-non-highlight-non-pair">
dhms2fsec  (class=time #args=1) Recovers floating-point seconds as in dhms2fsec("5d18h53m20.250000s") = 500000.250000
</pre>


### dhms2sec
<pre class="pre-non-highlight-non-pair">
dhms2sec  (class=time #args=1) Recovers integer seconds as in dhms2sec("5d18h53m20s") = 500000
</pre>


### fsec2dhms
<pre class="pre-non-highlight-non-pair">
fsec2dhms  (class=time #args=1) Formats floating-point seconds as in fsec2dhms(500000.25) = "5d18h53m20.250000s"
</pre>


### fsec2hms
<pre class="pre-non-highlight-non-pair">
fsec2hms  (class=time #args=1) Formats floating-point seconds as in fsec2hms(5000.25) = "01:23:20.250000"
</pre>


### gmt2sec
<pre class="pre-non-highlight-non-pair">
gmt2sec  (class=time #args=1) Parses GMT timestamp as integer seconds since the epoch.
</pre>


### hms2fsec
<pre class="pre-non-highlight-non-pair">
hms2fsec  (class=time #args=1) Recovers floating-point seconds as in hms2fsec("01:23:20.250000") = 5000.250000
</pre>


### hms2sec
<pre class="pre-non-highlight-non-pair">
hms2sec  (class=time #args=1) Recovers integer seconds as in hms2sec("01:23:20") = 5000
</pre>


### sec2dhms
<pre class="pre-non-highlight-non-pair">
sec2dhms  (class=time #args=1) Formats integer seconds as in sec2dhms(500000) = "5d18h53m20s"
</pre>


### sec2gmt
<pre class="pre-non-highlight-non-pair">
sec2gmt  (class=time #args=1,2) Formats seconds since epoch (integer part)
as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
Leaves non-numbers as-is. With second integer argument n, includes n decimal places
for the seconds part
</pre>


### sec2gmtdate
<pre class="pre-non-highlight-non-pair">
sec2gmtdate  (class=time #args=1) Formats seconds since epoch (integer part)
as GMT timestamp with year-month-date, e.g. sec2gmtdate(1440768801.7) = "2015-08-28".
Leaves non-numbers as-is.
</pre>


### sec2hms
<pre class="pre-non-highlight-non-pair">
sec2hms  (class=time #args=1) Formats integer seconds as in sec2hms(5000) = "01:23:20"
</pre>


### strftime
<pre class="pre-non-highlight-non-pair">
strftime  (class=time #args=2)  Formats seconds since the epoch as timestamp, e.g.
	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%SZ") = "2015-08-28T13:33:21Z", and
	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%3SZ") = "2015-08-28T13:33:21.700Z".
	Format strings are as in the C library (please see "man strftime" on your system),
	with the Miller-specific addition of "%1S" through "%9S" which format the seconds
	with 1 through 9 decimal places, respectively. ("%S" uses no decimal places.)
	See also strftime_local.
</pre>


### strptime
<pre class="pre-non-highlight-non-pair">
strptime  (class=time #args=2) strptime: Parses timestamp as floating-point seconds since the epoch,
	e.g. strptime("2015-08-28T13:33:21Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.000000,
	and  strptime("2015-08-28T13:33:21.345Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.345000.
	See also strptime_local.
</pre>


### systime
<pre class="pre-non-highlight-non-pair">
systime  (class=time #args=0) help string will go here
</pre>


### systimeint
<pre class="pre-non-highlight-non-pair">
systimeint  (class=time #args=0) help string will go here
</pre>


### uptime
<pre class="pre-non-highlight-non-pair">
uptime  (class=time #args=0) help string will go here
</pre>

## Typing functions


### asserting_absent
<pre class="pre-non-highlight-non-pair">
asserting_absent  (class=typing #args=1) Aborts with an error if is_absent on the argument returns false,
else returns its argument.
</pre>


### asserting_array
<pre class="pre-non-highlight-non-pair">
asserting_array  (class=typing #args=1) Aborts with an error if is_array on the argument returns false,
else returns its argument.
</pre>


### asserting_bool
<pre class="pre-non-highlight-non-pair">
asserting_bool  (class=typing #args=1) Aborts with an error if is_bool on the argument returns false,
else returns its argument.
</pre>


### asserting_boolean
<pre class="pre-non-highlight-non-pair">
asserting_boolean  (class=typing #args=1) Aborts with an error if is_boolean on the argument returns false,
else returns its argument.
</pre>


### asserting_empty
<pre class="pre-non-highlight-non-pair">
asserting_empty  (class=typing #args=1) Aborts with an error if is_empty on the argument returns false,
else returns its argument.
</pre>


### asserting_empty_map
<pre class="pre-non-highlight-non-pair">
asserting_empty_map  (class=typing #args=1) Aborts with an error if is_empty_map on the argument returns false,
else returns its argument.
</pre>


### asserting_error
<pre class="pre-non-highlight-non-pair">
asserting_error  (class=typing #args=1) Aborts with an error if is_error on the argument returns false,
else returns its argument.
</pre>


### asserting_float
<pre class="pre-non-highlight-non-pair">
asserting_float  (class=typing #args=1) Aborts with an error if is_float on the argument returns false,
else returns its argument.
</pre>


### asserting_int
<pre class="pre-non-highlight-non-pair">
asserting_int  (class=typing #args=1) Aborts with an error if is_int on the argument returns false,
else returns its argument.
</pre>


### asserting_map
<pre class="pre-non-highlight-non-pair">
asserting_map  (class=typing #args=1) Aborts with an error if is_map on the argument returns false,
else returns its argument.
</pre>


### asserting_nonempty_map
<pre class="pre-non-highlight-non-pair">
asserting_nonempty_map  (class=typing #args=1) Aborts with an error if is_nonempty_map on the argument returns false,
else returns its argument.
</pre>


### asserting_not_array
<pre class="pre-non-highlight-non-pair">
asserting_not_array  (class=typing #args=1) Aborts with an error if is_not_array on the argument returns false,
else returns its argument.
</pre>


### asserting_not_empty
<pre class="pre-non-highlight-non-pair">
asserting_not_empty  (class=typing #args=1) Aborts with an error if is_not_empty on the argument returns false,
else returns its argument.
</pre>


### asserting_not_map
<pre class="pre-non-highlight-non-pair">
asserting_not_map  (class=typing #args=1) Aborts with an error if is_not_map on the argument returns false,
else returns its argument.
</pre>


### asserting_not_null
<pre class="pre-non-highlight-non-pair">
asserting_not_null  (class=typing #args=1) Aborts with an error if is_not_null on the argument returns false,
else returns its argument.
</pre>


### asserting_null
<pre class="pre-non-highlight-non-pair">
asserting_null  (class=typing #args=1) Aborts with an error if is_null on the argument returns false,
else returns its argument.
</pre>


### asserting_numeric
<pre class="pre-non-highlight-non-pair">
asserting_numeric  (class=typing #args=1) Aborts with an error if is_numeric on the argument returns false,
else returns its argument.
</pre>


### asserting_present
<pre class="pre-non-highlight-non-pair">
asserting_present  (class=typing #args=1) Aborts with an error if is_present on the argument returns false,
else returns its argument.
</pre>


### asserting_string
<pre class="pre-non-highlight-non-pair">
asserting_string  (class=typing #args=1) Aborts with an error if is_string on the argument returns false,
else returns its argument.
</pre>


### is_absent
<pre class="pre-non-highlight-non-pair">
is_absent  (class=typing #args=1) False if field is present in input, true otherwise
</pre>


### is_array
<pre class="pre-non-highlight-non-pair">
is_array  (class=typing #args=1) True if argument is an array.
</pre>


### is_bool
<pre class="pre-non-highlight-non-pair">
is_bool  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_boolean.
</pre>


### is_boolean
<pre class="pre-non-highlight-non-pair">
is_boolean  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_bool.
</pre>


### is_empty
<pre class="pre-non-highlight-non-pair">
is_empty  (class=typing #args=1) True if field is present in input with empty string value, false otherwise.
</pre>


### is_empty_map
<pre class="pre-non-highlight-non-pair">
is_empty_map  (class=typing #args=1) True if argument is a map which is empty.
</pre>


### is_error
<pre class="pre-non-highlight-non-pair">
is_error  (class=typing #args=1) True if if argument is an error, such as taking string length of an integer.
</pre>


### is_float
<pre class="pre-non-highlight-non-pair">
is_float  (class=typing #args=1) True if field is present with value inferred to be float
</pre>


### is_int
<pre class="pre-non-highlight-non-pair">
is_int  (class=typing #args=1) True if field is present with value inferred to be int
</pre>


### is_map
<pre class="pre-non-highlight-non-pair">
is_map  (class=typing #args=1) True if argument is a map.
</pre>


### is_nonempty_map
<pre class="pre-non-highlight-non-pair">
is_nonempty_map  (class=typing #args=1) True if argument is a map which is non-empty.
</pre>


### is_not_array
<pre class="pre-non-highlight-non-pair">
is_not_array  (class=typing #args=1) True if argument is not an array.
</pre>


### is_not_empty
<pre class="pre-non-highlight-non-pair">
is_not_empty  (class=typing #args=1) False if field is present in input with empty value, true otherwise
</pre>


### is_not_map
<pre class="pre-non-highlight-non-pair">
is_not_map  (class=typing #args=1) True if argument is not a map.
</pre>


### is_not_null
<pre class="pre-non-highlight-non-pair">
is_not_null  (class=typing #args=1) False if argument is null (empty or absent), true otherwise.
</pre>


### is_null
<pre class="pre-non-highlight-non-pair">
is_null  (class=typing #args=1) True if argument is null (empty or absent), false otherwise.
</pre>


### is_numeric
<pre class="pre-non-highlight-non-pair">
is_numeric  (class=typing #args=1) True if field is present with value inferred to be int or float
</pre>


### is_present
<pre class="pre-non-highlight-non-pair">
is_present  (class=typing #args=1) True if field is present in input, false otherwise.
</pre>


### is_string
<pre class="pre-non-highlight-non-pair">
is_string  (class=typing #args=1) True if field is present with string (including empty-string) value
</pre>


### typeof
<pre class="pre-non-highlight-non-pair">
typeof  (class=typing #args=1) Convert argument to type of argument (e.g. "str"). For debug.
</pre>

