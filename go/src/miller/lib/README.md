This contains the implementation of the `lib.Mlrval` datatype which is used for record values, as well as expression/variable values in the Miller `put`/`filter` DSL.

The `lib.Mlrval` structure includes **string, int, float, boolean, void, absent, and error** types (not unlike PHP's `zval`) as well as type-conversion logic for various operators.

* Miller's `absent` type is like Javascript's `undefined` -- it's for times when there is no such key, as in a DSL expression `$out = $foo` when the input record is `$x=3,y=4` -- there is no `$foo` so `$foo` has `absent` type. Nothing is written to the `$out` field in this case. See also [here](http://johnkerl.org/miller/doc/reference.html#Null_data:_empty_and_absent) for more information.
* Miller's `void` type is like Javascript's `null` -- it's for times when there is a key with no value, as in `$out = $x` when the input record is `$x=,$y=4`. This is an overlap with `string` type, since a void value looks like an empty string. I've gone back and forth on this (including when I was writing the C implementation) -- whether to retain `void` as a distinct type from empty-string, or not. I ended up keeping it as it made the `Mlrval` logic easier to understand.
* Miller's `error` type is for things like doing type-uncoerced addition of strings. Data-dependent errors are intended to result in `(error)`-valued output, rather than crashing Miller. See also [here](http://johnkerl.org/miller/doc/reference.html#Data_types) for more information.
* Miller's number handling makes auto-overflow from int to float transparent, while preserving the possibility of 64-bit bitwise arithmetic.
  * Using `$a * $b` will auto-overflow to float. Using `$a .* $b` will stick with 64-bit integers (if `$a` and $b` are already 64-bit integers).
  * More generally:
    * Bitwise operators such as `|`, `&`, and `^` map ints to ints.
    * Auto-overflowing math operators `+`, `-`, etc. map ints to ints unless they overflow in which case float is produced.
    * Int-preserving math operators `.+`, `.-`, etc. map ints to ints even if they overflow.
  * This is different from JavaScript, which has only double-precision floats and thus no support for 64-bit numbers (note however that there is now [`BigInt`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/BigInt)).
  * This is also different from C and Go, wherein casts are necessary -- without which int arithmetic overflows.
  * See also [here](http://johnkerl.org/miller/doc/reference.html#Arithmetic) for the semantics of Miller arithmetic, which the `Mlrval` class implements.
