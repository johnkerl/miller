<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# DSL reference: built-in functions

## Summary

mlr: option "--list-all-functions-as-table" not recognized.
Please run "mlr --help" for usage information.

## List of functions

Each function takes a specific number of arguments, as shown below, except for functions marked as variadic such as `min` and `max`. (The latter compute min and max of any number of numerical arguments.) There is no notion of optional or default-on-absent arguments. All argument-passing is positional rather than by name; arguments are passed by value, not by reference.

You can get a list of all functions using **mlr -f**, with details using **mlr -F**.


<a id=colon />
## \!

<pre class="pre-non-highlight">
!  (class=boolean #args=1) Logical negation.
</pre>


## !=

<pre class="pre-non-highlight">
!=  (class=boolean #args=2) String/numeric inequality. Mixing number and string results in string compare.
</pre>


## !=~

<pre class="pre-non-highlight">
!=~  (class=boolean #args=2) String (left-hand side) does not match regex (right-hand side), e.g. '$name !=~ "^a.*b$"'.
</pre>


## %

<pre class="pre-non-highlight">
%  (class=arithmetic #args=2) Remainder; never negative-valued (pythonic).
</pre>


## &

<pre class="pre-non-highlight">
&  (class=arithmetic #args=2) Bitwise AND.
</pre>


## &&

<pre class="pre-non-highlight">
&&  (class=boolean #args=2) Logical AND.
</pre>


<a id=times />
## \*

<pre class="pre-non-highlight">
*  (class=arithmetic #args=2) Multiplication, with integer*integer overflow to float.
</pre>


<a id=exponentiation />
## \**

<pre class="pre-non-highlight">
**  (class=arithmetic #args=2) Exponentiation. Same as pow, but as an infix operator.
</pre>


<a id=plus />
## \+

<pre class="pre-non-highlight">
+  (class=arithmetic #args=1,2) Addition as binary operator; unary plus operator.
</pre>


<a id=minus />
## \-

<pre class="pre-non-highlight">
-  (class=arithmetic #args=1,2) Subtraction as binary operator; unary negation operator.
</pre>


## .

<pre class="pre-non-highlight">
.  (class=string #args=2) String concatenation.
</pre>


## .*

<pre class="pre-non-highlight">
.*  (class=arithmetic #args=2) Multiplication, with integer-to-integer overflow.
</pre>


## .+

<pre class="pre-non-highlight">
.+  (class=arithmetic #args=2) Addition, with integer-to-integer overflow.
</pre>


## .-

<pre class="pre-non-highlight">
.-  (class=arithmetic #args=2) Subtraction, with integer-to-integer overflow.
</pre>


## ./

<pre class="pre-non-highlight">
./  (class=arithmetic #args=2) Integer division; not pythonic.
</pre>


## /

<pre class="pre-non-highlight">
/  (class=arithmetic #args=2) Division. Integer / integer is floating-point.
</pre>


## //

<pre class="pre-non-highlight">
//  (class=arithmetic #args=2) Pythonic integer division, rounding toward negative.
</pre>


## <

<pre class="pre-non-highlight">
<  (class=boolean #args=2) String/numeric less-than. Mixing number and string results in string compare.
</pre>


## <<

<pre class="pre-non-highlight">
<<  (class=arithmetic #args=2) Bitwise left-shift.
</pre>


## <=

<pre class="pre-non-highlight">
<=  (class=boolean #args=2) String/numeric less-than-or-equals. Mixing number and string results in string compare.
</pre>


## ==

<pre class="pre-non-highlight">
==  (class=boolean #args=2) String/numeric equality. Mixing number and string results in string compare.
</pre>


## =~

<pre class="pre-non-highlight">
=~  (class=boolean #args=2) String (left-hand side) matches regex (right-hand side), e.g. '$name =~ "^a.*b$"'.
</pre>


## >

<pre class="pre-non-highlight">
>  (class=boolean #args=2) String/numeric greater-than. Mixing number and string results in string compare.
</pre>


## >=

<pre class="pre-non-highlight">
>=  (class=boolean #args=2) String/numeric greater-than-or-equals. Mixing number and string results in string compare.
</pre>


<a id=srsh />
## \>\>

<pre class="pre-non-highlight">
>>  (class=arithmetic #args=2) Bitwise signed right-shift.
</pre>


<a id=ursh />
## \>\>\>

<pre class="pre-non-highlight">
>>>  (class=arithmetic #args=2) Bitwise unsigned right-shift.
</pre>


<a id=question-mark-colon />
## \?

<pre class="pre-non-highlight">
?:  (class=boolean #args=3) Standard ternary operator.
</pre>


## ??

<pre class="pre-non-highlight">
??  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record.
</pre>


## ???

<pre class="pre-non-highlight">
???  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record, or has empty value.
</pre>


## ^

<pre class="pre-non-highlight">
^  (class=arithmetic #args=2) Bitwise XOR.
</pre>


## ^^

<pre class="pre-non-highlight">
^^  (class=boolean #args=2) Logical XOR.
</pre>


<a id=bitwise-or />
## \|

<pre class="pre-non-highlight">
|  (class=arithmetic #args=2) Bitwise OR.
</pre>


## ||

<pre class="pre-non-highlight">
||  (class=boolean #args=2) Logical OR.
</pre>


## ~

<pre class="pre-non-highlight">
~  (class=arithmetic #args=1) Bitwise NOT. Beware '$y=~$x' since =~ is the
regex-match operator: try '$y = ~$x'.
</pre>


## abs

<pre class="pre-non-highlight">
abs  (class=math #args=1) Absolute value.
</pre>


## acos

<pre class="pre-non-highlight">
acos  (class=math #args=1) Inverse trigonometric cosine.
</pre>


## acosh

<pre class="pre-non-highlight">
acosh  (class=math #args=1) Inverse hyperbolic cosine.
</pre>


## append

<pre class="pre-non-highlight">
append  (class=maps/arrays #args=2) Appends second argument to end of first argument, which must be an array.
</pre>


## arrayify

<pre class="pre-non-highlight">
arrayify  (class=maps/arrays #args=1) Walks through a nested map/array, converting any map with consecutive keys
"1", "2", ... into an array. Useful to wrap the output of unflatten.
</pre>


## asin

<pre class="pre-non-highlight">
asin  (class=math #args=1) Inverse trigonometric sine.
</pre>


## asinh

<pre class="pre-non-highlight">
asinh  (class=math #args=1) Inverse hyperbolic sine.
</pre>


## asserting_absent

<pre class="pre-non-highlight">
asserting_absent  (class=typing #args=1) Aborts with an error if is_absent on the argument returns false,
else returns its argument.
</pre>


## asserting_array

<pre class="pre-non-highlight">
asserting_array  (class=typing #args=1) Aborts with an error if is_array on the argument returns false,
else returns its argument.
</pre>


## asserting_bool

<pre class="pre-non-highlight">
asserting_bool  (class=typing #args=1) Aborts with an error if is_bool on the argument returns false,
else returns its argument.
</pre>


## asserting_boolean

<pre class="pre-non-highlight">
asserting_boolean  (class=typing #args=1) Aborts with an error if is_boolean on the argument returns false,
else returns its argument.
</pre>


## asserting_empty

<pre class="pre-non-highlight">
asserting_empty  (class=typing #args=1) Aborts with an error if is_empty on the argument returns false,
else returns its argument.
</pre>


## asserting_empty_map

<pre class="pre-non-highlight">
asserting_empty_map  (class=typing #args=1) Aborts with an error if is_empty_map on the argument returns false,
else returns its argument.
</pre>


## asserting_error

<pre class="pre-non-highlight">
asserting_error  (class=typing #args=1) Aborts with an error if is_error on the argument returns false,
else returns its argument.
</pre>


## asserting_float

<pre class="pre-non-highlight">
asserting_float  (class=typing #args=1) Aborts with an error if is_float on the argument returns false,
else returns its argument.
</pre>


## asserting_int

<pre class="pre-non-highlight">
asserting_int  (class=typing #args=1) Aborts with an error if is_int on the argument returns false,
else returns its argument.
</pre>


## asserting_map

<pre class="pre-non-highlight">
asserting_map  (class=typing #args=1) Aborts with an error if is_map on the argument returns false,
else returns its argument.
</pre>


## asserting_nonempty_map

<pre class="pre-non-highlight">
asserting_nonempty_map  (class=typing #args=1) Aborts with an error if is_nonempty_map on the argument returns false,
else returns its argument.
</pre>


## asserting_not_array

<pre class="pre-non-highlight">
asserting_not_array  (class=typing #args=1) Aborts with an error if is_not_array on the argument returns false,
else returns its argument.
</pre>


## asserting_not_empty

<pre class="pre-non-highlight">
asserting_not_empty  (class=typing #args=1) Aborts with an error if is_not_empty on the argument returns false,
else returns its argument.
</pre>


## asserting_not_map

<pre class="pre-non-highlight">
asserting_not_map  (class=typing #args=1) Aborts with an error if is_not_map on the argument returns false,
else returns its argument.
</pre>


## asserting_not_null

<pre class="pre-non-highlight">
asserting_not_null  (class=typing #args=1) Aborts with an error if is_not_null on the argument returns false,
else returns its argument.
</pre>


## asserting_null

<pre class="pre-non-highlight">
asserting_null  (class=typing #args=1) Aborts with an error if is_null on the argument returns false,
else returns its argument.
</pre>


## asserting_numeric

<pre class="pre-non-highlight">
asserting_numeric  (class=typing #args=1) Aborts with an error if is_numeric on the argument returns false,
else returns its argument.
</pre>


## asserting_present

<pre class="pre-non-highlight">
asserting_present  (class=typing #args=1) Aborts with an error if is_present on the argument returns false,
else returns its argument.
</pre>


## asserting_string

<pre class="pre-non-highlight">
asserting_string  (class=typing #args=1) Aborts with an error if is_string on the argument returns false,
else returns its argument.
</pre>


## atan

<pre class="pre-non-highlight">
atan  (class=math #args=1) One-argument arctangent.
</pre>


## atan2

<pre class="pre-non-highlight">
atan2  (class=math #args=2) Two-argument arctangent.
</pre>


## atanh

<pre class="pre-non-highlight">
atanh  (class=math #args=1) Inverse hyperbolic tangent.
</pre>


## bitcount

<pre class="pre-non-highlight">
bitcount  (class=arithmetic #args=1) Count of 1-bits.
</pre>


## boolean

<pre class="pre-non-highlight">
boolean  (class=conversion #args=1) Convert int/float/bool/string to boolean.
</pre>


## capitalize

<pre class="pre-non-highlight">
capitalize  (class=string #args=1) Convert string's first character to uppercase.
</pre>


## cbrt

<pre class="pre-non-highlight">
cbrt  (class=math #args=1) Cube root.
</pre>


## ceil

<pre class="pre-non-highlight">
ceil  (class=math #args=1) Ceiling: nearest integer at or above.
</pre>


## clean_whitespace

<pre class="pre-non-highlight">
clean_whitespace  (class=string #args=1) Same as collapse_whitespace and strip.
</pre>


## collapse_whitespace

<pre class="pre-non-highlight">
collapse_whitespace  (class=string #args=1) Strip repeated whitespace from string.
</pre>


## cos

<pre class="pre-non-highlight">
cos  (class=math #args=1) Trigonometric cosine.
</pre>


## cosh

<pre class="pre-non-highlight">
cosh  (class=math #args=1) Hyperbolic cosine.
</pre>


## depth

<pre class="pre-non-highlight">
depth  (class=maps/arrays #args=1) Prints maximum depth of map/array. Scalars have depth 0.
</pre>


## dhms2fsec

<pre class="pre-non-highlight">
dhms2fsec  (class=time #args=1) Recovers floating-point seconds as in dhms2fsec("5d18h53m20.250000s") = 500000.250000
</pre>


## dhms2sec

<pre class="pre-non-highlight">
dhms2sec  (class=time #args=1) Recovers integer seconds as in dhms2sec("5d18h53m20s") = 500000
</pre>


## erf

<pre class="pre-non-highlight">
erf  (class=math #args=1) Error function.
</pre>


## erfc

<pre class="pre-non-highlight">
erfc  (class=math #args=1) Complementary error function.
</pre>


## exp

<pre class="pre-non-highlight">
exp  (class=math #args=1) Exponential function e**x.
</pre>


## expm1

<pre class="pre-non-highlight">
expm1  (class=math #args=1) e**x - 1.
</pre>


## flatten

<pre class="pre-non-highlight">
flatten  (class=maps/arrays #args=3) Flattens multi-level maps to single-level ones. Examples:
flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
Two-argument version: flatten($*, ".") is the same as flatten("", ".", $*).
Useful for nested JSON-like structures for non-JSON file formats like CSV.
</pre>


## float

<pre class="pre-non-highlight">
float  (class=conversion #args=1) Convert int/float/bool/string to float.
</pre>


## floor

<pre class="pre-non-highlight">
floor  (class=math #args=1) Floor: nearest integer at or below.
</pre>


## fmtnum

<pre class="pre-non-highlight">
fmtnum  (class=conversion #args=2) Convert int/float/bool to string using
printf-style format string, e.g. '$s = fmtnum($n, "%06lld")'.
</pre>


## fsec2dhms

<pre class="pre-non-highlight">
fsec2dhms  (class=time #args=1) Formats floating-point seconds as in fsec2dhms(500000.25) = "5d18h53m20.250000s"
</pre>


## fsec2hms

<pre class="pre-non-highlight">
fsec2hms  (class=time #args=1) Formats floating-point seconds as in fsec2hms(5000.25) = "01:23:20.250000"
</pre>


## get_keys

<pre class="pre-non-highlight">
get_keys  (class=maps/arrays #args=1) Returns array of keys of map or array
</pre>


## get_values

<pre class="pre-non-highlight">
get_values  (class=maps/arrays #args=1) Returns array of keys of map or array -- in the latter case, returns a copy of the array
</pre>


## gmt2sec

<pre class="pre-non-highlight">
gmt2sec  (class=time #args=1) Parses GMT timestamp as integer seconds since the epoch.
</pre>


## gsub

<pre class="pre-non-highlight">
gsub  (class=string #args=3) Example: '$name=gsub($name, "old", "new")' (replace all).
</pre>


## haskey

<pre class="pre-non-highlight">
haskey  (class=maps/arrays #args=2) True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or
'haskey(mymap, mykey)', or true/false if array index is in bounds / out of bounds.
Error if 1st argument is not a map or array. Note -n..-1 alias to 1..n in Miller arrays.
</pre>


## hexfmt

<pre class="pre-non-highlight">
hexfmt  (class=conversion #args=1) Convert int to hex string, e.g. 255 to "0xff".
</pre>


## hms2fsec

<pre class="pre-non-highlight">
hms2fsec  (class=time #args=1) Recovers floating-point seconds as in hms2fsec("01:23:20.250000") = 5000.250000
</pre>


## hms2sec

<pre class="pre-non-highlight">
hms2sec  (class=time #args=1) Recovers integer seconds as in hms2sec("01:23:20") = 5000
</pre>


## hostname

<pre class="pre-non-highlight">
hostname  (class=system #args=0) Returns the hostname as a string.
</pre>


## int

<pre class="pre-non-highlight">
int  (class=conversion #args=1) Convert int/float/bool/string to int.
</pre>


## invqnorm

<pre class="pre-non-highlight">
invqnorm  (class=math #args=1) Inverse of normal cumulative distribution function.
Note that invqorm(urand()) is normally distributed.
</pre>


## is_absent

<pre class="pre-non-highlight">
is_absent  (class=typing #args=1) False if field is present in input, true otherwise
</pre>


## is_array

<pre class="pre-non-highlight">
is_array  (class=typing #args=1) True if argument is an array.
</pre>


## is_bool

<pre class="pre-non-highlight">
is_bool  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_boolean.
</pre>


## is_boolean

<pre class="pre-non-highlight">
is_boolean  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_bool.
</pre>


## is_empty

<pre class="pre-non-highlight">
is_empty  (class=typing #args=1) True if field is present in input with empty string value, false otherwise.
</pre>


## is_empty_map

<pre class="pre-non-highlight">
is_empty_map  (class=typing #args=1) True if argument is a map which is empty.
</pre>


## is_error

<pre class="pre-non-highlight">
is_error  (class=typing #args=1) True if if argument is an error, such as taking string length of an integer.
</pre>


## is_float

<pre class="pre-non-highlight">
is_float  (class=typing #args=1) True if field is present with value inferred to be float
</pre>


## is_int

<pre class="pre-non-highlight">
is_int  (class=typing #args=1) True if field is present with value inferred to be int
</pre>


## is_map

<pre class="pre-non-highlight">
is_map  (class=typing #args=1) True if argument is a map.
</pre>


## is_nonempty_map

<pre class="pre-non-highlight">
is_nonempty_map  (class=typing #args=1) True if argument is a map which is non-empty.
</pre>


## is_not_array

<pre class="pre-non-highlight">
is_not_array  (class=typing #args=1) True if argument is not an array.
</pre>


## is_not_empty

<pre class="pre-non-highlight">
is_not_empty  (class=typing #args=1) False if field is present in input with empty value, true otherwise
</pre>


## is_not_map

<pre class="pre-non-highlight">
is_not_map  (class=typing #args=1) True if argument is not a map.
</pre>


## is_not_null

<pre class="pre-non-highlight">
is_not_null  (class=typing #args=1) False if argument is null (empty or absent), true otherwise.
</pre>


## is_null

<pre class="pre-non-highlight">
is_null  (class=typing #args=1) True if argument is null (empty or absent), false otherwise.
</pre>


## is_numeric

<pre class="pre-non-highlight">
is_numeric  (class=typing #args=1) True if field is present with value inferred to be int or float
</pre>


## is_present

<pre class="pre-non-highlight">
is_present  (class=typing #args=1) True if field is present in input, false otherwise.
</pre>


## is_string

<pre class="pre-non-highlight">
is_string  (class=typing #args=1) True if field is present with string (including empty-string) value
</pre>


## joink

<pre class="pre-non-highlight">
joink  (class=conversion #args=2) Makes string from map/array keys. Examples:
joink({"a":3,"b":4,"c":5}, ",") = "a,b,c"
joink([1,2,3], ",") = "1,2,3".
</pre>


## joinkv

<pre class="pre-non-highlight">
joinkv  (class=conversion #args=3) Makes string from map/array key-value pairs. Examples:
joinkv([3,4,5], "=", ",") = "1=3,2=4,3=5"
joinkv({"a":3,"b":4,"c":5}, "=", ",") = "a=3,b=4,c=5"
</pre>


## joinv

<pre class="pre-non-highlight">
joinv  (class=conversion #args=2) Makes string from map/array values.
joinv([3,4,5], ",") = "3,4,5"
joinv({"a":3,"b":4,"c":5}, ",") = "3,4,5"
</pre>


## json_parse

<pre class="pre-non-highlight">
json_parse  (class=maps/arrays #args=1) Converts value from JSON-formatted string.
</pre>


## json_stringify

<pre class="pre-non-highlight">
json_stringify  (class=maps/arrays #args=1,2) Converts value to JSON-formatted string. Default output is single-line.
With optional second boolean argument set to true, produces multiline output.
</pre>


## leafcount

<pre class="pre-non-highlight">
leafcount  (class=maps/arrays #args=1) Counts total number of terminal values in map/array. For single-level
map/array, same as length.
</pre>


## length

<pre class="pre-non-highlight">
length  (class=maps/arrays #args=1) Counts number of top-level entries in array/map. Scalars have length 1.
</pre>


## log

<pre class="pre-non-highlight">
log  (class=math #args=1) Natural (base-e) logarithm.
</pre>


## log10

<pre class="pre-non-highlight">
log10  (class=math #args=1) Base-10 logarithm.
</pre>


## log1p

<pre class="pre-non-highlight">
log1p  (class=math #args=1) log(1-x).
</pre>


## logifit

<pre class="pre-non-highlight">
logifit  (class=math #args=3)  Given m and b from logistic regression, compute fit:
$yhat=logifit($x,$m,$b).
</pre>


## lstrip

<pre class="pre-non-highlight">
lstrip  (class=string #args=1) Strip leading whitespace from string.
</pre>


## madd

<pre class="pre-non-highlight">
madd  (class=arithmetic #args=3) a + b mod m (integers)
</pre>


## mapdiff

<pre class="pre-non-highlight">
mapdiff  (class=maps/arrays #args=variadic) With 0 args, returns empty map. With 1 arg, returns copy of arg.
With 2 or more, returns copy of arg 1 with all keys from any of remaining
argument maps removed.
</pre>


## mapexcept

<pre class="pre-non-highlight">
mapexcept  (class=maps/arrays #args=variadic) Returns a map with keys from remaining arguments, if any, unset.
Remaining arguments can be strings or arrays of string.
E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'
and  'mapexcept({1:2,3:4,5:6}, [1, 5, 7])' is '{3:4}'.
</pre>


## mapselect

<pre class="pre-non-highlight">
mapselect  (class=maps/arrays #args=variadic) Returns a map with only keys from remaining arguments set.
Remaining arguments can be strings or arrays of string.
E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'
and  'mapselect({1:2,3:4,5:6}, [1, 5, 7])' is '{1:2,5:6}'.
</pre>


## mapsum

<pre class="pre-non-highlight">
mapsum  (class=maps/arrays #args=variadic) With 0 args, returns empty map. With >= 1 arg, returns a map with
key-value pairs from all arguments. Rightmost collisions win, e.g.
'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.
</pre>


## max

<pre class="pre-non-highlight">
max  (class=math #args=variadic) Max of n numbers; null loses.
</pre>


## md5

<pre class="pre-non-highlight">
md5  (class=hashing #args=1) MD5 hash.
</pre>


## mexp

<pre class="pre-non-highlight">
mexp  (class=arithmetic #args=3) a ** b mod m (integers)
</pre>


## min

<pre class="pre-non-highlight">
min  (class=math #args=variadic) Min of n numbers; null loses.
</pre>


## mmul

<pre class="pre-non-highlight">
mmul  (class=arithmetic #args=3) a * b mod m (integers)
</pre>


## msub

<pre class="pre-non-highlight">
msub  (class=arithmetic #args=3) a - b mod m (integers)
</pre>


## os

<pre class="pre-non-highlight">
os  (class=system #args=0) Returns the operating-system name as a string.
</pre>


## pow

<pre class="pre-non-highlight">
pow  (class=arithmetic #args=2) Exponentiation. Same as **, but as a function.
</pre>


## qnorm

<pre class="pre-non-highlight">
qnorm  (class=math #args=1) Normal cumulative distribution function.
</pre>


## regextract

<pre class="pre-non-highlight">
regextract  (class=string #args=2) Example: '$name=regextract($name, "[A-Z]{3}[0-9]{2}")'
</pre>


## regextract_or_else

<pre class="pre-non-highlight">
regextract_or_else  (class=string #args=3) Example: '$name=regextract_or_else($name, "[A-Z]{3}[0-9]{2}", "default")'
</pre>


## round

<pre class="pre-non-highlight">
round  (class=math #args=1) Round to nearest integer.
</pre>


## roundm

<pre class="pre-non-highlight">
roundm  (class=math #args=2) Round to nearest multiple of m: roundm($x,$m) is
the same as round($x/$m)*$m.
</pre>


## rstrip

<pre class="pre-non-highlight">
rstrip  (class=string #args=1) Strip trailing whitespace from string.
</pre>


## sec2dhms

<pre class="pre-non-highlight">
sec2dhms  (class=time #args=1) Formats integer seconds as in sec2dhms(500000) = "5d18h53m20s"
</pre>


## sec2gmt

<pre class="pre-non-highlight">
sec2gmt  (class=time #args=1,2) Formats seconds since epoch (integer part)
as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
Leaves non-numbers as-is. With second integer argument n, includes n decimal places
for the seconds part
</pre>


## sec2gmtdate

<pre class="pre-non-highlight">
sec2gmtdate  (class=time #args=1) Formats seconds since epoch (integer part)
as GMT timestamp with year-month-date, e.g. sec2gmtdate(1440768801.7) = "2015-08-28".
Leaves non-numbers as-is.
</pre>


## sec2hms

<pre class="pre-non-highlight">
sec2hms  (class=time #args=1) Formats integer seconds as in sec2hms(5000) = "01:23:20"
</pre>


## sgn

<pre class="pre-non-highlight">
sgn  (class=math #args=1)  +1, 0, -1 for positive, zero, negative input respectively.
</pre>


## sha1

<pre class="pre-non-highlight">
sha1  (class=hashing #args=1) SHA1 hash.
</pre>


## sha256

<pre class="pre-non-highlight">
sha256  (class=hashing #args=1) SHA256 hash.
</pre>


## sha512

<pre class="pre-non-highlight">
sha512  (class=hashing #args=1) SHA512 hash.
</pre>


## sin

<pre class="pre-non-highlight">
sin  (class=math #args=1) Trigonometric sine.
</pre>


## sinh

<pre class="pre-non-highlight">
sinh  (class=math #args=1) Hyperbolic sine.
</pre>


## splita

<pre class="pre-non-highlight">
splita  (class=conversion #args=2) Splits string into array with type inference. Example:
splita("3,4,5", ",") = [3,4,5]
</pre>


## splitax

<pre class="pre-non-highlight">
splitax  (class=conversion #args=2) Splits string into array without type inference. Example:
splita("3,4,5", ",") = ["3","4","5"]
</pre>


## splitkv

<pre class="pre-non-highlight">
splitkv  (class=conversion #args=3) Splits string by separators into map with type inference. Example:
splitkv("a=3,b=4,c=5", "=", ",") = {"a":3,"b":4,"c":5}
</pre>


## splitkvx

<pre class="pre-non-highlight">
splitkvx  (class=conversion #args=3) Splits string by separators into map without type inference (keys and
values are strings). Example:
splitkvx("a=3,b=4,c=5", "=", ",") = {"a":"3","b":"4","c":"5"}
</pre>


## splitnv

<pre class="pre-non-highlight">
splitnv  (class=conversion #args=2) Splits string by separator into integer-indexed map with type inference. Example:
splitnv("a,b,c", ",") = {"1":"a","2":"b","3":"c"}
</pre>


## splitnvx

<pre class="pre-non-highlight">
splitnvx  (class=conversion #args=2) Splits string by separator into integer-indexed map without type
inference (values are strings). Example:
splitnvx("3,4,5", ",") = {"1":"3","2":"4","3":"5"}
</pre>


## sqrt

<pre class="pre-non-highlight">
sqrt  (class=math #args=1) Square root.
</pre>


## ssub

<pre class="pre-non-highlight">
ssub  (class=string #args=3) Like sub but does no regexing. No characters are special.
</pre>


## strftime

<pre class="pre-non-highlight">
strftime  (class=time #args=2)  Formats seconds since the epoch as timestamp, e.g.
	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%SZ") = "2015-08-28T13:33:21Z", and
	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%3SZ") = "2015-08-28T13:33:21.700Z".
	Format strings are as in the C library (please see "man strftime" on your system),
	with the Miller-specific addition of "%1S" through "%9S" which format the seconds
	with 1 through 9 decimal places, respectively. ("%S" uses no decimal places.)
	See also strftime_local.
</pre>


## string

<pre class="pre-non-highlight">
string  (class=conversion #args=1) Convert int/float/bool/string/array/map to string.
</pre>


## strip

<pre class="pre-non-highlight">
strip  (class=string #args=1) Strip leading and trailing whitespace from string.
</pre>


## strlen

<pre class="pre-non-highlight">
strlen  (class=string #args=1) String length.
</pre>


## strptime

<pre class="pre-non-highlight">
strptime  (class=time #args=2) strptime: Parses timestamp as floating-point seconds since the epoch,
	e.g. strptime("2015-08-28T13:33:21Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.000000,
	and  strptime("2015-08-28T13:33:21.345Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.345000.
	See also strptime_local.
</pre>


## sub

<pre class="pre-non-highlight">
sub  (class=string #args=3) Example: '$name=sub($name, "old", "new")' (replace once).
</pre>


## substr

<pre class="pre-non-highlight">
substr  (class=string #args=3) substr is an alias for substr0. See also substr1. Miller is generally 1-up
with all array indices, but, this is a backward-compatibility issue with Miller 5 and below.
Arrays are new in Miller 6; the substr function is older.
</pre>


## substr0

<pre class="pre-non-highlight">
substr0  (class=string #args=3) substr0(s,m,n) gives substring of s from 0-up position m to n
inclusive. Negative indices -len .. -1 alias to 0 .. len-1.
</pre>


## substr1

<pre class="pre-non-highlight">
substr1  (class=string #args=3) substr1(s,m,n) gives substring of s from 1-up position m to n
inclusive. Negative indices -len .. -1 alias to 1 .. len.
</pre>


## system

<pre class="pre-non-highlight">
system  (class=system #args=1) Run command string, yielding its stdout minus final carriage return.
</pre>


## systime

<pre class="pre-non-highlight">
systime  (class=time #args=0) help string will go here
</pre>


## systimeint

<pre class="pre-non-highlight">
systimeint  (class=time #args=0) help string will go here
</pre>


## tan

<pre class="pre-non-highlight">
tan  (class=math #args=1) Trigonometric tangent.
</pre>


## tanh

<pre class="pre-non-highlight">
tanh  (class=math #args=1) Hyperbolic tangent.
</pre>


## tolower

<pre class="pre-non-highlight">
tolower  (class=string #args=1) Convert string to lowercase.
</pre>


## toupper

<pre class="pre-non-highlight">
toupper  (class=string #args=1) Convert string to uppercase.
</pre>


## truncate

<pre class="pre-non-highlight">
truncate  (class=string #args=2) Truncates string first argument to max length of int second argument.
</pre>


## typeof

<pre class="pre-non-highlight">
typeof  (class=typing #args=1) Convert argument to type of argument (e.g. "str"). For debug.
</pre>


## unflatten

<pre class="pre-non-highlight">
unflatten  (class=maps/arrays #args=2) Reverses flatten. Example:
unflatten({"a.b.c" : 4}, ".") is {"a": "b": { "c": 4 }}.
Useful for nested JSON-like structures for non-JSON file formats like CSV.
See also arrayify.
</pre>


## uptime

<pre class="pre-non-highlight">
uptime  (class=time #args=0) help string will go here
</pre>


## urand

<pre class="pre-non-highlight">
urand  (class=math #args=0) Floating-point numbers uniformly distributed on the unit interval.
Int-valued example: '$n=floor(20+urand()*11)'.
</pre>


## urand32

<pre class="pre-non-highlight">
urand32  (class=math #args=0) Integer uniformly distributed 0 and 2**32-1 inclusive.
</pre>


## urandint

<pre class="pre-non-highlight">
urandint  (class=math #args=2) Integer uniformly distributed between inclusive integer endpoints.
</pre>


## urandrange

<pre class="pre-non-highlight">
urandrange  (class=math #args=2) Floating-point numbers uniformly distributed on the interval [a, b).
</pre>


## version

<pre class="pre-non-highlight">
version  (class=system #args=0) Returns the Miller version as a string.
</pre>

