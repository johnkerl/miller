<!---  PLEASE DO NOT EDIT DIRECTLY. EDIT THE .md.in FILE PLEASE. --->
# DSL reference: built-in functions

## Summary

mlr: option "--list-all-functions-as-table" not recognized.
Please run "mlr --help" for usage information.

## List of functions

Each function takes a specific number of arguments, as shown below, except for functions marked as variadic such as `min` and `max`. (The latter compute min and max of any number of numerical arguments.) There is no notion of optional or default-on-absent arguments. All argument-passing is positional rather than by name; arguments are passed by value, not by reference.

You can get a list of all functions using **mlr -f**, with details using **mlr -F**.


## \!

<pre>
!  (class=boolean #args=1) Logical negation.
</pre>


## !=

<pre>
!=  (class=boolean #args=2) String/numeric inequality. Mixing number and string results in string compare.
</pre>


## !=~

<pre>
!=~  (class=boolean #args=2) String (left-hand side) does not match regex (right-hand side), e.g. '$name !=~ "^a.*b$"'.
</pre>


## %

<pre>
%  (class=arithmetic #args=2) Remainder; never negative-valued (pythonic).
</pre>


## &

<pre>
&  (class=arithmetic #args=2) Bitwise AND.
</pre>


## &&

<pre>
&&  (class=boolean #args=2) Logical AND.
</pre>


## \*

<pre>
*  (class=arithmetic #args=2) Multiplication, with integer*integer overflow to float.
</pre>


## \**

<pre>
**  (class=arithmetic #args=2) Exponentiation. Same as pow, but as an infix operator.
</pre>


## \+

<pre>
+  (class=arithmetic #args=1,2) Addition as binary operator; unary plus operator.
</pre>


## \-

<pre>
-  (class=arithmetic #args=1,2) Subtraction as binary operator; unary negation operator.
</pre>


## .

<pre>
.  (class=string #args=2) String concatenation.
</pre>


## .*

<pre>
.*  (class=arithmetic #args=2) Multiplication, with integer-to-integer overflow.
</pre>


## .+

<pre>
.+  (class=arithmetic #args=2) Addition, with integer-to-integer overflow.
</pre>


## .-

<pre>
.-  (class=arithmetic #args=2) Subtraction, with integer-to-integer overflow.
</pre>


## ./

<pre>
./  (class=arithmetic #args=2) Integer division; not pythonic.
</pre>


## /

<pre>
/  (class=arithmetic #args=2) Division. Integer / integer is floating-point.
</pre>


## //

<pre>
//  (class=arithmetic #args=2) Pythonic integer division, rounding toward negative.
</pre>


## <

<pre>
<  (class=boolean #args=2) String/numeric less-than. Mixing number and string results in string compare.
</pre>


## <<

<pre>
<<  (class=arithmetic #args=2) Bitwise left-shift.
</pre>


## <=

<pre>
<=  (class=boolean #args=2) String/numeric less-than-or-equals. Mixing number and string results in string compare.
</pre>


## ==

<pre>
==  (class=boolean #args=2) String/numeric equality. Mixing number and string results in string compare.
</pre>


## =~

<pre>
=~  (class=boolean #args=2) String (left-hand side) matches regex (right-hand side), e.g. '$name =~ "^a.*b$"'.
</pre>


## >

<pre>
>  (class=boolean #args=2) String/numeric greater-than. Mixing number and string results in string compare.
</pre>


## >=

<pre>
>=  (class=boolean #args=2) String/numeric greater-than-or-equals. Mixing number and string results in string compare.
</pre>


## \>\>

<pre>
>>  (class=arithmetic #args=2) Bitwise signed right-shift.
</pre>


## \>\>\>

<pre>
>>>  (class=arithmetic #args=2) Bitwise unsigned right-shift.
</pre>


## \?

<pre>
?:  (class=boolean #args=3) Standard ternary operator.
</pre>


## ??

<pre>
??  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record.
</pre>


## ???

<pre>
???  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record, or has empty value.
</pre>


## ^

<pre>
^  (class=arithmetic #args=2) Bitwise XOR.
</pre>


## ^^

<pre>
^^  (class=boolean #args=2) Logical XOR.
</pre>


## \|

<pre>
|  (class=arithmetic #args=2) Bitwise OR.
</pre>


## ||

<pre>
||  (class=boolean #args=2) Logical OR.
</pre>


## ~

<pre>
~  (class=arithmetic #args=1) Bitwise NOT. Beware '$y=~$x' since =~ is the
regex-match operator: try '$y = ~$x'.
</pre>


## abs

<pre>
abs  (class=math #args=1) Absolute value.
</pre>


## acos

<pre>
acos  (class=math #args=1) Inverse trigonometric cosine.
</pre>


## acosh

<pre>
acosh  (class=math #args=1) Inverse hyperbolic cosine.
</pre>


## append

<pre>
append  (class=maps/arrays #args=2) Appends second argument to end of first argument, which must be an array.
</pre>


## arrayify

<pre>
arrayify  (class=maps/arrays #args=1) Walks through a nested map/array, converting any map with consecutive keys
"1", "2", ... into an array. Useful to wrap the output of unflatten.
</pre>


## asin

<pre>
asin  (class=math #args=1) Inverse trigonometric sine.
</pre>


## asinh

<pre>
asinh  (class=math #args=1) Inverse hyperbolic sine.
</pre>


## asserting_absent

<pre>
asserting_absent  (class=typing #args=1) Aborts with an error if is_absent on the argument returns false,
else returns its argument.
</pre>


## asserting_array

<pre>
asserting_array  (class=typing #args=1) Aborts with an error if is_array on the argument returns false,
else returns its argument.
</pre>


## asserting_bool

<pre>
asserting_bool  (class=typing #args=1) Aborts with an error if is_bool on the argument returns false,
else returns its argument.
</pre>


## asserting_boolean

<pre>
asserting_boolean  (class=typing #args=1) Aborts with an error if is_boolean on the argument returns false,
else returns its argument.
</pre>


## asserting_empty

<pre>
asserting_empty  (class=typing #args=1) Aborts with an error if is_empty on the argument returns false,
else returns its argument.
</pre>


## asserting_empty_map

<pre>
asserting_empty_map  (class=typing #args=1) Aborts with an error if is_empty_map on the argument returns false,
else returns its argument.
</pre>


## asserting_error

<pre>
asserting_error  (class=typing #args=1) Aborts with an error if is_error on the argument returns false,
else returns its argument.
</pre>


## asserting_float

<pre>
asserting_float  (class=typing #args=1) Aborts with an error if is_float on the argument returns false,
else returns its argument.
</pre>


## asserting_int

<pre>
asserting_int  (class=typing #args=1) Aborts with an error if is_int on the argument returns false,
else returns its argument.
</pre>


## asserting_map

<pre>
asserting_map  (class=typing #args=1) Aborts with an error if is_map on the argument returns false,
else returns its argument.
</pre>


## asserting_nonempty_map

<pre>
asserting_nonempty_map  (class=typing #args=1) Aborts with an error if is_nonempty_map on the argument returns false,
else returns its argument.
</pre>


## asserting_not_array

<pre>
asserting_not_array  (class=typing #args=1) Aborts with an error if is_not_array on the argument returns false,
else returns its argument.
</pre>


## asserting_not_empty

<pre>
asserting_not_empty  (class=typing #args=1) Aborts with an error if is_not_empty on the argument returns false,
else returns its argument.
</pre>


## asserting_not_map

<pre>
asserting_not_map  (class=typing #args=1) Aborts with an error if is_not_map on the argument returns false,
else returns its argument.
</pre>


## asserting_not_null

<pre>
asserting_not_null  (class=typing #args=1) Aborts with an error if is_not_null on the argument returns false,
else returns its argument.
</pre>


## asserting_null

<pre>
asserting_null  (class=typing #args=1) Aborts with an error if is_null on the argument returns false,
else returns its argument.
</pre>


## asserting_numeric

<pre>
asserting_numeric  (class=typing #args=1) Aborts with an error if is_numeric on the argument returns false,
else returns its argument.
</pre>


## asserting_present

<pre>
asserting_present  (class=typing #args=1) Aborts with an error if is_present on the argument returns false,
else returns its argument.
</pre>


## asserting_string

<pre>
asserting_string  (class=typing #args=1) Aborts with an error if is_string on the argument returns false,
else returns its argument.
</pre>


## atan

<pre>
atan  (class=math #args=1) One-argument arctangent.
</pre>


## atan2

<pre>
atan2  (class=math #args=2) Two-argument arctangent.
</pre>


## atanh

<pre>
atanh  (class=math #args=1) Inverse hyperbolic tangent.
</pre>


## bitcount

<pre>
bitcount  (class=arithmetic #args=1) Count of 1-bits.
</pre>


## boolean

<pre>
boolean  (class=conversion #args=1) Convert int/float/bool/string to boolean.
</pre>


## capitalize

<pre>
capitalize  (class=string #args=1) Convert string's first character to uppercase.
</pre>


## cbrt

<pre>
cbrt  (class=math #args=1) Cube root.
</pre>


## ceil

<pre>
ceil  (class=math #args=1) Ceiling: nearest integer at or above.
</pre>


## clean_whitespace

<pre>
clean_whitespace  (class=string #args=1) Same as collapse_whitespace and strip.
</pre>


## collapse_whitespace

<pre>
collapse_whitespace  (class=string #args=1) Strip repeated whitespace from string.
</pre>


## cos

<pre>
cos  (class=math #args=1) Trigonometric cosine.
</pre>


## cosh

<pre>
cosh  (class=math #args=1) Hyperbolic cosine.
</pre>


## depth

<pre>
depth  (class=maps/arrays #args=1) Prints maximum depth of map/array. Scalars have depth 0.
</pre>


## dhms2fsec

<pre>
dhms2fsec  (class=time #args=1) Recovers floating-point seconds as in dhms2fsec("5d18h53m20.250000s") = 500000.250000
</pre>


## dhms2sec

<pre>
dhms2sec  (class=time #args=1) Recovers integer seconds as in dhms2sec("5d18h53m20s") = 500000
</pre>


## erf

<pre>
erf  (class=math #args=1) Error function.
</pre>


## erfc

<pre>
erfc  (class=math #args=1) Complementary error function.
</pre>


## exp

<pre>
exp  (class=math #args=1) Exponential function e**x.
</pre>


## expm1

<pre>
expm1  (class=math #args=1) e**x - 1.
</pre>


## flatten

<pre>
flatten  (class=maps/arrays #args=3) Flattens multi-level maps to single-level ones. Examples:
flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
Two-argument version: flatten($*, ".") is the same as flatten("", ".", $*).
Useful for nested JSON-like structures for non-JSON file formats like CSV.
</pre>


## float

<pre>
float  (class=conversion #args=1) Convert int/float/bool/string to float.
</pre>


## floor

<pre>
floor  (class=math #args=1) Floor: nearest integer at or below.
</pre>


## fmtnum

<pre>
fmtnum  (class=conversion #args=2) Convert int/float/bool to string using
printf-style format string, e.g. '$s = fmtnum($n, "%06lld")'.
</pre>


## fsec2dhms

<pre>
fsec2dhms  (class=time #args=1) Formats floating-point seconds as in fsec2dhms(500000.25) = "5d18h53m20.250000s"
</pre>


## fsec2hms

<pre>
fsec2hms  (class=time #args=1) Formats floating-point seconds as in fsec2hms(5000.25) = "01:23:20.250000"
</pre>


## get_keys

<pre>
get_keys  (class=maps/arrays #args=1) Returns array of keys of map or array
</pre>


## get_values

<pre>
get_values  (class=maps/arrays #args=1) Returns array of keys of map or array -- in the latter case, returns a copy of the array
</pre>


## gmt2sec

<pre>
gmt2sec  (class=time #args=1) Parses GMT timestamp as integer seconds since the epoch.
</pre>


## gsub

<pre>
gsub  (class=string #args=3) Example: '$name=gsub($name, "old", "new")' (replace all).
</pre>


## haskey

<pre>
haskey  (class=maps/arrays #args=2) True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or
'haskey(mymap, mykey)', or true/false if array index is in bounds / out of bounds.
Error if 1st argument is not a map or array. Note -n..-1 alias to 1..n in Miller arrays.
</pre>


## hexfmt

<pre>
hexfmt  (class=conversion #args=1) Convert int to hex string, e.g. 255 to "0xff".
</pre>


## hms2fsec

<pre>
hms2fsec  (class=time #args=1) Recovers floating-point seconds as in hms2fsec("01:23:20.250000") = 5000.250000
</pre>


## hms2sec

<pre>
hms2sec  (class=time #args=1) Recovers integer seconds as in hms2sec("01:23:20") = 5000
</pre>


## hostname

<pre>
hostname  (class=system #args=0) Returns the hostname as a string.
</pre>


## int

<pre>
int  (class=conversion #args=1) Convert int/float/bool/string to int.
</pre>


## invqnorm

<pre>
invqnorm  (class=math #args=1) Inverse of normal cumulative distribution function.
Note that invqorm(urand()) is normally distributed.
</pre>


## is_absent

<pre>
is_absent  (class=typing #args=1) False if field is present in input, true otherwise
</pre>


## is_array

<pre>
is_array  (class=typing #args=1) True if argument is an array.
</pre>


## is_bool

<pre>
is_bool  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_boolean.
</pre>


## is_boolean

<pre>
is_boolean  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_bool.
</pre>


## is_empty

<pre>
is_empty  (class=typing #args=1) True if field is present in input with empty string value, false otherwise.
</pre>


## is_empty_map

<pre>
is_empty_map  (class=typing #args=1) True if argument is a map which is empty.
</pre>


## is_error

<pre>
is_error  (class=typing #args=1) True if if argument is an error, such as taking string length of an integer.
</pre>


## is_float

<pre>
is_float  (class=typing #args=1) True if field is present with value inferred to be float
</pre>


## is_int

<pre>
is_int  (class=typing #args=1) True if field is present with value inferred to be int
</pre>


## is_map

<pre>
is_map  (class=typing #args=1) True if argument is a map.
</pre>


## is_nonempty_map

<pre>
is_nonempty_map  (class=typing #args=1) True if argument is a map which is non-empty.
</pre>


## is_not_array

<pre>
is_not_array  (class=typing #args=1) True if argument is not an array.
</pre>


## is_not_empty

<pre>
is_not_empty  (class=typing #args=1) False if field is present in input with empty value, true otherwise
</pre>


## is_not_map

<pre>
is_not_map  (class=typing #args=1) True if argument is not a map.
</pre>


## is_not_null

<pre>
is_not_null  (class=typing #args=1) False if argument is null (empty or absent), true otherwise.
</pre>


## is_null

<pre>
is_null  (class=typing #args=1) True if argument is null (empty or absent), false otherwise.
</pre>


## is_numeric

<pre>
is_numeric  (class=typing #args=1) True if field is present with value inferred to be int or float
</pre>


## is_present

<pre>
is_present  (class=typing #args=1) True if field is present in input, false otherwise.
</pre>


## is_string

<pre>
is_string  (class=typing #args=1) True if field is present with string (including empty-string) value
</pre>


## joink

<pre>
joink  (class=conversion #args=2) Makes string from map/array keys. Examples:
joink({"a":3,"b":4,"c":5}, ",") = "a,b,c"
joink([1,2,3], ",") = "1,2,3".
</pre>


## joinkv

<pre>
joinkv  (class=conversion #args=3) Makes string from map/array key-value pairs. Examples:
joinkv([3,4,5], "=", ",") = "1=3,2=4,3=5"
joinkv({"a":3,"b":4,"c":5}, "=", ",") = "a=3,b=4,c=5"
</pre>


## joinv

<pre>
joinv  (class=conversion #args=2) Makes string from map/array values.
joinv([3,4,5], ",") = "3,4,5"
joinv({"a":3,"b":4,"c":5}, ",") = "3,4,5"
</pre>


## json_parse

<pre>
json_parse  (class=maps/arrays #args=1) Converts value from JSON-formatted string.
</pre>


## json_stringify

<pre>
json_stringify  (class=maps/arrays #args=1,2) Converts value to JSON-formatted string. Default output is single-line.
With optional second boolean argument set to true, produces multiline output.
</pre>


## leafcount

<pre>
leafcount  (class=maps/arrays #args=1) Counts total number of terminal values in map/array. For single-level
map/array, same as length.
</pre>


## length

<pre>
length  (class=maps/arrays #args=1) Counts number of top-level entries in array/map. Scalars have length 1.
</pre>


## log

<pre>
log  (class=math #args=1) Natural (base-e) logarithm.
</pre>


## log10

<pre>
log10  (class=math #args=1) Base-10 logarithm.
</pre>


## log1p

<pre>
log1p  (class=math #args=1) log(1-x).
</pre>


## logifit

<pre>
logifit  (class=math #args=3)  Given m and b from logistic regression, compute fit:
$yhat=logifit($x,$m,$b).
</pre>


## lstrip

<pre>
lstrip  (class=string #args=1) Strip leading whitespace from string.
</pre>


## madd

<pre>
madd  (class=arithmetic #args=3) a + b mod m (integers)
</pre>


## mapdiff

<pre>
mapdiff  (class=maps/arrays #args=variadic) With 0 args, returns empty map. With 1 arg, returns copy of arg.
With 2 or more, returns copy of arg 1 with all keys from any of remaining
argument maps removed.
</pre>


## mapexcept

<pre>
mapexcept  (class=maps/arrays #args=variadic) Returns a map with keys from remaining arguments, if any, unset.
Remaining arguments can be strings or arrays of string.
E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'
and  'mapexcept({1:2,3:4,5:6}, [1, 5, 7])' is '{3:4}'.
</pre>


## mapselect

<pre>
mapselect  (class=maps/arrays #args=variadic) Returns a map with only keys from remaining arguments set.
Remaining arguments can be strings or arrays of string.
E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'
and  'mapselect({1:2,3:4,5:6}, [1, 5, 7])' is '{1:2,5:6}'.
</pre>


## mapsum

<pre>
mapsum  (class=maps/arrays #args=variadic) With 0 args, returns empty map. With >= 1 arg, returns a map with
key-value pairs from all arguments. Rightmost collisions win, e.g.
'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.
</pre>


## max

<pre>
max  (class=math #args=variadic) Max of n numbers; null loses.
</pre>


## md5

<pre>
md5  (class=hashing #args=1) MD5 hash.
</pre>


## mexp

<pre>
mexp  (class=arithmetic #args=3) a ** b mod m (integers)
</pre>


## min

<pre>
min  (class=math #args=variadic) Min of n numbers; null loses.
</pre>


## mmul

<pre>
mmul  (class=arithmetic #args=3) a * b mod m (integers)
</pre>


## msub

<pre>
msub  (class=arithmetic #args=3) a - b mod m (integers)
</pre>


## os

<pre>
os  (class=system #args=0) Returns the operating-system name as a string.
</pre>


## pow

<pre>
pow  (class=arithmetic #args=2) Exponentiation. Same as **, but as a function.
</pre>


## qnorm

<pre>
qnorm  (class=math #args=1) Normal cumulative distribution function.
</pre>


## regextract

<pre>
regextract  (class=string #args=2) Example: '$name=regextract($name, "[A-Z]{3}[0-9]{2}")'
</pre>


## regextract_or_else

<pre>
regextract_or_else  (class=string #args=3) Example: '$name=regextract_or_else($name, "[A-Z]{3}[0-9]{2}", "default")'
</pre>


## round

<pre>
round  (class=math #args=1) Round to nearest integer.
</pre>


## roundm

<pre>
roundm  (class=math #args=2) Round to nearest multiple of m: roundm($x,$m) is
the same as round($x/$m)*$m.
</pre>


## rstrip

<pre>
rstrip  (class=string #args=1) Strip trailing whitespace from string.
</pre>


## sec2dhms

<pre>
sec2dhms  (class=time #args=1) Formats integer seconds as in sec2dhms(500000) = "5d18h53m20s"
</pre>


## sec2gmt

<pre>
sec2gmt  (class=time #args=1,2) Formats seconds since epoch (integer part)
as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
Leaves non-numbers as-is. With second integer argument n, includes n decimal places
for the seconds part
</pre>


## sec2gmtdate

<pre>
sec2gmtdate  (class=time #args=1) Formats seconds since epoch (integer part)
as GMT timestamp with year-month-date, e.g. sec2gmtdate(1440768801.7) = "2015-08-28".
Leaves non-numbers as-is.
</pre>


## sec2hms

<pre>
sec2hms  (class=time #args=1) Formats integer seconds as in sec2hms(5000) = "01:23:20"
</pre>


## sgn

<pre>
sgn  (class=math #args=1)  +1, 0, -1 for positive, zero, negative input respectively.
</pre>


## sha1

<pre>
sha1  (class=hashing #args=1) SHA1 hash.
</pre>


## sha256

<pre>
sha256  (class=hashing #args=1) SHA256 hash.
</pre>


## sha512

<pre>
sha512  (class=hashing #args=1) SHA512 hash.
</pre>


## sin

<pre>
sin  (class=math #args=1) Trigonometric sine.
</pre>


## sinh

<pre>
sinh  (class=math #args=1) Hyperbolic sine.
</pre>


## splita

<pre>
splita  (class=conversion #args=2) Splits string into array with type inference. Example:
splita("3,4,5", ",") = [3,4,5]
</pre>


## splitax

<pre>
splitax  (class=conversion #args=2) Splits string into array without type inference. Example:
splita("3,4,5", ",") = ["3","4","5"]
</pre>


## splitkv

<pre>
splitkv  (class=conversion #args=3) Splits string by separators into map with type inference. Example:
splitkv("a=3,b=4,c=5", "=", ",") = {"a":3,"b":4,"c":5}
</pre>


## splitkvx

<pre>
splitkvx  (class=conversion #args=3) Splits string by separators into map without type inference (keys and
values are strings). Example:
splitkvx("a=3,b=4,c=5", "=", ",") = {"a":"3","b":"4","c":"5"}
</pre>


## splitnv

<pre>
splitnv  (class=conversion #args=2) Splits string by separator into integer-indexed map with type inference. Example:
splitnv("a,b,c", ",") = {"1":"a","2":"b","3":"c"}
</pre>


## splitnvx

<pre>
splitnvx  (class=conversion #args=2) Splits string by separator into integer-indexed map without type
inference (values are strings). Example:
splitnvx("3,4,5", ",") = {"1":"3","2":"4","3":"5"}
</pre>


## sqrt

<pre>
sqrt  (class=math #args=1) Square root.
</pre>


## ssub

<pre>
ssub  (class=string #args=3) Like sub but does no regexing. No characters are special.
</pre>


## strftime

<pre>
strftime  (class=time #args=2)  Formats seconds since the epoch as timestamp, e.g.
	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%SZ") = "2015-08-28T13:33:21Z", and
	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%3SZ") = "2015-08-28T13:33:21.700Z".
	Format strings are as in the C library (please see "man strftime" on your system),
	with the Miller-specific addition of "%1S" through "%9S" which format the seconds
	with 1 through 9 decimal places, respectively. ("%S" uses no decimal places.)
	See also strftime_local.
</pre>


## string

<pre>
string  (class=conversion #args=1) Convert int/float/bool/string/array/map to string.
</pre>


## strip

<pre>
strip  (class=string #args=1) Strip leading and trailing whitespace from string.
</pre>


## strlen

<pre>
strlen  (class=string #args=1) String length.
</pre>


## strptime

<pre>
strptime  (class=time #args=2) strptime: Parses timestamp as floating-point seconds since the epoch,
	e.g. strptime("2015-08-28T13:33:21Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.000000,
	and  strptime("2015-08-28T13:33:21.345Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.345000.
	See also strptime_local.
</pre>


## sub

<pre>
sub  (class=string #args=3) Example: '$name=sub($name, "old", "new")' (replace once).
</pre>


## substr

<pre>
substr  (class=string #args=3) substr is an alias for substr0. See also substr1. Miller is generally 1-up
with all array indices, but, this is a backward-compatibility issue with Miller 5 and below.
Arrays are new in Miller 6; the substr function is older.
</pre>


## substr0

<pre>
substr0  (class=string #args=3) substr0(s,m,n) gives substring of s from 0-up position m to n
inclusive. Negative indices -len .. -1 alias to 0 .. len-1.
</pre>


## substr1

<pre>
substr1  (class=string #args=3) substr1(s,m,n) gives substring of s from 1-up position m to n
inclusive. Negative indices -len .. -1 alias to 1 .. len.
</pre>


## system

<pre>
system  (class=system #args=1) Run command string, yielding its stdout minus final carriage return.
</pre>


## systime

<pre>
systime  (class=time #args=0) help string will go here
</pre>


## systimeint

<pre>
systimeint  (class=time #args=0) help string will go here
</pre>


## tan

<pre>
tan  (class=math #args=1) Trigonometric tangent.
</pre>


## tanh

<pre>
tanh  (class=math #args=1) Hyperbolic tangent.
</pre>


## tolower

<pre>
tolower  (class=string #args=1) Convert string to lowercase.
</pre>


## toupper

<pre>
toupper  (class=string #args=1) Convert string to uppercase.
</pre>


## truncate

<pre>
truncate  (class=string #args=2) Truncates string first argument to max length of int second argument.
</pre>


## typeof

<pre>
typeof  (class=typing #args=1) Convert argument to type of argument (e.g. "str"). For debug.
</pre>


## unflatten

<pre>
unflatten  (class=maps/arrays #args=2) Reverses flatten. Example:
unflatten({"a.b.c" : 4}, ".") is {"a": "b": { "c": 4 }}.
Useful for nested JSON-like structures for non-JSON file formats like CSV.
See also arrayify.
</pre>


## uptime

<pre>
uptime  (class=time #args=0) help string will go here
</pre>


## urand

<pre>
urand  (class=math #args=0) Floating-point numbers uniformly distributed on the unit interval.
Int-valued example: '$n=floor(20+urand()*11)'.
</pre>


## urand32

<pre>
urand32  (class=math #args=0) Integer uniformly distributed 0 and 2**32-1 inclusive.
</pre>


## urandint

<pre>
urandint  (class=math #args=2) Integer uniformly distributed between inclusive integer endpoints.
</pre>


## urandrange

<pre>
urandrange  (class=math #args=2) Floating-point numbers uniformly distributed on the interval [a, b).
</pre>


## version

<pre>
version  (class=system #args=0) Returns the Miller version as a string.
</pre>

