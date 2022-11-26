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
# Strings

## Essentials

Miller string literals are always written with double quotes, like `"abcde"`; single quotes
are not part of the grammar of [Miller's programming language](miller-programming-language.md).
Single quotes are used for wrapping `put`/`filter` statements, as in `mlr put '$b=$a.".suffix"' myfile.csv'`:
the single-quotes are consumed by the shell and Miller gets `$b=$a.".suffix"`. (See however the
[Miller on Windows page](miller-on-windows.md).)

A basic string operation is the `.` (concatenation) operator:

<pre class="pre-highlight-in-pair">
<b>mlr --c2p --from example.csv put '$output = $color . ":" . $shape'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
color  shape    flag  k  index quantity rate   output
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
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

Also see the [list of string-related built-in functions](reference-dsl-builtin-functions.md#string-functions).

## 1-up indexing

The most important difference between Miller's strings and strings in other
languages is that indices start with 1, not 0.  (The same is true for [Miller
arrays](reference-main-arrays.md).) This is intentional.

1-up indices may feel like a thing of the past, belonging to Fortran and Matlab,
say; or R and Julia as well, which are more modern.  But the overall trend is
decidedly toward 0-up. This means that if Miller does 1-up indices, it should
do so for good reasons.

Miller strings are indexed 1-up simply because Miller arrays are indexed 1-up.
See [this section](reference-main-arrays.md#1-up-indexing) for the reasoning.

Strings have been in Miller since the beginning, but they weren't accessible
using indices or slices until [Miller 6](new-in-miller-6.md). Also, the
[`substr`](reference-dsl-builtin-functions.md#substr) function predates Miller
6. This function was implemented to take 0-up indices.  When Miller 6 was
implemented, this became inconsistent.  As a result, there are
[`substr0`](reference-dsl-builtin-functions.md#substr0) and
[`substr1`](reference-dsl-builtin-functions.md#substr1) functions. For backward
compatibility with existing Miller scripts, `substr` is the same as `substr0`.
But users starting out with Miller 6 will probably want `substr1`.

## Negative-index aliasing

Imitating Python and other languages, you can use negative indices to read
backward from the end of the string, while positive indices read forward from
the start. If a string has length `n` then `-n..-1` are aliases for `1..n`,
respectively; 0 is never a valid string index in Miller.

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = "abcde";</b>
<b>    print x[1];</b>
<b>    print x[-1];</b>
<b>    print x[1:2];</b>
<b>    print x[-2:-1];</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a
e
ab
de
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

## Slicing

Miller supports slicing using `[lo:hi]` syntax.  Either or both of the indices
in a slice can be negatively aliased as described above.  Unlike in Python,
Miller string-slice indices are inclusive on both sides: `x[3:5]` means `x[3] . x[4] . x[5]`.

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = "abcde";</b>
<b>    print x[3:4];</b>
<b>    print x[:2];</b>
<b>    print x[3:];</b>
<b>    print x[1:-1];</b>
<b>    print x[2:-2];</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
cd
ab
cde
abcde
bcd
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

## Out-of-bounds indexing

Out-of-bounds index accesses are [errors](reference-main-data-types.md), but out-of-bounds slice
accesses result in trimming the indices, resulting in a short string or even the empty string.
(This behavior intentionally imitates Python.)

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = "abcde";</b>
<b>    print x[1];</b>
<b>    print x[5];</b>
<b>    print x[6];</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
a
e
(error)
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = "abcde";</b>
<b>    print "\"" . x[1:2] . "\"";</b>
<b>    print "\"" . x[1:6] . "\"";</b>
<b>    print "\"" . x[10:20] . "\"";</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
"ab"
"abcde"
""
Memory profile started.
Memory profile finished.
go tool pprof -http=:8080 foo-stream
</pre>

## Escape sequences for string literals

You can use the following backslash escapes for strings such as between the double quotes in contexts such as `mlr filter '$name =~ "..."'`, `mlr put '$name = $othername . "..."'`, `mlr put '$name = sub($name, "...", "...")`, etc.:

* `\a`: ASCII code 0x07 (alarm/bell)
* `\b`: ASCII code 0x08 (backspace)
* `\f`: ASCII code 0x0c (formfeed)
* `\n`: ASCII code 0x0a (LF/linefeed/newline)
* `\r`: ASCII code 0x0d (CR/carriage return)
* `\t`: ASCII code 0x09 (tab)
* `\v`: ASCII code 0x0b (vertical tab)
* `\\`: backslash
* `\"`: double quote
* `\123`: Octal 123, etc. for `\000` up to `\377`
* `\x7f`: Hexadecimal 7f, etc. for `\x00` up to `\xff`
* `\u2766`, `\U00010877:`: Unicode literals. For technical reasons, you must supply four hex digits after `\u` and eight hex digits after `\U`.

<pre class="pre-highlight-in-pair">
<b>mlr repl</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[mlr] "a\nb"
"a
b"

[mlr] "a\tb"
"a	b"

[mlr] "a\x62c"
"abc"

[mlr] "\u2766\U00010877"
"❦𐡷"
</pre>

See also [https://en.wikipedia.org/wiki/Escape_sequences_in_C](https://en.wikipedia.org/wiki/Escape_sequences_in_C).

These replacements apply only to strings you key in for the DSL expressions for `filter` and `put`: that is, if you type `\t` in a string literal for a `filter`/`put` expression, it will be turned into a tab character. If you want a backslash followed by a `t`, then please type `\\t`.

However, these replacements are done automatically only for string literals within DSL expressions -- they are not done automatically to fields within your data stream.  If you wish to make these replacements, you can do (for example) `mlr put '$field = gsub($field, "\\t", "\t")'`. If you need to make such a replacement for all fields in your data, you should probably use the system `sed` command instead. 
