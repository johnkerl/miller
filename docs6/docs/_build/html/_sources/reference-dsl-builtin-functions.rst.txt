..
    PLEASE DO NOT EDIT DIRECTLY. EDIT THE .rst.in FILE PLEASE.

DSL reference: built-in functions
=====================================

Summary
----------------------------------------------------------------

mlr: option "--list-all-functions-as-table" not recognized.
Please run "mlr --help" for usage information.

List of functions
----------------------------------------------------------------

Each function takes a specific number of arguments, as shown below, except for functions marked as variadic such as ``min`` and ``max``. (The latter compute min and max of any number of numerical arguments.) There is no notion of optional or default-on-absent arguments. All argument-passing is positional rather than by name; arguments are passed by value, not by reference.

You can get a list of all functions using **mlr -f**, with details using **mlr -F**.


.. _reference-dsl-colon:

\!
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    !  (class=boolean #args=1) Logical negation.



.. _reference-dsl-!=:

!=
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    !=  (class=boolean #args=2) String/numeric inequality. Mixing number and string results in string compare.



.. _reference-dsl-!=~:

!=~
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    !=~  (class=boolean #args=2) String (left-hand side) does not match regex (right-hand side), e.g. '$name !=~ "^a.*b$"'.



.. _reference-dsl-%:

%
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    %  (class=arithmetic #args=2) Remainder; never negative-valued (pythonic).



.. _reference-dsl-&:

&
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    &  (class=arithmetic #args=2) Bitwise AND.



.. _reference-dsl-&&:

&&
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    &&  (class=boolean #args=2) Logical AND.



.. _reference-dsl-times:

\*
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    *  (class=arithmetic #args=2) Multiplication, with integer*integer overflow to float.



.. _reference-dsl-exponentiation:

\**
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    **  (class=arithmetic #args=2) Exponentiation. Same as pow, but as an infix operator.



.. _reference-dsl-plus:

\+
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    +  (class=arithmetic #args=1,2) Addition as binary operator; unary plus operator.



.. _reference-dsl-minus:

\-
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    -  (class=arithmetic #args=1,2) Subtraction as binary operator; unary negation operator.



.. _reference-dsl-.:

.
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    .  (class=string #args=2) String concatenation.



.. _reference-dsl-.*:

.*
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    .*  (class=arithmetic #args=2) Multiplication, with integer-to-integer overflow.



.. _reference-dsl-.+:

.+
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    .+  (class=arithmetic #args=2) Addition, with integer-to-integer overflow.



.. _reference-dsl-.-:

.-
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    .-  (class=arithmetic #args=2) Subtraction, with integer-to-integer overflow.



.. _reference-dsl-./:

./
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ./  (class=arithmetic #args=2) Integer division; not pythonic.



.. _reference-dsl-/:

/
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    /  (class=arithmetic #args=2) Division. Integer / integer is floating-point.



.. _reference-dsl-//:

//
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    //  (class=arithmetic #args=2) Pythonic integer division, rounding toward negative.



.. _reference-dsl-<:

<
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    <  (class=boolean #args=2) String/numeric less-than. Mixing number and string results in string compare.



.. _reference-dsl-<<:

<<
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    <<  (class=arithmetic #args=2) Bitwise left-shift.



.. _reference-dsl-<=:

<=
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    <=  (class=boolean #args=2) String/numeric less-than-or-equals. Mixing number and string results in string compare.



.. _reference-dsl-==:

==
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ==  (class=boolean #args=2) String/numeric equality. Mixing number and string results in string compare.



.. _reference-dsl-=~:

=~
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    =~  (class=boolean #args=2) String (left-hand side) matches regex (right-hand side), e.g. '$name =~ "^a.*b$"'.



.. _reference-dsl->:

>
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    >  (class=boolean #args=2) String/numeric greater-than. Mixing number and string results in string compare.



.. _reference-dsl->=:

>=
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    >=  (class=boolean #args=2) String/numeric greater-than-or-equals. Mixing number and string results in string compare.



.. _reference-dsl-srsh:

\>\>
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    >>  (class=arithmetic #args=2) Bitwise signed right-shift.



.. _reference-dsl-ursh:

\>\>\>
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    >>>  (class=arithmetic #args=2) Bitwise unsigned right-shift.



.. _reference-dsl-question-mark-colon:

\?
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ?:  (class=boolean #args=3) Standard ternary operator.



.. _reference-dsl-??:

??
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ??  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record.



.. _reference-dsl-???:

???
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ???  (class=boolean #args=2) Absent-coalesce operator. $a ?? 1 evaluates to 1 if $a isn't defined in the current record, or has empty value.



.. _reference-dsl-^:

^
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ^  (class=arithmetic #args=2) Bitwise XOR.



.. _reference-dsl-^^:

^^
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ^^  (class=boolean #args=2) Logical XOR.



.. _reference-dsl-bitwise-or:

\|
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    |  (class=arithmetic #args=2) Bitwise OR.



.. _reference-dsl-||:

||
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ||  (class=boolean #args=2) Logical OR.



.. _reference-dsl-~:

~
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ~  (class=arithmetic #args=1) Bitwise NOT. Beware '$y=~$x' since =~ is the
    regex-match operator: try '$y = ~$x'.



.. _reference-dsl-abs:

abs
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    abs  (class=math #args=1) Absolute value.



.. _reference-dsl-acos:

acos
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    acos  (class=math #args=1) Inverse trigonometric cosine.



.. _reference-dsl-acosh:

acosh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    acosh  (class=math #args=1) Inverse hyperbolic cosine.



.. _reference-dsl-append:

append
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    append  (class=maps/arrays #args=2) Appends second argument to end of first argument, which must be an array.



.. _reference-dsl-arrayify:

arrayify
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    arrayify  (class=maps/arrays #args=1) Walks through a nested map/array, converting any map with consecutive keys
    "1", "2", ... into an array. Useful to wrap the output of unflatten.



.. _reference-dsl-asin:

asin
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asin  (class=math #args=1) Inverse trigonometric sine.



.. _reference-dsl-asinh:

asinh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asinh  (class=math #args=1) Inverse hyperbolic sine.



.. _reference-dsl-asserting_absent:

asserting_absent
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_absent  (class=typing #args=1) Aborts with an error if is_absent on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_array:

asserting_array
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_array  (class=typing #args=1) Aborts with an error if is_array on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_bool:

asserting_bool
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_bool  (class=typing #args=1) Aborts with an error if is_bool on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_boolean:

asserting_boolean
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_boolean  (class=typing #args=1) Aborts with an error if is_boolean on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_empty:

asserting_empty
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_empty  (class=typing #args=1) Aborts with an error if is_empty on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_empty_map:

asserting_empty_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_empty_map  (class=typing #args=1) Aborts with an error if is_empty_map on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_error:

asserting_error
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_error  (class=typing #args=1) Aborts with an error if is_error on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_float:

asserting_float
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_float  (class=typing #args=1) Aborts with an error if is_float on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_int:

asserting_int
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_int  (class=typing #args=1) Aborts with an error if is_int on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_map:

asserting_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_map  (class=typing #args=1) Aborts with an error if is_map on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_nonempty_map:

asserting_nonempty_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_nonempty_map  (class=typing #args=1) Aborts with an error if is_nonempty_map on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_not_array:

asserting_not_array
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_not_array  (class=typing #args=1) Aborts with an error if is_not_array on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_not_empty:

asserting_not_empty
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_not_empty  (class=typing #args=1) Aborts with an error if is_not_empty on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_not_map:

asserting_not_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_not_map  (class=typing #args=1) Aborts with an error if is_not_map on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_not_null:

asserting_not_null
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_not_null  (class=typing #args=1) Aborts with an error if is_not_null on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_null:

asserting_null
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_null  (class=typing #args=1) Aborts with an error if is_null on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_numeric:

asserting_numeric
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_numeric  (class=typing #args=1) Aborts with an error if is_numeric on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_present:

asserting_present
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_present  (class=typing #args=1) Aborts with an error if is_present on the argument returns false,
    else returns its argument.



.. _reference-dsl-asserting_string:

asserting_string
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    asserting_string  (class=typing #args=1) Aborts with an error if is_string on the argument returns false,
    else returns its argument.



.. _reference-dsl-atan:

atan
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    atan  (class=math #args=1) One-argument arctangent.



.. _reference-dsl-atan2:

atan2
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    atan2  (class=math #args=2) Two-argument arctangent.



.. _reference-dsl-atanh:

atanh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    atanh  (class=math #args=1) Inverse hyperbolic tangent.



.. _reference-dsl-bitcount:

bitcount
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    bitcount  (class=arithmetic #args=1) Count of 1-bits.



.. _reference-dsl-boolean:

boolean
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    boolean  (class=conversion #args=1) Convert int/float/bool/string to boolean.



.. _reference-dsl-capitalize:

capitalize
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    capitalize  (class=string #args=1) Convert string's first character to uppercase.



.. _reference-dsl-cbrt:

cbrt
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    cbrt  (class=math #args=1) Cube root.



.. _reference-dsl-ceil:

ceil
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ceil  (class=math #args=1) Ceiling: nearest integer at or above.



.. _reference-dsl-clean_whitespace:

clean_whitespace
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    clean_whitespace  (class=string #args=1) Same as collapse_whitespace and strip.



.. _reference-dsl-collapse_whitespace:

collapse_whitespace
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    collapse_whitespace  (class=string #args=1) Strip repeated whitespace from string.



.. _reference-dsl-cos:

cos
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    cos  (class=math #args=1) Trigonometric cosine.



.. _reference-dsl-cosh:

cosh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    cosh  (class=math #args=1) Hyperbolic cosine.



.. _reference-dsl-depth:

depth
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    depth  (class=maps/arrays #args=1) Prints maximum depth of map/array. Scalars have depth 0.



.. _reference-dsl-dhms2fsec:

dhms2fsec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    dhms2fsec  (class=time #args=1) Recovers floating-point seconds as in dhms2fsec("5d18h53m20.250000s") = 500000.250000



.. _reference-dsl-dhms2sec:

dhms2sec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    dhms2sec  (class=time #args=1) Recovers integer seconds as in dhms2sec("5d18h53m20s") = 500000



.. _reference-dsl-erf:

erf
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    erf  (class=math #args=1) Error function.



.. _reference-dsl-erfc:

erfc
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    erfc  (class=math #args=1) Complementary error function.



.. _reference-dsl-exp:

exp
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    exp  (class=math #args=1) Exponential function e**x.



.. _reference-dsl-expm1:

expm1
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    expm1  (class=math #args=1) e**x - 1.



.. _reference-dsl-flatten:

flatten
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    flatten  (class=maps/arrays #args=3) Flattens multi-level maps to single-level ones. Examples:
    flatten("a", ".", {"b": { "c": 4 }}) is {"a.b.c" : 4}.
    flatten("", ".", {"a": { "b": 3 }}) is {"a.b" : 3}.
    Two-argument version: flatten($*, ".") is the same as flatten("", ".", $*).
    Useful for nested JSON-like structures for non-JSON file formats like CSV.



.. _reference-dsl-float:

float
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    float  (class=conversion #args=1) Convert int/float/bool/string to float.



.. _reference-dsl-floor:

floor
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    floor  (class=math #args=1) Floor: nearest integer at or below.



.. _reference-dsl-fmtnum:

fmtnum
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    fmtnum  (class=conversion #args=2) Convert int/float/bool to string using
    printf-style format string, e.g. '$s = fmtnum($n, "%06lld")'.



.. _reference-dsl-fsec2dhms:

fsec2dhms
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    fsec2dhms  (class=time #args=1) Formats floating-point seconds as in fsec2dhms(500000.25) = "5d18h53m20.250000s"



.. _reference-dsl-fsec2hms:

fsec2hms
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    fsec2hms  (class=time #args=1) Formats floating-point seconds as in fsec2hms(5000.25) = "01:23:20.250000"



.. _reference-dsl-get_keys:

get_keys
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    get_keys  (class=maps/arrays #args=1) Returns array of keys of map or array



.. _reference-dsl-get_values:

get_values
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    get_values  (class=maps/arrays #args=1) Returns array of keys of map or array -- in the latter case, returns a copy of the array



.. _reference-dsl-gmt2sec:

gmt2sec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    gmt2sec  (class=time #args=1) Parses GMT timestamp as integer seconds since the epoch.



.. _reference-dsl-gsub:

gsub
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    gsub  (class=string #args=3) Example: '$name=gsub($name, "old", "new")' (replace all).



.. _reference-dsl-haskey:

haskey
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    haskey  (class=maps/arrays #args=2) True/false if map has/hasn't key, e.g. 'haskey($*, "a")' or
    'haskey(mymap, mykey)', or true/false if array index is in bounds / out of bounds.
    Error if 1st argument is not a map or array. Note -n..-1 alias to 1..n in Miller arrays.



.. _reference-dsl-hexfmt:

hexfmt
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    hexfmt  (class=conversion #args=1) Convert int to hex string, e.g. 255 to "0xff".



.. _reference-dsl-hms2fsec:

hms2fsec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    hms2fsec  (class=time #args=1) Recovers floating-point seconds as in hms2fsec("01:23:20.250000") = 5000.250000



.. _reference-dsl-hms2sec:

hms2sec
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    hms2sec  (class=time #args=1) Recovers integer seconds as in hms2sec("01:23:20") = 5000



.. _reference-dsl-hostname:

hostname
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    hostname  (class=system #args=0) Returns the hostname as a string.



.. _reference-dsl-int:

int
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    int  (class=conversion #args=1) Convert int/float/bool/string to int.



.. _reference-dsl-invqnorm:

invqnorm
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    invqnorm  (class=math #args=1) Inverse of normal cumulative distribution function.
    Note that invqorm(urand()) is normally distributed.



.. _reference-dsl-is_absent:

is_absent
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_absent  (class=typing #args=1) False if field is present in input, true otherwise



.. _reference-dsl-is_array:

is_array
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_array  (class=typing #args=1) True if argument is an array.



.. _reference-dsl-is_bool:

is_bool
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_bool  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_boolean.



.. _reference-dsl-is_boolean:

is_boolean
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_boolean  (class=typing #args=1) True if field is present with boolean value. Synonymous with is_bool.



.. _reference-dsl-is_empty:

is_empty
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_empty  (class=typing #args=1) True if field is present in input with empty string value, false otherwise.



.. _reference-dsl-is_empty_map:

is_empty_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_empty_map  (class=typing #args=1) True if argument is a map which is empty.



.. _reference-dsl-is_error:

is_error
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_error  (class=typing #args=1) True if if argument is an error, such as taking string length of an integer.



.. _reference-dsl-is_float:

is_float
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_float  (class=typing #args=1) True if field is present with value inferred to be float



.. _reference-dsl-is_int:

is_int
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_int  (class=typing #args=1) True if field is present with value inferred to be int



.. _reference-dsl-is_map:

is_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_map  (class=typing #args=1) True if argument is a map.



.. _reference-dsl-is_nonempty_map:

is_nonempty_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_nonempty_map  (class=typing #args=1) True if argument is a map which is non-empty.



.. _reference-dsl-is_not_array:

is_not_array
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_not_array  (class=typing #args=1) True if argument is not an array.



.. _reference-dsl-is_not_empty:

is_not_empty
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_not_empty  (class=typing #args=1) False if field is present in input with empty value, true otherwise



.. _reference-dsl-is_not_map:

is_not_map
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_not_map  (class=typing #args=1) True if argument is not a map.



.. _reference-dsl-is_not_null:

is_not_null
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_not_null  (class=typing #args=1) False if argument is null (empty or absent), true otherwise.



.. _reference-dsl-is_null:

is_null
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_null  (class=typing #args=1) True if argument is null (empty or absent), false otherwise.



.. _reference-dsl-is_numeric:

is_numeric
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_numeric  (class=typing #args=1) True if field is present with value inferred to be int or float



.. _reference-dsl-is_present:

is_present
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_present  (class=typing #args=1) True if field is present in input, false otherwise.



.. _reference-dsl-is_string:

is_string
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    is_string  (class=typing #args=1) True if field is present with string (including empty-string) value



.. _reference-dsl-joink:

joink
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    joink  (class=conversion #args=2) Makes string from map/array keys. Examples:
    joink({"a":3,"b":4,"c":5}, ",") = "a,b,c"
    joink([1,2,3], ",") = "1,2,3".



.. _reference-dsl-joinkv:

joinkv
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    joinkv  (class=conversion #args=3) Makes string from map/array key-value pairs. Examples:
    joinkv([3,4,5], "=", ",") = "1=3,2=4,3=5"
    joinkv({"a":3,"b":4,"c":5}, "=", ",") = "a=3,b=4,c=5"



.. _reference-dsl-joinv:

joinv
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    joinv  (class=conversion #args=2) Makes string from map/array values.
    joinv([3,4,5], ",") = "3,4,5"
    joinv({"a":3,"b":4,"c":5}, ",") = "3,4,5"



.. _reference-dsl-json_parse:

json_parse
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    json_parse  (class=maps/arrays #args=1) Converts value from JSON-formatted string.



.. _reference-dsl-json_stringify:

json_stringify
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    json_stringify  (class=maps/arrays #args=1,2) Converts value to JSON-formatted string. Default output is single-line.
    With optional second boolean argument set to true, produces multiline output.



.. _reference-dsl-leafcount:

leafcount
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    leafcount  (class=maps/arrays #args=1) Counts total number of terminal values in map/array. For single-level
    map/array, same as length.



.. _reference-dsl-length:

length
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    length  (class=maps/arrays #args=1) Counts number of top-level entries in array/map. Scalars have length 1.



.. _reference-dsl-log:

log
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    log  (class=math #args=1) Natural (base-e) logarithm.



.. _reference-dsl-log10:

log10
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    log10  (class=math #args=1) Base-10 logarithm.



.. _reference-dsl-log1p:

log1p
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    log1p  (class=math #args=1) log(1-x).



.. _reference-dsl-logifit:

logifit
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    logifit  (class=math #args=3)  Given m and b from logistic regression, compute fit:
    $yhat=logifit($x,$m,$b).



.. _reference-dsl-lstrip:

lstrip
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    lstrip  (class=string #args=1) Strip leading whitespace from string.



.. _reference-dsl-madd:

madd
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    madd  (class=arithmetic #args=3) a + b mod m (integers)



.. _reference-dsl-mapdiff:

mapdiff
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    mapdiff  (class=maps/arrays #args=variadic) With 0 args, returns empty map. With 1 arg, returns copy of arg.
    With 2 or more, returns copy of arg 1 with all keys from any of remaining
    argument maps removed.



.. _reference-dsl-mapexcept:

mapexcept
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    mapexcept  (class=maps/arrays #args=variadic) Returns a map with keys from remaining arguments, if any, unset.
    Remaining arguments can be strings or arrays of string.
    E.g. 'mapexcept({1:2,3:4,5:6}, 1, 5, 7)' is '{3:4}'
    and  'mapexcept({1:2,3:4,5:6}, [1, 5, 7])' is '{3:4}'.



.. _reference-dsl-mapselect:

mapselect
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    mapselect  (class=maps/arrays #args=variadic) Returns a map with only keys from remaining arguments set.
    Remaining arguments can be strings or arrays of string.
    E.g. 'mapselect({1:2,3:4,5:6}, 1, 5, 7)' is '{1:2,5:6}'
    and  'mapselect({1:2,3:4,5:6}, [1, 5, 7])' is '{1:2,5:6}'.



.. _reference-dsl-mapsum:

mapsum
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    mapsum  (class=maps/arrays #args=variadic) With 0 args, returns empty map. With >= 1 arg, returns a map with
    key-value pairs from all arguments. Rightmost collisions win, e.g.
    'mapsum({1:2,3:4},{1:5})' is '{1:5,3:4}'.



.. _reference-dsl-max:

max
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    max  (class=math #args=variadic) Max of n numbers; null loses.



.. _reference-dsl-md5:

md5
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    md5  (class=hashing #args=1) MD5 hash.



.. _reference-dsl-mexp:

mexp
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    mexp  (class=arithmetic #args=3) a ** b mod m (integers)



.. _reference-dsl-min:

min
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    min  (class=math #args=variadic) Min of n numbers; null loses.



.. _reference-dsl-mmul:

mmul
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    mmul  (class=arithmetic #args=3) a * b mod m (integers)



.. _reference-dsl-msub:

msub
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    msub  (class=arithmetic #args=3) a - b mod m (integers)



.. _reference-dsl-os:

os
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    os  (class=system #args=0) Returns the operating-system name as a string.



.. _reference-dsl-pow:

pow
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    pow  (class=arithmetic #args=2) Exponentiation. Same as **, but as a function.



.. _reference-dsl-qnorm:

qnorm
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    qnorm  (class=math #args=1) Normal cumulative distribution function.



.. _reference-dsl-regextract:

regextract
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    regextract  (class=string #args=2) Example: '$name=regextract($name, "[A-Z]{3}[0-9]{2}")'



.. _reference-dsl-regextract_or_else:

regextract_or_else
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    regextract_or_else  (class=string #args=3) Example: '$name=regextract_or_else($name, "[A-Z]{3}[0-9]{2}", "default")'



.. _reference-dsl-round:

round
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    round  (class=math #args=1) Round to nearest integer.



.. _reference-dsl-roundm:

roundm
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    roundm  (class=math #args=2) Round to nearest multiple of m: roundm($x,$m) is
    the same as round($x/$m)*$m.



.. _reference-dsl-rstrip:

rstrip
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    rstrip  (class=string #args=1) Strip trailing whitespace from string.



.. _reference-dsl-sec2dhms:

sec2dhms
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sec2dhms  (class=time #args=1) Formats integer seconds as in sec2dhms(500000) = "5d18h53m20s"



.. _reference-dsl-sec2gmt:

sec2gmt
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sec2gmt  (class=time #args=1,2) Formats seconds since epoch (integer part)
    as GMT timestamp, e.g. sec2gmt(1440768801.7) = "2015-08-28T13:33:21Z".
    Leaves non-numbers as-is. With second integer argument n, includes n decimal places
    for the seconds part



.. _reference-dsl-sec2gmtdate:

sec2gmtdate
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sec2gmtdate  (class=time #args=1) Formats seconds since epoch (integer part)
    as GMT timestamp with year-month-date, e.g. sec2gmtdate(1440768801.7) = "2015-08-28".
    Leaves non-numbers as-is.



.. _reference-dsl-sec2hms:

sec2hms
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sec2hms  (class=time #args=1) Formats integer seconds as in sec2hms(5000) = "01:23:20"



.. _reference-dsl-sgn:

sgn
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sgn  (class=math #args=1)  +1, 0, -1 for positive, zero, negative input respectively.



.. _reference-dsl-sha1:

sha1
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sha1  (class=hashing #args=1) SHA1 hash.



.. _reference-dsl-sha256:

sha256
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sha256  (class=hashing #args=1) SHA256 hash.



.. _reference-dsl-sha512:

sha512
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sha512  (class=hashing #args=1) SHA512 hash.



.. _reference-dsl-sin:

sin
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sin  (class=math #args=1) Trigonometric sine.



.. _reference-dsl-sinh:

sinh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sinh  (class=math #args=1) Hyperbolic sine.



.. _reference-dsl-splita:

splita
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    splita  (class=conversion #args=2) Splits string into array with type inference. Example:
    splita("3,4,5", ",") = [3,4,5]



.. _reference-dsl-splitax:

splitax
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    splitax  (class=conversion #args=2) Splits string into array without type inference. Example:
    splita("3,4,5", ",") = ["3","4","5"]



.. _reference-dsl-splitkv:

splitkv
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    splitkv  (class=conversion #args=3) Splits string by separators into map with type inference. Example:
    splitkv("a=3,b=4,c=5", "=", ",") = {"a":3,"b":4,"c":5}



.. _reference-dsl-splitkvx:

splitkvx
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    splitkvx  (class=conversion #args=3) Splits string by separators into map without type inference (keys and
    values are strings). Example:
    splitkvx("a=3,b=4,c=5", "=", ",") = {"a":"3","b":"4","c":"5"}



.. _reference-dsl-splitnv:

splitnv
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    splitnv  (class=conversion #args=2) Splits string by separator into integer-indexed map with type inference. Example:
    splitnv("a,b,c", ",") = {"1":"a","2":"b","3":"c"}



.. _reference-dsl-splitnvx:

splitnvx
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    splitnvx  (class=conversion #args=2) Splits string by separator into integer-indexed map without type
    inference (values are strings). Example:
    splitnvx("3,4,5", ",") = {"1":"3","2":"4","3":"5"}



.. _reference-dsl-sqrt:

sqrt
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sqrt  (class=math #args=1) Square root.



.. _reference-dsl-ssub:

ssub
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    ssub  (class=string #args=3) Like sub but does no regexing. No characters are special.



.. _reference-dsl-strftime:

strftime
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    strftime  (class=time #args=2)  Formats seconds since the epoch as timestamp, e.g.
    	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%SZ") = "2015-08-28T13:33:21Z", and
    	strftime(1440768801.7,"%Y-%m-%dT%H:%M:%3SZ") = "2015-08-28T13:33:21.700Z".
    	Format strings are as in the C library (please see "man strftime" on your system),
    	with the Miller-specific addition of "%1S" through "%9S" which format the seconds
    	with 1 through 9 decimal places, respectively. ("%S" uses no decimal places.)
    	See also strftime_local.



.. _reference-dsl-string:

string
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    string  (class=conversion #args=1) Convert int/float/bool/string/array/map to string.



.. _reference-dsl-strip:

strip
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    strip  (class=string #args=1) Strip leading and trailing whitespace from string.



.. _reference-dsl-strlen:

strlen
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    strlen  (class=string #args=1) String length.



.. _reference-dsl-strptime:

strptime
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    strptime  (class=time #args=2) strptime: Parses timestamp as floating-point seconds since the epoch,
    	e.g. strptime("2015-08-28T13:33:21Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.000000,
    	and  strptime("2015-08-28T13:33:21.345Z","%Y-%m-%dT%H:%M:%SZ") = 1440768801.345000.
    	See also strptime_local.



.. _reference-dsl-sub:

sub
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    sub  (class=string #args=3) Example: '$name=sub($name, "old", "new")' (replace once).



.. _reference-dsl-substr:

substr
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    substr  (class=string #args=3) substr is an alias for substr0. See also substr1. Miller is generally 1-up
    with all array indices, but, this is a backward-compatibility issue with Miller 5 and below.
    Arrays are new in Miller 6; the substr function is older.



.. _reference-dsl-substr0:

substr0
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    substr0  (class=string #args=3) substr0(s,m,n) gives substring of s from 0-up position m to n
    inclusive. Negative indices -len .. -1 alias to 0 .. len-1.



.. _reference-dsl-substr1:

substr1
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    substr1  (class=string #args=3) substr1(s,m,n) gives substring of s from 1-up position m to n
    inclusive. Negative indices -len .. -1 alias to 1 .. len.



.. _reference-dsl-system:

system
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    system  (class=system #args=1) Run command string, yielding its stdout minus final carriage return.



.. _reference-dsl-systime:

systime
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    systime  (class=time #args=0) help string will go here



.. _reference-dsl-systimeint:

systimeint
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    systimeint  (class=time #args=0) help string will go here



.. _reference-dsl-tan:

tan
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    tan  (class=math #args=1) Trigonometric tangent.



.. _reference-dsl-tanh:

tanh
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    tanh  (class=math #args=1) Hyperbolic tangent.



.. _reference-dsl-tolower:

tolower
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    tolower  (class=string #args=1) Convert string to lowercase.



.. _reference-dsl-toupper:

toupper
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    toupper  (class=string #args=1) Convert string to uppercase.



.. _reference-dsl-truncate:

truncate
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    truncate  (class=string #args=2) Truncates string first argument to max length of int second argument.



.. _reference-dsl-typeof:

typeof
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    typeof  (class=typing #args=1) Convert argument to type of argument (e.g. "str"). For debug.



.. _reference-dsl-unflatten:

unflatten
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    unflatten  (class=maps/arrays #args=2) Reverses flatten. Example:
    unflatten({"a.b.c" : 4}, ".") is {"a": "b": { "c": 4 }}.
    Useful for nested JSON-like structures for non-JSON file formats like CSV.
    See also arrayify.



.. _reference-dsl-uptime:

uptime
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    uptime  (class=time #args=0) help string will go here



.. _reference-dsl-urand:

urand
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    urand  (class=math #args=0) Floating-point numbers uniformly distributed on the unit interval.
    Int-valued example: '$n=floor(20+urand()*11)'.



.. _reference-dsl-urand32:

urand32
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    urand32  (class=math #args=0) Integer uniformly distributed 0 and 2**32-1 inclusive.



.. _reference-dsl-urandint:

urandint
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    urandint  (class=math #args=2) Integer uniformly distributed between inclusive integer endpoints.



.. _reference-dsl-urandrange:

urandrange
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    urandrange  (class=math #args=2) Floating-point numbers uniformly distributed on the interval [a, b).



.. _reference-dsl-version:

version
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

.. code-block:: none

    version  (class=system #args=0) Returns the Miller version as a string.


