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
# Number formatting

## The --ofmt flag

The command-line option `--ofmt {format string}` is the global number format for all numeric fields.  Examples:

<pre class="pre-non-highlight-non-pair">
--ofmt %.9e --ofmt %.6f --ofmt %.0f
</pre>

These are just familiar `printf` formats -- as of Miller 6.0.0, supported options are those at
[https://pkg.go.dev/fmt](https://pkg.go.dev/fmt).  Additionally, if you use leading width (e.g.
`%18.12f`) then the output will contain embedded whitespace, which may not be what you want if you
pipe the output to something else, particularly CSV. I use Miller's pretty-print format (`mlr
--opprint`) to column-align numerical data.

<pre class="pre-highlight-in-pair">
<b>echo 'x=3.1,y=4.3' | mlr --ofmt '%8.3f' cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=   3.100,y=   4.300
</pre>

<pre class="pre-highlight-in-pair">
<b>echo 'x=3.1,y=4.3' | mlr --ofmt '%11.8e' cat</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=3.10000000e+00,y=4.30000000e+00
</pre>

## The format-values verb

To separately specify formatting for string, int, and float fields, you can use
the [`format-values`](reference-verbs.md#format-values) verb -- see that section for examples.

## The fmtnum and hexfmt functions

To apply formatting to a single field, you can also use
[`fmtnum`](reference-dsl-builtin-functions.md#fmtnum) function within `mlr
put`. For example:

<pre class="pre-highlight-in-pair">
<b>echo 'x=3.1,y=4.3' | mlr put '$z=fmtnum($x*$y,"%08f")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=3.1,y=4.3,z=13.330000
</pre>

<pre class="pre-highlight-in-pair">
<b>echo 'x=0xffff,y=0xff' | mlr put '$z=fmtnum(int($x*$y),"%08x")'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=0xffff,y=0xff,z=00feff01
</pre>

Input conversion from hexadecimal is done automatically on fields handled by `mlr put` and `mlr filter` as long as the field value begins with `0x`.  To apply output conversion to hexadecimal on a single column, you may use `fmtnum`, or the keystroke-saving [`hexfmt`](reference-dsl-builtin-functions.md#hexfmt) function. Example:

<pre class="pre-highlight-in-pair">
<b>echo 'x=0xffff,y=0xff' | mlr put '$z=$x*$y'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=0xffff,y=0xff,z=16711425
</pre>

<pre class="pre-highlight-in-pair">
<b>echo 'x=0xffff,y=0xff' | mlr put '$z=hexfmt($x*$y)'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x=0xffff,y=0xff,z=0xfeff01
</pre>

## Thousands separators

Some tools accept a `printf`-style apostrophe flag, e.g. `%'d`, to insert
locale-dependent thousands separators, formatting `12345` as `12,345`. Since
Miller's number formatting uses Go's [fmt](https://pkg.go.dev/fmt) package,
which doesn't support that flag, formats like `--ofmt "%'d"`,
`fmtnum($x, "%'d")`, and `format-values -i "%'d"` won't work. Likewise, the
well-known `sed`/`perl` recipe for inserting separators relies on regular-expression
lookaheads such as `(?=...)`, which [Go's regular-expression engine](reference-main-regular-expressions.md)
doesn't support.

You can, however, insert thousands separators with a small user-defined
[DSL function](reference-dsl-user-defined-functions.md) which applies
[`sub`](reference-dsl-builtin-functions.md#sub) until there's nothing left to
do. Digits are grouped from the left, so the decimal part (if any) is
untouched, and non-numeric values pass through unchanged:

<pre class="pre-highlight-in-pair">
<b>echo 'x=12345,y=1234567.8901,z=-987654321,s=hello' | mlr --opprint put '</b>
<b>  func commify(n) {</b>
<b>    s = string(n);</b>
<b>    while (s =~ "^(-?[0-9]+)([0-9]{3})") {</b>
<b>      s = sub(s, "^(-?[0-9]+)([0-9]{3})", "\1,\2");</b>
<b>    }</b>
<b>    return s;</b>
<b>  }</b>
<b>  for (k, v in $*) {</b>
<b>    $[k] = commify(v);</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
x      y              z            s
12,345 1,234,567.8901 -987,654,321 hello
</pre>

To use a separator other than comma -- say, a space or a period -- simply
change the `,` in the `"\1,\2"` replacement string.
