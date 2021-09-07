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
# Arrays

Miller data types are listed on the [Data types](reference-main-data-types.md)
page; here we focus specifically on arrays.

Arrays are supported [as of Miller 6](new-in-miller-6.md), and constitute one
of the major advantages of Miller 6.

## Array literals

Array literals are written in square brackets braces with integer indices. Array slots can be any [Miller data type](reference-main-data-types.md) (including other arrays, or maps).

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = [ "a", 1, "b", {"x": 2, "y": [3,4,5]}, 99, true];</b>
<b>    print x;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[
  "a",
  1,
  "b",
  {
    "x": 2,
    "y": [3, 4, 5]
  },
  99,
  true
]
</pre>

As with maps and argument-lists, trailing commas are supported:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = [</b>
<b>      "a",</b>
<b>      "b",</b>
<b>      "c",</b>
<b>    ];</b>
<b>    print x;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
["a", "b", "c"]
</pre>

Also note that several [built-in functions](reference-dsl-builtin-functions.md) operate on arrays and/or return arrays.

## 1-up indexing

The most important difference between Miller's arrays and arrays in other
languages is that indices start with 1, not 0. This is intentional.

1-up array indices may feel like a thing of the past, belonging to Fortran and
Matlab, say; or R and Julia as well, which are more modern.  But the overall
trend is decidedly toward 0-up. This means that if Miller does 1-up array
indices, it should do so for good reasons.

When arrays were introduced into [Miller 6](new-in-miller-6.md), it quickly became
clear that 1-up indexing is the right thing for Miller.  So many other things
are already 1-up in Miller, and always have been, mostly inherited from AWK:

* The `awk`-like [built-in variables](reference-dsl-variables.md#built-in-variables) `NF`, `NR`, and `FNR` are 1-up in Miller. So for idioms like `@records[NR] = $*` it's natural to index from 1; `@records[NR-1] = $*` would be error-prone and would result in frequent off-by-one errors.
* In particular, fields have always been indexed 1-up for [NIDX and DKVP formats](file-formats.md).
* [Regex captures](reference-main-regular-expressions.md) run from `"\1"` to `"\9"` (`"\0"` is the entire match substring).

## Negative-index aliasing

Imitating Python and other languages, you can use negative indices to read backward from the end of the array,
while positive indices read forward from the start. If an array has length `n` then `-n..-1` are aliases for `1..n`, respectively; 0 is never a valid array index in Miller.

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = [10, 20, 30, 40, 50];</b>
<b>    print x[1];</b>
<b>    print x[-1];</b>
<b>    print x[1:2];</b>
<b>    print x[-2:-1];</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
10
50
[10, 20]
[40, 50]
</pre>

## Auto-create results in maps

As noted on the [maps page](reference-main-maps.md), indexing any
as-yet-assigned local variable or out-of-stream variable results in
**auto-create** of that variable as a map variable:

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from example.csv put -q '</b>
<b>  # You can do this but you do not need to:</b>
<b>  # begin { @last_rates = {} }</b>
<b>  @last_rates[$shape] = $rate;</b>
<b>  end {</b>
<b>    dump @last_rates;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "triangle": 5.8240,
  "square": 8.2430,
  "circle": 8.3350
}
</pre>

*This also means that auto-create results in maps, not arrays, even if keys are integers.*
If you want to auto-extend an [array](reference-main-arrays.md), initialize it explicitly to `[]`.

<pre class="pre-highlight-in-pair">
<b>mlr --csv --from example.csv head -n 4 then put -q '</b>
<b>  begin {</b>
<b>    @my_array = [];</b>
<b>  }</b>
<b>  @my_array[NR] = $quantity;</b>
<b>  @my_map[NR] = $rate;</b>
<b>  end {</b>
<b>    dump</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
{
  "my_array": [43.6498, 79.2778, 13.8103, 77.5542],
  "my_map": {
    "1": 9.8870,
    "2": 0.0130,
    "3": 2.9010,
    "4": 7.4670
  }
}
</pre>

## Auto-extend and null-gaps

Once an array is initialized, it can be extended by assigning to indices beyond
its length.  If each write is one past the end of the array, the array will
grow by one. (Memory management, handled for you, is careful handled here in
Miller: not to worry, capacity is doubled so performance doesn't suffer a
rellocate on every single extend.)

This is important in Miller so you can do things like `@records[NR] = $*` with
a minimum of keystrokes without worrying about explicitly resizing arrays. In
particular, you can iteratively populate arrays as you read your data files,
without having to first know how many records they have.

However, if an array is written to more than one past its end, [values of type
JSON-null](reference-main-data-types.md) are used to fill in the gaps. These
are called **null-gaps**.

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    no_gaps = [];</b>
<b>    no_gaps[1] = "a";</b>
<b>    no_gaps[2] = "b";</b>
<b></b>
<b>    gaps = [];</b>
<b>    gaps[1] = "a";</b>
<b>    gaps[5] = "e";</b>
<b></b>
<b>    print no_gaps;</b>
<b>    print gaps;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
["a", "b"]
["a", null, null, null, "e"]
</pre>

## Unset as shift

Unsetting an array index results in shifting all higher-index elements down by one:

<pre class="pre-highlight-in-pair">
<b>mlr -n put '</b>
<b>  end {</b>
<b>    x = [ "a", "b", "c", "d", "e"];</b>
<b>    print x;</b>
<b>    unset x[2];</b>
<b>    print x;</b>
<b>  }</b>
<b>'</b>
</pre>
<pre class="pre-non-highlight-in-pair">
["a", "b", "c", "d", "e"]
["a", "c", "d", "e"]
</pre>

More generally, you can get shift and pop operations by unsetting indices 1 and -1:

<pre class="pre-highlight-in-pair">
<b>$ mlr repl -q</b>
</pre>
<pre class="pre-non-highlight-in-pair">
[mlr] x=[1,2,3,4,5]
[mlr] unset x[-1]
[mlr] x
[1, 2, 3, 4]
[mlr] unset x[-1]
[mlr] x
[1, 2, 3]
[mlr]
[mlr] x=[1,2,3,4,5]
[mlr] unset x[1]
[mlr] x
[2, 3, 4, 5]
[mlr] unset x[1]
[mlr] x
[3, 4, 5]
[mlr]
</pre>

## Looping

See [single-variable for-loops](reference-dsl-control-structures.md#single-variable-for-loops) and [key-value for-loops](reference-dsl-control-structures.md#key-value-for-loops).
